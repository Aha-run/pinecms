import Vuex from 'vuex'
import Vue from 'vue'
// import {getToken, setToken, removeToken} from '@/request/token'
import {login, getUserInfo, logout, register} from '@/api/login'
import Cookie from 'js-cookie'

Vue.use(Vuex);

const createStore = () => {
  return new Vuex.Store({
    state: {
      id: '',
      account: '',
      name: '',
      avatar: '',
      token: '',
    },
    mutations: {
      SET_TOKEN: (state, token) => {
        state.token = token;
        if (token) {
          Cookie.set('token', token);
        } else {
          Cookie.remove('token');
        }
      },
      SET_ACCOUNT: (state, account) => {
        state.account = account
      },
      SET_NAME: (state, name) => {
        state.name = name
      },
      SET_AVATAR: (state, avatar) => {
        state.avatar = avatar
      },
      SET_ID: (state, id) => {
        state.id = id
      }
    },
    actions: {
      nuxtServerInit ({ commit }, { req }) {
      },
      // 登录
      login({commit}, user) {
        return new Promise((resolve, reject) => {
          login(user.account, user.password, user.token).then(data => {
            commit('SET_TOKEN', data.data.data['sign_token'])
            resolve()
          }).catch(error => {
            reject(error)
          })
        })
      },
      // 获取用户信息
      getUserInfo({commit, state}) {
        let that = this
        return new Promise((resolve, reject) => {
          getUserInfo().then(data => {
            if (data.data) {
              commit('SET_ACCOUNT', data.data.account)
              commit('SET_NAME', data.data.nickname)
              commit('SET_AVATAR', data.data.avatar)
              commit('SET_ID', data.data.id)
            } else {
              commit('SET_ACCOUNT', '')
              commit('SET_NAME', '')
              commit('SET_AVATAR', '')
              commit('SET_ID', '')
              // removeToken()
            }
            resolve(data)
          }).catch(error => {
            reject(error)
          })
        })
      },
      // 退出
      logout({commit, state}) {
        return new Promise((resolve, reject) => {
          // logout().then(data => {
            commit('SET_TOKEN', '')
            commit('SET_ACCOUNT', '')
            commit('SET_NAME', '')
            commit('SET_AVATAR', '')
            commit('SET_ID', '')

            // removeToken()
            resolve()
          // }).catch(error => {
          //   reject(error)
          // })
        })
      },
      // 前端 登出 todo 后端验证token是否过期
      fedLogOut({commit}) {
        return new Promise(resolve => {
          commit('SET_TOKEN', '')
          commit('SET_ACCOUNT', '')
          commit('SET_NAME', '')
          commit('SET_AVATAR', '')
          commit('SET_ID', '')
          // removeToken()
          resolve()
        }).catch(error => {
          reject(error)
        })
      },
      // 注册
      register({commit}, user) {
        return new Promise((resolve, reject) => {
          register(user.account, user.email, user.password,user.token).then((data) => {
            commit('SET_TOKEN', data.data.data['sign_token'])
            // setToken(data.data['sign_token'])
            resolve(data)
          }).catch((error) => {
            reject(error)
          })
        })
      }
    }
  })
}

export default createStore
