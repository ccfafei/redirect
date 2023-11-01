import request from '@/utils/system/request'

// 登录api
export function loginApi(data) {
  return request({
    url: '/share_login',
    method: 'post',
    data
  })
}