package crypto

import (
	"crypto/rsa"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSA(t *testing.T) {
	t.Run("密钥对生成", func(t *testing.T) {
		privateKey, publicKey, err := GenerateKeyPair(2048)
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)
		assert.NotNil(t, publicKey)
		assert.Equal(t, 0, privateKey.N.Cmp(publicKey.N), "私钥和公钥的模数应该匹配")
	})

	t.Run("密钥PEM格式转换", func(t *testing.T) {
		privateKey, publicKey, err := GenerateKeyPair(2048)
		assert.NoError(t, err)

		// 测试私钥PEM转换
		privateKeyPEM := ExportPrivateKeyToPEM(privateKey)
		assert.NotEmpty(t, privateKeyPEM, "私钥PEM格式不应为空")

		importedPrivateKey, err := ImportPrivateKeyFromPEM(privateKeyPEM)
		assert.NoError(t, err)
		assert.Equal(t, 0, importedPrivateKey.N.Cmp(privateKey.N), "导入的私钥应与原始私钥匹配")

		// 测试公钥PEM转换
		publicKeyPEM := ExportPublicKeyToPEM(publicKey)
		assert.NotEmpty(t, publicKeyPEM, "公钥PEM格式不应为空")

		importedPublicKey, err := ImportPublicKeyFromPEM(publicKeyPEM)
		assert.NoError(t, err)
		assert.Equal(t, 0, importedPublicKey.N.Cmp(publicKey.N), "导入的公钥应与原始公钥匹配")
	})

	t.Run("加密解密", func(t *testing.T) {
		// 准备测试数据
		testCases := []struct {
			name string
			data []byte
		}{
			{"空数据", []byte{}},
			{"短文本", []byte("Hello")},
			{"中文字符", []byte("你好，世界")},
			{"长文本", []byte("这是一段很长的文本，用于测试RSA加密解密功能。This is a long text for testing RSA encryption and decryption.")},
		}

		privateKey, publicKey, err := GenerateKeyPair(2048)
		assert.NoError(t, err)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// 加密
				encryptedData, err := Encrypt(publicKey, tc.data)
				assert.NoError(t, err)
				assert.NotEmpty(t, encryptedData, "加密后的数据不应为空")

				// 验证base64编码
				encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedData)
				assert.NotEmpty(t, encryptedBase64, "加密后的Base64字符串不应为空")

				// 解密
				decryptedData, err := Decrypt(privateKey, encryptedData)
				assert.NoError(t, err)
				assert.Equal(t, tc.data, decryptedData, "解密后的数据应与原始数据匹配")
			})
		}
	})

	t.Run("签名验证", func(t *testing.T) {
		// 准备测试数据
		testCases := []struct {
			name string
			data []byte
		}{
			{"空数据", []byte{}},
			{"短文本", []byte("Hello")},
			{"中文字符", []byte("你好，世界")},
			{"长文本", []byte("这是一段很长的文本，用于测试RSA签名验证功能。This is a long text for testing RSA signature verification.")},
		}

		privateKey, publicKey, err := GenerateKeyPair(2048)
		assert.NoError(t, err)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// 签名
				signature, err := Sign(privateKey, tc.data)
				assert.NoError(t, err)
				assert.NotEmpty(t, signature, "签名不应为空")

				// 验证签名
				err = Verify(publicKey, tc.data, signature)
				assert.NoError(t, err, "签名验证应该成功")

				// 测试篡改数据
				if len(tc.data) > 0 {
					tamperedData := make([]byte, len(tc.data))
					copy(tamperedData, tc.data)
					tamperedData[0] = ^tamperedData[0] // 修改第一个字节
					err = Verify(publicKey, tamperedData, signature)
					assert.Error(t, err, "篡改数据后签名验证应该失败")
				}

				// 测试篡改签名
				if len(signature) > 0 {
					tamperedSignature := make([]byte, len(signature))
					copy(tamperedSignature, signature)
					tamperedSignature[0] = ^tamperedSignature[0] // 修改第一个字节
					err = Verify(publicKey, tc.data, tamperedSignature)
					assert.Error(t, err, "篡改签名后验证应该失败")
				}
			})
		}
	})

	t.Run("签名错误处理", func(t *testing.T) {
		_, publicKey, err := GenerateKeyPair(2048)
		assert.NoError(t, err)

		// 测试无效的私钥
		invalidPrivateKey := &rsa.PrivateKey{}
		_, err = Sign(invalidPrivateKey, []byte("test"))
		assert.Error(t, err, "使用无效私钥签名应该失败")

		// 测试无效的公钥
		invalidPublicKey := &rsa.PublicKey{}
		err = Verify(invalidPublicKey, []byte("test"), []byte("signature"))
		assert.Error(t, err, "使用无效公钥验证应该失败")

		// 测试空签名
		err = Verify(publicKey, []byte("test"), nil)
		assert.Error(t, err, "验证空签名应该失败")
	})

	t.Run("错误处理", func(t *testing.T) {
		// 测试无效的PEM格式
		invalidPEM := "invalid pem data"
		_, err := ImportPrivateKeyFromPEM(invalidPEM)
		assert.Error(t, err, "导入无效PEM时应返回错误")

		_, err = ImportPublicKeyFromPEM(invalidPEM)
		assert.Error(t, err, "导入无效PEM时应返回错误")

		// 测试过大的数据
		_, publicKey, err := GenerateKeyPair(2048)
		assert.NoError(t, err)

		// 生成一个超过RSA密钥长度限制的数据
		largeData := make([]byte, 300)
		_, err = Encrypt(publicKey, largeData)
		assert.Error(t, err, "加密过大数据时应返回错误")
	})
}
