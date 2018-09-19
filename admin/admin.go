package admin

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sunface/talent"

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

	a.InitAdminService()

	// 查询所有service的服务节点状态
	g.ETCD.QueryAll(misc.Conf.Etcd.Addrs)

	a.sessions = &sync.Map{}
	a.listen()
}

func (a *Admin) Shutdown() {
	g.L.Info("shutdown tfe..")
}

func (a *Admin) listen() {
	go func() {
		e := echo.New()
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}))

		e.Static("/images", "./images")

		e.Static("/adminui", "./ui/dist")
		e.Static("/lab", "../labsite/dist")

		// 这里把static目录单独独立出来，是为了底下的rewrite不报错
		e.Static("/adminStatic", "./ui/dist")
		e.Static("/labStatic", "../labsite/dist")

		// 为了让vuejs的mode: history模式下，我们可以访问子路径，这里需要做重写
		e.Pre(middleware.Rewrite(map[string]string{
			"/lab/*":     "/lab",
			"/adminui/*": "/adminui",
		}))

		// 其它路径都路由到/lab
		e.GET("/", func(c echo.Context) error {
			fmt.Println(c.Request().Host)
			return c.Redirect(http.StatusMovedPermanently, "/lab")
		})

		e.GET("/:any", func(c echo.Context) error {
			return c.Redirect(http.StatusMovedPermanently, "/lab")
		})
		e.Logger.Fatal(e.Start(":" + misc.Conf.Static.Port))
	}()

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
	e.POST("/tools/testApi", a.testApi)

	//service管理
	e.GET("/ops/service/query", a.QueryServices, timing)
	e.POST("/ops/service/create", a.CreateService, timing)
	e.POST("/ops/service/update", a.UpdateService, timing)
	e.GET("/ops/service/queryUser", a.QueryServiceUser, timing)
	e.POST("/ops/service/addUser", a.AddServiceUser, timing)
	e.POST("/ops/service/delUser", a.DelServiceUser, timing)
	// application manage
	e.POST("/ops/app/create", a.CreateAPP, timing)
	e.POST("/ops/app/update", a.UpdateAPP, timing)
	e.GET("/ops/app/query", a.QueryApps, timing)
	// server manage
	e.POST("/ops/server/define", a.DefineServer, timing)
	e.GET("/ops/server/query", a.QueryServers, timing)
	e.GET("/ops/server/queryPW", a.QueryServerPW, timing)

	// 通用管理接口，只做代理，不做逻辑
	e.Any("/admin", a.proxy, timing)

	e.Logger.Fatal(e.Start(":" + misc.Conf.Common.Port))
}

func timing(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ts := time.Now()
		rid := (ts.UnixNano()/10)*10 + misc.Conf.Common.ServerID

		g.L.Debug("New request accepted", zap.Int64("rid", rid), zap.String("path", c.Path()), zap.String("ip", c.RealIP()))
		c.Set("rid", rid)
		defer func() {
			g.L.Debug("Request done", zap.Int64("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000000))
		}()

		return f(c)
	}
}

func (a *Admin) InitAdminService() {
	g.DB.Exec(fmt.Sprintf("insert into service (name,owner,description,create_date) values ('admin','admin','admin service','%s')", talent.Time2StringSecond(time.Now())))
	g.DB.Exec(fmt.Sprintf("insert into application (name,service,description,create_date) values ('admin','admin','admin application','%s')", talent.Time2StringSecond(time.Now())))

	g.DB.Exec(fmt.Sprintf("insert into service (name,owner,description,create_date) values ('juz','admin','juz api-gateway service','%s')", talent.Time2StringSecond(time.Now())))
	g.DB.Exec(fmt.Sprintf("insert into application (name,service,description,create_date) values ('juzManage','juz','juz api-gateway manager application','%s')", talent.Time2StringSecond(time.Now())))
}
