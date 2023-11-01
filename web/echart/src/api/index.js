import request from '@/utils/system/request'

// 获取统计数据
export function totalApi(params) {
  return request({
    url: '/user/stats/total',
    method: 'get',
    params
  })
}



// 获取图表数据
export function dataEchartApi(params) {
  return request({
    url: '/user/stats/chart',
    method: 'get',
    params
  })
}

// 获取来源域名
export function dataDomainApi(params) {
  return request({
    url: '/user/stats/rank',
    method: 'get',
    params
  })
}

export function getTitle(params) {
  return request({
    url: '/user/rule/info',
    method: 'get',
    params
  })
}

