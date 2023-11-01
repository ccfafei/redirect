import axios from 'axios'
import store from '@/store'
import { ElMessage } from 'element-plus'
import { routerKey } from 'vue-router'
const baseURL = import.meta.env.VITE_BASE_URL

const service = axios.create({
  baseURL: baseURL,
  timeout: 5000
})

// 请求前的统一处理
service.interceptors.request.use(
  (config) => {
    // JWT鉴权处理
    if (store.getters['user/token']) {
      config.headers['token'] = store.state.user.token
    }
    return config
  },
  (error) => {
    showError(error)
    // console.log(error) // for debug
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code === 200) {
      return res
    } else {
      showError(res)
      return Promise.reject(res)
    }
  },
  (error) => {
    const errRes=error.response.data  //错误信息
    showError(errRes)
    return Promise.reject(error)
  }
)

function showError (error) {
  if (error.code === 408||error.code === 401) { //token 408过期
    // to re-login
    store.dispatch('user/loginOut')

  } else {
    ElMessage({
      message: error.msg || error.message || '服务异常',
      type: 'error',
      duration: 3 * 1000
    })
  }
  
}

export default service