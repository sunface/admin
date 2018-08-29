package service

import (
	"net/http"
	"time"

	"github.com/mafanr/g"

	"go.uber.org/zap"

	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
)

var defaultRole = "normal"

type Session struct {
	User       *UserInfo
	CreateTime time.Time
}

type UserInfo struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Avatar   string   `json:"avatar"`
	Role     []string `json:"role"`
	SsoToken string   `json:"ssoToken"`
}

func (a *Admin) login(c echo.Context) error {
	role := getGlobalRole("13269")
	user := &UserInfo{
		ID:       "13269",
		Name:     "孙飞",
		Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Role:     []string{role},
		SsoToken: "",
	}

	//sub token验证成功，保存session
	a.sessions.Store(user.SsoToken, &Session{
		User:       user,
		CreateTime: time.Now(),
	})

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (a *Admin) loginMock(c echo.Context) error {
	role := getGlobalRole("13269")
	user := &UserInfo{
		ID:       "13269",
		Name:     "孙飞",
		Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Role:     []string{role},
		SsoToken: "0af8c18eed353af38b2e2524f4850f76",
	}

	//sub token验证成功，保存session
	a.sessions.Store(user.SsoToken, &Session{
		User:       user,
		CreateTime: time.Now(),
	})

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (a *Admin) logout(c echo.Context) error {
	token := getToken(c)
	// 删除用户的session
	a.sessions.Delete(token)

	return c.JSON(http.StatusOK, g.Result{
		Status:  http.StatusOK,
		Message: "退出登陆成功",
	})
}

func (a *Admin) userInfo(c echo.Context) error {
	token := getToken(c)
	sess, ok := a.sessions.Load(token)
	if !ok {
		// 用户未登陆或者session失效
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   sess.(*Session).User,
	})
}

//subToken认证，返回的用户信息
type TokenInfo struct {
	Code int `json:"code"`
	Data struct {
		SubTokenObj struct {
			Message     string `json:"message"`
			SsoToken    string `json:"ssoToken"`
			Status      int    `json:"status"`
			UserSession struct {
				Facility        string `json:"facility"`
				HeadImgUrl      string `json:"headImgUrl"`
				SysDepartmentId string `json:"sysDepartmentId"`
				UserId          string `json:"userId"`
				UserName        string `json:"userName"`
				UserType        string `json:"userType"`
			} `json:"userSession"`
		} `json:"subTokenObj"`
	} `json:"data"`
	Message string `json:"message"`
}

func requestToSso(body string, url string) []byte {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBodyString(body)

	req.SetRequestURI(url)
	var cli = &fasthttp.Client{}
	err := cli.DoTimeout(req, resp, 15*time.Second)
	if err != nil {
		g.L.Info("获取sso用户信息失败", zap.Error(err))
		return nil
	}

	return resp.Body()
}

func getToken(c echo.Context) string {
	return c.Request().Header.Get("X-Token")
}

func getGlobalRole(u string) string {
	role := defaultRole
	_, ok := GlobalAdmins[u]
	if ok {
		role = "admin"
	}

	return role
}
