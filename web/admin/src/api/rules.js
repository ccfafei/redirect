import request from '@/utils/system/request'



// 获取所有规则
export function getRulesList(params) {
  return request({
    url: '/api/rule/all',
    method: 'get',
    params
  })
}



// 添加用户
export function addRule(data) {
  return request({
    url: '/api/rule/add',
    method: 'post',
    data
  })
}

// 修改用户
export function updateRule(data) {
  return request({
    url: '/api/rule/update',
    method: 'post',
    data
  })
}

// 删除用户
export function delRule(data) {
  return request({
    url: '/api/rule/delete',
    method: 'post',
    data
  })
}
//
export function urlPassword(params) {
  return request({
    url: '/api/share/info',
    method: 'get',
    params
  })
}

export function updatePassword(data) {
  return request({
    url: '/api/share/update',
    method: 'post',
    data
  })
}
 