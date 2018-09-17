package admin

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
	args.Set("username", sess.LoginInfo.Username)

	// 设置admin token
	args.Set("admin_token", misc.Conf.Common.Token)
	method := c.Request().Method

	// 查询目标节点
	srv := c.FormValue("target_app")
	s := g.GetServer(srv)
	if s == nil {
		g.L.Info("No application node founded", zap.String("targe_app", srv))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusNotFound,
			ErrCode: g.TargetServiceNotFountC,
			Message: g.TargetServiceNotFountE,
		})
	}

	service := c.FormValue("target_service")
	// 设置应用访问角色,是否是admin
	var appPriv string
	if isServiceOwner(service, sess.LoginInfo.Username) {
		appPriv = g.PRIV_OWNER
	} else if sess.LoginInfo.Priv == g.PRIV_ADMIN || sess.LoginInfo.Priv == g.PRIV_SUPER_ADMIN || isServiceAdmin(service, sess.LoginInfo.Username) {
		appPriv = g.PRIV_ADMIN
	} else if isServiceUser(service, sess.LoginInfo.Username) {
		appPriv = g.PRIV_NORMAL
	} else {
		appPriv = g.PRIV_GUEST
	}

	args.Set("app_priv", appPriv)

	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.Header.SetMethod(method)

	ip := s.IP
	path := c.FormValue("target_path")
	url := "http://" + ip + path

	rid := c.Get("rid").(int64)
	g.L.Info("Start requesting", zap.Int64("rid", rid), zap.String("target_app", srv), zap.String("target_path", path))

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
		g.L.Info("Request error", zap.Error(err), zap.String("url", url))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadGateway,
			ErrCode: g.ReqFailedC,
			Message: g.ReqFailedE,
		})
	}

	return c.String(http.StatusOK, talent.Bytes2String(resp.Body()))
}
