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

// DBStorage implements Storage interface
type DBStorage struct {
	conn *pgx.Conn
	db   *pgxpool.Pool
}

// check that DBRepository implements all required methods
var _ storage.Storage = (*DBStorage)(nil)

// NewDBStorage returns a new DBStorage.
func NewDBStorage(connection *pgx.Conn, pool *pgxpool.Pool) *DBStorage {
	return &DBStorage{conn: connection, db: pool}
}

func (d *DBStorage) RegisterUser(ctx context.Context, user *models.User) error {
	err := d.conn.QueryRow(ctx,
		"INSERT INTO users (id, login, password) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING login",
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

	log.Debug().Msg("User registered")
	return nil
}

func (d *DBStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var users []models.User
	err := pgxscan.Select(ctx, d.db, &users, "SELECT id, login, password FROM users WHERE login=$1",
		login)
	if err != nil {
		log.Error().Msgf("GetUserByLogin error %s", err)
		return nil, err
	}

	if len(users) == 0 {
		log.Error().Msg("User doesn't exist")
		return nil, storage.ErrorUserNotFound
	}

	log.Debug().Msg("User loaded")
	return &users[0], nil
}

func (d *DBStorage) AddData(ctx context.Context, data *models.Data) error {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorage) GetDataByUserID(ctx context.Context, s string) ([]*models.Data, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorage) DeleteDataByDataID(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DBStorage) ReleaseStorage(ctx context.Context) error {
	err := d.conn.Close(ctx)
	if err != nil {
		return err
	}
	log.Info().Msg("Storage released")
	return nil
}
