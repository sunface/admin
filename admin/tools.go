package admin

import (
	"encoding/json"
	"net/http"

	"github.com/mafanr/g"

	"time"

	"github.com/labstack/echo"
	"github.com/sunface/talent"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (a *Admin) testApi(c echo.Context) error {
	paramsS := c.FormValue("params")
	var temp map[string]interface{}
	json.Unmarshal([]byte(paramsS), &temp)

	// 解析参数
	args := &fasthttp.Args{}
	for k, vi := range temp {
		v, err := talent.Interface2String(vi)
		if err != nil {
			g.L.Info("invalid test param", zap.String("params", paramsS))
			continue
		}
		args.Set(k, v)
	}

	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.Header.SetMethod("POST")

	tfeAddr := c.FormValue("tfe_addr")

	args.WriteTo(req.BodyWriter())
	req.SetRequestURI(tfeAddr)

	err := g.Cli.DoTimeout(req, resp, 15*time.Second)
	if err != nil {
		g.L.Info("请求tfe出错", zap.Error(err), zap.String("url", tfeAddr))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusServiceUnavailable,
			Message: "请求TFE出错",
		})
	}

	if resp.StatusCode() != http.StatusOK {
		g.L.Info("请求tfe异常", zap.Int("code", resp.StatusCode()), zap.String("url", tfeAddr))
		return c.JSON(http.StatusOK, g.Result{
			Status:  resp.StatusCode(),
			Message: talent.Bytes2String(resp.Body()),
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   talent.Bytes2String(resp.Body()),
	})
}
