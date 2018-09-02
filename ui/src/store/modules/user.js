/* eslint-disable */
import {logout } from '@/api/login'
import { getToken, setToken, removeToken } from '@/utils/auth'

const user = {
  state: {
    user: '',
    status: '',
    code: '',
    token: getToken(),
    name: '',
    avatar: '',
    introduction: '',
    roles: [],
    setting: {
      articlePlatform: []
    }
  },

  mutations: {
    SET_USER: (state,userID) => {
      state.user = userID
    },
    SET_CODE: (state, code) => {
      state.code = code
    },
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_INTRODUCTION: (state, introduction) => {
      state.introduction = introduction
    },
    SET_SETTING: (state, setting) => {
      state.setting = setting
    },
    SET_STATUS: (state, status) => {
      state.status = status
    },
    SET_NAME: (state, name) => {
      state.name = name
    },
    SET_AVATAR: (state, avatar) => {
      state.avatar = avatar
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    }
  },

  actions: {   
    // SSO登陆成功，保存信息
    SetUserInfo({ commit,state},userInfo) {
      return new Promise(resolve => {
        setToken(userInfo.token)
        commit('SET_TOKEN', userInfo.token)
        commit('SET_ROLES', userInfo.priv)
        commit('SET_USER', userInfo.id)
        commit('SET_NAME', userInfo.username)
        commit('SET_AVATAR', userInfo.avatar)
        resolve()
      })
    },
    // 登出
    LogOut({ commit, state }) {
      return new Promise((resolve, reject) => {
        logout(state.token).then(() => {
          commit('SET_TOKEN', '')
          commit('SET_ROLES', [])
          removeToken()
          resolve()
        }).catch(error => {
          commit('SET_TOKEN', '')
          commit('SET_ROLES', [])
          removeToken()
          reject(error)
        })
      })
    }
  }
}

export default user
