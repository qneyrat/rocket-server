package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Crypto struct {
	Key []byte
}

func (c *Crypto) Encrypt(message []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(message))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], message)

	return c.b64encode(cipherText), nil
}

func (c *Crypto) Decrypt(message []byte) ([]byte, error) {
	cipherText, err := c.b64decode(message)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("cipherText block size is too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func (c *Crypto) b64encode(b []byte) []byte {
	encodedData := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, encodedData)
	encoder.Write(b)
	encoder.Close()
	return encodedData.Bytes()
}

func (c *Crypto) b64decode(b []byte) ([]byte, error) {
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(b))
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, dec)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
