package service

import (
	"sync"
	"time"

	"github.com/mafanr/admin/misc"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mafanr/g"
	"go.uber.org/zap"
)

type Admin struct {
	sessions *sync.Map
}

func (a *Admin) Start() {
	g.L.Info("start mafanr admin")

	// 初始化mysql连接
	g.InitMysql(misc.Conf.Mysql.Acc, misc.Conf.Mysql.Pw, misc.Conf.Mysql.Addr, misc.Conf.Mysql.Port, misc.Conf.Mysql.Database)

	// 初始化超级管理员
	a.InitSuperAdmin()

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
	e.POST("/logout", a.logout, timing)
	e.GET("/login/info", a.loginInfo, timing)

	// 用户相关
	e.POST("/user/add", a.AddUser, timing)
	e.POST("/user/edit", a.EditUser, timing)
	e.POST("/user/delete", a.DeleteUser, timing)
	e.GET("/user/load", a.LoadUsers, timing)
	e.POST("/user/setPW", a.SetPW, timing)
	e.POST("/user/setUserPW", a.SetUserPW, timing)
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
