/* eslint-disable */

import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css'// progress bar style
import { getToken } from '@/utils/auth' // getToken from cookie
import {getLoginInfo} from '@/api/login'

NProgress.configure({ showSpinner: false })// NProgress Configuration

// permission judge function
function hasPermission(roles, permissionRoles) {
  if (roles.indexOf('admin') >= 0) return true // admin permission passed directly
  if (!permissionRoles) return true
  return roles.some(role => permissionRoles.indexOf(role) >= 0)
}

const whiteList = ['/adminui/login']// no redirect whitelist

router.beforeEach((to, _, next) => {
  NProgress.start() // start progress bar
  // 只有在ssoToken不存在时，才能访问到login和callback页面

  if (getToken()) { // determine if there has token
    /* has token*/
    if (to.path === '/adminui/login') {
      next({ path: '/adminui' })
      NProgress.done() // if current page is dashboard will not trigger	afterEach hook, so manually handle it
    } else {
      if (store.getters.name == '') {
        // 刷新后导致state信息不存在，重新初始化
        getLoginInfo().then((res) => {
          // 设置用户信息
          store.dispatch('SetUserInfo', res.data.data)
            // 判断用户信息是否存在，不存在则需要加载
          if (store.getters.addRouters.length==0) {
            var role = store.getters.roles
            // 路由还有没有加载
            store.dispatch('GenerateRoutes', {role}).then(() => { // 根据roles权限生成可访问的路由表
              router.addRoutes(store.getters.addRouters) // 动态添加可访问路由表
            })
          }
            // 没有动态改变权限的需求可直接next() 删除下方权限判断 ↓
            if (hasPermission(store.getters.roles, to.meta.roles)) {
              next()//
            } else {
              next({ path: '/adminui/401', replace: true, query: { noGoBack: true }})
            }
        }).catch(error => {
          // 发生错误，需要重新登陆
          store.dispatch('LogOut')
        })
      }   else {
        // 判断用户信息是否存在，不存在则需要加载
        if (store.getters.addRouters.length==0) {
        var role = store.getters.roles
        // 路由还有没有加载
        store.dispatch('GenerateRoutes', {role}).then(() => { // 根据roles权限生成可访问的路由表
          router.addRoutes(store.getters.addRouters) // 动态添加可访问路由表
        })
      }
        // 没有动态改变权限的需求可直接next() 删除下方权限判断 ↓
        if (hasPermission(store.getters.roles, to.meta.roles)) {
          next()//
        } else {
          next({ path: '/adminui/401', replace: true, query: { noGoBack: true }})
        }
      }
     
    
    }
  } else {
    /* has no token*/
    if (whiteList.indexOf(to.path) !== -1) { // 在免登录白名单，直接进入
      next()
    } else {
      next('/adminui/login') // 否则全部重定向到登录页
      NProgress.done() // if current page is login will not trigger afterEach hook, so manually handle it
    }
  }
})

router.afterEach(() => {
  NProgress.done() // finish progress bar
})
