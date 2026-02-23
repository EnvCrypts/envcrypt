package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
	"github.com/vijayvenkatj/envcrypt/internal/helpers/auth"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type UserService struct {
	q     *database.Queries
	audit *AuditService
}

func NewUserService(q *database.Queries) *UserService {
	return &UserService{q: q}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]database.User, error) {
	return s.q.GetUsers(ctx)
}

func (s *UserService) Create(ctx context.Context, createBody config.CreateRequestBody) (*config.UserBody, error) {

	passwordHash, err := auth.HashPassword(createBody.Password)
	if err != nil {
		return nil, err
	}

	paramsJson, err := json.Marshal(passwordHash.Argon2idParam)
	if err != nil {
		return nil, err
	}

	user, err := s.q.CreateUser(ctx, database.CreateUserParams{
		Email:                   createBody.Email,
		PasswordHash:            passwordHash.Hash,
		PasswordSalt:            passwordHash.Salt,
		UserPublicKey:           createBody.PublicKey,
		EncryptedUserPrivateKey: createBody.EncryptedUserPrivateKey,
		PrivateKeySalt:          createBody.PrivateKeySalt,
		PrivateKeyNonce:         createBody.PrivateKeyNonce,
		ArgonParams:             paramsJson,
	})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionRegister, ActorType: config.ActorTypeUser, ActorID: "unknown", ActorEmail: createBody.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		if dberrors.IsUniqueViolation(err) {
			return nil, errors.New("user already exists")
		}
		return nil, err
	}

	var argonParams auth.Argon2idParams
	err = json.Unmarshal(user.ArgonParams, &argonParams)
	if err != nil {
		return nil, err
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
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	var argonParams auth.Argon2idParams
	err = json.Unmarshal(user.ArgonParams, &argonParams)
	if err != nil {
		return nil, err
	}

	stored := auth.PasswordHash{
		Hash:          user.PasswordHash,
		Salt:          user.PasswordSalt,
		Argon2idParam: argonParams,
	}

	if auth.VerifyPassword(password, &stored) == false {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogin, ActorType: config.ActorTypeUser, ActorID: user.ID.String(), ActorEmail: email, Status: config.StatusFailure, ErrMsg: helpers.Ptr("invalid password")})
		return nil, errors.New("invalid password")
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
			return uuid.Nil, nil, errors.New("user not found")
		}
		return uuid.Nil, nil, err
	}

	return user.ID, user.UserPublicKey, nil
}

func (s *UserService) Logout(ctx context.Context, userId uuid.UUID) error {

	user, err := s.q.GetUserByID(ctx, userId)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), Status: config.StatusFailure, ErrMsg: helpers.Ptr("user not found")})
		if dberrors.IsNoRows(err) {
			return errors.New("user not found")
		}
		return err
	}

	err = s.q.DeleteRefreshTokens(ctx, userId)
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), ActorEmail: user.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return err
	}
	err = s.q.DeleteUserAccessTokens(ctx, uuid.NullUUID{UUID: userId, Valid: true})
	if err != nil {
		s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), ActorEmail: user.Email, Status: config.StatusFailure, ErrMsg: helpers.Ptr(err.Error())})
		return err
	}

	s.audit.Log(ctx, AuditEntry{Action: config.ActionLogout, ActorType: config.ActorTypeUser, ActorID: userId.String(), ActorEmail: user.Email, Status: config.StatusSuccess})
	return nil
}
