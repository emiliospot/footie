package gorm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormdb "gorm.io/driver/sqlite"
	gormlib "gorm.io/gorm"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gormlib.DB {
	db, err := gormlib.Open(gormdb.Open(":memory:"), &gormlib.Config{})
	require.NoError(t, err, "Failed to setup test database")

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err, "Failed to migrate database")

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		user    *models.User
		name    string
		wantErr bool
	}{
		{
			name: "valid user",
			user: &models.User{
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				FirstName:    "John",
				LastName:     "Doe",
				Role:         "user",
				IsActive:     true,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &models.User{
				Email:        "test@example.com",
				PasswordHash: "hashed_password2",
				FirstName:    "Jane",
				LastName:     "Doe",
				Role:         "user",
				IsActive:     true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.user.ID)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	// Create test user
	user := &models.User{
		Email:        "findme@example.com",
		PasswordHash: "hash",
		FirstName:    "Find",
		LastName:     "Me",
		Role:         "user",
		IsActive:     true,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      user.ID,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			found, err := repo.FindByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, found)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, found)
				assert.Equal(t, tt.id, found.ID)
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	email := "unique@example.com"
	user := &models.User{
		Email:        email,
		PasswordHash: "hash",
		FirstName:    "Test",
		LastName:     "User",
		Role:         "user",
		IsActive:     true,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	found, err := repo.FindByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, email, found.Email)

	// Test non-existent email
	notFound, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Nil(t, notFound)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		Email:        "update@example.com",
		PasswordHash: "hash",
		FirstName:    "Old",
		LastName:     "Name",
		Role:         "user",
		IsActive:     true,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Update user
	user.FirstName = "New"
	user.LastName = "Name"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "New", updated.FirstName)
	assert.Equal(t, "Name", updated.LastName)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		Email:        "delete@example.com",
		PasswordHash: "hash",
		FirstName:    "Delete",
		LastName:     "Me",
		Role:         "user",
		IsActive:     true,
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Delete user
	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	// Create multiple users
	for i := 0; i < 5; i++ {
		user := &models.User{
			Email:        fmt.Sprintf("user%d@example.com", i),
			PasswordHash: "hash",
			FirstName:    fmt.Sprintf("User%d", i),
			LastName:     "Test",
			Role:         "user",
			IsActive:     true,
		}
		err := repo.Create(ctx, user)
		require.NoError(t, err)
	}

	// Test pagination
	users, total, err := repo.List(ctx, 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 3)

	// Test second page
	users, total, err = repo.List(ctx, 3, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, users, 2)
}

// Benchmark tests.
func BenchmarkUserRepository_Create(b *testing.B) {
	db := setupTestDB(&testing.T{})
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &models.User{
			Email:        fmt.Sprintf("bench%d@example.com", i),
			PasswordHash: "hash",
			FirstName:    "Benchmark",
			LastName:     "User",
			Role:         "user",
			IsActive:     true,
		}
		_ = repo.Create(ctx, user)
	}
}

func BenchmarkUserRepository_FindByID(b *testing.B) {
	db := setupTestDB(&testing.T{})
	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		Email:        "bench@example.com",
		PasswordHash: "hash",
		FirstName:    "Benchmark",
		LastName:     "User",
		Role:         "user",
		IsActive:     true,
	}
	_ = repo.Create(ctx, user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.FindByID(ctx, user.ID)
	}
}
