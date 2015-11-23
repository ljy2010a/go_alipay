package alipay

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

func AlipayWebRequestForm(alipayConfig *AlipayConfig, r *AlipayWebRequest, w io.Writer) error {
	p := Kvpairs{
		Kvpair{`total_fee`, fmt.Sprintf("%.2f", r.TotalFee)},
		Kvpair{`subject`, r.Subject},
		Kvpair{`body`, r.Body},
		Kvpair{`show_url`, r.ShowUrl},
		Kvpair{`out_trade_no`, r.OutTradeNo},
		Kvpair{`service`, alipayConfig.Service},
		Kvpair{`partner`, alipayConfig.Partner},
		Kvpair{`payment_type`, alipayConfig.Payment_type},
		Kvpair{`notify_url`, alipayConfig.Notify_url},
		Kvpair{`return_url`, alipayConfig.Return_url},
		Kvpair{`seller_email`, alipayConfig.Seller_id},
		Kvpair{`_input_charset`, alipayConfig.Input_charset},
	}

	paraFilter(&p)
	argSort(&p)
	sign := md5Sign(createLinkStringNoUrl(&p), alipayConfig.Key)

	p = append(p, Kvpair{`sign`, sign})
	p = append(p, Kvpair{`sign_type`, `MD5`})

	fmt.Fprintln(w, `<html><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head><body>`)
	fmt.Fprintf(w, `<form name='alipaysubmit' action='%s_input_charset=utf-8' method='post'> `, alipayGatewayNew)
	for _, kv := range p {
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
func VerifyWebReturn(r *http.Request, alipayConfig *AlipayConfig) error {
	log.Println("VerifyWebReturn begin")

	p := &Kvpairs{}
	sign := ""
	sign_type := ""
	for k := range r.Form {
		v := r.PostForm.Get(k)
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

	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.Println("VerifyWebReturn prestr is : %v ", prestr)
	log.Println("VerifyWebReturn sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, alipayConfig.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.Println("VerifyWebReturn success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, alipayConfig)
	if err != nil {
		return err
	}
	log.Println("VerifyWebReturn responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.Println("VerifyWebReturn responseTxt verify success ")
	return nil

}

/**
 * 针对notify_url验证消息是否是支付宝发出的合法消息
 * @return 验证结果
 */
func VerifyWebNotify(r *http.Request, alipayConfig *AlipayConfig) error {
	log.Println("VerifyWebNotify begin")

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

	//把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串
	prestr := createLinkStringNoUrl(p)

	log.Println("VerifyWebNotify prestr is : %v ", prestr)
	log.Println("VerifyWebNotify sign is : %v  , sign_type is %v", sign, sign_type)

	switch sign_type {
	case "MD5":
		if md5Sign(prestr, alipayConfig.Key) != sign {
			return fmt.Errorf("sign invalid")
		}
		break
	default:
		return fmt.Errorf("no right sign_type")
	}

	log.Println("VerifyWebNotify success")

	notify_id := r.FormValue("notify_id")
	//获取支付宝远程服务器ATN结果（验证是否是支付宝发来的消息）(1分钟认证)
	responseTxt, err := getResponse(notify_id, alipayConfig)
	if err != nil {
		return err
	}
	log.Println("VerifyWebNotify responseTxt is: %v", responseTxt)

	reg := regexp.MustCompile(`true`)
	if 0 == len(reg.FindAllString(responseTxt, -1)) {
		log.Println("responseTxt verify fail ")
		return fmt.Errorf("responseTxt is wrong")
	}
	log.Println("VerifyWebNotify responseTxt verify success ")
	return nil
}
