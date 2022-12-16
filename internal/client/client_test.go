package client

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/server"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage/postgres/testhelpers"
)

func startGrpcServer(t *testing.T) {
	cfg, err := config.ReadConfig()
	assert.NoError(t, err)

	err = server.RunServer(cfg)
	assert.NoError(t, err)
}

func startGrpcClient() (*AuthClient, *SecretClient, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	authClient, secretClient, err := startClient(cfg)
	if err != nil {
		return nil, nil, err
	}
	return authClient, secretClient, nil
}

func TestGophkeeperClient_Positive_Negative(t *testing.T) {
	// run docker with postgres
	storageContainer := testhelpers.NewTestDatabase(t)
	dsn := storageContainer.ConnectionString(t)
	log.Printf("DATABASE_DSN: %v", dsn)
	t.Setenv("DATABASE_DSN", dsn)
	t.Setenv("SERVER_ADDRESS", "localhost:3202")

	// connect to db only for migrations
	db, err := testhelpers.ConnectDB(context.Background(), dsn)
	assert.NoError(t, err)

	// do migration
	err = testhelpers.MigrateTables(db)
	assert.NoError(t, err)

	// close connection (it will be opened during start of server)
	db.Close()

	// start grpc server
	go startGrpcServer(t)

	// start client
	authClient, secretClient, err := startGrpcClient()
	assert.NoError(t, err)

	ctx := context.Background()
	user := models.User{
		ID:       "",
		Login:    "clientUser",
		Password: "password",
	}
	authClient.SetUser(user)

	// test Register user
	err = authClient.Register(ctx)
	assert.NoError(t, err)

	// test Login user
	token, err := authClient.Login(ctx)
	assert.NoError(t, err)

	// set jwt token
	authClient.SetAccessToken(token)

	secretCredentials := models.NewCredentials("credentials description", "login", "password")
	binaryCredentials, err := secretCredentials.GetJSON()
	assert.NoError(t, err)
	dataCredentials := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secretCredentials.GetType(),
		DataBinary: binaryCredentials,
	}
	err = secretClient.AddData(ctx, dataCredentials)
	assert.NoError(t, err)

	secretText := models.NewText("text description", "some text")
	binaryText, err := secretText.GetJSON()
	assert.NoError(t, err)
	dataText := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secretText.GetType(),
		DataBinary: binaryText,
	}
	err = secretClient.AddData(ctx, dataText)
	assert.NoError(t, err)

	secretCard := models.NewCard("card description", "ivanov ivan", "5555 5555 5555 5555", "01/24", "000")
	binaryCard, err := secretCard.GetJSON()
	assert.NoError(t, err)
	dataCard := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secretCard.GetType(),
		DataBinary: binaryCard,
	}
	err = secretClient.AddData(ctx, dataCard)
	assert.NoError(t, err)

	secretBinary := models.NewBinary("binary description", []byte("some binary value"))
	binaryBinary, err := secretBinary.GetJSON()
	assert.NoError(t, err)
	dataBinary := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secretBinary.GetType(),
		DataBinary: binaryBinary,
	}
	err = secretClient.AddData(ctx, dataBinary)
	assert.NoError(t, err)

	data, err := secretClient.GetData(ctx)
	for _, secret := range data {
		log.Printf("secret from storage %v", secret)
		switch secret.DataType {
		case models.CREDENTIALS_TYPE:
			assert.Equal(t, binaryCredentials, secret.DataBinary)
			break
		case models.TEXT_TYPE:
			assert.Equal(t, binaryText, secret.DataBinary)
			break
		case models.BINARY_TYPE:
			assert.Equal(t, binaryBinary, secret.DataBinary)
			break
		case models.CARD_TYPE:
			assert.Equal(t, binaryCard, secret.DataBinary)
			break
		default:
			assert.Equal(t, "correct_type", "invalid_type_from_server")
		}
	}

	err = secretClient.DeleteData(ctx, data[0].ID)
	assert.NoError(t, err)
}
