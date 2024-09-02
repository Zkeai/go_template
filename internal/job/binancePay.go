package binancepay

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// BinancePay 是一个包含API密钥和基础URL的结构体
type BinancePay struct {
	APIKey        string
	APISecret     string
	BaseURL       string
	CertificateSN string
}

// NewBinancePay 是一个构造函数，返回一个新的BinancePay实例
func NewBinancePay(apiKey, apiSecret, certificateSN string) *BinancePay {
	return &BinancePay{
		APIKey:        apiKey,
		APISecret:     apiSecret,
		BaseURL:       "https://bpay.binanceapi.com",
		CertificateSN: certificateSN,
	}
}

// generateSignature 用于生成HMAC SHA256签名
func (bp *BinancePay) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(bp.APISecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// sendRequest 发送HTTP请求并返回响应
func (bp *BinancePay) sendRequest(endpoint string, method string, payload map[string]interface{}) ([]byte, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// 将payload序列化为JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest(method, bp.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	nonce := fmt.Sprintf("%d", time.Now().UnixNano()) // 生成唯一随机字符串
	signature := bp.generateSignature(string(jsonData))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("BinancePay-Timestamp", timestamp)
	req.Header.Set("BinancePay-Nonce", nonce)
	req.Header.Set("BinancePay-Certificate-SN", bp.CertificateSN)
	req.Header.Set("BinancePay-Signature", signature)
	req.Header.Set("BinancePay-API-Key", bp.APIKey)

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", body)
	}

	return body, nil
}

// CreateOrder 创建支付订单并返回响应
func (bp *BinancePay) CreateOrder(orderId string, orderAmount string, currency string) ([]byte, error) {
	endpoint := "/binancepay/openapi/v3/order"
	method := "POST"

	payload := map[string]interface{}{
		"merchantId":  "你的商户ID", // 替换为实际的商户ID
		"orderAmount": orderAmount,
		"currency":    currency,
		"orderId":     orderId,
		"goods": map[string]string{
			"goodsType":        "01",   // 商品类型
			"goodsCategory":    "D000", // 商品类别
			"referenceGoodsId": "商品ID", // 商品ID
			"goodsName":        "商品名称", // 商品名称
		},
	}

	response, err := bp.sendRequest(endpoint, method, payload)
	if err != nil {
		return nil, err
	}

	return response, nil
}
