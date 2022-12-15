package server

import (
	"context"
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

func TestGophkeeperServer_Positive_Negative(t *testing.T) {
	// run docker with postgres
	storageContainer := testhelpers.NewTestDatabase(t)
	dsn := storageContainer.ConnectionString(t)
	log.Printf("DATABASE_DSN: %v", dsn)
	t.Setenv("DATABASE_DSN", dsn)
	t.Setenv("SERVER_ADDRESS", "localhost:3201")

	// connect to db only for migrations
	db, err := testhelpers.ConnectDB(context.Background(), dsn)
	if err != nil {
		log.Fatalf("ConnectDB error: %s", err)
	}
	assert.NoError(t, err)

	// do migration
	err = testhelpers.MigrateTables(db)
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
		Login:    "serverUser",
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

	// add jwt token for authorization
	md := metadata.New(map[string]string{auth.AccessToken: loginResponse.GetToken().GetToken(),
		auth.UserCtx: loginResponse.GetToken().GetUserId()})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	// add text data
	textSecret, err := models.NewText("text description", "text value").GetJSON()
	assert.NoError(t, err)

	textData := &pb.AddDataRequest{Data: &pb.Data{
		DataType:   pb.DataType_TEXT_TYPE,
		DataBinary: textSecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, textData)
	assert.NoError(t, err)

	// add binary data
	binarySecret, err := models.NewBinary("binary description", []byte("binary value")).GetJSON()
	assert.NoError(t, err)

	binaryData := &pb.AddDataRequest{Data: &pb.Data{
		DataType:   pb.DataType_BINARY_TYPE,
		DataBinary: binarySecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, binaryData)
	assert.NoError(t, err)

	// add card data
	cardSecret, err := models.NewCard("card description", "CARDHOLDER NAME", "6666 6666 6666 6666", "12/23", "000").GetJSON()
	assert.NoError(t, err)

	cardData := &pb.AddDataRequest{Data: &pb.Data{
		DataType:   pb.DataType_CARD_TYPE,
		DataBinary: cardSecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, cardData)
	assert.NoError(t, err)

	// add credentials data
	credentialsSecret, err := models.NewCredentials("credentials description", "login", "password").GetJSON()
	assert.NoError(t, err)

	credentialsData := &pb.AddDataRequest{Data: &pb.Data{
		DataType:   pb.DataType_CREDENTIALS_TYPE,
		DataBinary: credentialsSecret,
	}}
	_, err = gophkeeperClient.AddData(ctx, credentialsData)
	assert.NoError(t, err)

	// Get all secret data for current user
	getDataResponse, err := gophkeeperClient.GetData(ctx, &pb.GetDataRequest{})
	log.Printf("getDataResponse: %v", getDataResponse)

	for _, secret := range getDataResponse.Data {
		log.Printf("secret from storage %v", secret)
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
			assert.Equal(t, "correct_type", "invalid_type_from_server")
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

	_, err = gophkeeperClient.DeleteData(ctx, &pb.DeleteDataRequest{DataId: "invalid_dataid"})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	// reset metadata and context to get authorization error
	md = metadata.New(map[string]string{"InvalidAccessToken": loginResponse.GetToken().GetToken()})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	_, err = gophkeeperClient.GetData(ctx, &pb.GetDataRequest{})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())

	ctx = context.Background()
	_, err = gophkeeperClient.GetData(ctx, &pb.GetDataRequest{})
	assert.NotNil(t, err)
	log.Printf("err : %v", err.Error())
}
