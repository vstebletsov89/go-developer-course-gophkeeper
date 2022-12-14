package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"reflect"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test",
			args:    args{pwd: "password"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPassword(tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEqual(t, tt.args.pwd, got)
		})
	}
}

func TestIsUserAuthorized(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		userDB   string
		password string
		want     bool
		wantErr  bool
	}{
		{
			name:     "positive test",
			user:     "user",
			userDB:   "user",
			password: "password",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "negative test with different login",
			user:     "user",
			userDB:   "anotherUser",
			password: "password",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "negative test with different password",
			user:     "user",
			userDB:   "user",
			password: "password",
			want:     false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			encryptedPassword, err := EncryptPassword(tt.password)
			assert.NoError(t, err)

			password := tt.password
			if tt.wantErr {
				password = "invalid"
			}

			user := &models.User{
				ID:       "",
				Login:    tt.user,
				Password: password,
			}

			userDB := &models.User{
				ID:       "id",
				Login:    tt.userDB,
				Password: encryptedPassword,
			}

			got, err := IsUserAuthorized(user, userDB)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserAuthorized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUserAuthorized() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWTManager_GenerateToken(t *testing.T) {
	type fields struct {
		secretKey string
	}
	type args struct {
		user string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "positive test",
			fields:  fields{secretKey: "some_secret_key"},
			args:    args{user: "user"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			J := &JWTManager{
				secretKey: tt.fields.secretKey,
			}
			got, err := J.GenerateToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got)
		})
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	type fields struct {
		secretKey string
	}
	tests := []struct {
		name             string
		fields           fields
		user             string
		wantUserClaimsID string
		wantErr          bool
	}{
		{
			name:             "positive test",
			fields:           fields{secretKey: "some_key"},
			user:             "user",
			wantUserClaimsID: "user",
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			J := &JWTManager{
				secretKey: tt.fields.secretKey,
			}
			token, err := J.GenerateToken(tt.user)
			assert.NoError(t, err)

			got, err := J.ValidateToken(token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.user, got.ID)
		})
	}
}

func TestNewJWTManager(t *testing.T) {
	type args struct {
		secretKey string
	}
	tests := []struct {
		name string
		args args
		want *JWTManager
	}{
		{
			name: "positive test",
			args: args{secretKey: "some_secret_key"},
			want: &JWTManager{secretKey: "some_secret_key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJWTManager(tt.args.secretKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJWTManager() = %v, want %v", got, tt.want)
			}
		})
	}
}
