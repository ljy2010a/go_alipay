package alipay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func AlipayWapRequest(alipayConfig *AlipayConfig, r *AlipayWebRequest, w io.Writer) error {

	var format = "xml"
	var v = "2.0"
	args := []interface{}{}
	args = append(args, alipayConfig.Notify_url)
	args = append(args, alipayConfig.Wap_callback_url)
	args = append(args, alipayConfig.Seller_id)
	args = append(args, r.OutTradeNo)
	args = append(args, r.Subject)
	args = append(args, fmt.Sprintf("%.2f", r.TotalFee))
	args = append(args, alipayConfig.Wap_merchant_url)

	var token_req_data = fmt.Sprintf("<direct_trade_create_req><notify_url>%v</notify_url><call_back_url>%v</call_back_url><seller_account_name>%v</seller_account_name><out_trade_no>%v</out_trade_no><subject>%v</subject><total_fee>%v</total_fee><merchant_url>%v</merchant_url></direct_trade_create_req>", args...)

	log.Println("token_req_data = %v", token_req_data)
	p := Kvpairs{
		Kvpair{`service`, alipayConfig.Wap_Service},
		Kvpair{`partner`, alipayConfig.Partner},
		Kvpair{`sec_id`, alipayConfig.Sign_type},
		Kvpair{`format`, format},
		Kvpair{`v`, v},
		Kvpair{`req_id`, r.OutTradeNo},
		Kvpair{`req_data`, token_req_data},
		Kvpair{`_input_charset`, alipayConfig.Input_charset},
	}

	paraFilter(&p)

	argSort(&p)

	buildRequestPara(&p, alipayConfig)

	log.Println("pararms is  = %v", p)

	origin_token, err := getHttpResponsePOST(wap_alipayGatewayNew, "", alipayConfig.Input_charset, &p)
	if nil != err {
		log.Println("get token fail : %v", err)
		return err
	}

	log.Println("origin_token = %v", origin_token)
	origin_token, _ = url.QueryUnescape(origin_token)
	log.Println("origin_token url.QueryUnescape  = %v", origin_token)

	token := ParseOriginTokenMsg(origin_token, alipayConfig)
	if token == "" {
		fmt.Fprintln(w, origin_token)
		return errors.New("解析token失败")
	}

	var trade_req_data = fmt.Sprintf("<auth_and_execute_req><request_token>%v</request_token></auth_and_execute_req>", token)

	p2 := Kvpairs{
		Kvpair{`service`, alipayConfig.Service},
		Kvpair{`partner`, alipayConfig.Partner},
		Kvpair{`sec_id`, alipayConfig.Sign_type},
		Kvpair{`format`, format},
		Kvpair{`v`, v},
		Kvpair{`req_id`, r.OutTradeNo},
		Kvpair{`req_data`, trade_req_data},
		Kvpair{`_input_charset`, alipayConfig.Input_charset},
	}
	paraFilter(&p2)

	argSort(&p2)

	buildRequestPara(&p2, alipayConfig)

	// log.Println("pararm2 = %v", p2)

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=utf-8' method='get'> `, wap_alipayGatewayNew)
	for _, kv := range p2 {
		fmt.Fprintf(w, `<input type='hidden' name='%s' value='%s' />`, kv.K, kv.V)
	}
	fmt.Fprintln(w, `<script>document.forms['alipaysubmit'].submit();</script>`)
	fmt.Fprintln(w, `</body></html>`)
	return nil
}

/**
 * 验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWapNotify(r *http.Request, alipayConfig *AlipayConfig) (Notify, error) {
	log.Println("VerifyWapNotify begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	notify_data := ""
	var alipayNotify Notify

	for k := range r.Form {
		v := r.Form.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		case "notify_data":
			notify_data = v
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	if notify_data == "" {
		log.Println("notify_data is null")
		return alipayNotify, errors.New("notify_data is null")
	}

	err := xml.Unmarshal([]byte(notify_data), &alipayNotify)
	if err != nil {
		log.Println("xml Unmarshal err : %v", err)
		return alipayNotify, err
	}

	// log.Println("p = %v", p)
	p2 := sortNotifyPara(p)
	// log.Println("p2 = %v", p2)
	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(&p2)

	log.Println("VerifyWapNotify prestr is : %v ", prestr)
	log.Println("VerifyWapNotify sign is : %v  , sign_type is %v", sign, sign_type)

	switch alipayConfig.Sign_type {
	case "MD5":
		if md5Sign(prestr, alipayConfig.Key) != sign {
			return alipayNotify, fmt.Errorf("sign invalid")
		}
		break
	default:
		return alipayNotify, fmt.Errorf("no right sign_type")
	}

	log.Println("VerifyWapNotify Notify is %v", alipayNotify)

	notify_id := alipayNotify.Notify_id
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, alipayConfig)
	if err != nil {
		return alipayNotify, err
	}
	log.Println("VerifyWapNotify responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("VerifyWapNotify responseTxt verify fail ")
		return alipayNotify, fmt.Errorf("VerifyWapNotify responseTxt is wrong")
	}
	log.Println("VerifyWapNotify responseTxt verify success ")
	return alipayNotify, nil
}

/**
 * 验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWapCallback(r *http.Request, alipayConfig *AlipayConfig) error {
	log.Println("VerifyWapCallback begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.Form.Get(k)
		switch k {
		case "sign":
			sign = v
			continue
		case "sign_type":
			sign_type = v
			continue
		}
		*p = append(*p, Kvpair{k, v})
	}
	//除去待签名参数数组中的空值和签名参数
	paraFilter(p)

	//对待签名参数数组排序
	argSort(p)

	// log.Println("VerifyWapCallback p = %v", p)
	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.Println("VerifyWapCallback prestr is : %v ", prestr)
	log.Println("VerifyWapCallback sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, alipayConfig.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.Println("VerifyWapCallback success")
	return nil
}
