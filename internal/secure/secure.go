package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/models"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/proto"
	"github.com/vstebletsov89/go-developer-course-gophkeeper/internal/rand"
	"sync"
)

type cipherData struct {
	key    []byte
	nonce  []byte
	aesGCM cipher.AEAD
}

var cipherInstance *cipherData
var once sync.Once

func cipherInit() error {
	var e error
	once.Do(func() {
		key := rand.GenerateRandom(2 * aes.BlockSize)

		aesblock, err := aes.NewCipher(key)
		if err != nil {
			e = err
		}

		aesgcm, err := cipher.NewGCM(aesblock)
		if err != nil {
			e = err
		}

		nonce := rand.GenerateRandom(aesgcm.NonceSize())
		cipherInstance = &cipherData{key: key, aesGCM: aesgcm, nonce: nonce}
	})
	return e
}

func Encrypt(data []byte) ([]byte, error) {
	if err := cipherInit(); err != nil {
		return nil, err
	}
	log.Debug().Msgf("Data to be encrypted: %s", string(data))

	encrypted := cipherInstance.aesGCM.Seal(nil, cipherInstance.nonce, data, nil)
	dst := make([]byte, hex.EncodedLen(len(encrypted)))
	hex.Encode(dst, encrypted)

	log.Debug().Msgf("Encrypted data: %s", string(dst))
	return dst, nil
}

func Decrypt(data []byte) ([]byte, error) {
	if err := cipherInit(); err != nil {
		return nil, err
	}
	log.Debug().Msgf("Data to be decrypted: %s", string(data))

	dst := make([]byte, hex.DecodedLen(len(data)))
	_, err := hex.Decode(dst, data)
	if err != nil {
		return nil, err
	}

	decrypted, err := cipherInstance.aesGCM.Open(nil, cipherInstance.nonce, dst, nil)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Decrypted data: %s", string(decrypted))
	return decrypted, nil
}

func EncryptPrivateData(data *proto.Data, userID string) (models.Data, error) {
	var securedData models.Data
	encryptedBinary, err := Encrypt(data.GetDataBinary())
	if err != nil {
		return models.Data{}, err
	}

	securedData.ID = uuid.NewString()
	securedData.UserID = userID
	securedData.DataType = models.DataType(data.GetDataType())
	securedData.DataBinary = encryptedBinary

	return securedData, nil
}

func DecryptPrivateData(data models.Data) (*proto.Data, error) {
	var securedData proto.Data
	decryptedBinary, err := Decrypt(data.DataBinary)
	if err != nil {
		return nil, err
	}

	securedData.DataId = data.ID
	securedData.DataType = proto.DataType(data.DataType)
	securedData.DataBinary = decryptedBinary

	return &securedData, nil
}
