//go:build integration
// +build integration

package gorm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/emiliospot/footie/api/internal/domain/models"
	"github.com/emiliospot/footie/api/internal/repository/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	gormdb "gorm.io/driver/postgres"
	gormlib "gorm.io/gorm"
)

// setupPostgresContainer creates a real Postgres container for integration tests.
func setupPostgresContainer(t *testing.T) (*gormlib.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "footie_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to start Postgres container")

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=footie_test sslmode=disable",
		host, port.Port())

	db, err := gormlib.Open(gormdb.Open(dsn), &gormlib.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Run migrations
	err = db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Player{},
		&models.Match{},
		&models.MatchEvent{},
		&models.PlayerStatistics{},
		&models.TeamStatistics{},
	)
	require.NoError(t, err, "Failed to run migrations")

	cleanup := func() {
		container.Terminate(ctx)
	}

	return db, cleanup
}

func TestIntegration_UserRepository_WithRealPostgres(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := gorm.NewUserRepository(db)
	ctx := context.Background()

	// Test full CRUD cycle
	user := &models.User{
		Email:        "integration@example.com",
		PasswordHash: "hash",
		FirstName:    "Integration",
		LastName:     "Test",
		Role:         "user",
		IsActive:     true,
	}

	// Create
	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// Read
	found, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)

	// Update
	user.FirstName = "Updated"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	updated, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.FirstName)

	// Delete
	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	_, err = repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
}

func TestIntegration_RepositoryManager_Transactions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repoManager := gorm.NewRepositoryManager(db)
	ctx := context.Background()

	// Test transaction rollback
	t.Run("rollback on error", func(t *testing.T) {
		txManager, err := repoManager.BeginTx(ctx)
		require.NoError(t, err)

		user := &models.User{
			Email:        "transaction@example.com",
			PasswordHash: "hash",
			FirstName:    "Trans",
			LastName:     "Action",
			Role:         "user",
			IsActive:     true,
		}

		err = txManager.User().Create(ctx, user)
		assert.NoError(t, err)

		// Rollback
		err = txManager.Rollback()
		assert.NoError(t, err)

		// Verify user was not created
		_, err = repoManager.User().FindByID(ctx, user.ID)
		assert.Error(t, err)
	})

	// Test transaction commit
	t.Run("commit transaction", func(t *testing.T) {
		txManager, err := repoManager.BeginTx(ctx)
		require.NoError(t, err)

		user := &models.User{
			Email:        "committed@example.com",
			PasswordHash: "hash",
			FirstName:    "Commit",
			LastName:     "Test",
			Role:         "user",
			IsActive:     true,
		}

		err = txManager.User().Create(ctx, user)
		assert.NoError(t, err)

		// Commit
		err = txManager.Commit()
		assert.NoError(t, err)

		// Verify user was created
		found, err := repoManager.User().FindByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
	})
}
