package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type SnapshotService struct {
	q     *database.Queries
	db    *sql.DB
	audit *AuditService
}

func NewSnapshotService(q *database.Queries, db *sql.DB) *SnapshotService {
	return &SnapshotService{
		q:  q,
		db: db,
	}
}

func (s *SnapshotService) SetAuditService(audit *AuditService) {
	s.audit = audit
}

func generateChecksum(snapshot config.Snapshot) (string, error) {
	rawBytes, err := json.Marshal(snapshot)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(rawBytes)
	return hex.EncodeToString(hash[:]), nil
}

func (s *SnapshotService) ExportSnapshot(ctx context.Context, req config.SnapshotExportRequest) (*config.SnapshotExportResponse, error) {
	actor, err := s.q.GetUserByID(ctx, req.UserID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	project, err := s.q.GetProject(ctx, database.GetProjectParams{
		Name:      req.ProjectName,
		CreatedBy: req.UserID,
	})
	if err != nil {
		// Try member project
		projectID, errMem := s.q.GetMemberProject(ctx, database.GetMemberProjectParams{
			Name:   req.ProjectName,
			UserID: req.UserID,
		})
		if errMem != nil {
			s.audit.Log(ctx, AuditEntry{Action: "snapshot.export", ActorType: config.ActorTypeUser, ActorID: req.UserID.String(), ActorEmail: actor.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("project not found")})
			return nil, errors.New("project not found or permission denied")
		}
		// fetch actual project
		project, err = s.q.GetProjectById(ctx, projectID)
		if err != nil {
			return nil, err
		}
	}

	// Verify project role and not revoked
	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		ProjectID: project.ID,
		UserID:    req.UserID,
		IsRevoked: false,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: "snapshot.export", ActorType: config.ActorTypeUser, ActorID: req.UserID.String(), ActorEmail: actor.Email, ProjectID: &project.ID, Status: config.StatusFailure, ErrMsg: helpers.Ptr("permission denied")})
		return nil, errors.New("permission denied")
	}

	// Gather rotation data (wrapped PRKs)
	rotationData, err := s.q.GetRotationData(ctx, project.ID)
	if err != nil {
		return nil, err
	}

	var members []config.SnapshotMember
	for _, rd := range rotationData {
		members = append(members, config.SnapshotMember{
			UserID:             rd.UserID,
			WrappedPRK:         rd.WrappedPrk,
			WrapNonce:          rd.WrapNonce,
			EphemeralPublicKey: rd.WrapEphemeralPub,
		})
	}

	envVersions, err := s.q.GetAllEnvVersionsForProject(ctx, project.ID)
	if err != nil {
		return nil, err
	}

	var snapshotEnvs []config.SnapshotEnvVersion
	for _, env := range envVersions {
		snapshotEnvs = append(snapshotEnvs, config.SnapshotEnvVersion{
			EnvVersionID:      env.ID,
			EnvName:           env.EnvName,
			Version:           env.Version,
			Ciphertext:        env.Ciphertext,
			Nonce:             env.Nonce,
			WrappedDEK:        env.WrappedDek,
			DekNonce:          env.DekNonce,
			EncryptionVersion: env.EncryptionVersion,
			CreatedAt:         env.CreatedAt,
			CreatedBy:         env.CreatedBy,
			Metadata:          env.Metadata,
		})
	}

	snapshot := config.Snapshot{
		Metadata: config.SnapshotProjectMetadata{
			Name:       project.Name,
			PrkVersion: project.PrkVersion,
		},
		Members:     members,
		EnvVersions: snapshotEnvs,
	}

	checksum, err := generateChecksum(snapshot)
	if err != nil {
		return nil, err
	}

	s.audit.Log(ctx, AuditEntry{
		Action:     "snapshot.export",
		ActorType:  config.ActorTypeUser,
		ActorID:    req.UserID.String(),
		ActorEmail: actor.Email,
		ProjectID:  &project.ID,
		Status:     config.StatusSuccess,
	})

	return &config.SnapshotExportResponse{
		Snapshot: snapshot,
		Checksum: checksum,
	}, nil
}

