package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// GenerateKeyPair 生成RSA密钥对
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// ExportPrivateKeyToPEM 将私钥导出为PEM格式
func ExportPrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

// ExportPublicKeyToPEM 将公钥导出为PEM格式
func ExportPublicKeyToPEM(publicKey *rsa.PublicKey) string {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM)
}

// ImportPrivateKeyFromPEM 从PEM格式导入私钥
func ImportPrivateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ImportPublicKeyFromPEM 从PEM格式导入公钥
func ImportPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// Encrypt 使用公钥加密数据
func Encrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
}

// Decrypt 使用私钥解密数据
func Decrypt(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
}

// Sign 使用私钥对数据进行签名
func Sign(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := crypto.SHA256.New()
	hash.Write(data)
	hashed := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
}

// Verify 使用公钥验证签名
func Verify(publicKey *rsa.PublicKey, data []byte, signature []byte) error {
	hash := crypto.SHA256.New()
	hash.Write(data)
	hashed := hash.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, signature)
}
