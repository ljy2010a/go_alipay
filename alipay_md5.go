package alipay

import (
	"crypto/md5"
	"fmt"
	"io"
)

/**
 * 签名字符串
 * @param $prestr 需要签名的字符串
 * @param $key 私钥
 * return 签名结果
 */
func md5Sign(str, key string) string {
	h := md5.New()
	io.WriteString(h, str)
	io.WriteString(h, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
 * 验证签名
 * @param $prestr 需要签名的字符串
 * @param $sign 签名结果
 * @param $key 私钥
 * return 签名结果
 */
//func md5Verify(prestr, sign, key string) (bool) {
//	h := md5.New()
//	io.WriteString(h, prestr)
//	io.WriteString(h, key)
//	mysgin := fmt.Sprintf("%x", h.Sum(nil))
//
//	if mysgin == sign {
//		return true
//	}
//	return false
//}
