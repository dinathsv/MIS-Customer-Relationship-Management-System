import { createRouter, createWebHistory } from 'vue-router'
import Register from '../views/Register.vue'
import AdminDashboard from '../views/AdminDashboard.vue'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/register', component: Register },
  { 
    path: '/dashboard', 
    component: AdminDashboard,
    meta: { requiresAuth: true }
  },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation Guard (Global Auth Guard)
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
  
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!token) {
      alert("Please login to the portal first.")
      window.location.href = '/' // Redirect to portal root
    } else if (!isAdmin) {
      alert("Access Denied: User Management module is restricted to administrators only.")
      window.location.href = '/' // Redirect to portal root
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router