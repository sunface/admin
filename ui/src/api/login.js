/* eslint-disable */
import request from '@/utils/request'


export function logout(token) {
  return request({
    url: '/logout',
    method: 'post',
    params:{
      token : token
    }
  })
}

export function getUserInfo(token) {
  return request({
    url: '/user/info',
    method: 'get',
    params: { token },
    withCredentials: true
  })
}

 