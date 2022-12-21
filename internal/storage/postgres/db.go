// Package postgres contains business logic to work with postgres database.
package postgres

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage"
)

// DBStorage implements Storage interface.
type DBStorage struct {
	db *pgxpool.Pool
}

// check that DBStorage implements all required methods.
var _ storage.Storage = (*DBStorage)(nil)

// NewDBStorage returns a new DBStorage.
func NewDBStorage(pool *pgxpool.Pool) *DBStorage {
	return &DBStorage{db: pool}
}

// RegisterUser adds a new unique user to the storage.
func (d *DBStorage) RegisterUser(ctx context.Context, user models.User) error {
	err := d.db.QueryRow(ctx,
		`INSERT INTO users (id, login, password)
			 VALUES ($1, $2, $3) ON CONFLICT
			 DO NOTHING RETURNING login`,
		user.ID,
		user.Login,
		user.Password,
	).Scan(&user.Login)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Error().Msg("User already exist")
		return storage.ErrorUserAlreadyExist
	}

	if err != nil {
		log.Error().Msgf("RegisterUser error %s", err)
		return err
	}

	log.Debug().Msgf("User registered %+v", user)
	return nil
}

// GetUserByLogin gets current user data from storage.
func (d *DBStorage) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var users []models.User
	err := pgxscan.Select(ctx, d.db, &users, "SELECT id, login, password FROM users WHERE login=$1",
		login)
	if err != nil {
		log.Error().Msgf("GetUserByLogin error %s", err)
		return models.User{}, err
	}

	if len(users) == 0 {
		log.Error().Msg("User doesn't exist")
		return models.User{}, storage.ErrorUserNotFound
	}

	log.Debug().Msg("User loaded")
	return users[0], nil
}

// AddData adds private data to storage.
func (d *DBStorage) AddData(ctx context.Context, data models.Data) error {
	log.Debug().Msgf("AddData (postgres): %v", data)
	_, err := d.db.Exec(ctx,
		`INSERT INTO data (id, user_id, data_type, data_binary) 
			 VALUES ($1, $2, $3, $4) ON CONFLICT(id) 
			 DO UPDATE SET data_type = EXCLUDED.data_type,
			               data_binary = EXCLUDED.data_binary`,
		data.ID,
		data.UserID,
		data.DataType,
		data.DataBinary,
	)

	if err != nil {
		log.Error().Msgf("AddData error %s", err)
		return err
	}

	log.Debug().Msg("DataBinary added")
	return nil
}

// GetDataByUserID gets all related user data from storage.
func (d *DBStorage) GetDataByUserID(ctx context.Context, userID string) ([]models.Data, error) {
	var data []models.Data
	err := pgxscan.Select(ctx, d.db, &data,
		"SELECT id, user_id, data_type, data_binary FROM data WHERE user_id=$1",
		userID)
	if err != nil {
		log.Error().Msgf("GetDataByUserID error %s", err)
		return nil, err
	}

	if len(data) == 0 {
		log.Error().Msg("Data doesn't exist")
		return nil, storage.ErrorPrivateDataNotFound
	}

	log.Debug().Msg("Data loaded")
	return data, nil
}

// DeleteDataByDataID deletes private data from storage.
func (d *DBStorage) DeleteDataByDataID(ctx context.Context, id string) error {
	_, err := d.db.Exec(ctx,
		`DELETE from data WHERE id = $1`,
		id)

	if err != nil {
		return err
	}

	log.Info().Msg("DataBinary deleted")
	return nil
}

// ReleaseStorage closes database connection.
func (d *DBStorage) ReleaseStorage() {
	d.db.Close()
	log.Info().Msg("Storage released")
}