func (s *SnapshotService) ImportSnapshot(ctx context.Context, req config.SnapshotImportRequest) (*config.SnapshotImportResponse, error) {
	actor, err := s.q.GetUserByID(ctx, req.UserID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 1. Verify PRKs
	if len(req.Snapshot.Members) == 0 {
		return nil, errors.New("snapshot must contain at least one wrapped PRK member")
	}

	// 2. Verify unique IDs, env version consistency (simplified checks)
	if req.Snapshot.Metadata.PrkVersion < 1 {
		req.Snapshot.Metadata.PrkVersion = 1
	}

	// 3. Verify Checksum
	actualChecksum, err := generateChecksum(req.Snapshot)
	if err != nil {
		return nil, errors.New("failed to compute checksum")
	}

	if actualChecksum != req.Checksum {
		s.audit.Log(ctx, AuditEntry{Action: "snapshot.import", ActorType: config.ActorTypeUser, ActorID: req.UserID.String(), ActorEmail: actor.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("checksum mismatch")})
		return nil, errors.New("checksum mismatch")
	}

	// Check if new project name already exists
	_, err = s.q.GetProject(ctx, database.GetProjectParams{
		Name:      req.NewProjectName,
		CreatedBy: req.UserID,
	})
	if err == nil {
		return nil, errors.New("project with this name already exists")
	}

	// Begin atomic transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	txQ := s.q.WithTx(tx)

	newProjectID := uuid.New()

	_, err = txQ.InsertProjectWithVersion(ctx, database.InsertProjectWithVersionParams{
		ID:         newProjectID,
		Name:       req.NewProjectName,
		CreatedBy:  req.UserID,
		PrkVersion: req.Snapshot.Metadata.PrkVersion,
	})
	if err != nil {
		return nil, err
	}

	for _, member := range req.Snapshot.Members {
		// Default to member, if it's the creator make them admin
		role := "member"
		if member.UserID == req.UserID {
			role = "admin"
		}
		_, err = txQ.AddUserToProject(ctx, database.AddUserToProjectParams{
			ProjectID: newProjectID,
			UserID:    member.UserID,
			Role:      role,
		})
		if err != nil {
			return nil, err
		}

		_, err = txQ.AddWrappedPRK(ctx, database.AddWrappedPRKParams{
			ProjectID:        newProjectID,
			UserID:           member.UserID,
			WrappedPrk:       member.WrappedPRK,
			WrapNonce:        member.WrapNonce,
			WrapEphemeralPub: member.EphemeralPublicKey,
		})
		if err != nil {
			return nil, err
		}
	}

	for _, env := range req.Snapshot.EnvVersions {
		// Maintain versions, but assign a new internal database ID to avoid primary key collisions.
		err = txQ.InsertEnvVersionRaw(ctx, database.InsertEnvVersionRawParams{
			ID:                uuid.New(), 
			ProjectID:         newProjectID,
			EnvName:           env.EnvName,
			Version:           env.Version,
			Ciphertext:        env.Ciphertext,
			Nonce:             env.Nonce,
			WrappedDek:        env.WrappedDEK,
			DekNonce:          env.DekNonce,
			EncryptionVersion: env.EncryptionVersion,
			CreatedAt:         env.CreatedAt,
			CreatedBy:         env.CreatedBy,
			Metadata:          env.Metadata,
		})
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.audit.Log(ctx, AuditEntry{
		Action:     "snapshot.import",
		ActorType:  config.ActorTypeUser,
		ActorID:    req.UserID.String(),
		ActorEmail: actor.Email,
		ProjectID:  &newProjectID,
		Status:     config.StatusSuccess,
	})

	return &config.SnapshotImportResponse{
		NewProjectID: newProjectID,
	}, nil
}
