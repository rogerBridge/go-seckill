/*
使用fasthttp制造http请求
*/
package pressuremaker

import (
	//"encoding/json"

	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type Order struct {
	Token       string
	ProductID   int
	PurchaseNum int
}

func (o *Order) CreateOrder(w *sync.WaitGroup, timeStatistics chan float64, errChan chan error) (bool, error) {
	// 首先, 构造client
	client := FastHttpClient
	host := "http://go-seckill:4000"
	var URL = host + "/api/v0/user/order/buy"

	req := &fasthttp.Request{}

	// req := fasthttp.AcquireRequest()
	// defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(http.MethodPost)
	req.Header.Set("Authorization", o.Token)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(URL)
	// 构造request body里面的值
	r := ReqBuy{
		ProductId:   o.ProductID,
		PurchaseNum: o.PurchaseNum,
	}
	reqBody, err := json.Marshal(r)
	if err != nil {
		logger.Fatalf("Marshal struct to []byte error message %v", err)
		return false, err
	}
	req.SetBody(reqBody)

	resp := &fasthttp.Response{}

	// resp := fasthttp.AcquireResponse()
	// defer fasthttp.ReleaseResponse(resp)

	//resp := fasthttp.AcquireResponse()
	// 开始发送请求
	t0 := time.Now() // 客户端开始发起请求的时间

	err = client.Do(req, resp)
	if err != nil {
		errChan <- fmt.Errorf("client do error %v", err)
		w.Done()
		logger.Warnf("发送请求时: %v", err)
		return false, err
	}
	// if resp.StatusCode() != 200 || len(resp.Body()) == 0
	if resp.StatusCode() != 200 {
		errChan <- fmt.Errorf("server response status code error")
		w.Done()
		return false, err
	}
	t1 := time.Since(t0)           // 客户端结束请求的时间
	timeStatistics <- t1.Seconds() // 将客户端整个请求的时间段发送给timeStatistics
	// 请求结束了
	// see this: [fasthttp resp.Body not use any io.Reader](https://github.com/valyala/fasthttp/issues/411#issuecomment-420857389)
	w.Done()
	return true, nil
}
