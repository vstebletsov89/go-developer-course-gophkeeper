package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage/postgres/testhelpers"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestStorage interface {
	storage.Storage
}

type StorageTestSuite struct {
	suite.Suite
	TestStorage
	container *testhelpers.TestDatabase
}

const PostgreSQLTables = `
CREATE TABLE IF NOT EXISTS users (
    id       uuid NOT NULL PRIMARY KEY,
    login    text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS data (
    id          uuid        NOT NULL PRIMARY KEY,
    user_id     uuid        REFERENCES users(id) ON DELETE CASCADE,
    data_type   text        NOT NULL,
    data_binary bytea       NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  timestamptz
);
`

func migrateTables(connection *pgx.Conn) error {
	log.Info().Msg("Migration started..")
	_, err := connection.Exec(context.Background(), PostgreSQLTables)
	if err != nil {
		return err
	}
	log.Info().Msg("Migration done")
	return nil
}

func connectDB(ctx context.Context, databaseURL string) (*pgx.Conn, *pgxpool.Pool, error) {
	log.Debug().Msg("Connect to DB...")
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Error().Msgf("Failed to connect to database. Error: %v", err.Error())
		return nil, nil, err
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Error().Msgf("Failed to create pgx pool. Error: %v", err.Error())
		return nil, nil, err
	}
	return conn, pool, nil
}

func (sts *StorageTestSuite) SetupTest() {
	// init global logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// set log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	storageContainer := testhelpers.NewTestDatabase(sts.T())

	dsn := storageContainer.ConnectionString(sts.T())
	log.Printf("DATABASE_DSN: %v", dsn)
	sts.T().Setenv("DATABASE_DSN", dsn)

	cfg, err := config.ReadConfig()
	require.NoError(sts.T(), err)

	conn, pool, err := connectDB(context.Background(), cfg.DatabaseDsn)
	require.NoError(sts.T(), err)

	// migrations
	err = migrateTables(conn)
	require.NoError(sts.T(), err)

	db := NewDBStorage(conn, pool)
	require.NoError(sts.T(), err)

	sts.TestStorage = db
	sts.container = storageContainer
}

func (sts *StorageTestSuite) TearDownTest() {
	sts.container.Close(sts.T())
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (sts *StorageTestSuite) TestDBStorage_RegisterUser() {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "positive test",
			user: &models.User{
				ID:       uuid.NewString(),
				Login:    "login",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "negative test",
			user: &models.User{
				ID:       uuid.NewString(),
				Login:    "login",
				Password: "password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorage

			err := s.RegisterUser(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.ErrorIs(sts.T(), err, storage.ErrorUserAlreadyExist)
			}
		})
	}
}

func (sts *StorageTestSuite) TestDBStorage_GetUserByLogin() {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "positive test",
			user: &models.User{
				ID:       uuid.NewString(),
				Login:    "login",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "negative test",
			user: &models.User{
				ID:       uuid.NewString(),
				Login:    "login2",
				Password: "password2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorage

			err := s.RegisterUser(context.Background(), tt.user)
			if err != nil {
				sts.T().Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var login string
			if !tt.wantErr {
				login = tt.user.Login
			} else {
				login = "invalid_login"
			}
			user, err := s.GetUserByLogin(context.Background(), login)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("GetUserByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(sts.T(), user.ID, tt.user.ID)
				assert.Equal(sts.T(), user.Login, tt.user.Login)
				assert.Equal(sts.T(), user.Password, tt.user.Password)
			} else {
				assert.ErrorIs(sts.T(), err, storage.ErrorUserNotFound)
			}
		})
	}
}
