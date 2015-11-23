package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/Centny/gwf/log"
)

/**
 * RSA签名
 * @param $data 待签名数据
 * @param $private_key_path 商户私钥文件路径
 * return 签名结果
 */
func rsaSign(origData string, privateKey *rsa.PrivateKey) (string, error) {

	log.I("rsaSign for origData :%v ", origData)
	// log.I("rsaSign for privateKey :%v ", privateKey)

	h := sha1.New()
	h.Write([]byte(origData))
	digest := h.Sum(nil)

	s, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, digest)
	if err != nil {
		log.E("rsaSign SignPKCS1v15 error")
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

/**
 * RSA验签
 * @param $data 待签名数据
 * @param $ali_public_key_path 支付宝的公钥文件路径
 * @param $sign 要校对的的签名结果
 * return 验证结果
 */
func rsaVerify(data, sign string, public_key *rsa.PublicKey) error {
	log.I("rsaVerify \n sign ： %v \n data : %v", sign, data)
	h := sha1.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)

	err := rsa.VerifyPKCS1v15(public_key, crypto.SHA1, digest, []byte(sign))
	if err != nil {
		log.E("VerifyPKCS1v15 fail : %v\n", err)
		return err
	}
	return nil
}

/**
 * RSA解密
 * @param $content 需要解密的内容，密文
 * @param $private_key_path 商户私钥文件路径
 * return 解密后内容，明文
 */
func rsaDecrypt(content string, privateKey *rsa.PrivateKey) (string, error) {
	log.I("rsaDecrypt for content")
	rsaBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(content))
	if err != nil {
		log.E("rsaSign DecryptPKCS1v15 error")
		return "", err
	}
	return string(rsaBytes), nil
}

func base64EnCode(sign string) (string, error) {
	base64, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		fmt.Errorf("base64 DecodeString error \n")
		return "", err
	}
	return string(base64), nil
}
