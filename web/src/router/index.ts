import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import PasswordRecordsDetail from '../views/password/PasswordRecordsDetail.vue' // 假设详情页面组件路径

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    // component: () => import('../views/LoginView.vue'),
    component: () => import('../views/LoginCasdoor.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/callback',
    name: 'callback',
    component: () => import('../views/CallbackPage.vue'),
  },

  {
    path: '/',
    name: 'layout',
    component: () => import('../views/LayoutView.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'passwordList',
        component: () => import('../views/password/ListView.vue'),
      },
      {
        path: 'password/create',
        name: 'passwordCreate',
        component: () => import('../views/password/CreateView.vue'),
      },
      {
        path: 'password/:id',
        name: 'passwordDetail',
        component: () => import('../views/password/DetailView.vue'),
      },
      {
        path: 'password/edit/:id',
        name: 'passwordEdit',
        component: () => import('../views/password/EditView.vue'),
      },
      {
        path: 'password/records/detail',
        name: 'passwordRecordsDetail',
        component: () => import('../views/password/PasswordRecordsDetail.vue'),
      },
    ],
  },
  {
    path: '/userinfo',
    name: 'userinfo',
    component: () => import('../views/UserInfoView.vue'),
  },
  {
    path: '/session', // 你可以根据需求修改路径
    name: 'Session', // 路由名称
    component: () => import('../components/SomeComponent.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
