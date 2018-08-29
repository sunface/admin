package service

import (
	"sync"
	"time"

	"github.com/mafanr/admin/misc"

	"github.com/mafanr/g"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

type Admin struct {
	sessions *sync.Map
}

func (a *Admin) Start() {
	// 查询所有service的服务节点状态
	go queryServices()

	a.sessions = &sync.Map{}
	a.listen()
}

func (a *Admin) Shutdown() {
	g.L.Info("shutdown tfe..")
}

func (a *Admin) listen() {
	e := echo.New()
	// 设置跨域
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Token"},
		AllowCredentials: true,
	}))

	// sso登陆和session相关
	e.POST("/login", a.login, timing)
	e.POST("/login/mock", a.loginMock, timing)
	e.POST("/logout", a.logout, timing)
	e.GET("/user/info", a.userInfo, timing)

	// tools
	e.POST("/tools/praise", a.genPraise, timing)
	e.POST("/tools/testApi", a.testApi)

	// 通用管理接口，只做代理，不做逻辑
	e.Any("/admin", a.proxy, timing)
	e.Logger.Fatal(e.Start(":" + misc.Conf.Common.Port))
}

func timing(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ts := time.Now()
		rid := (ts.UnixNano()/10)*10 + misc.Conf.Common.ServerID

		g.L.Info("收到新的请求", zap.Int64("rid", rid), zap.String("path", c.Path()), zap.String("ip", c.RealIP()))
		c.Set("rid", rid)
		defer func() {
			g.L.Info("请求完成", zap.Int64("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000000))
		}()

		return f(c)
	}
}

var services map[string][]*g.ServerInfo

func queryServices() {
	g.EtcdCli = g.InitEtcd(misc.Conf.Etcd.Addrs)
	var err error

	go func() {
		for {
			services, err = g.QueryAllService(g.EtcdCli)
			if err != nil {
				g.L.Error("获取服务节点状态出错", zap.Error(err))
			}

			time.Sleep(g.ServiceQueryInterval * time.Second)
		}
	}()

}
