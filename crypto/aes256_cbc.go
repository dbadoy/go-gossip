package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
)

type AES256_CBC struct{}

func (a AES256_CBC) Encrypt(key string, plain []byte) ([]byte, error) {
	k := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := k[:aes.BlockSize]
	enc := cipher.NewCBCEncrypter(block, iv)
	padded := padPKCS7(plain, block.BlockSize())
	cipherText := make([]byte, len(padded))
	enc.CryptBlocks(cipherText, padded)
	return cipherText, nil
}

func (a AES256_CBC) Decrypt(key string, data []byte) ([]byte, error) {
	k := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := k[:aes.BlockSize]
	dec := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(data))
	dec.CryptBlocks(plainText, data)
	return trimPKCS5(plainText), nil
}

func (a AES256_CBC) EncryptWithPublicKey(key *rsa.PublicKey, plain []byte) ([]byte, error) {
	panic("not supported")
}

func (a AES256_CBC) DecryptWithPrivateKey(key *rsa.PrivateKey, plain []byte) ([]byte, error) {
	panic("not supported")
}

func padPKCS7(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padText...)
}

func trimPKCS5(text []byte) []byte {
	padding := text[len(text)-1]
	return text[:len(text)-int(padding)]
}
