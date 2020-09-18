package goutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func CreateHashMd5(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateSha256(key string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hasher.Sum(nil)
}

func EncryptAes(data []byte, passphrase string) (string, error) {
	block, err := aes.NewCipher([]byte(CreateSha256(passphrase)))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	base64Str := base64.StdEncoding.EncodeToString(ciphertext)
	return base64Str, nil
}

func DecryptAes(param []byte, passphrase string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(param))
	if err != nil {
		return "", err
	}

	data := []byte(decoded)
	key := []byte(CreateSha256(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func EncryptAesFile(filename string, data []byte, passphrase string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	dataFile, err := EncryptAes(data, passphrase)
	if err != nil {
		return err
	}
	f.Write([]byte(dataFile))
	return nil
}

func DecryptAesFile(filename string, passphrase string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return DecryptAes(data, passphrase)
}

func ValidateSHA1(param string) bool {
	match, err := regexp.MatchString(`\b[0-9a-f]{5,40}\b`, param)
	if err != nil {
		return false
	}
	return match
}

func TgEnc(text, key string) string {
	var result strings.Builder
	for i := 0; i < len(text); i++ {
		result.WriteString(fmt.Sprintf("%c", []rune(text)[i]^[]rune(key)[(i%len(key))]))
	}
	return base64.StdEncoding.EncodeToString([]byte(result.String()))
}

func TgDec(text, key string) string {
	param, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "--empty--"
	}
	paramStr := string(param)
	var result strings.Builder
	for i := 0; i < len(paramStr); i++ {
		result.WriteString(fmt.Sprintf("%c", []rune(paramStr)[i]^[]rune(key)[(i%len(key))]))
	}
	return result.String()
}
