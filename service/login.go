package service

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mafanr/g"

	"go.uber.org/zap"

	"github.com/labstack/echo"
)

/* user login module*/

var defaultRole = "normal"

type Session struct {
	LoginInfo  *LoginInfo
	CreateTime time.Time
}

type LoginInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
	Priv   string `json:"priv"`
}

// admin自带login接口
func (a *Admin) login(c echo.Context) error {
	username := c.FormValue("user")
	pw := c.FormValue("pw")
	if username == "" || pw == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	query := fmt.Sprintf("select id,priv from users where username = '%s' and pw ='%s'", username, pw)
	rows, err := g.DB.Query(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	if !rows.Next() {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusNotFound,
			ErrCode: g.UserNotExistC,
			Message: g.UserNotExistE,
		})
	}

	li := &LoginInfo{
		Name:   username,
		Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Token:  strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	rows.Scan(&li.ID, &li.Priv)
	//sub token验证成功，保存session
	a.sessions.Store(li.Token, &Session{
		LoginInfo:  li,
		CreateTime: time.Now(),
	})

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   li,
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

func (a *Admin) loginInfo(c echo.Context) error {
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
		Data:   sess.(*Session).LoginInfo,
	})
}

func (a *Admin) getLoginInfo(c echo.Context) *LoginInfo {
	token := getToken(c)
	sess, ok := a.sessions.Load(token)
	if !ok {
		// 用户未登陆或者session失效
		return nil
	}

	return sess.(*Session).LoginInfo
}

func getToken(c echo.Context) string {
	return c.Request().Header.Get("X-Token")
}
