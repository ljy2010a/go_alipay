package alipay

import (
	"net/url"
	"sort"
	"strings"
)

type Kvpair struct {
	K, V string
}

type Kvpairs []Kvpair

func (t Kvpairs) Less(i, j int) bool {
	return t[i].K < t[j].K
}

func (t Kvpairs) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Kvpairs) Len() int {
	return len(t)
}

func (t Kvpairs) Sort() {
	sort.Sort(t)
}

/**
 * 除去数组中的空值和签名参数
 * @param $para 签名参数组
 * return 去掉空值与签名参数后的新签名参数组
 */
func paraFilter(para *Kvpairs) {
	new := Kvpairs{}
	for _, kv := range *para {
		if kv.V != "" {
			new = append(new, kv)
		}
	}
	*para = new
}

/**
 * 对数组排序
 * @param $para 排序前的数组
 * return 排序后的数组
 */
func argSort(para *Kvpairs) {
	para.Sort()
}

/**
 * 把数组所有元素，按照“参数=参数值”的模式用“&”字符拼接成字符串，并对字符串做urlencode编码
 * @param $para 需要拼接的数组
 * return 拼接完成以后的字符串
 */
func createLinkstringUrlencode(para *Kvpairs) string {
	var strs []string
	for _, kv := range *para {
		if kv.K == "notify_url" || kv.K == "sign" {
			//移动支付这里要做URL转换，加双引号  坑
			strs = append(strs, kv.K+"=\""+url.QueryEscape(kv.V)+"\"")
		} else {
			strs = append(strs, kv.K+"=\""+kv.V+"\"")
		}

	}
	return strings.Join(strs, "&")
}

func createLinkstringForPost(para *Kvpairs) string {
	var strs []string
	for _, kv := range *para {
		if kv.K == "notify_url" || kv.K == "sign" {
			//移动支付这里要做URL转换，加双引号  坑
			strs = append(strs, kv.K+"="+url.QueryEscape(kv.V)+"")
		} else {
			strs = append(strs, kv.K+"="+kv.V+"")
		}

	}
	return strings.Join(strs, "&")
}

func createLinkStringNoUrl(para *Kvpairs) string {
	var strs []string
	for _, kv := range *para {
		//验签这里不用转 坑
		strs = append(strs, kv.K+"="+kv.V)
	}
	return strings.Join(strs, "&")
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
