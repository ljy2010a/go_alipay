package alipay

import (
	"encoding/xml"
	"log"
	"strings"
)

var alipayGatewayNew = `https://mapi.alipay.com/gateway.do?`

/**
 * WAP形式消息验证地址
 */
var wap_alipayGatewayNew = `http://wappaygw.alipay.com/service/rest.htm?`

/**
 * HTTPS形式消息验证地址
 */
var https_verify_url = "https://mapi.alipay.com/gateway.do?service=notify_verify&"

/**
 * HTTP形式消息验证地址
 */
var http_verify_url = "http://notify.alipay.com/trade/notify_query.do?"

/**
 * 异步通知时，对参数做固定排序
 * @param $para 排序前的参数组
 * @return 排序后的参数组
 */
func sortNotifyPara(para *Kvpairs) Kvpairs {
	new := Kvpairs{}
	for _, kv := range *para {
		if kv.K == "service" {
			new = append(new, kv)
		}
	}
	for _, kv := range *para {
		if kv.K == "v" {
			new = append(new, kv)
		}
	}
	for _, kv := range *para {
		if kv.K == "sec_id" {
			new = append(new, kv)
		}
	}
	for _, kv := range *para {
		if kv.K == "notify_data" {
			new = append(new, kv)
		}
	}
	return new
}

func ParseOriginTokenMsg(origin_token string, alipayConfig *AlipayConfig) string {

	tokenArray := strings.Split(origin_token, "&")
	kvs := Kvpairs{}
	tokenkv := Kvpair{}
	for _, v := range tokenArray {

		lenght := len(v)
		index := strings.Index(v, "=")
		if index != -1 {
			str1 := Substr(v, index+1, lenght-index-1)
			str2 := Substr(v, 0, index)
			subkv := Kvpair{str2, str1}
			kvs = append(kvs, subkv)
			if str2 == "res_data" {
				tokenkv = Kvpair{str2, str1}
			}
		}
	}

	//	token_data_Decrypt := ""

	//	for _, v := range kvs {
	//		if v.K == "sign_type" && v.V == "MD5" {
	//			token_data_Decrypt, err := rsaDecrypt(tokenkv.V, config.Private_key)
	//			if nil != err {
	//				return token_data_Decrypt
	//			}
	//		}
	//	}

	// log.Println("kvs = %v", kvs)

	// log.Println("tokenkv.V = %v", tokenkv.V)

	if tokenkv.V == "" {
		return ""
	}

	type Direct_trade_create_res struct {
		Token string `xml:"request_token"`
	}

	var alipay Direct_trade_create_res

	err := xml.Unmarshal([]byte(tokenkv.V), &alipay)
	if err != nil {
		log.Println("xml Unmarshal err : %v", err)
		return ""
	}

	log.Println("alipay xml : %v", alipay)

	if alipay.Token == "" {
		log.Println("token is null  ")
		return ""
	}

	return alipay.Token
}
