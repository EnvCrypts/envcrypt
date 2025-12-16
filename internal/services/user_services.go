package services

import (
	"context"
	"encoding/json"
	"errors"

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

func (s *UserService) Create(ctx context.Context, createBody config.CreateRequestBody) error {

	passwordHash, err := auth.HashPassword(createBody.Password)
	if err != nil {
		return err
	}

	paramsJson, err := json.Marshal(passwordHash.Argon2idParam)
	if err != nil {
		return err
	}

	_, err = s.q.CreateUser(ctx, database.CreateUserParams{
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
			return errors.New("user already exists")
		}
		return err
	}

	return nil
}
