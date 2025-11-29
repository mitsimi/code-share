package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
	"mitsimi.dev/codeShare/internal/repository"
)

func setupUserTestDB(t *testing.T) (*sql.DB, *UserRepository) {
	err := logger.Init(logger.Config{
		Environment: "development",
		Level:       "debug",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	storage, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	repo := NewUserRepository(storage.DB())
	return storage.DB(), repo
}

func TestUserRepository_Create(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		user := &domain.UserCreation{
			ID:           "test-id",
			Username:     "test-user",
			Email:        "test@example.com",
			PasswordHash: "test-hash",
		}

		createdUser, err := repo.Create(context.Background(), user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.ID, createdUser.ID)
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("duplicate username", func(t *testing.T) {
		user := &domain.UserCreation{
			ID:           "test-id-2",
			Username:     "test-user", // Same username as before
			Email:        "test2@example.com",
			PasswordHash: "test-hash",
		}

		_, err := repo.Create(context.Background(), user)
		assert.ErrorIs(t, err, repository.ErrAlreadyExists)
	})

	t.Run("duplicate email", func(t *testing.T) {
		user := &domain.UserCreation{
			ID:           "test-id-3",
			Username:     "test-user-3",
			Email:        "test@example.com", // Same email as before
			PasswordHash: "test-hash",
		}

		_, err := repo.Create(context.Background(), user)
		assert.ErrorIs(t, err, repository.ErrAlreadyExists)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	// Seed the database with a user
	user := &domain.UserCreation{
		ID:           "test-id",
		Username:     "test-user",
		Email:        "test@example.com",
		PasswordHash: "test-hash",
	}
	_, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		foundUser, err := repo.GetByID(context.Background(), "test-id")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByID(context.Background(), "non-existent-id")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	// Seed the database with a user
	user := &domain.UserCreation{
		ID:           "test-id",
		Username:     "test-user",
		Email:        "test@example.com",
		PasswordHash: "test-hash",
	}
	_, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		foundUser, err := repo.GetByUsername(context.Background(), "test-user")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByUsername(context.Background(), "non-existent-user")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	// Seed the database with a user
	user := &domain.UserCreation{
		ID:           "test-id",
		Username:     "test-user",
		Email:        "test@example.com",
		PasswordHash: "test-hash",
	}
	_, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		foundUser, err := repo.GetByEmail(context.Background(), "test@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetByEmail(context.Background(), "non-existent-email@example.com")
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestUserRepository_Update(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	// Seed the database with a user
	user := &domain.UserCreation{
		ID:           "test-id",
		Username:     "test-user",
		Email:        "test@example.com",
		PasswordHash: "test-hash",
	}
	createdUser, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		updatedUser := &domain.User{
			ID:       createdUser.ID,
			Username: "updated-user",
			Email:    "updated@example.com",
		}
		result, err := repo.Update(context.Background(), updatedUser)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedUser.Username, result.Username)
		assert.Equal(t, updatedUser.Email, result.Email)

		// Verify the update in the database
		foundUser, err := repo.GetByID(context.Background(), createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, updatedUser.Username, foundUser.Username)
		assert.Equal(t, updatedUser.Email, foundUser.Email)
	})
}

func TestUserRepository_UpdateAvatar(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	// Seed the database with a user
	user := &domain.UserCreation{
		ID:           "test-id",
		Username:     "test-user",
		Email:        "test@example.com",
		PasswordHash: "test-hash",
	}
	createdUser, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		avatarURL := "https://example.com/avatar.png"
		err := repo.UpdateAvatar(context.Background(), createdUser.ID, avatarURL)
		assert.NoError(t, err)

		// Verify the update in the database
		foundUser, err := repo.GetByID(context.Background(), createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, avatarURL, *foundUser.Avatar)
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, repo := setupUserTestDB(t)
	defer db.Close()

	// Seed the database with a user
	user := &domain.UserCreation{
		ID:           "test-id",
		Username:     "test-user",
		Email:        "test@example.com",
		PasswordHash: "old-password-hash",
	}
	createdUser, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		newPasswordHash := "new-password-hash"
		err := repo.UpdatePassword(context.Background(), createdUser.ID, newPasswordHash)
		assert.NoError(t, err)

		// Verify the update in the database
		var passwordHash string
		err = db.QueryRow("SELECT password_hash FROM users WHERE id = ?", createdUser.ID).Scan(&passwordHash)
		assert.NoError(t, err)
		assert.Equal(t, newPasswordHash, passwordHash)
	})
}
