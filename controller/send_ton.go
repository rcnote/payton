package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"payton/pkg/blockchain"
)

var Ctrl = &Resp{}

type Resp struct{}

type (
	Request struct {
		ReceiverAddress string `json:"receiver_address"`
		Amount          uint64 `json:"amount"`
		Month           int    `json:"month"`
		Ref             string `json:"ref"`
	}
	Response struct {
		StatusCode int         `json:"status_code"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
		RequestID  string      `json:"request_id"`
	}
)

func (r *Resp) JSONS(ctx echo.Context, data interface{}, message ...string) error {
	rp := new(Response)
	rp.StatusCode = http.StatusOK
	if len(message) == 0 {
		rp.Message = "success"
	} else {
		for _, m := range message {
			rp.Message += "," + m
		}
	}
	rp.Data = data
	rp.RequestID = ctx.Request().Header.Get(echo.HeaderXRequestID)
	return ctx.JSON(http.StatusOK, rp)
}

func (r *Resp) SendTon(ctx echo.Context) (err error) {
	var s Request
	if err := ctx.Bind(&s); err != nil {
		d := &Response{
			StatusCode: 0,
		}
		return r.JSONS(ctx, d, "绑定参数失败")
	}

	// 验证api_key是否正确
	apiKey := ctx.Request().Header.Get("APIKEY")
	if apiKey != viper.GetString("api.key") {
		d := &Response{
			StatusCode: 0,
		}
		return r.JSONS(ctx, d, "APIKEY认证失败")
	}

	receiverAddress := s.ReceiverAddress
	if receiverAddress == "" {
		d := &Response{
			StatusCode: 0,
		}
		return r.JSONS(ctx, d, "参数错误, 未传入收款地址")
	}

	amount := s.Amount
	if amount < 100_000_000 {
		d := &Response{
			StatusCode: 0,
		}
		return r.JSONS(ctx, d, "参数错误, 数量单位是100000000")
	}

	month := s.Month
	orderSN := s.Ref
	if len(orderSN) < 9 {
		d := &Response{
			StatusCode: 0,
		}
		return r.JSONS(ctx, d, "参数错误, 当前Ref# 字符小于9，请检测是否解析错误")
	}

	var tonCommentFormats = map[int]string{
		3:  "Telegram Premium for 3 months \n\nRef#%s",
		6:  "Telegram Premium for 6 months \n\nRef#%s",
		12: "Telegram Premium for 1 year \n\nRef#%s",
	}
	common := fmt.Sprintf(tonCommentFormats[month], orderSN)

	err = blockchain.Transfer(receiverAddress, amount, common)
	if err != nil {
		d := &Response{
			StatusCode: 0,
		}
		return r.JSONS(ctx, d, "TON付款失败")
	}
	d := &Response{
		StatusCode: 1,
	}
	return r.JSONS(ctx, d, "TON付款成功")
}
