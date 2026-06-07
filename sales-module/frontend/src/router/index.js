import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Dashboard', component: () => import('../views/DashboardView.vue') },
  { path: '/sales', name: 'Sales', component: () => import('../views/SalesView.vue') },
  { path: '/reports', name: 'Reports', component: () => import('../views/ReportsView.vue') },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory('/sales/'),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('crm_token')
  let isAdmin = false;
  
  try {
    const userStr = localStorage.getItem('crm_user')
    if (userStr) {
      const user = JSON.parse(userStr)
      if (user.role_id === 1 || user.role === 'admin') {
        isAdmin = true;
      }
    }
  } catch (e) {
    console.error('Error parsing user data', e)
  }

  if (!token) {
    window.location.href = '/' // Redirect to portal root
  } else if (!isAdmin) {
    alert("Access Denied: Sales module is restricted to administrators only.")
    window.location.href = '/' // Redirect to portal root
  } else {
    next()
  }
})

export default router
