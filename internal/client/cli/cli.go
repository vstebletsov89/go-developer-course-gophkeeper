// Package cli implements cli for client.
package cli

import (
	"context"
	"errors"
	"github.com/c-bata/go-prompt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/client/service"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"os"
	"strings"
)

// CLI represents a structure for cli communication with user.
type CLI struct {
	authClient   *service.AuthClient
	secretClient *service.SecretClient
}

// NewCLI returns an instance of CLI.
func NewCLI(authClient *service.AuthClient, secretClient *service.SecretClient) *CLI {
	return &CLI{authClient: authClient, secretClient: secretClient}
}

// Completer is a menu items for the Gophkeeper UI.
func (c *CLI) Completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "register", Description: "Register new user for gophkeeper application. Example: register <user> <password>"},
		{Text: "login", Description: "Sign-in into gophkeeper application. Example: Login <user> <password>"},
		{Text: "add-text", Description: "Add new private text data. Example: add-text <description> <text>"},
		{Text: "add-card", Description: "Add new private card data. Example: add-card <description> <name> <number> <date> <cvv>"},
		{Text: "add-binary", Description: "Add new private binary data. Example: add-binary <description> <value>"},
		{Text: "add-credentials", Description: "Add new private credentials data. Example: add-credentials <user> <password>"},
		{Text: "get-data", Description: "Get all private data for the user. Example: get-data"},
		{Text: "delete-data", Description: "Delete private data. Example: delete-data <data_id>"},
		{Text: "exit", Description: "Exit from gophkeeper application. Example: exit"},
	}
	return prompt.FilterContains(s, d.CurrentLine(), true)
}

// Register add new user for gophkeeper application.
func (c *CLI) Register(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}
	c.setCurrentUser(args)
	return c.authClient.Register(ctx)
}

func (c *CLI) setCurrentUser(args []string) {
	user := models.User{
		ID:       "",
		Login:    args[0],
		Password: args[1],
	}
	c.authClient.SetUser(user)
}

// Login sign-in into gophkeeper application.
func (c *CLI) Login(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}

	c.setCurrentUser(args)

	token, err := c.authClient.Login(ctx)
	if err != nil {
		log.Error().Msgf("Failed to Login: %v", err)
		return err
	}

	// set jwt token
	c.authClient.SetAccessToken(token)
	return nil
}

// DeleteData deletes private data from storage.
func (c *CLI) DeleteData(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return errors.New("invalid arguments")
	}

	return c.secretClient.DeleteData(ctx, args[0])
}

// GetData gets all private data from the storage.
func (c *CLI) GetData(ctx context.Context) ([]models.Data, error) {
	data, err := c.secretClient.GetData(ctx)
	if err != nil {
		log.Error().Msgf("Failed to get private data: %v", err)
		return nil, err
	}
	return data, nil
}

// AddBinary add binary data to the storage.
func (c *CLI) AddBinary(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}

	secret := models.NewBinary(args[0], []byte(args[1]))
	binary, err := secret.GetJSON()
	if err != nil {
		log.Error().Msgf("Failed to convert binary data: %v", err)
		return err
	}

	data := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secret.GetType(),
		DataBinary: binary,
	}

	return c.secretClient.AddData(ctx, data)
}

// AddCredentials add credentials data to the storage.
func (c *CLI) AddCredentials(ctx context.Context, args []string) error {
	if len(args) != 3 {
		return errors.New("invalid arguments")
	}

	secret := models.NewCredentials(args[0], args[1], args[2])
	binary, err := secret.GetJSON()
	if err != nil {
		log.Error().Msgf("Failed to convert credentials data: %v", err)
		return err
	}

	data := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secret.GetType(),
		DataBinary: binary,
	}

	return c.secretClient.AddData(ctx, data)
}

// AddText add text data to the storage.
func (c *CLI) AddText(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("invalid arguments")
	}

	secret := models.NewText(args[0], args[1])
	binary, err := secret.GetJSON()
	if err != nil {
		log.Error().Msgf("Failed to convert text data: %v", err)
		return err
	}

	data := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secret.GetType(),
		DataBinary: binary,
	}

	return c.secretClient.AddData(ctx, data)
}

// AddCard add card data to the storage.
func (c *CLI) AddCard(ctx context.Context, args []string) error {
	if len(args) != 5 {
		return errors.New("invalid arguments")
	}

	secret := models.NewCard(args[0], args[1], args[2], args[3], args[4])
	binary, err := secret.GetJSON()
	if err != nil {
		log.Error().Msgf("Failed to convert card data: %v", err)
		return err
	}

	data := models.Data{
		ID:         uuid.NewString(),
		UserID:     "",
		DataType:   secret.GetType(),
		DataBinary: binary,
	}

	return c.secretClient.AddData(ctx, data)
}

// Executor runs actions for option items.
func (c *CLI) Executor(input string) {
	log.Debug().Msgf("Option selected: " + input)
	args := strings.Split(input, " ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch args[0] {
	case "register":
		err := c.Register(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to register new user: %v", err)
			return
		}
		log.Info().Msg("User was registered. Use login command to sign-in.")
	case "login":
		err := c.Login(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to login: %v", err)
			return
		}
		log.Info().Msg("UI login done.")
	case "add-text":
		err := c.AddText(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to add text data: %v", err)
			return
		}
		log.Info().Msg("Text data was added.")
	case "add-card":
		err := c.AddCard(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to add card data: %v", err)
			return
		}
		log.Info().Msg("Card data was added.")
	case "add-binary":
		err := c.AddBinary(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to add binary data: %v", err)
			return
		}
		log.Info().Msg("Binary data was added.")
	case "add-credentials":
		err := c.AddCredentials(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to add credentials data: %v", err)
			return
		}
		log.Info().Msg("Credentials data was added.")
	case "get-data":
		data, err := c.GetData(ctx)
		if err != nil {
			log.Error().Msgf("Failed to get data: %v", err)
			return
		}
		c.LogData(data)
		log.Info().Msg("All user data was received.")
	case "delete-data":
		err := c.DeleteData(ctx, args[1:])
		if err != nil {
			log.Error().Msgf("Failed to delete data: %v", err)
			return
		}
		log.Info().Msg("Data was deleted.")
	case "exit":
		log.Debug().Msg("Client shutdown.")
		os.Exit(0)
	default:
		log.Info().Msg("Invalid option.")
	}
}

// LogData prints formatted private data.
func (c *CLI) LogData(data []models.Data) {
	log.Info().Msg("Private data for current user:")
	for _, secret := range data {
		switch secret.DataType {
		case models.CREDENTIALS_TYPE:
			log.Info().Msgf("ID: %s type: CREDENTIALS data: %s",
				secret.ID, string(secret.DataBinary))
		case models.TEXT_TYPE:
			log.Info().Msgf("ID: %s type: TEXT data: %s",
				secret.ID, string(secret.DataBinary))
		case models.BINARY_TYPE:
			log.Info().Msgf("ID: %s type: BINARY data: %s",
				secret.ID, string(secret.DataBinary))
		case models.CARD_TYPE:
			log.Info().Msgf("ID: %s type: CARD data: %s",
				secret.ID, string(secret.DataBinary))
		}
	}
}
