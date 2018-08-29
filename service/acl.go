package service

import "github.com/mafanr/g"
 
/* 管理后台的整体权限控制，集中在这里管理 */
// 权限分为两种：admin自身权限，应用权限
// admin服务自身权限分为超级管理员、管理员和普通用户
// 应用权限分为创建者、管理员、普通、体验者

// #admin服务
// 超级管理员能看到所有页面

// #应用权限
// 创建者、管理员、普通、体验者

var (
	GlobalAdmins = map[string]string{
		"13269": g.ROLE_SUPER_ADMIN,
		"14929": g.ROLE_ADMIN,
	}

	// 应用管理员
	AppAdmins = map[string]map[string]string{
		g.TFEManage: map[string]string{
			"13269": g.ROLE_SUPER_ADMIN,
			"14929": g.ROLE_VIEWER,
		},
	}
)
