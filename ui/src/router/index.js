/* eslint-disable */
import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/views/layout/Layout'

/** note: submenu only apppear when children.length>=1
*   detail see  https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
**/

/**
* hidden: true                   if `hidden:true` will not show in the sidebar(default is false)
* alwaysShow: true               if set true, will always show the root menu, whatever its child routes length
*                                if not set alwaysShow, only more than one route under the children
*                                it will becomes nested mode, otherwise not show the root menu
* redirect: noredirect           if `redirect:noredirect` will no redirct in the breadcrumb
* name:'router-name'             the name is used by <keep-alive> (must set!!!)
* meta : {
    roles: ['admin','editor']     will control the page roles (you can set multiple roles)
    title: 'title'               the name show in submenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar,
    noCache: true                if true ,the page will no be cached(default is false)
  }
**/
export const constantRouterMap = [
  { path: '/login', component: () => import('@/views/login/index'), hidden: true },
  { path: '/404', component: () => import('@/views/errorPage/404'), hidden: true },
  { path: '/401', component: () => import('@/views/errorPage/401'), hidden: true },
  {
    path: '',
    component: Layout,
    redirect: 'dashboard',
    children: [{
      path: 'dashboard',
      component: () => import('@/views/dashboard/index'),
      name: 'dashboard',
      meta: { title: 'dashboard', icon: 'dashboard', noCache: true }
    }]
  }
]

export default new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRouterMap,
  mode: 'history'
})

export const asyncRouterMap = [
  {
    path: '/services',
    name: 'services',
    component: Layout,
    meta: { title: 'Services', icon: 'table', noCache: true },
    children: [
      {
        path: 'tfe',
        name: 'tfe',
        meta: { title: 'TFE', noCache: true },
        component: () => import('@/views/services/tfe/index'),
        redirect: '/services/tfe/dashboard',
        children: [
          {
            path: 'dashboard',
            component: () => import('@/views/services/tfe/dashboard'),
            name: 'tfe-dashboard',
            hidden:true,
            meta: { title: 'Dashboard',  noCache: true }
          },
          {
            path: 'service',
            component: () => import('@/views/services/tfe/service'),
            name: 'service',
            meta: { title: 'Service管理',  noCache: true}
          },
          {
            path: 'api',
            component: () => import('@/views/services/tfe/api'),
            name: 'api',
            meta: { title: 'API管理', noCache: true}
          },
          {
            path: 'strategy',
            component: () => import('@/views/services/tfe/strategy'),
            name: 'strategy',
            meta: { title: '策略管理', noCache: true}
          },
          {
            path: 'label',
            component: () => import('@/views/services/tfe/label'),
            name: 'label',
            meta: { title: '标签分组', noCache: true}
          },
          { 
            path: 'audit',
            component: () => import('@/views/services/tfe/auditLogs'),
            name: 'auditLogs',
            hidden:false,
            meta: { title: '审计日志',  noCache: true}
          },
          { 
            path: 'import',
            component: () => import('@/views/services/tfe/import'),
            name: 'import',
            hidden:false,
            meta: { title: '配置导入',  noCache: true}
          },
          {
            path: 'log',
            component: () => import('@/views/services/tfe/log'),
            name: 'log',
            hidden:true,
            meta: { title: '日志查询',  noCache: true}
          }
        ]
      }
    ]
  },
  {
    path: '/users',
    component: Layout,
    children: [{
      path: '',
      component: () => import('@/views/users'),
      name: 'users',
      meta: { title: 'Users', icon: 'user', noCache: true }
    }]
  },
  {
    path: '/setting',
    component: Layout,
    children: [{
      path: '',
      component: () => import('@/views/setting'),
      name: 'setting',
      meta: { title: 'Setting', icon: 'eye', noCache: true }
    }]
  },
  // {
  //   path: '/tools',
  //   name: 'tools',
  //   component: Layout,
  //   meta: { title: '小工具', icon: 'table', noCache: true },
  //   children: [
  //     {
  //       path: 'code',
  //       component: () => import('@/views/tools/code'),
  //       name: 'code',
  //       meta: { title: '编程常用',  noCache: true }
  //     },
  //     {
  //       path: 'praise',
  //       component: () => import('@/views/tools/praise'),
  //       name: 'praise',
  //       meta: { title: '抽奖工具',  noCache: true }
  //     }
  //   ]
  // },
  { path: '*', redirect: '/404', hidden: true }
]
