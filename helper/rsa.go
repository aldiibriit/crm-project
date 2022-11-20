package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func GenRsaKeyForBE(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "private key",
		Bytes: derStream,
	}
	file, err := os.Create("privateBE.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "public key",
		Bytes: derPkix,
	}
	file, err = os.Create("publicBE.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func GenRsaKeyForFE(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "private key",
		Bytes: derStream,
	}
	file, err := os.Create("privateFE.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "public key",
		Bytes: derPkix,
	}
	file, err = os.Create("publicFE.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// encryption process when FE send request to BE using BE's publicKey
func RsaEncryptFEToBE(origData []byte) (string, error) {
	publicKeyBE, err := ioutil.ReadFile("publicBE.pem")
	if err != nil {
		log.Panicln("error in read publicBE")
	}

	block, _ := pem.Decode(publicKeyBE)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	if err != nil {
		log.Println("error in encrypt Data", err.Error())
	}

	result := base64.StdEncoding.EncodeToString([]byte(encryptedData))
	return result, nil
}

// decryption process when BE get request from BE using BE's privateKey
func RsaDecryptFromFEInBE(ciphertext []byte) (string, error) {
	privateKeyBE, err := ioutil.ReadFile("privateBE.pem")
	if err != nil {
		log.Panicln("error in read privateKeyBE")
	}
	block, _ := pem.Decode(privateKeyBE)
	if block == nil {
		return "", errors.New("error in pem decode privateKeyBE")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// privateKey := priv.(*rsa.PrivateKey)
	plainData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		log.Fatalln("error in decrypt data in BE", err.Error(), string(ciphertext))
	}
	return string(plainData), nil
}

func RsaDecryptFromFEInBEJava(ciphertext []byte) (string, error) {
	privateKeyBEJava, err := ioutil.ReadFile("privateBEJava.pem")
	if err != nil {
		log.Panicln("error in read privateKeyBEJava")
	}
	block, _ := pem.Decode(privateKeyBEJava)
	if block == nil {
		return "", errors.New("error in pem decode privateKeyBEJava")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	privateKey := priv.(*rsa.PrivateKey)
	plainData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		log.Fatalln("error in decrypt data in BEJava")
	}
	return string(plainData), nil
}

func RsaEncryptBEToFE(origData []byte) (string, error) {
	publicKeyFE, err := ioutil.ReadFile("publicFE.pem")
	if err != nil {
		log.Panicln("error in read publicFE")
	}

	block, _ := pem.Decode(publicKeyFE)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	if err != nil {
		log.Panicln("error in encrypt Data", err.Error())
	}

	result := base64.StdEncoding.EncodeToString([]byte(encryptedData))
	return result, nil
}

func RsaDecryptFromBEInFE(ciphertext []byte) (string, error) {
	privateKeyFE, err := ioutil.ReadFile("privateFE.pem")
	if err != nil {
		log.Panicln("error in read privateKeyFE")
	}
	block, _ := pem.Decode(privateKeyFE)
	if block == nil {
		return "", errors.New("error in pem decode privateKeyFE")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// privateKey := priv.(*rsa.PrivateKey)
	plainData, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		log.Fatalln("error in decrypt data in FE")
	}
	return string(plainData), nil
}
