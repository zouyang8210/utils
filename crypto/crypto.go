package crypto

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

/*
功能:生成32位md5字串
参数:
	str:待加密的字符串
返回:签名
*/
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

/*
功能:RSA2签名
参数:
	data:待签名字符串
	privateKetStr:私钥
返回:签名,错误信息
*/
func Rsa2Sign(data string, privateKeyStr string) (sign string, err error) {
	var privateKey *rsa.PrivateKey
	var bSign []byte
	mPrivate := "-----BEGIN PRIVATE KEY-----\r\n" + privateKeyStr + "\r\n-----END PRIVATE KEY-----"
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	blockPri, _ := pem.Decode([]byte(mPrivate))
	if blockPri != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(blockPri.Bytes)
		if err == nil {
			bSign, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
			if err == nil {
				sign = base64.StdEncoding.EncodeToString(bSign)
			}
		}
	} else {
		err = errors.New("密钥不能还原")
	}
	return
}

/*
功能:RSA2验签
参数:
	data:待验签字符串
	publicKey:公钥
返回:验签是否成功(验签失败,查看错误信息),错误信息
*/
func VerifyRas2Sign(data, sign, publicKey string) (result bool, err error) {
	var pubInterface interface{}
	result = false
	mPublic := "-----BEGIN PUBLIC KEY-----\r\n" + publicKey + "\r\n-----END PUBLIC KEY-----"
	block, _ := pem.Decode([]byte(mPublic))

	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	if block != nil {
		pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err == nil {
			var bSign []byte
			pub := pubInterface.(*rsa.PublicKey)
			bSign, err = base64.StdEncoding.DecodeString(sign)
			if err == nil {
				err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, bSign)
				if err == nil {
					result = true
				}
			}
		}
	} else {
		err = errors.New("公钥不能还原")
	}

	return
}

func EncryptBase64(source []byte) string {
	return base64.StdEncoding.EncodeToString(source)
}

func DecodeBase64(base64Str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64Str)

}
