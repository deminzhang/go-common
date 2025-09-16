package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
)

func RsaParsePublicKey(key string) (*rsa.PublicKey, error) {
	b, _ := pem.Decode([]byte(key))
	if b == nil {
		return nil, errors.New("rsa public key wrong")
	}
	p, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	if pubKey, ok := p.(*rsa.PublicKey); ok {
		return pubKey, nil
	}
	return nil, errors.New("rsa public key wrong")
}

func RsaParsePrivateKey(key string) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode([]byte(key))
	if b == nil {
		return nil, errors.New("rsa private key wrong")
	}
	p, err := x509.ParsePKCS8PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}
	if priKey, ok := p.(*rsa.PrivateKey); ok {
		return priKey, nil
	}
	return nil, errors.New("rsa private key wrong")
}

func RsaEncrypt(src []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	random := rand.Reader
	srcLen := len(src)
	step := publicKey.Size() - 11
	var encryptedBytes []byte

	for start := 0; start < srcLen; start += step {
		finish := start + step
		if finish > srcLen {
			finish = srcLen
		}

		encryptedBlockBytes, err := rsa.EncryptPKCS1v15(random, publicKey, src[start:finish])
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func RsaDecrypt(src []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	random := rand.Reader
	srcLen := len(src)
	step := privateKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < srcLen; start += step {
		finish := start + step
		if finish > srcLen {
			finish = srcLen
		}

		decryptedBlockBytes, err := rsa.DecryptPKCS1v15(random, privateKey, src[start:finish])
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func RsaEncryptOAEP(src []byte, publicKey *rsa.PublicKey, hash hash.Hash) ([]byte, error) {
	random := rand.Reader
	srcLen := len(src)
	step := publicKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < srcLen; start += step {
		finish := start + step
		if finish > srcLen {
			finish = srcLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, publicKey, src[start:finish], nil)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func RsaDecryptOAEP(src []byte, privateKey *rsa.PrivateKey, hash hash.Hash) ([]byte, error) {
	random := rand.Reader
	srcLen := len(src)
	step := privateKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < srcLen; start += step {
		finish := start + step
		if finish > srcLen {
			finish = srcLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, privateKey, src[start:finish], nil)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func RsaSign(content []byte, privateKey *rsa.PrivateKey, hashType crypto.Hash) ([]byte, error) {
	hash := hashType.New()
	hash.Write([]byte(content))
	hashed := hash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hashType, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func RsaVerify(content []byte, signature []byte, pubKey *rsa.PublicKey, hashType crypto.Hash) error {
	hash := hashType.New()
	hash.Write([]byte(content))
	hashed := hash.Sum(nil)

	err := rsa.VerifyPKCS1v15(pubKey, hashType, hashed, signature)
	if err != nil {
		return err
	}
	return nil
}
