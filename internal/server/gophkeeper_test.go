package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	pb "github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/service/auth"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage/postgres/testhelpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"
)

func startGrpcServer() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read server configuration. Error: %v", err.Error())
	}

	if err := RunServer(cfg); err != nil {
		log.Fatalf("Failed to read server configuration. Error: %v", err.Error())
	}
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
    data_type   integer     NOT NULL,
    data_binary bytea       NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func migrateTables(pool *pgxpool.Pool) error {
	log.Println("Migration started..")
	_, err := pool.Exec(context.Background(), PostgreSQLTables)
	if err != nil {
		return err
	}
	log.Println("Migration done")
	return nil
}

func TestGophkeeperServer_Positive_Negative(t *testing.T) {
	// run docker with postgres
	storageContainer := testhelpers.NewTestDatabase(t)
	dsn := storageContainer.ConnectionString(t)
	log.Printf("DATABASE_DSN: %v", dsn)
	t.Setenv("DATABASE_DSN", dsn)
	t.Setenv("SERVER_ADDRESS", "localhost:3201")

	// connect to db only for migrations
	db, err := connectDB(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connectDB error: %s", err)
	}
	assert.NoError(t, err)

	// do migration
	err = migrateTables(db)
	assert.NoError(t, err)

	// close connection (it will be opened during start of server)
	db.Close()

	// start grpc server
	go startGrpcServer()

	// start client
	conn, err := grpc.Dial("localhost:3201", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		assert.NoError(t, err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			assert.NoError(t, err)
		}
	}(conn)
	authClient := pb.NewAuthClient(conn)
	gophkeeperClient := pb.NewGophkeeperClient(conn)
	ctx := context.Background()

	user := &pb.User{
		Login:    "testUser",
		Password: "password",
	}

	// Register
	_, err = authClient.Register(ctx, &pb.RegisterRequest{User: user})
	assert.NoError(t, err)

	// Login
	loginResponse, err := authClient.Login(ctx, &pb.LoginRequest{User: user})
	assert.NoError(t, err)
	log.Printf("loginResponse: %v", loginResponse)
	assert.Equal(t, user.Login, loginResponse.GetUser().GetLogin())
	assert.Equal(t, user.Password, loginResponse.GetUser().GetPassword())
	assert.NotNil(t, loginResponse.GetToken().GetUserId())
	assert.NotNil(t, loginResponse.GetToken().GetToken())

	userID := loginResponse.GetToken().GetUserId()

	// add jwt token for authorization
	md := metadata.New(map[string]string{auth.AccessToken: loginResponse.GetToken().GetToken()})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// add text data
	textSecret, err := models.NewText("text description", "text value").GetJSON()
	assert.NoError(t, err)

	textData := &pb.AddDataRequest{Data: &pb.Data{
		UserId:     userID,
		DataType:   pb.DataType_TEXT_TYPE,
		DataBinary: textSecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, textData)
	assert.NoError(t, err)

	// add binary data
	binarySecret, err := models.NewBinary("binary description", []byte("binary value")).GetJSON()
	assert.NoError(t, err)

	binaryData := &pb.AddDataRequest{Data: &pb.Data{
		UserId:     userID,
		DataType:   pb.DataType_BINARY_TYPE,
		DataBinary: binarySecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, binaryData)
	assert.NoError(t, err)

	// add card data
	cardSecret, err := models.NewCard("card description", "CARDHOLDER NAME", "6666 6666 6666 6666", "12/23", "000").GetJSON()
	assert.NoError(t, err)

	cardData := &pb.AddDataRequest{Data: &pb.Data{
		UserId:     userID,
		DataType:   pb.DataType_CARD_TYPE,
		DataBinary: cardSecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, cardData)
	assert.NoError(t, err)

	// add credentials data
	credentialsSecret, err := models.NewCredentials("credentials description", "login", "password").GetJSON()
	assert.NoError(t, err)

	credentialsData := &pb.AddDataRequest{Data: &pb.Data{
		UserId:     userID,
		DataType:   pb.DataType_CREDENTIALS_TYPE,
		DataBinary: credentialsSecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, credentialsData)
	assert.NoError(t, err)

	// Get all secret data for current user
	getDataResponse, err := gophkeeperClient.GetData(ctx, &pb.GetDataRequest{UserId: userID})
	log.Printf("getDataResponse: %v", getDataResponse)

	for _, secret := range getDataResponse.Data {
		log.Printf("secret from storage %v", secret)
		assert.Equal(t, userID, secret.UserId)
		switch secret.GetDataType() {
		case pb.DataType_CREDENTIALS_TYPE:
			assert.Equal(t, credentialsSecret, secret.GetDataBinary())
			break
		case pb.DataType_TEXT_TYPE:
			assert.Equal(t, textSecret, secret.GetDataBinary())
			break
		case pb.DataType_BINARY_TYPE:
			assert.Equal(t, binarySecret, secret.GetDataBinary())
			break
		case pb.DataType_CARD_TYPE:
			assert.Equal(t, cardSecret, secret.GetDataBinary())
			break
		default:
			errors.New("invalid data type")
		}
	}

	// Delete one secret from storage
	secret := getDataResponse.Data[0]
	_, err = gophkeeperClient.DeleteData(ctx, &pb.DeleteDataRequest{DataId: secret.DataId})

	// negative tests for authClient
	_, err = authClient.Register(ctx, nil)
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	invalidUser := &pb.User{
		Login:    "testUser",
		Password: "invalid_password",
	}
	_, err = authClient.Login(ctx, &pb.LoginRequest{User: invalidUser})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	// negative tests for gophkeeperClient
	_, err = gophkeeperClient.AddData(ctx, nil)
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	_, err = gophkeeperClient.GetData(ctx, &pb.GetDataRequest{UserId: uuid.NewString()})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	_, err = gophkeeperClient.GetData(ctx, &pb.GetDataRequest{UserId: "invalid_userid"})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	_, err = gophkeeperClient.DeleteData(ctx, &pb.DeleteDataRequest{DataId: "invalid_dataid"})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	// reset metadata and context to get authorization error
	md = metadata.New(map[string]string{"InvalidAccessToken": loginResponse.GetToken().GetToken()})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	_, err = gophkeeperClient.GetData(ctx, &pb.GetDataRequest{UserId: uuid.NewString()})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	ctx = context.Background()
	_, err = gophkeeperClient.GetData(ctx, &pb.GetDataRequest{UserId: uuid.NewString()})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())
}
