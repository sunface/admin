package admin

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mafanr/g"

	"github.com/labstack/echo"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

/*----------------------------Service相关----------------------------*/
type Service struct {
	Name       string  `json:"name"`
	Owner      string  `json:"owner"`
	Desc       *string `json:"desc"`
	ModifyDate string  `json:"modify_date"`
}

type Services []*Service

func (a *Admin) QueryServices(c echo.Context) error {
	services := make(Services, 0)

	// get username and global priv
	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	query := fmt.Sprintf("select name,owner,description,modify_date from service")
	rows, err := g.DB.Query(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	for rows.Next() {
		s := &Service{}
		rows.Scan(&s.Name, &s.Owner, &s.Desc, &s.ModifyDate)
		services = append(services, s)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   services,
	})
}

func (a *Admin) CreateService(c echo.Context) error {
	service := c.FormValue("service")
	owner := c.FormValue("owner")
	desc := c.FormValue("desc")
	if service == "" || !talent.OnlyAlphaAndNum(service) || owner == "" {
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
	if !a.isRole(li.Priv, []string{g.PRIV_SUPER_ADMIN, g.PRIV_ADMIN}) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// check the owner username is exist
	if !a.userExist(owner) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusNotFound,
			ErrCode: g.UserNotExistC,
			Message: g.UserNotExistE,
		})
	}

	now := talent.Time2StringSecond(time.Now())
	query := fmt.Sprintf("insert into service (name,owner,description,create_date) values ('%s','%s','%s','%s')",
		service, owner, desc, now)
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

func (a *Admin) UpdateService(c echo.Context) error {
	service := c.FormValue("service")
	owner := c.FormValue("owner")
	desc := c.FormValue("desc")
	if service == "" || desc == "" || owner == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
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

	// check the owner username is exist
	if !a.userExist(owner) {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusNotFound,
			ErrCode: g.UserNotExistC,
			Message: g.UserNotExistE,
		})
	}

	if li.Priv == g.PRIV_ADMIN || li.Priv == g.PRIV_SUPER_ADMIN || isServiceOwner(service, li.Username) {
		query := fmt.Sprintf("update service set owner='%s',description='%s' where name='%s'", owner, desc, service)
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
				Status:  http.StatusNotFound,
				ErrCode: g.NotExistC,
				Message: g.NotExistE,
			})
		}

		return c.JSON(http.StatusOK, g.Result{
			Status: http.StatusOK,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status:  http.StatusForbidden,
		ErrCode: g.ForbiddenC,
		Message: g.ForbiddenE,
	})
}

type ServiceUser struct {
	Username string `json:"username"`
	Priv     string `json:"priv"`
}

func (a *Admin) QueryServiceUser(c echo.Context) error {
	service := c.FormValue("service")
	if service == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	sus := make([]*ServiceUser, 0)
	query := fmt.Sprintf("select username,priv from service_priv where service='%s'", service)
	rows, err := g.DB.Query(query)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	for rows.Next() {
		su := &ServiceUser{}
		rows.Scan(&su.Username, &su.Priv)
		sus = append(sus, su)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   sus,
	})
}
func (a *Admin) AddServiceUser(c echo.Context) error {
	service := c.FormValue("service")
	username := c.FormValue("username")
	priv := c.FormValue("priv")
	if service == "" || username == "" || priv == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	// operator must be Global admin or service owner
	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	if li.Priv == g.PRIV_ADMIN || li.Priv == g.PRIV_SUPER_ADMIN || isServiceOwner(service, li.Username) {
		query := fmt.Sprintf("insert into service_priv (username,service,priv) values ('%s','%s','%s')",
			username, service, priv)
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

	return c.JSON(http.StatusOK, g.Result{
		Status:  http.StatusForbidden,
		ErrCode: g.ForbiddenC,
		Message: g.ForbiddenE,
	})
}

func (a *Admin) DelServiceUser(c echo.Context) error {
	service := c.FormValue("service")
	username := c.FormValue("username")
	if service == "" || username == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamEmptyC,
			Message: g.ParamEmptyE,
		})
	}

	// operator must be Global admin or service owner
	li := a.getLoginInfo(c)
	if li == nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	if li.Priv == g.PRIV_ADMIN || li.Priv == g.PRIV_SUPER_ADMIN || isServiceOwner(service, li.Username) {
		query := fmt.Sprintf("delete from service_priv where username='%s' and service='%s'",
			username, service)
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

	return c.JSON(http.StatusOK, g.Result{
		Status:  http.StatusForbidden,
		ErrCode: g.ForbiddenC,
		Message: g.ForbiddenE,
	})
}

func isServiceOwner(service string, username string) bool {
	var temp interface{}
	query := fmt.Sprintf("select id from service where name='%s' and owner ='%s'", service, username)
	err := g.DB.Get(&temp, query)
	if err != nil {
		return false
	}

	return true
}

func isServiceAdmin(service, username string) bool {
	var temp interface{}
	query := fmt.Sprintf("select id from service_priv where service='%s' and username='%s' and priv='%s'", service, username, g.PRIV_ADMIN)
	err := g.DB.Get(&temp, query)
	if err == nil {
		return true
	}
	return false
}

func isServiceUser(service string, username string) bool {
	var temp interface{}
	query := fmt.Sprintf("select id from service where name='%s' and owner ='%s'", service, username)
	err := g.DB.Get(&temp, query)
	if err == nil {
		return true
	}

	query = fmt.Sprintf("select id from service_priv where service='%s' and username='%s'", service, username)
	err = g.DB.Get(&temp, query)
	if err == nil {
		return true
	}

	return false
}
