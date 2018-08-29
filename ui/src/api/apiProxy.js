/* eslint-disable */
import request from '@/utils/request'
export function proxy(method,params) {
    return request({
      url: '/admin',
      method: method,
      params: params,
      withCredentials: true
    })
}
 