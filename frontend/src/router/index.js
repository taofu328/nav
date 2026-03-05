import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { requiresAuth: false, isPublic: true }
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false, isPublic: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/Admin.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const adminToken = localStorage.getItem('admin_token')
  const userToken = localStorage.getItem('token')
  
  if (to.meta.requiresAdmin) {
    if (!adminToken) {
      next('/admin/login')
    } else {
      next()
    }
  } else if (to.meta.requiresAuth) {
    if (!userToken) {
      next('/login')
    } else {
      next()
    }
  } else if (to.path === '/admin/login' && adminToken) {
    next('/admin')
  } else if ((to.path === '/login' || to.path === '/register') && userToken) {
    next('/admin')
  } else {
    next()
  }
})

export default router
