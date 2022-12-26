package postgres

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
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

	pool, err := ConnectDB(context.Background(), cfg.DatabaseDsn)
	require.NoError(sts.T(), err)

	// migrations
	conn, err := ConnectDBForMigration(cfg.DatabaseDsn)
	require.NoError(sts.T(), err)

	err = RunMigrations(conn)
	require.NoError(sts.T(), err)

	db := NewDBStorage(pool)
	require.NoError(sts.T(), err)

	sts.TestStorage = db
	sts.container = storageContainer
}

func (sts *StorageTestSuite) TearDownTest() {
	sts.container.Close(sts.T())
}

func TestStorageTestSuite(t *testing.T) {
	if testhelpers.IsGithubActions() {
		// skip testcontainers for github actions
		return
	}
	suite.Run(t, new(StorageTestSuite))
}

func (sts *StorageTestSuite) TestDBStorage_RegisterUser() {
	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "positive test",
			user: models.User{
				ID:       uuid.NewString(),
				Login:    "login",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "negative test",
			user: models.User{
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
		user    models.User
		wantErr bool
	}{
		{
			name: "positive test",
			user: models.User{
				ID:       uuid.NewString(),
				Login:    "login",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "negative test",
			user: models.User{
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

func (sts *StorageTestSuite) TestDBStorage_AddData() {
	tests := []struct {
		name    string
		id      string
		data    *models.Text
		wantErr bool
	}{
		{
			name:    "positive test",
			id:      uuid.NewString(),
			data:    models.NewText("description", "some text here"),
			wantErr: false,
		},
	}

	user := models.User{
		ID:       uuid.NewString(),
		Login:    "login",
		Password: "password",
	}

	err := sts.TestStorage.RegisterUser(context.Background(), user)
	if err != nil {
		sts.T().Errorf("RegisterUser() error = %v", err)
		return
	}

	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorage

			binary, err := json.Marshal(tt.data)
			assert.NoError(sts.T(), err)

			err = s.AddData(context.Background(),
				models.Data{
					ID:         tt.id,
					UserID:     user.ID,
					DataType:   tt.data.GetType(),
					DataBinary: binary,
				})
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("AddData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (sts *StorageTestSuite) TestDBStorage_GetDataByUserID() {
	tests := []struct {
		name    string
		data    *models.Text
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    models.NewText("description", "some text here"),
			wantErr: false,
		},
		{
			name:    "positive test (update)",
			data:    models.NewText("description updated", "text updated"),
			wantErr: false,
		},
		{
			name:    "negative test",
			data:    models.NewText("test", "test"),
			wantErr: true,
		},
	}

	id := uuid.NewString()
	user := models.User{
		ID:       uuid.NewString(),
		Login:    "login",
		Password: "password",
	}

	err := sts.TestStorage.RegisterUser(context.Background(), user)
	if err != nil {
		sts.T().Errorf("RegisterUser() error = %v", err)
		return
	}

	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorage

			binary, err := json.Marshal(tt.data)
			assert.NoError(sts.T(), err)

			err = s.AddData(context.Background(),
				models.Data{
					ID:         id,
					UserID:     user.ID,
					DataType:   tt.data.GetType(),
					DataBinary: binary,
				})
			if err != nil {
				sts.T().Errorf("AddData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			userID := user.ID
			if tt.wantErr {
				userID = uuid.NewString()
			}

			data, err := s.GetDataByUserID(context.Background(), userID)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("GetDataByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.ErrorIs(sts.T(), err, storage.ErrorPrivateDataNotFound)
				return
			}

			var extracted models.Text
			err = json.Unmarshal(data[0].DataBinary, &extracted)
			assert.NoError(sts.T(), err)
			log.Debug().Msgf("data %v", extracted)

			assert.Equal(sts.T(), data[0].ID, id)
			assert.Equal(sts.T(), data[0].UserID, user.ID)
			assert.Equal(sts.T(), data[0].DataType, tt.data.GetType())
			assert.Equal(sts.T(), data[0].DataBinary, binary)
		})
	}
}

func (sts *StorageTestSuite) TestDBStorage_DeleteDataByDataID() {
	tests := []struct {
		name    string
		id      string
		data    *models.Text
		wantErr bool
	}{
		{
			name:    "positive test",
			id:      uuid.NewString(),
			data:    models.NewText("description", "some text here"),
			wantErr: false,
		},
	}

	user := models.User{
		ID:       uuid.NewString(),
		Login:    "login",
		Password: "password",
	}

	err := sts.TestStorage.RegisterUser(context.Background(), user)
	if err != nil {
		sts.T().Errorf("RegisterUser() error = %v", err)
		return
	}

	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorage

			binary, err := json.Marshal(tt.data)
			assert.NoError(sts.T(), err)

			err = s.AddData(context.Background(),
				models.Data{
					ID:         tt.id,
					UserID:     user.ID,
					DataType:   tt.data.GetType(),
					DataBinary: binary,
				})
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("AddData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = s.DeleteDataByDataID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				sts.T().Errorf("DeleteDataByDataID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (sts *StorageTestSuite) TestDBStorage_NegativeAll() {
	tests := []struct {
		name    string
		user    models.User
		id      string
		data    *models.Text
		wantErr bool
	}{
		{
			name: "negative tests for all methods",
			user: models.User{
				ID:       uuid.NewString(),
				Login:    "login",
				Password: "password",
			},
			id:      uuid.NewString(),
			data:    models.NewText("description", "some text here"),
			wantErr: true,
		},
	}

	sts.TestStorage.ReleaseStorage()

	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorage

			binary, err := json.Marshal(tt.data)
			assert.NoError(sts.T(), err)

			err = s.RegisterUser(context.Background(), tt.user)
			assert.NotNil(sts.T(), err)

			_, err = s.GetUserByLogin(context.Background(), "login")
			assert.NotNil(sts.T(), err)

			err = s.AddData(context.Background(),
				models.Data{
					ID:         tt.id,
					UserID:     uuid.NewString(),
					DataType:   tt.data.GetType(),
					DataBinary: binary,
				})
			assert.NotNil(sts.T(), err)

			_, err = s.GetDataByUserID(context.Background(), tt.user.ID)
			assert.NotNil(sts.T(), err)

			err = s.DeleteDataByDataID(context.Background(), tt.id)
			assert.NotNil(sts.T(), err)
		})
	}
}
