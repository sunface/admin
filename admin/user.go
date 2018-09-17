package admin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"

	"github.com/mafanr/admin/misc"
	"github.com/mafanr/g"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

/* user manage module */

func (a *Admin) InitSuperAdmin() {
	query := fmt.Sprintf("insert into users (username,pw,priv,create_date) values ('admin','admin','%s','%s')",
		g.PRIV_SUPER_ADMIN, talent.Time2StringSecond(time.Now()))
	_, err := g.DB.Exec(query)
	if err != nil {
		if !strings.Contains(err.Error(), g.DUP_KEY_ERR) {
			g.L.Fatal("add super admin error", zap.Error(err))
		}
	}
}

func (a *Admin) AddUser(c echo.Context) error {
	userName := c.FormValue("username")
	priv := c.FormValue("priv")
	if userName == "" || priv == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	// user name must consist of alphabet and numeric
	if !talent.OnlyAlphaAndNum(userName) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// validate the priv
	if priv != g.PRIV_ADMIN && priv != g.PRIV_NORMAL && priv != g.PRIV_OUTSOURCE {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	// check operator has the permission
	if !a.isRole(li.Priv, []string{g.PRIV_SUPER_ADMIN}) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	query := fmt.Sprintf("insert into users (username,pw,priv,create_date) values ('%s','%s','%s','%s')",
		userName, misc.DefaultPassword, priv, talent.Time2StringSecond(time.Now()))
	_, err := g.DB.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), g.DUP_KEY_ERR) {
			return c.JSON(http.StatusOK, g.Result{
				Status:  http.StatusConflict,
				ErrCode: g.AlreadyExistC,
				Message: g.AlreadyExistE,
			})
		}
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (a *Admin) EditUser(c echo.Context) error {
	userName := c.FormValue("username")
	priv := c.FormValue("priv")
	if userName == "" || priv == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	// validate the priv
	if priv != g.PRIV_ADMIN && priv != g.PRIV_NORMAL && priv != g.PRIV_OUTSOURCE {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	// check operator has the permission
	if !a.isRole(li.Priv, []string{g.PRIV_SUPER_ADMIN}) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	query := fmt.Sprintf("update  users  set priv='%s' where username='%s'", priv, userName)
	_, err := g.DB.Exec(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (a *Admin) LoadUsers(c echo.Context) error {
	users := make([]*misc.User, 0)
	query := fmt.Sprintf("select * from users")

	err := g.DB.Select(&users, query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   users,
	})
}

func (a *Admin) DeleteUser(c echo.Context) error {
	userName := c.FormValue("username")
	if userName == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}
	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	// check operator has the permission
	if !a.isRole(li.Priv, []string{g.PRIV_SUPER_ADMIN}) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// cant operator on self
	if li.Username == userName {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.OperateOnSelfC,
			Message: g.OperateOnSelfE,
		})
	}
	query := fmt.Sprintf("delete from users where username='%s'", userName)
	_, err := g.DB.Exec(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (a *Admin) SetPW(c echo.Context) error {
	oldPW := talent.FormValue(c, "old_pw")
	newPW := talent.FormValue(c, "new_pw")
	if oldPW == "" || newPW == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	if !talent.OnlyAlphaAndNum(newPW) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.PWAlphabetAndNumC,
			Message: g.PWAlphabetAndNumE,
		})
	}

	// get login info
	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	query := fmt.Sprintf("update users set pw='%s' where username='%s' and pw='%s'", newPW, li.Username, oldPW)
	res, err := g.DB.Exec(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.OldPWErrorC,
			Message: g.OldPWErrorE,
		})
	}
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}
func (a *Admin) SetUserPW(c echo.Context) error {
	username := talent.FormValue(c, "username")
	newPW := talent.FormValue(c, "new_pw")
	if username == "" || newPW == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	if !talent.OnlyAlphaAndNum(newPW) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.PWAlphabetAndNumC,
			Message: g.PWAlphabetAndNumE,
		})
	}

	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	if !a.isRole(li.Priv, []string{g.PRIV_SUPER_ADMIN, g.PRIV_ADMIN}) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	query := fmt.Sprintf("update users set pw='%s' where username='%s' ", newPW, username)
	res, err := g.DB.Exec(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.OldPWErrorC,
			Message: g.OldPWErrorE,
		})
	}
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}
func (a *Admin) isRole(priv string, privs []string) bool {
	can := false
	for _, p := range privs {
		if priv == p {
			can = true
			break
		}
	}
	return can
}

func (a *Admin) userExist(username string) bool {
	// 验证service是否存在
	var temp interface{}
	query := fmt.Sprintf("select id from users where username ='%s'", username)
	err := g.DB.Get(&temp, query)
	if err != nil {
		return false
	}

	return true
}
