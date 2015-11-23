package alipay

type AlipayWebRequest struct {
	OutTradeNo string
	Subject    string
	TotalFee   float64
	Body       string
	ShowUrl    string
}

type Notify struct {
	Payment_type        string `xml:"payment_type"`
	Subject             string `xml:"subject"`
	Trade_no            string `xml:"trade_no"`
	Buyer_email         string `xml:"buyer_email"`
	Gmt_create          string `xml:"gmt_create"`
	Notify_type         string `xml:"notify_type"`
	Quantity            string `xml:"quantity"`
	Out_trade_no        string `xml:"out_trade_no"`
	Seller_id           string `xml:"seller_id"`
	Trade_status        string `xml:"trade_status"`
	Is_total_fee_adjust string `xml:"is_total_fee_adjust"`
	Gmt_payment         string `xml:"gmt_payment"`
	Seller_email        string `xml:"seller_email"`
	Price               string `xml:"price"`
	Buyer_id            string `xml:"buyer_id"`
	Notify_id           string `xml:"notify_id"`
	Use_coupon          string `xml:"use_coupon"`
}

type AlipayMobileRequest struct {
	OutTradeNo string
	Subject    string
	TotalFee   float64
	Body       string
	ShowUrl    string
}
