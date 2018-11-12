package admin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

type App struct {
	Name       string  `json:"name"`
	Service    string  `json:"service"`
	Desc       *string `json:"desc"`
	ModifyDate string  `json:"modify_date"`
}

type Apps []*App

func (a *Admin) QueryApps(c echo.Context) error {
	service := c.FormValue("service")
	if service == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	apps := make(Apps, 0)

	query := fmt.Sprintf("select name,service,description,modify_date from application where service='%s'", service)
	rows, err := g.DB.Query(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}
	defer rows.Close()

	for rows.Next() {
		s := &App{}
		rows.Scan(&s.Name, &s.Service, &s.Desc, &s.ModifyDate)
		apps = append(apps, s)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   apps,
	})
}

func (a *Admin) CreateAPP(c echo.Context) error {
	service := c.FormValue("service")
	app := c.FormValue("app")
	desc := c.FormValue("desc")
	if service == "" || !talent.OnlyAlphaAndNum(app) || app == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
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

	if !isServiceUser(service, li.Username) && li.Priv != g.PRIV_SUPER_ADMIN && li.Priv != g.PRIV_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	now := talent.Time2StringSecond(time.Now())
	query := fmt.Sprintf("insert into application (name,service,description,create_date) values ('%s','%s','%s','%s')",
		app, service, desc, now)
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

func (a *Admin) UpdateAPP(c echo.Context) error {
	service := c.FormValue("service")
	app := c.FormValue("app")
	desc := c.FormValue("desc")
	if service == "" || !talent.OnlyAlpha(app) || app == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
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

	// only global admin can crate service
	if !isServiceUser(service, li.Username) && li.Priv != g.PRIV_SUPER_ADMIN && li.Priv != g.PRIV_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	query := fmt.Sprintf("update application set description='%s' where name = '%s'",
		desc, app)
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
