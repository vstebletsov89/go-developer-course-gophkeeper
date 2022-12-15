package secure

import (
	"github.com/stretchr/testify/assert"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
	"reflect"
	"testing"
)

func TestDecrypt(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "positive test",
			args:    args{data: []byte("data to be encrypted")},
			want:    []byte("data to be encrypted"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.args.data)
			assert.NoError(t, err)

			got, err := Decrypt(encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecryptPrivateData(t *testing.T) {
	type args struct {
		data *proto.Data
	}
	tests := []struct {
		name     string
		args     args
		want     *proto.Data
		wantData []byte
		wantErr  bool
	}{
		{
			name: "positive test",
			args: args{data: &proto.Data{
				DataType:   0,
				DataBinary: []byte("some binary data"),
			}},
			wantData: []byte("some binary data"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptPrivateData(tt.args.data, "userId")
			assert.NoError(t, err)

			got, err := DecryptPrivateData(encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptPrivateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantData, got.DataBinary)
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive test",
			args:    args{data: []byte("data to be encrypted")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestEncryptPrivateData(t *testing.T) {
	type args struct {
		data *proto.Data
	}
	tests := []struct {
		name    string
		args    args
		want    models.Data
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{data: &proto.Data{
				DataType:   0,
				DataBinary: []byte("some binary data"),
			}},
			want:    models.Data{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPrivateData(tt.args.data, "userID")
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptPrivateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func Test_cipherInit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "positive test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cipherInit(); (err != nil) != tt.wantErr {
				t.Errorf("cipherInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
