package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers/reqcontext"
)

type AuditService struct {
	q *database.Queries
}

func NewAuditService(q *database.Queries) *AuditService {
	return &AuditService{q: q}
}

type AuditEntry struct {
	Action      string
	ActorType   string
	ActorID     string
	ActorEmail  string
	ProjectID   *uuid.UUID
	Environment *string
	TargetID    *string
	Status      string
	ErrMsg      *string
	Metadata    json.RawMessage
}

func (s *AuditService) Log(ctx context.Context, e AuditEntry) {
	if s == nil {
		return
	}

	reqID, ip, ua := reqcontext.GetRequestDetails(ctx)
	if reqID == "" {
		reqID = "unknown_request"
	}
	if e.ActorType == "" {
		e.ActorType = config.ActorTypeSystem
	}
	if e.ActorID == "" {
		e.ActorID = "system"
	}
	if e.ActorEmail == "" {
		e.ActorEmail = "system@envcrypt"
	}

	auditLog := buildAuditLog(e, reqID, ip, ua)

	if err := s.create(ctx, auditLog); err != nil {
		log.Printf("failed to write audit log: %v", err)
	}
}

func buildAuditLog(e AuditEntry, reqID string, ip *string, ua *string) *config.AuditLog {
	return &config.AuditLog{
		ID:           uuid.New(),
		Timestamp:    time.Now().UTC(),
		RequestID:    reqID,
		ActorType:    e.ActorType,
		ActorID:      e.ActorID,
		ActorEmail:   e.ActorEmail,
		Action:       e.Action,
		ProjectID:    e.ProjectID,
		Environment:  e.Environment,
		TargetID:     e.TargetID,
		IPAddress:    ip,
		UserAgent:    ua,
		Status:       e.Status,
		ErrorMessage: e.ErrMsg,
		Metadata:     e.Metadata,
	}
}

func (s *AuditService) create(ctx context.Context, auditLog *config.AuditLog) error {
	var projectID uuid.NullUUID
	if auditLog.ProjectID != nil {
		projectID = uuid.NullUUID{UUID: *auditLog.ProjectID, Valid: true}
	}

	var env sql.NullString
	if auditLog.Environment != nil {
		env = sql.NullString{String: *auditLog.Environment, Valid: true}
	}

	var targetID sql.NullString
	if auditLog.TargetID != nil {
		targetID = sql.NullString{String: *auditLog.TargetID, Valid: true}
	}

	var ipAddr pqtype.Inet
	if auditLog.IPAddress != nil && *auditLog.IPAddress != "" {
		ip := net.ParseIP(*auditLog.IPAddress)
		if ip != nil {
			var mask net.IPMask
			if ip.To4() != nil {
				mask = net.CIDRMask(32, 32)
			} else {
				mask = net.CIDRMask(128, 128)
			}
			ipAddr = pqtype.Inet{
				IPNet: net.IPNet{IP: ip, Mask: mask},
				Valid: true,
			}
		}
	}

	var userAgent sql.NullString
	if auditLog.UserAgent != nil {
		userAgent = sql.NullString{String: *auditLog.UserAgent, Valid: true}
	}

	var errMsg sql.NullString
	if auditLog.ErrorMessage != nil {
		errMsg = sql.NullString{String: *auditLog.ErrorMessage, Valid: true}
	}

	var meta pqtype.NullRawMessage
	if auditLog.Metadata != nil {
		meta = pqtype.NullRawMessage{RawMessage: json.RawMessage(auditLog.Metadata), Valid: true}
	}

	return s.q.CreateAuditLog(ctx, database.CreateAuditLogParams{
		ID:           auditLog.ID,
		Timestamp:    auditLog.Timestamp,
		RequestID:    auditLog.RequestID,
		ActorType:    auditLog.ActorType,
		ActorID:      auditLog.ActorID,
		ActorEmail:   auditLog.ActorEmail,
		Action:       auditLog.Action,
		ProjectID:    projectID,
		Environment:  env,
		TargetID:     targetID,
		IpAddress:    ipAddr,
		UserAgent:    userAgent,
		Status:       auditLog.Status,
		ErrorMessage: errMsg,
		Metadata:     meta,
	})
}
