package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBinary(t *testing.T) {
	type args struct {
		description string
		value       []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test",
			args: args{
				description: "test description",
				value:       []byte("raw bytes"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBinary(tt.args.description, tt.args.value)
			assert.NotNil(t, got)
		})
	}
}

func TestNewCard(t *testing.T) {
	type args struct {
		description string
		name        string
		number      string
		date        string
		CVV         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test",
			args: args{
				description: "test description",
				name:        "DIGITAL CARD",
				number:      "5555 5555 5555 5555",
				date:        "01/2027",
				CVV:         "000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCard(tt.args.description, tt.args.name, tt.args.number, tt.args.date, tt.args.CVV)
			assert.NotNil(t, got)
		})
	}
}

func TestNewCredentials(t *testing.T) {
	type args struct {
		description string
		login       string
		password    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test",
			args: args{
				description: "test description",
				login:       "login",
				password:    "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCredentials(tt.args.description, tt.args.login, tt.args.password)
			assert.NotNil(t, got)
		})
	}
}

func TestNewText(t *testing.T) {
	type args struct {
		description string
		value       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test",
			args: args{
				description: "test description",
				value:       "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewText(tt.args.description, tt.args.value)
			assert.NotNil(t, got)
		})
	}
}

func Test_binary_GetJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    *binary
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    NewBinary("description", []byte("raw data")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.data.GetJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, "{\"description\":\"description\",\"value\":\"cmF3IGRhdGE=\"}", string(got))
		})
	}
}

func Test_binary_GetType(t *testing.T) {
	tests := []struct {
		name string
		data *binary
		want DataType
	}{
		{
			name: "positive test",
			data: NewBinary("description", []byte("raw data")),
			want: BinaryType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.data.GetType())
		})
	}
}

func Test_card_GetJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    *card
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    NewCard("description", "DIGITAL CARD", "5555 5555 5555 5555", "03/24", "000"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.data.GetJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, "{\"cvv\":\"000\", \"date\":\"03/24\", \"description\":\"description\", \"name\":\"DIGITAL CARD\", \"number\":\"5555 5555 5555 5555\"}", string(got))
		})
	}
}

func Test_card_GetType(t *testing.T) {
	tests := []struct {
		name string
		data *card
		want DataType
	}{
		{
			name: "positive test",
			data: NewCard("description", "DIGITAL CARD", "5555 5555 5555 5555", "03/24", "000"),
			want: CardType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.data.GetType())
		})
	}
}

func Test_credentials_GetJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    *credentials
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    NewCredentials("description", "user", "password"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.data.GetJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, "{\"description\":\"description\", \"login\":\"user\", \"password\":\"password\"}", string(got))
		})
	}
}

func Test_credentials_GetType(t *testing.T) {
	tests := []struct {
		name string
		data *credentials
		want DataType
	}{
		{
			name: "positive test",
			data: NewCredentials("description", "user", "password"),
			want: CredentialsType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.data.GetType())
		})
	}
}

func Test_text_GetJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    *text
		wantErr bool
	}{
		{
			name:    "positive test",
			data:    NewText("description", "some text here"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.data.GetJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, "{\"description\":\"description\", \"value\":\"some text here\"}", string(got))
		})
	}
}

func Test_text_GetType(t *testing.T) {
	tests := []struct {
		name string
		data *text
		want DataType
	}{
		{
			name: "positive test",
			data: NewText("description", "some text here"),
			want: TextType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.data.GetType())
		})
	}
}
