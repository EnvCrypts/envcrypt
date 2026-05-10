package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	"github.com/vijayvenkatj/envcrypt/internal/helpers/auth"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type UserService struct {
	q     *database.Queries
	db    *sql.DB
	audit *AuditService
}

func NewUserService(q *database.Queries, db *sql.DB) *UserService {
	return &UserService{q: q, db: db}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]database.User, error) {
	return s.q.GetUsers(ctx)
}

func (s *UserService) Create(ctx context.Context, createBody config.CreateRequestBody) (*config.UserBody, error) {

	passwordHash, err := auth.HashPassword(createBody.Password)
	if err != nil {
		return nil, errors.InternalMessage("Failed to hash password", err)
	}

	paramsJson, err := json.Marshal(passwordHash.Argon2idParam)
	if err != nil {
		return nil, errors.InternalMessage("Failed to serialize password parameters", err)
	}

	user, err := s.q.CreateUser(ctx, database.CreateUserParams{
		ID:                          uuid.New(),
		Email:                       createBody.Email,
		PasswordHash:                passwordHash.Hash,
		PasswordSalt:                passwordHash.Salt,
		UserPublicKey:               createBody.PublicKey,
		EncryptedUserPrivateKey:     createBody.EncryptedUserPrivateKey,
		PrivateKeySalt:              createBody.PrivateKeySalt,
		PrivateKeyNonce:             createBody.PrivateKeyNonce,
		RecoveryEncryptedPrivateKey: createBody.RecoveryPrivateKey,
		RecoveryKdfSalt:             createBody.RecoverySalt,
		RecoveryNonce:               createBody.RecoveryNonce,
		ArgonParams:                 paramsJson,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionRegister, ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: createBody.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) {
			return nil, errors.Conflict("User already exists", "Try a different email address")
		}
		return nil, errors.Internal(err)
	}

	var argonParams auth.Argon2idParams
	err = json.Unmarshal(user.ArgonParams, &argonParams)
	if err != nil {
		return nil, errors.InternalMessage("Failed to parse password parameters", err)
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionRegister, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: user.Email, Status: config.StatusSuccess})

	return &config.UserBody{
		Id:                      user.ID,
		Email:                   user.Email,
		PublicKey:               user.UserPublicKey,
		EncryptedUserPrivateKey: user.EncryptedUserPrivateKey,
		PrivateKeySalt:          user.PrivateKeySalt,
		PrivateKeyNonce:         user.PrivateKeyNonce,
		ArgonParams:             argonParams,
	}, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*config.UserBody, error) {

	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return nil, errors.Unauthorized("INVALID_CREDENTIALS", "Invalid email or password", "Check your credentials and try again")
		}
		return nil, errors.Internal(err)
	}

	var argonParams auth.Argon2idParams
	err = json.Unmarshal(user.ArgonParams, &argonParams)
	if err != nil {
		return nil, errors.InternalMessage("Failed to parse password parameters", err)
	}

	stored := auth.PasswordHash{
		Hash:          user.PasswordHash,
		Salt:          user.PasswordSalt,
		Argon2idParam: argonParams,
	}

	if auth.VerifyPassword(password, &stored) == false {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("invalid password")})
		return nil, errors.Unauthorized("INVALID_CREDENTIALS", "Invalid email or password", "Check your credentials and try again")
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: email, Status: config.StatusSuccess})

	return &config.UserBody{
		Id:                      user.ID,
		Email:                   user.Email,
		PublicKey:               user.UserPublicKey,
		EncryptedUserPrivateKey: user.EncryptedUserPrivateKey,
		PrivateKeySalt:          user.PrivateKeySalt,
		PrivateKeyNonce:         user.PrivateKeyNonce,
		ArgonParams:             argonParams,
	}, nil
}

func (s *UserService) GetUserPublicKey(ctx context.Context, email string) (uuid.UUID, []byte, error) {

	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if dberrors.IsNoRows(err) {
			return uuid.Nil, nil, errors.NotFound("User", "Check the email address")
		}
		return uuid.Nil, nil, errors.Internal(err)
	}

	return user.ID, user.UserPublicKey, nil
}

func (s *UserService) Logout(ctx context.Context, userId uuid.UUID) error {

	user, err := s.q.GetUserByID(ctx, userId)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return errors.NotFound("User", "")
		}
		return errors.Internal(err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.InternalMessage("Unable to begin logout transaction", err)
	}
	defer tx.Rollback()

	txQ := s.q.WithTx(tx)

	err = txQ.DeleteRefreshTokens(ctx, userId)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), ActorEmail: user.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return errors.Internal(err)
	}
	err = txQ.DeleteUserAccessTokens(ctx, uuid.NullUUID{UUID: userId, Valid: true})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), ActorEmail: user.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return errors.Internal(err)
	}

	if err = tx.Commit(); err != nil {
		return errors.InternalMessage("Unable to commit logout transaction", err)
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), ActorEmail: user.Email, Status: config.StatusSuccess})
	return nil
}

func (s *UserService) RecoveryInit(ctx context.Context, email string) (*config.RecoveryInitResponseBody, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: "Recovery Init", ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return nil, errors.NotFound("User", "Check the email address")
		}
		return nil, errors.Internal(err)
	}

	s.audit.Log(ctx, AuditEntry{Action: "Recovery Init", ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: email, Status: config.StatusSuccess})
	return &config.RecoveryInitResponseBody{
		RecoveryPrivateKey: user.RecoveryEncryptedPrivateKey,
		RecoverySalt:       user.RecoveryKdfSalt,
		RecoveryNonce:      user.RecoveryNonce,
	}, nil
}

func (s *UserService) RecoveryComplete(ctx context.Context, req config.RecoveryCompleteRequestBody) error {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: "Recovery Complete", ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: req.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return errors.NotFound("User", "Check the email address")
		}
		return errors.Internal(err)
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return errors.InternalMessage("Failed to hash password", err)
	}

	paramsJson, err := json.Marshal(passwordHash.Argon2idParam)
	if err != nil {
		return errors.InternalMessage("Failed to serialize password parameters", err)
	}

	_, err = s.q.UpdateUserCredentials(ctx, database.UpdateUserCredentialsParams{
		Email:                   req.Email,
		PasswordHash:            passwordHash.Hash,
		PasswordSalt:            passwordHash.Salt,
		ArgonParams:             paramsJson,
		EncryptedUserPrivateKey: req.EncryptedUserPrivateKey,
		PrivateKeyNonce:         req.PrivateKeyNonce,
		PrivateKeySalt:          req.PrivateKeySalt,
	})

	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: "Recovery Complete", ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: req.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return errors.Internal(err)
	}

	s.audit.Log(ctx, AuditEntry{Action: "Recovery Complete", ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: req.Email, Status: config.StatusSuccess})
	return nil
}
