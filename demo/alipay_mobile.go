package demo

import (
	"encoding/json"
	"fmt"
	"github.com/ljy2010a/go_alipay"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func NewOrderNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), RandInt(10000, 99999))
}

func RandInt(min int, max int) int {
	if max-min <= 0 {
		return min
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

//POST
//支付宝回调处理
func AlipayMobileNotify(w http.ResponseWriter, r *http.Request) {
	log.Panicln("alipay Notify Begin")

	var callbackMsg = "fail"
	defer func() {
		log.Panicln("alipay Notify End")
		log.Panicln("callbackMsg to alipay : %v", callbackMsg)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, callbackMsg)
	}()

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.PostForm = nil
	r.ParseForm()

	log.Panicln("==========================================================")
	log.Panicln("Request :%v", r)
	log.Panicln("==========================================================")

	if err := alipay.VerifyMobileNotify(r, alipay.AMobileConfig); err != nil {
		//验证失败
		log.Panicln("verify notify fail")
		return
	}

	trade_status := r.FormValue("trade_status")
	out_trade_no := r.FormValue("out_trade_no")
	buyer_email := r.FormValue("buyer_email")
	subject := r.FormValue("subject")

	log.Panicln("trade_status is : %v ", trade_status)
	log.Panicln("out_trade_no is : %v ", out_trade_no)
	log.Panicln("buyer_email is : %v ", buyer_email)
	log.Panicln("subject is : %v ", subject)

	var total_fee float64
	fmt.Sscanf(r.FormValue("total_fee"), "%f", &total_fee)

	//判断该笔订单是否在商户网站中已经做过处理
	//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
	//如果有做过处理，不执行商户的业务程序

	//注意：
	//该种交易状态只在一种情况下出现——开通了高级即时到账，买家付款成功后。

	if trade_status == "TRADE_SUCCESS" {

	}

	if trade_status == "TRADE_FINISHED" {
		//判断是否已做操作

		//判断该笔订单是否在商户网站中已经做过处理
		//如果没有做过处理，根据订单号（out_trade_no）在商户网站的订单系统中查到该笔订单的详细，并执行商户的业务程序
		//如果有做过处理，不执行商户的业务程序

		//注意：
		//该种交易状态只在两种情况下出现
		//1、开通了普通即时到账，买家付款成功后。
		//2、开通了高级即时到账，从该笔交易成功时间算起，过了签约时的可退款时限（如：三个月以内可退款、一年以内可退款等）后。

		log.Println("在这处理订单")
	}
	//	echo "success";		//请不要修改或删除
	callbackMsg = "success"
	return
}

//获取支付宝签名
func GetRsaSign(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	amr := alipay.AlipayMobileRequest{}
	amr.OutTradeNo = NewOrderNo()
	amr.Subject = "测试"
	amr.Body = "测试"
	amr.TotalFee = 0.01

	orderinfo := alipay.AlipayMobileRsaSign(amr, alipay.AMobileConfig)

	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	rs["data"] = orderinfo //createLinkString(&p)

	b, _ := json.Marshal(rs)
	// http.Error(w, string(b), 200)
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	fmt.Fprint(w, string(b))
}
