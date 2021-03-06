/* eslint-disable */
import Vue from 'vue'
import Router from 'vue-router'

// Vue.use(Router)

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
  {path: '/', redirect: '/adminui/dashboard'},
  { path: '/adminui/login', component: () => import('@/views/login/index'), hidden: true },
  { path: '/adminui/404', component: () => import('@/views/errorPage/404'), hidden: true },
  { path: '/adminui/401', component: () => import('@/views/errorPage/401'), hidden: true },
  {
    path: '/adminui',
    component: Layout,
    redirect: '/adminui/dashboard',
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
    path: '/adminui/ops',
    component: Layout,
    redirect: '/adminui/ops/service',
    meta: { title: 'ops',icon: 'component',noCache: true},
    children: [{
      path: 'service',
      component: () => import('@/views/ops/service'),
      name: 'service',
      meta: {
        title: 'service',
        noCache: true
      }
    }, 
    {
      path: 'application',
      component: () => import('@/views/ops/application'),
      name: 'application',
      meta: {
        title: 'application',
        noCache: true
      }
    },
    {
      path: 'server',
      component: () => import('@/views/ops/server'),
      name: 'server',
      hidden:true,
      meta: {
        title: 'cloudserver',
        noCache: true
      }
    }
  ]
  },
  {
    path: '/adminui/services',
    name: 'services',
    component: Layout,
    meta: { title: 'services', icon: 'table', noCache: true },
    children: [
      {
        path: 'juz',
        name: 'juz',
        meta: { title: 'Api网关', noCache: true },
        component: () => import('@/views/services/juz/index'),
        redirect: '/services/juz/dashboard',
        children: [
          {
            path: 'dashboard',
            component: () => import('@/views/services/juz/dashboard'),
            name: 'tfe-dashboard',
            hidden:true,
            meta: { title: 'Dashboard',  noCache: true }
          },
          {
            path: 'api',
            component: () => import('@/views/services/juz/api'),
            name: 'api',
            meta: { title: 'api', noCache: true}
          },
          {
            path: 'strategy',
            component: () => import('@/views/services/juz/strategy'),
            name: 'strategy',
            meta: { title: 'strategy', noCache: true}
          },
          { 
            path: 'audit',
            component: () => import('@/views/services/juz/auditLogs'),
            name: 'auditLogs',
            hidden:false,
            meta: { title: 'auditLogs',  noCache: true}
          },
          { 
            path: 'import',
            component: () => import('@/views/services/juz/import'),
            name: 'import',
            hidden:false,
            meta: { title: 'configImport',  noCache: true}
          },
          { 
            path: 'metrics',
            component: () => import('@/views/services/juz/metrics'),
            name: 'metrics',
            hidden:false,
            meta: { title: 'metrics',  noCache: true}
          }
        ]
      }
    ]
  },
  {
    path: '/adminui/users',
    component: Layout,
    children: [{
      path: '',
      component: () => import('@/views/users'),
      name: 'users',
      meta: { title: 'users', icon: 'user', noCache: true }
    }]
  },
  {
    path: '/adminui/setting',
    component: Layout,
    children: [{
      path: '',
      component: () => import('@/views/setting'),
      name: 'setting',
      meta: { title: 'setting', icon: 'eye', noCache: true }
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
  { path: '/adminui/*', redirect: '/404', hidden: true }
]
