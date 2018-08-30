import request from '@/utils/request'


export function loadUsers(token) {
    return request({
        url: '/user/load',
        method: 'GET'
    })
}