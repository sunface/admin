package service

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sort"
	"strings"

	"github.com/mafanr/g"

	"time"

	"github.com/labstack/echo"
	"github.com/sunface/talent"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (a *Admin) testApi(c echo.Context) error {
	method := c.FormValue("method")

	paramsS := c.FormValue("params")
	var temp map[string]string
	json.Unmarshal([]byte(paramsS), &temp)

	// 解析参数
	args := &fasthttp.Args{}
	for k, v := range temp {
		args.Set(k, v)
	}

	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.Header.SetMethod(method)

	tfeAddr := c.FormValue("tfe_addr")

	switch method {
	case "GET":
		// 拼接url
		tfeAddr = tfeAddr + "?" + args.String()
	default:
		args.WriteTo(req.BodyWriter())
	}
	req.SetRequestURI(tfeAddr)

	err := g.Cli.DoTimeout(req, resp, 15*time.Second)
	if err != nil {
		g.L.Info("请求tfe出错", zap.Error(err), zap.String("url", tfeAddr))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusServiceUnavailable,
			Message: "请求TFE出错",
		})
	}

	if resp.StatusCode() != http.StatusOK {
		g.L.Info("请求tfe异常", zap.Int("code", resp.StatusCode()), zap.String("url", tfeAddr))
		return c.JSON(http.StatusOK, g.Result{
			Status:  resp.StatusCode(),
			Message: talent.Bytes2String(resp.Body()),
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   talent.Bytes2String(resp.Body()),
	})
}

var userList = []string{
	"沈工",

	"会员-徐永宁", "会员-廖冰", "会员-王林", "会员-江媛媛", "会员-汤敏捷", "会员-蔡缤纷",

	"测试-周嵩", "测试-邵磊", "测试-章彬彬", "测试-占志滔", "测试-张碧琳",
	"测试-孙泽波", "测试-章立星", "测试-戴靖", "测试-梁东东", "测试-张颍",

	"基础组件-孙飞",

	"前端-杨光远", "前端-陈欢斌", "前端-鲁延涛", "前端-王丽波", "前端-朱平雷", "前端-朱昊泽", "前端-陈国庆",
	"前端-方熹", "前端-刘泽威", "前端-南雷蒙", "前端-曾文平",

	"基础组件-邵聪聪",

	"中间件-王乃江", "中间件-王浩阳", "中间件-莫凌峰", "中间件-赵一鸣", "中间件-陈志玉", "中间件-黄江帆",

	"基础组件-任江林",

	"测试-贾壁宁", "测试-卢丹", "测试-严蕾兵", "测试-周磊", "测试-张园庆", "测试-王青青", "测试-章琴", "测试-张立杰",
	"测试-张春耀", "测试-王旭松", "测试-柴树军", "测试-刘利军", "测试-王昌林", "测试-李相花", "测试-杜敏涛", "测试-章津",
	"测试-冯旭妍", "测试-陈小丽", "测试-林爱萍", "测试-秦丹", "测试-徐丹",

	"中间件-王厚达", "中间件-周成刚", "中间件-夏光", "中间件-孙振维", "中间件-孙鹭易", "中间件-李国喜"}

func (a *Admin) genPraise(c echo.Context) error {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	punisheeRaw := c.FormValue("punishee")

	// 被惩罚者列表
	punishee := strings.Split(punisheeRaw, " ")

	auth := true
	for _, p := range punishee {
		valid := false
		for _, u := range userList {
			if strings.Contains(u, p) {
				valid = true
			}
		}

		if !valid {
			auth = false
		}
	}

	if !auth {
		return c.String(200, "输入不合法")
	}
	// 计算奖励者列表
	rewarder := getRewarder(punishee)

	return c.JSON(200, rewarder)
}

// 每个惩罚者对应N个奖励者
const PerPubnishee = 10

// 根据惩罚者列表获取奖励者列表
func getRewarder(punishee []string) []string {
	// 奖励者列表
	var rewarder []string

	// 获取惩罚者之外的剩余用户列表
	var restUsers []string
	for _, u := range userList {
		isPun := false
		for _, p := range punishee {
			if strings.Contains(u, p) {
				// 该用户是惩罚者，不加入奖励名单
				isPun = true
			}
		}
		if !isPun {
			restUsers = append(restUsers, u)
		}
	}

	rn := len(punishee) * PerPubnishee
	if rn > len(userList)-len(punishee) {
		// 若 (被惩罚者*10) >= (总列表 - 被惩罚者),则奖励被惩罚者之外的所有人
		rewarder = restUsers
		return rewarder
	}

	// 随机出rn长度的用户
	for i := 0; i < rn; i++ {
		// 随机一个用户[0:n)

		index := rand.Intn(len(restUsers))

		rewarder = append(rewarder, restUsers[index])

		// 从剩余列表，删除该用户
		restUsers = append(restUsers[:index], restUsers[index+1:]...)
	}

	sort.Strings(rewarder)
	return rewarder
}
