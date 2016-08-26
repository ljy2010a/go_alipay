package demo

import (
	"fmt"

	alipay "github.com/ljy2010a/go_alipay"
)

var private_key = `
-----BEGIN RSA PRIVATE KEY-----
xxxxxxxxx
-----END RSA PRIVATE KEY-----
`

var public_key = `
-----BEGIN PUBLIC KEY-----
xxxxxxxxx
-----END PUBLIC KEY-----
`

var partner = "xxxxxxxxx"

var key = "xxxxxxxxx"

var seller = "xxxxxxxxx@gmail.com"

func initAlipayConfig() {

	host := "xxxxxxxxx"

	alipay.AWebConfig = &alipay.AlipayConfig{
		Partner:        partner,
		Key:            key,
		Sign_type:      "MD5",
		Input_charset:  "utf-8",
		Cacert:         "Cacert",
		Transport:      "http",
		Service:        "create_direct_pay_by_user",
		Seller_id:      seller,
		Payment_type:   "1",
		Show_order_url: "/paymentStatus.html",
	}

	alipay.AMobileConfig = &alipay.AlipayConfig{
		Partner:             partner,
		Key:                 key,
		Sign_type:           "RSA",
		Private_key_path:    []byte(private_key),
		Ali_public_key_path: []byte(public_key),
		Input_charset:       "UTF-8",
		Cacert:              "Cacert",
		Transport:           "http",
		Service:             "mobile.securitypay.pay",
		Seller_id:           seller,
		Payment_type:        "1",
	}

	alipay.AWapConfig = &alipay.AlipayConfig{
		Partner:             partner,
		Key:                 key,
		Sign_type:           "MD5",
		Private_key_path:    []byte(private_key),
		Ali_public_key_path: []byte(public_key),
		Input_charset:       "utf-8",
		Transport:           "http",
		Service:             "alipay.wap.auth.authAndExecute",
		Wap_Service:         "alipay.wap.trade.create.direct",
		Seller_id:           seller,
		Show_order_url:      "/paymentStatus.html",
	}

	alipay.AMobileConfig.Notify_url = fmt.Sprintf("http://%v/alipay-mobile-notify", host)
	alipay.InitKeys(alipay.AMobileConfig)

	alipay.AWapConfig.Notify_url = fmt.Sprintf(
		"http://%v/alipay-wap-notify", host)
	alipay.AWapConfig.Wap_merchant_url = fmt.Sprintf(
		"http://%v/merchant", host)
	alipay.AWapConfig.Wap_callback_url = fmt.Sprintf(
		"http://%v/alipay-wap-callback", host)
	alipay.AWapConfig.Show_order_url = fmt.Sprintf(
		"http://%v/orderDetail", host)

	alipay.InitKeys(alipay.AWapConfig)

	// !!!!!! web no need to init !!!!!!

	alipay.AWebConfig.Notify_url = fmt.Sprintf(
		"http://%v/alipay-web-notify", host)
	alipay.AWebConfig.Return_url = fmt.Sprintf(
		"http://%v/alipay-web-return", host)
	alipay.AWebConfig.Show_order_url = fmt.Sprintf(
		"http://%v/orderDetail", host)

}
