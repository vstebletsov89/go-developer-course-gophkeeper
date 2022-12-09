// Package models defines common structures and interfaces for server and client.
package models

import (
	"encoding/json"
)

// User represents a structure for user data.
type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"credentials"`
}

// DataType enum type for data types.
type DataType string

const (
	CredentialsType DataType = "credentials"
	TextType        DataType = "text"
	BinaryType      DataType = "binary"
	CardType        DataType = "card"
)

// Data represents a structure for data type.
type Data struct {
	ID       string
	UserID   string
	DataType DataType
	Data     []byte
}

// PrivateData is the interface that must be implemented by specific data type (credentials, text, binary, card).
type PrivateData interface {
	GetType() DataType
	GetJSON() ([]byte, error)
}

// check that credentials implements all required methods
var _ PrivateData = (*credentials)(nil)

// credentials represents a structure for login/password data.
type credentials struct {
	Description string `json:"description"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

func NewCredentials(description string, login string, password string) *credentials {
	return &credentials{Description: description, Login: login, Password: password}
}

func (p credentials) GetType() DataType {
	return CredentialsType
}

func (p credentials) GetJSON() ([]byte, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// check that text implements all required methods
var _ PrivateData = (*text)(nil)

// text represents a structure for text data.
type text struct {
	Description string `json:"description"`
	Value       string `json:"value"`
}

func NewText(description string, value string) *text {
	return &text{Description: description, Value: value}
}

func (t text) GetType() DataType {
	return TextType
}

func (t text) GetJSON() ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// check that text implements all required methods
var _ PrivateData = (*binary)(nil)

// binary represents a structure for binary data.
type binary struct {
	Description string `json:"description"`
	Value       []byte `json:"value"`
}

func NewBinary(description string, value []byte) *binary {
	return &binary{Description: description, Value: value}
}

func (b binary) GetType() DataType {
	return BinaryType
}

func (b binary) GetJSON() ([]byte, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// check that card implements all required methods
var _ PrivateData = (*card)(nil)

// card represents a structure for credit card data.
type card struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	Date        string `json:"date"`
	CVV         string `json:"cvv"`
}

func NewCard(description string, name string, number string, date string, CVV string) *card {
	return &card{Description: description, Name: name, Number: number, Date: date, CVV: CVV}
}

func (c card) GetType() DataType {
	return CardType
}

func (c card) GetJSON() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return data, nil
}
