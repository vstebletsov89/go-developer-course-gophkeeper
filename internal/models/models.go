// Package models defines common structures and interfaces for server and client.
package models

import (
	"encoding/json"
)

// User represents a structure for user data.
type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// DataType enum type for data types (same as in grpc).
type DataType int32

const (
	CREDENTIALS_TYPE DataType = 0
	TEXT_TYPE        DataType = 1
	BINARY_TYPE      DataType = 2
	CARD_TYPE        DataType = 3
)

// Data represents a structure for data type.
type Data struct {
	ID         string
	UserID     string
	DataType   DataType
	DataBinary []byte
}

// PrivateData is the interface that must be implemented by specific data type (credentials, text, binary, card).
type PrivateData interface {
	GetType() DataType
	GetJSON() ([]byte, error)
}

// check that Credentials implements all required methods.
var _ PrivateData = (*Credentials)(nil)

// Credentials represents a structure for login/password data.
type Credentials struct {
	Description string `json:"description"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

// NewCredentials returns an instance of Credentials.
func NewCredentials(description string, login string, password string) *Credentials {
	return &Credentials{Description: description, Login: login, Password: password}
}

// GetType getter for Credentials type.
func (p Credentials) GetType() DataType {
	return CREDENTIALS_TYPE
}

// GetJSON getter for Credentials binary data.
func (p Credentials) GetJSON() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// check that Text implements all required methods.
var _ PrivateData = (*Text)(nil)

// Text represents a structure for Text data.
type Text struct {
	Description string `json:"description"`
	Value       string `json:"value"`
}

// NewText returns an instance of Text.
func NewText(description string, value string) *Text {
	return &Text{Description: description, Value: value}
}

// GetType getter for Text type.
func (t Text) GetType() DataType {
	return TEXT_TYPE
}

// GetJSON getter for Text binary data.
func (t Text) GetJSON() ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// check that Text implements all required methods.
var _ PrivateData = (*Binary)(nil)

// Binary represents a structure for Binary data.
type Binary struct {
	Description string `json:"description"`
	Value       []byte `json:"value"`
}

// NewBinary returns an instance of Binary.
func NewBinary(description string, value []byte) *Binary {
	return &Binary{Description: description, Value: value}
}

// GetType getter for Binary type.
func (b Binary) GetType() DataType {
	return BINARY_TYPE
}

// GetJSON getter for Binary binary data.
func (b Binary) GetJSON() ([]byte, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// check that Card implements all required methods.
var _ PrivateData = (*Card)(nil)

// Card represents a structure for credit Card data.
type Card struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	Date        string `json:"date"`
	CVV         string `json:"cvv"`
}

// NewCard returns an instance of Card.
func NewCard(description string, name string, number string, date string, cvv string) *Card {
	return &Card{Description: description, Name: name, Number: number, Date: date, CVV: cvv}
}

// GetType getter for Card type.
func (c Card) GetType() DataType {
	return CARD_TYPE
}

// GetJSON getter for Card binary data.
func (c Card) GetJSON() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}
