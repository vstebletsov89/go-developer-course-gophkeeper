package client

import (
	"context"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/client/cli"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/config"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/server"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/storage/postgres/testhelpers"
)

func startGrpcServer(t *testing.T) {
	cfg, err := config.ReadConfig()
	assert.NoError(t, err)

	err = server.RunServer(cfg)
	assert.NoError(t, err)
}

func startGrpcClient() (*cli.CLI, error) {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	client, err := startClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestGophkeeperClient_Positive_Negative(t *testing.T) {
	if testhelpers.IsGithubActions() {
		// skip testcontainers for github actions
		return
	}

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
	client, err := startGrpcClient()
	assert.NoError(t, err)

	ctx := context.Background()

	// test Register user
	args := make([]string, 2)
	args[0] = "user"
	args[1] = "password"
	err = client.Register(ctx, args)
	assert.NoError(t, err)

	// test Login user
	err = client.Login(ctx, args)
	assert.NoError(t, err)

	// add different data types
	args[0] = "binary description"
	args[1] = "some binary value"
	err = client.AddBinary(ctx, args)
	assert.NoError(t, err)

	args[0] = "text description"
	args[1] = "some text"
	err = client.AddText(ctx, args)
	assert.NoError(t, err)

	args = make([]string, 3)
	args[0] = "credentials description"
	args[1] = "login"
	args[2] = "password"
	err = client.AddCredentials(ctx, args)
	assert.NoError(t, err)

	args = make([]string, 5)
	args[0] = "card description"
	args[1] = "ivanov ivan"
	args[2] = "5555 5555 5555 5555"
	args[3] = "01/24"
	args[4] = "000"
	err = client.AddCard(ctx, args)
	assert.NoError(t, err)

	// get all data
	data, err := client.GetData(ctx)
	assert.NoError(t, err)

	client.LogData(data)

	// delete data
	args = make([]string, 1)
	args[0] = data[0].ID
	err = client.DeleteData(ctx, args)
	assert.NoError(t, err)
}
