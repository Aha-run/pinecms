import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/Home'

/*import Index from '@/views/Index'
import Login from '@/views/Login'
import Register from '@/views/Register'
import Log from '@/views/Log'
import MessageBoard from '@/views/MessageBoard'
import BlogWrite from '@/views/blog/BlogWrite'
import BlogView from '@/views/blog/BlogView'
import BlogAllCategoryTag from '@/views/blog/BlogAllCategoryTag'
import BlogCategoryTag from '@/views/blog/BlogCategoryTag'*/

import {Message} from 'element-ui';
import {getToken} from '@/request/token'
import store from '@/store'
import {removeToken} from "../request/token";

Vue.use(Router)

const router = new Router({
  routes: [
    {
      path: '/write/:id?',
      component: r => require.ensure([], () => r(require('@/views/blog/BlogWrite')), 'blogwrite'),
      meta: {
        requireLogin: true
      },
    },
    {
      path: '/verify/:token',
      component: r => require.ensure([], () => r(require('@/views/Verify')), 'verify')
    },
    {
      path: '',
      name: 'Home',
      component: Home,
      children: [
        {
          path: '/user/center',
          component: r => require.ensure([], () => r(require('@/views/usercenter/Index')), 'usercenter'),
          meta: {
            requireLogin: true
          },
        },
        {
          path: '/user/modipwd',
          component: r => require.ensure([], () => r(require('@/views/usercenter/Password')), 'usercenter'),
          meta: {
            requireLogin: true
          },
        },
        {
          path: '/user/money',
          component: r => require.ensure([], () => r(require('@/views/usercenter/Money')), 'usercenter'),
          meta: {
            requireLogin: true
          },
        },
        {
          path: '/user/task',
          component: r => require.ensure([], () => r(require('@/views/usercenter/Task')), 'usercenter'),
          meta: {
            requireLogin: true
          },
        },
        {
          path: '/user/cart',
          component: r => require.ensure([], () => r(require('@/views/usercenter/Cart')), 'usercenter'),
          meta: {
            requireLogin: true
          },
        },
        {
          path: '/user/ref',
          component: r => require.ensure([], () => r(require('@/views/usercenter/Ref')), 'usercenter'),
          meta: {
            requireLogin: true
          },
        },

        {
          path: '/buy/:id/:paytype',
          component: r => require.ensure([], () => r(require('@/views/Buy')), 'buy'),
          meta: {
            requireLogin: true
          },
        },

        {
          path: '/',
          component: r => require.ensure([], () => r(require('@/views/Index')), 'index')
        },
        {
          path: '/log',
          component: r => require.ensure([], () => r(require('@/views/Log')), 'log')
        },
        {
          path: '/free/video/:id?',
          component: r => require.ensure([], () => r(require('@/views/blog/FreeVideoCategory')), 'freeVideoCategory')
        },
        {
          path: '/free/book/:id?',
          component: r => require.ensure([], () => r(require('@/views/blog/FreeBookCategory')), 'freeBookCategory')
        },
        {
          path: '/paid/video/:id?',
          component: r => require.ensure([], () => r(require('@/views/blog/PaidVideoCategory')), 'paidVideoCategory')
        },
        {
          path: '/paid/book/:id?',
          component: r => require.ensure([], () => r(require('@/views/blog/PaidBookCategory')), 'paidBookCategory')
        },
        {
          path: '/archives/:year?/:month?',
          component: r => require.ensure([], () => r(require('@/views/blog/BlogArchive')), 'archives')
        },
        {
          path: '/archives/:year?/:month?',
          component: r => require.ensure([], () => r(require('@/views/blog/BlogArchive')), 'archives')
        },
        {
          path: '/feedback',
          component: r => require.ensure([], () => r(require('@/views/MessageBoard')), 'messageboard')
        },
        {
          path: '/view/:id',
          component: r => require.ensure([], () => r(require('@/views/blog/BlogView')), 'blogview')
        },
        {
          path: '/:type/all',
          component: r => require.ensure([], () => r(require('@/views/blog/BlogAllCategoryTag')), 'blogallcategorytag')
        },
        {
          path: '/:type/:id',
          component: r => require.ensure([], () => r(require('@/views/blog/BlogCategoryTag')), 'blogcategorytag')
        }
      ]
    },
    {
      path: '/login',
      component: r => require.ensure([], () => r(require('@/views/Login')), 'login')
    },
    {
      path: '/register',
      component: r => require.ensure([], () => r(require('@/views/Register')), 'register')
    }
  ],
  scrollBehavior(to, from, savedPosition) {
    return {x: 0, y: 0}
  }
})
router.beforeEach((to, from, next) => {
  if (getToken()) { // 有token
    if (to.path === '/login') {
      next({path: '/'})
    } else {
      if (store.state.account.length === 0) {
        store.dispatch('getUserInfo').then(data => { //获取用户信息
          next()
        }).catch(() => {
          removeToken() //todo 移除token防止死循环一直跳
          next({path: '/'})
        })
      } else {
        next()
      }
    }
  } else { // 无token
    if (to.matched.some(r => r.meta.requireLogin)) {
      Message({
        type: 'warning',
        showClose: true,
        message: '请先登录账户哟!',
        onClose: () => {
          next({path: '/login'})
        }
      })
    } else {
      next();
    }
  }
})


export default router
