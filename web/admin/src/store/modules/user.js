import { loginApi, getInfoApi, loginOutApi } from '@/api/user'

const state = () => ({
  token: '', // 登录token
  info: {},  // 用户信息
})

// getters
const getters = {
  token(state) {
    return state.token
  },
  info(state) {
    return state.info
  }
}

// mutations
const mutations = {
  tokenChange(state, token) {
    state.token = token
  },
  infoChange(state, info) {
    state.info = info
  }
}

// actions
const actions = {
  // login by login.vue
  login ({ commit, dispatch }, params) {
    return new Promise((resolve, reject) => {
      loginApi(params)
        .then(res => {
        commit('tokenChange', res.result.token)
        dispatch('getInfo', {token: res.result.token ,id:res.result.id})
        .then(infoRes => {
          resolve(res.result.token)
        })
      })
    })
  },
  // get user info after user logined
  getInfo({ commit }, params) {
    return new Promise((resolve, reject) => {
      getInfoApi(params)
        .then(res => {
        commit('infoChange', res.result.account)  // 用户账号
        resolve(res.result.account)
      })
    })
  },

  // login out the system after user click the loginOut button
  loginOut({ commit }) {
    loginOutApi()
    .then(res => {
      commit('tokenChange', '')
    })
    .catch(error => {

    })
    .finally(() => {
      localStorage.removeItem('tabs')
      localStorage.removeItem('vuex')
      location.reload()
    })
  }
}

export default {
  namespaced: true,
  state,
  actions,
  getters,
  mutations
}