import { loginApi } from '@/api/user'

const state = () => ({
  token: '', // 登录token
  ruleId: '',  // ID
})

// getters
const getters = {
  token(state) {
    return state.token
  },
  ruleId(state) {
    return state.ruleId
  }
}

// mutations
const mutations = {
  tokenChange(state, token) {
    state.token = token
  },
  idChange(state, rule_id) {
    state.ruleId = rule_id
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
        dispatch('getId', {token: res.result.token ,id:res.result.rule_id})
        .then(infoRes => {
          resolve(res.result.token)
        })
      })
    })
  },
  // get user info after user logined
  getId({ commit }, params) {
    return new Promise((resolve, reject) => {
      commit('idChange', params.id)
      commit('tokenChange', params.token)
      resolve(params)
    })
  },

 // login out the system after user click the loginOut button
  loginOut({ commit }) {
    localStorage.removeItem('tabs')
    localStorage.removeItem('vuex')
    commit('tokenChange', '')
    location.reload()
   

  }
}

export default {
  namespaced: true,
  state,
  actions,
  getters,
  mutations
}