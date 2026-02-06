package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/envcrypt/database"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers/auth"
	dberrors "github.com/vijayvenkatj/envcrypt/internal/helpers/db"
)

type UserService struct {
	q *database.Queries
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
		if dberrors.IsNoRows(err) {
			return nil, errors.New("user not found")
		}
		log.Print(err.Error())
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
		return nil, errors.New("invalid password")
	}

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
