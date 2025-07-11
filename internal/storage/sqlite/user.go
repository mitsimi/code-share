package sqlite

import (
	"context"
	"database/sql"

	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"
)

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.UserCreation) (*domain.User, error) {
	logger.Log.Info("Creating new user", zap.String("username", user.Username), zap.String("email", user.Email), zap.String("password_hash", user.PasswordHash))
	newUser, err := r.q.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, repository.ErrAlreadyExists
		}
		return nil, repository.WrapError(err, "failed to create user")
	}
	return domain.ToDomainUser(newUser), nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := r.q.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get user by ID")
	}
	return domain.ToDomainUser(user), nil

}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get user by username")
	}
	return domain.ToDomainUser(user), nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to get user by email")
	}
	return domain.ToDomainUser(user), nil
}
func (r *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	updatedUser, err := r.q.UpdateUserInfo(ctx, db.UpdateUserInfoParams{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		return nil, repository.WrapError(err, "failed to update user")
	}
	return domain.ToDomainUser(updatedUser), nil
}

func (r *UserRepository) UpdateAvatar(ctx context.Context, userID, avatarURL string) error {
	_, err := r.q.UpdateUserAvatar(ctx, db.UpdateUserAvatarParams{
		ID:     userID,
		Avatar: sql.NullString{String: avatarURL, Valid: avatarURL != ""},
	})
	if err != nil {
		return repository.WrapError(err, "failed to update user avatar")
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	err := r.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return repository.WrapError(err, "failed to update user password")
	}
	return nil
}
