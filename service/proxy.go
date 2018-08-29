package service

import (
	"net/http"

	"github.com/mafanr/admin/misc"

	"github.com/mafanr/g"

	"time"

	"github.com/sunface/talent"

	"go.uber.org/zap"

	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
)

func (a *Admin) proxy(c echo.Context) error {
	// 解析参数
	args := &fasthttp.Args{}

	ps, _ := c.FormParams()
	for k, v := range ps {
		args.Set(k, v[0])
	}
	// 获取当前登陆用户的信息
	token := getToken(c)
	if token == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}
	sessI, ok := a.sessions.Load(token)
	if !ok {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}
	sess := sessI.(*Session)
	args.Set("user_id", sess.User.ID)
	args.Set("user_name", sess.User.Name)

	// 设置admin token
	args.Set("admin_token", misc.Conf.Common.Token)
	method := c.Request().Method

	// 查询节点ip
	srv := c.FormValue("target_service")
	service, ok := services[srv]
	if !ok || len(service) == 0 {
		g.L.Info("没有查询到对应后台接口的节点信息", zap.String("service", srv))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusNotFound,
			ErrCode: g.TargetServiceNotFountC,
			Message: g.TargetServiceNotFountE,
		})
	}

	// 设置应用访问角色,是否是admin
	args.Set("app_role", g.ROLE_NORMAL)
	aa, ok := AppAdmins[srv]
	if ok {
		role, ok1 := aa[sess.User.ID]
		if ok1 {
			args.Set("app_role", role)
		}
	}
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.Header.SetMethod(method)

	ip := service[0].IP
	path := c.FormValue("target_path")
	url := "http://" + ip + path

	rid := c.Get("rid").(int64)
	g.L.Info("开始请求", zap.Int64("rid", rid), zap.String("target_service", srv), zap.String("target_path", path))

	switch method {
	case "GET":
		// 拼接url
		url = url + "?" + args.String()
	default:
		args.WriteTo(req.BodyWriter())
	}
	req.SetRequestURI(url)

	err := g.Cli.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		g.L.Info("请求后台服务接口出错", zap.Error(err), zap.String("url", url))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadGateway,
			Message: "请求后台服务接口出错",
		})
	}

	return c.String(http.StatusOK, talent.Bytes2String(resp.Body()))
}
