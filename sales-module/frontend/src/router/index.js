import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/LoginView.vue'), meta: { public: true } },
  { path: '/', name: 'Dashboard', component: () => import('../views/DashboardView.vue') },
  { path: '/sales', name: 'Sales', component: () => import('../views/SalesView.vue') },
  { path: '/invoices', name: 'Invoices', component: () => import('../views/InvoicesView.vue') },
  { path: '/products', name: 'Products', component: () => import('../views/ProductsView.vue') },
  { path: '/reports', name: 'Reports', component: () => import('../views/ReportsView.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
