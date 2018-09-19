package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sunface/talent"

	"github.com/labstack/echo"
	"github.com/mafanr/admin/misc"
	"github.com/mafanr/g"
	"go.uber.org/zap"
)

type CloudServer struct {
	ID            int    `json:"id" db:"id"`
	PublicIP      string `json:"public_ip" db:"public_ip"`
	PrivateIP     string `json:"private_ip" db:"private_ip"`
	CloudProvider string `json:"cloud_provider" db:"cloud_provider"`
	SshUser       string `json:"ssh_user" db:"ssh_user"`
	SshPW         string `json:"ssh_pw" db:"ssh_pw"`
	Region        string `json:"region" db:"region"`
	Configure     string `json:"configure" db:"configure"`
	Expire        string `json:"expire" db:"expire"`
	ModifyDate    string `json:"modify_date" db:"modify_date"`
}

func (a *Admin) DefineServer(c echo.Context) error {
	serverR := c.FormValue("server")

	server := &CloudServer{}
	err := json.Unmarshal([]byte(serverR), &server)
	if err != nil {
		g.L.Info("parse server error", zap.Error(err), zap.String("api", serverR))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	if !talent.IsIP(server.PrivateIP) || server.CloudProvider == "" ||
		server.SshUser == "" || server.Configure == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	op := talent.FormValue(c, "op")
	// when create,you must provice ssh password
	// but when edit, if the password is empty, means that we dont need to change the password
	if op == "create" {
		if server.SshPW == "" {
			return c.JSON(http.StatusOK, g.Result{
				Status:  http.StatusBadRequest,
				ErrCode: g.ParamInvalidC,
				Message: g.ParamInvalidE,
			})
		}
	}

	if server.PublicIP != "" {
		if !talent.IsIP(server.PrivateIP) {
			return c.JSON(http.StatusOK, g.Result{
				Status:  http.StatusBadRequest,
				ErrCode: g.ParamInvalidC,
				Message: g.ParamInvalidE,
			})
		}
	}

	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	if li.Priv != g.PRIV_SUPER_ADMIN && li.Priv != g.PRIV_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	server.Configure = strings.ToUpper(server.Configure)
	var query string
	if op == "create" {
		query = fmt.Sprintf("insert into cloud_server (public_ip,private_ip,cloud_provider,ssh_user,ssh_pw,region,configure,expire) values ('%s','%s','%s','%s','%s','%s','%s','%s')",
			server.PublicIP, server.PrivateIP, server.CloudProvider, server.SshUser, server.SshPW, server.Region, server.Configure, server.Expire)
	} else {
		if server.SshPW == "" {
			query = fmt.Sprintf("update cloud_server set public_ip='%s',private_ip='%s',cloud_provider='%s',ssh_user='%s',region='%s',configure='%s',expire='%s' where id='%d'",
				server.PublicIP, server.PrivateIP, server.CloudProvider, server.SshUser, server.Region, server.Configure, server.Expire, server.ID)
		} else {
			query = fmt.Sprintf("update cloud_server set public_ip='%s',private_ip='%s',cloud_provider='%s',ssh_user='%s',ssh_pw='%s',region='%s',configure='%s',expire='%s' where id='%d'",
				server.PublicIP, server.PrivateIP, server.CloudProvider, server.SshUser, server.SshPW, server.Region, server.Configure, server.Expire, server.ID)
		}
	}
	_, err = g.DB.Exec(query)
	if err != nil {
		if strings.Contains(err.Error(), g.DUP_KEY_ERR) {
			return c.JSON(http.StatusOK, g.Result{
				Status:  http.StatusConflict,
				ErrCode: misc.IpAlreadyExistC,
				Message: misc.IpAlreadyExistE,
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

func (a *Admin) QueryServers(c echo.Context) error {
	servers := make([]*CloudServer, 0)
	query := fmt.Sprintf("select * from cloud_server")
	err := g.DB.Select(&servers, query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusInternalServerError, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	for _, s := range servers {
		// for safe,pw can only ge by method below
		s.SshPW = ""
	}
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   servers,
	})
}

func (a *Admin) QueryServerPW(c echo.Context) error {
	id := talent.FormValue(c, "id")
	if id == "" {
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

	if li.Priv != g.PRIV_SUPER_ADMIN && li.Priv != g.PRIV_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	query := fmt.Sprintf("select ssh_pw from cloud_server where id='%s'", id)
	rows, err := g.DB.Query(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusInternalServerError, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}
	defer rows.Close()

	if !rows.Next() {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusNotFound,
			ErrCode: g.ServerNotExistC,
			Message: g.ServerNotExistE,
		})
	}

	var pw string
	rows.Scan(&pw)

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   pw,
	})
}
