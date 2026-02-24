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
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
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

func (s *AuditService) GetProjectAuditLogs(
	ctx context.Context,
	req config.ProjectAuditRequest,
) (config.ProjectAuditResponse, error) {

	sessionID, ok := ctx.Value("session_id").(uuid.UUID)
	if !ok {
		return config.ProjectAuditResponse{}, helpers.ErrUnauthorized("SESSION_MISSING", "User not authenticated", "")
	}

	session, err := s.q.GetSession(ctx, sessionID)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return config.ProjectAuditResponse{}, helpers.ErrUnauthorized("SESSION_EXPIRED", "Session is invalid or expired", "Please log in again")
		}
		return config.ProjectAuditResponse{}, err
	}

	var userID uuid.UUID
	if session.UserID.Valid {
		userID = session.UserID.UUID
	} else {
		return config.ProjectAuditResponse{}, helpers.ErrUnauthorized("USER_NOT_FOUND", "No user associated with session", "")
	}

	// Validate project membership
	_, err = s.q.GetUserProjectRole(ctx, database.GetUserProjectRoleParams{
		ProjectID: req.ProjectID,
		UserID:    userID,
		IsRevoked: false,
	})
	if err != nil {
		if dberrors.IsNoRows(err) {
			return config.ProjectAuditResponse{}, helpers.ErrForbidden("User doesn't have permission for this project", "")
		}
		return config.ProjectAuditResponse{}, err
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 50
	} else if limit > 200 {
		limit = 200
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	var actorEmail sql.NullString
	if req.ActorEmail != nil {
		actorEmail = sql.NullString{String: *req.ActorEmail, Valid: true}
	}
	var action sql.NullString
	if req.Action != nil {
		action = sql.NullString{String: *req.Action, Valid: true}
	}
	var status sql.NullString
	if req.Status != nil {
		status = sql.NullString{String: *req.Status, Valid: true}
	}
	var fromTime sql.NullTime
	if req.From != nil {
		fromTime = sql.NullTime{Time: *req.From, Valid: true}
	}
	var toTime sql.NullTime
	if req.To != nil {
		toTime = sql.NullTime{Time: *req.To, Valid: true}
	}

	projectIDNull := uuid.NullUUID{UUID: req.ProjectID, Valid: true}

	logs, err := s.q.GetProjectAuditLogsPaginated(ctx, database.GetProjectAuditLogsPaginatedParams{
		ProjectID:  projectIDNull,
		ActorEmail: actorEmail,
		Action:     action,
		Status:     status,
		FromTime:   fromTime,
		ToTime:     toTime,
		LimitVal:   limit,
		OffsetVal:  offset,
	})
	if err != nil {
		return config.ProjectAuditResponse{}, err
	}

	total, err := s.q.CountProjectAuditLogs(ctx, database.CountProjectAuditLogsParams{
		ProjectID:  projectIDNull,
		ActorEmail: actorEmail,
		Action:     action,
		Status:     status,
		FromTime:   fromTime,
		ToTime:     toTime,
	})
	if err != nil {
		return config.ProjectAuditResponse{}, err
	}

	resp := config.ProjectAuditResponse{
		Logs: make([]config.AuditLog, len(logs)),
	}
	for i, log := range logs {
		var ipAddr *string
		if log.IpAddress.Valid {
			ipStr := log.IpAddress.IPNet.IP.String()
			ipAddr = &ipStr
		}
		var errStr *string
		if log.ErrorMessage.Valid {
			errStr = &log.ErrorMessage.String
		}
		var envStr *string
		if log.Environment.Valid {
			envStr = &log.Environment.String
		}
		var targetStr *string
		if log.TargetID.Valid {
			targetStr = &log.TargetID.String
		}
		var uaStr *string
		if log.UserAgent.Valid {
			uaStr = &log.UserAgent.String
		}
		var meta json.RawMessage
		if log.Metadata.Valid {
			meta = log.Metadata.RawMessage
		}
		var projID *uuid.UUID
		if log.ProjectID.Valid {
			projID = &log.ProjectID.UUID
		}
		resp.Logs[i] = config.AuditLog{
			ID:           log.ID,
			Timestamp:    log.Timestamp,
			RequestID:    log.RequestID,
			ActorType:    log.ActorType,
			ActorID:      log.ActorID,
			ActorEmail:   log.ActorEmail,
			Action:       log.Action,
			ProjectID:    projID,
			Environment:  envStr,
			TargetID:     targetStr,
			IPAddress:    ipAddr,
			UserAgent:    uaStr,
			Status:       log.Status,
			ErrorMessage: errStr,
			Metadata:     meta,
		}
	}
	resp.Pagination.Limit = limit
	resp.Pagination.Offset = offset
	resp.Pagination.Total = total

	return resp, nil
}
