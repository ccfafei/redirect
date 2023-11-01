import request from '@/utils/system/request'
// 获取验证码id,图片
export function getCaptchaId(params) {
  return request({
    url: '/captcha',
    method: 'get',
    params
  })
}
// 登录api
export function loginApi(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 获取用户信息Api
export function getInfoApi(params) {
  return request({
    url: '/api/account/info',
    method: 'get',
    params
  })
}

// 获取所有用户
export function getUserList(params) {
  return request({
    url: '/api/account/all',
    method: 'get',
    params
  })
}


// 添加用户
export function addUser(data) {
  return request({
    url: '/api/account/add',
    method: 'post',
    data
  })
}

// 修改用户
export function updateUser(data) {
  return request({
    url: '/api/account/update',
    method: 'post',
    data
  })
}

// 删除用户
export function delUser(data) {
  return request({
    url: '/api/account/del',
    method: 'post',
    data
  })
}
// 退出登录Api
export function loginOutApi() {
  return request({
    url: '/api/logout',
    method: 'post',
  })
}


// 获取日志
export function getLogsList(params) {
  return request({
    url: '/api/logs/all',
    method: 'get',
    params
  })
}
// 删除日志
export function delLogs(data) {
  return request({
    url: '/api/logs/delete',
    method: 'post',
    data
  })
}
