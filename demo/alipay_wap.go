package demo

import (
	"fmt"
	"github.com/ljy2010a/go_alipay"
	"log"
	"net/http"
)

/**
请求支付
*/
func AlipayWapRequest(w http.ResponseWriter, r *http.Request) {

	alipayR := &alipay.AlipayWebRequest{
		OutTradeNo: NewOrderNo(), // 订单号
		Subject:    `waptest中文`,  // 商品名称 不能加空格
		TotalFee:   0.01,         // 价格
	}

	// 输出的是 html 页面，会自动跳转到支付界面
	err := alipay.AlipayWapRequest(alipay.AWapConfig, alipayR, w)
	if err != nil {
		return
	}
	return
}

//支付宝异步通知处理
func AlipayWapNotify(w http.ResponseWriter, r *http.Request) {

	log.Println("AlipayWapNotify Begin")
	var callbackMsg = "fail"
	defer func() {
		log.Println("AlipayWapNotify Notify End")
		log.Println("callbackMsg to alipay : %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.PostForm = nil
	r.ParseForm()

	log.Println("==========================================================")
	log.Println("AlipayWapNotify Request :%v", r)
	log.Println("==========================================================")

	alipayNotify, err := alipay.VerifyWapNotify(r, alipay.AWapConfig)
	if err != nil {
		//验证失败
		log.Println("verify notify fail")
		return
	}

	trade_status := alipayNotify.Trade_status
	out_trade_no := alipayNotify.Out_trade_no
	buyer_email := alipayNotify.Buyer_email
	subject := alipayNotify.Subject

	log.Println("trade_status is : %v ", trade_status)
	log.Println("out_trade_no is : %v ", out_trade_no)
	log.Println("buyer_email is : %v ", buyer_email)
	log.Println("subject is : %v ", subject)

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//该种交易状态只在一种情况下出现——开通了高级即时到账，买家付款成功后。

	if trade_status == "TRADE_SUCCESS" {

	}

	//判断是否已做操作

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//1、开通了普通即时到账，买家付款成功后。
	//该种交易状态只在两种情况下出现
	//2、开通了高级即时到账，从该笔交易成功时间算起，过了签约时的可退款时限（如：三个月以内可退款、一年以内可退款等）后。

	if trade_status == "TRADE_FINISHED" {

		log.Println("在这处理订单")
	}
	//	echo "success";		//请不要修改或删除
	callbackMsg = "success"
	return
}

//支付宝 同步通知处理
func AlipayWapCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("AlipayWapCallback Begin")

	var callbackMsg = "fail"
	defer func() {
		log.Println("AlipayWapCallback End")
		log.Println("callbackMsg to alipay : %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	//	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	r.PostForm = nil
	r.ParseForm()

	log.Println("==========================================================")
	log.Println("AlipayWapCallback Request :%v", r)
	log.Println("==========================================================")

	if err := alipay.VerifyWapCallback(r, alipay.AWapConfig); err != nil {
		//验证失败
		log.Println("verify notify fail")
		return
	}

	trade_status := r.FormValue("result")
	out_trade_no := r.FormValue("out_trade_no")

	log.Println("trade_status is : %v ", trade_status)
	log.Println("out_trade_no is : %v ", out_trade_no)

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	if trade_status == "success" {

	}
	callbackMsg = "success"
	return
}
