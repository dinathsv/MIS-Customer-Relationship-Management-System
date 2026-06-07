<template>
  <div class="relative min-h-screen">
    <!-- Animated background -->
    <div class="app-bg"></div>

    <!-- Floating orbs for depth -->
    <div class="fixed pointer-events-none" style="
      width: 600px; height: 600px; border-radius: 50%;
      background: radial-gradient(circle, rgba(124,58,237,0.12) 0%, transparent 70%);
      top: -200px; left: -200px; animation: float 8s ease-in-out infinite;
      z-index: 0;
    "></div>
    <div class="fixed pointer-events-none" style="
      width: 500px; height: 500px; border-radius: 50%;
      background: radial-gradient(circle, rgba(0,229,255,0.08) 0%, transparent 70%);
      bottom: -150px; right: -150px; animation: float 10s ease-in-out infinite reverse;
      z-index: 0;
    "></div>

    <!-- Main layout -->
    <div class="relative z-10 min-h-screen flex flex-col">

      <!-- ── Top Navigation ──────────────────────────────── -->
      <header class="sticky top-0 z-50" style="
        background: rgba(10,10,26,0.7);
        backdrop-filter: blur(20px);
        -webkit-backdrop-filter: blur(20px);
        border-bottom: 1px solid rgba(255,255,255,0.08);
      ">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 py-4 flex items-center justify-between">

          <!-- Logo -->
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg flex items-center justify-center text-lg"
                 style="background: linear-gradient(135deg, #7c3aed, #4f46e5); box-shadow: 0 0 16px rgba(124,58,237,0.5)">
              ◈
            </div>
            <div>
              <p class="text-xs font-semibold" style="color: var(--text-secondary); letter-spacing: 0.1em;">CRM ENTERPRISE</p>
              <p class="text-sm font-bold font-display text-white leading-tight">Feedback Module</p>
            </div>
          </div>

          <!-- Nav Tabs -->
          <nav class="flex items-center gap-1 p-1 rounded-xl"
               style="background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.08)">
            <button
              v-for="tab in tabs" :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
                activeTab === tab.id
                  ? 'text-white'
                  : 'hover:text-white/80'
              ]"
              :style="activeTab === tab.id
                ? 'background: linear-gradient(135deg, rgba(124,58,237,0.7), rgba(79,70,229,0.7)); box-shadow: 0 2px 12px rgba(124,58,237,0.35);'
                : 'color: rgba(255,255,255,0.45);'"
            >
              <span class="text-base">{{ tab.icon }}</span>
              <span class="hidden sm:inline">{{ tab.label }}</span>
            </button>
          </nav>

          <!-- Right Actions -->
          <div class="flex items-center gap-4">
            <button @click="goToPortal" class="flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium hover:bg-white/10 transition-colors" style="color: rgba(255,255,255,0.8); border: 1px solid rgba(255,255,255,0.2);">
              <span class="text-base">🏠</span> Portal
            </button>
            <button @click="logout" class="flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors hover:bg-red-500/20" style="color: #ef4444; border: 1px solid rgba(239,68,68,0.3);">
              Logout
            </button>
            <!-- API Status -->
            <div class="flex items-center gap-2 hidden md:flex">
              <span class="w-2 h-2 rounded-full animate-pulse"
                    :style="{ background: apiStatus === 'healthy' ? '#10b981' : apiStatus === 'checking' ? '#f59e0b' : '#ef4444' }">
              </span>
              <span class="text-xs" style="color: var(--text-secondary)">
                {{ apiStatus === 'healthy' ? 'API Online' : apiStatus === 'checking' ? 'Connecting...' : 'API Offline' }}
              </span>
            </div>
          </div>
        </div>
      </header>

      <!-- ── Hero Banner ─────────────────────────────────── -->
      <div class="max-w-7xl mx-auto w-full px-4 sm:px-6 pt-8 pb-2">
        <Transition name="fade" mode="out-in">
          <div v-if="activeTab === 'form'" key="form-hero">
            <div class="flex items-center gap-3 mb-1">
              <span class="text-3xl">📝</span>
              <h1 class="text-3xl font-bold font-display">
                <span class="gradient-text">Customer</span> Feedback
              </h1>
            </div>
            <p style="color: var(--text-secondary); font-size: 0.9rem;">
              Submit a new feedback entry — linked directly to the CRM customer profile via Customer ID.
            </p>
          </div>
          <div v-else-if="activeTab === 'list'" key="list-hero">
            <div class="flex items-center gap-3 mb-1">
              <span class="text-3xl">📋</span>
              <h1 class="text-3xl font-bold font-display">
                <span class="gradient-text">Feedback</span> Reviews
              </h1>
            </div>
            <p style="color: var(--text-secondary); font-size: 0.9rem;">
              Review all submitted feedback from customers in real-time.
            </p>
          </div>
          <div v-else-if="activeTab === 'analytics' && isAdmin" key="analytics-hero">
            <div class="flex items-center gap-3 mb-1">
              <span class="text-3xl">📊</span>
              <h1 class="text-3xl font-bold font-display">
                <span class="gradient-text">Analytics</span> Dashboard
              </h1>
            </div>
            <p style="color: var(--text-secondary); font-size: 0.9rem;">
              Real-time insights across all customer feedback — powered by live PostgreSQL aggregations.
            </p>
          </div>
        </Transition>
      </div>

      <!-- ── Main Content ────────────────────────────────── -->
      <main class="flex-1 max-w-7xl mx-auto w-full px-4 sm:px-6 py-6">
        <Transition name="fade" mode="out-in">

          <!-- Feedback Form Tab -->
          <div v-if="activeTab === 'form'" key="form-tab"
               class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
            <div class="lg:col-span-2">
              <FeedbackForm @submitted="onFeedbackSubmitted" />
            </div>
            <!-- Side Info Panel -->
            <div class="space-y-4">
              <div class="glass-card p-5">
                <p class="section-label mb-3">API Integration</p>
                <p class="text-xs mb-3" style="color: var(--text-secondary)">
                  Feedback is linked to the CRM via <span class="font-mono text-cyan-400">customer_id</span>.
                  Query all feedback for a customer:
                </p>
                <div class="rounded-lg p-3 text-xs font-mono overflow-x-auto"
                     style="background: rgba(0,0,0,0.4); border: 1px solid rgba(255,255,255,0.08); color: #67e8f9;">
                  GET /api/feedback/customer/{'{'}id{'}'}
                </div>
              </div>

              <div class="glass-card p-5">
                <p class="section-label mb-3">Categories</p>
                <div class="flex flex-wrap gap-2">
                  <span v-for="cat in categories" :key="cat"
                        class="badge text-xs"
                        style="background: rgba(124,58,237,0.15); color: rgba(167,139,250,0.9); border: 1px solid rgba(124,58,237,0.25)">
                    {{ cat }}
                  </span>
                </div>
              </div>

              <div class="glass-card p-5">
                <p class="section-label mb-3">Rating Scale</p>
                <div class="space-y-2">
                  <div v-for="(label, i) in ratingScale" :key="i"
                       class="flex items-center justify-between text-xs">
                    <span class="text-amber-400">{{ '★'.repeat(5-i) }}<span class="text-white/15">{{ '★'.repeat(i) }}</span></span>
                    <span style="color: var(--text-secondary)">{{ label }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Feedback List Tab -->
          <div v-else-if="activeTab === 'list'" key="list-tab">
            <FeedbackList />
          </div>

          <!-- Analytics Tab -->
          <div v-else-if="activeTab === 'analytics' && isAdmin" key="analytics-tab">
            <AnalyticsDashboard ref="dashboardRef" />
          </div>

        </Transition>
      </main>

      <!-- ── Footer ──────────────────────────────────────── -->
      <footer class="text-center py-6 px-4"
              style="border-top: 1px solid rgba(255,255,255,0.06); color: var(--text-secondary)">
        <p class="text-xs">
          CRM Feedback Module &nbsp;·&nbsp; Vue 3 + Go + PostgreSQL &nbsp;·&nbsp;
          <span class="gradient-text font-medium">Glassmorphism UI</span>
        </p>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import FeedbackForm from './components/FeedbackForm.vue'
import FeedbackList from './components/FeedbackList.vue'
import AnalyticsDashboard from './components/AnalyticsDashboard.vue'
import { useApi } from './composables/useApi.js'

const { checkHealth } = useApi()

const userStr_init = localStorage.getItem('crm_user')
const user_init = userStr_init ? JSON.parse(userStr_init) : null
const isAdmin_init = user_init && user_init.role_id === 1

const activeTab = ref(isAdmin_init ? 'list' : 'form')
const apiStatus = ref('checking')
const dashboardRef = ref(null)

const userStr = localStorage.getItem('crm_user')
const user = userStr ? JSON.parse(userStr) : null
const isAdmin = user && user.role_id === 1

const tabs = computed(() => {
  if (isAdmin) {
    return [
      { id: 'list', label: 'Review Feedback', icon: '📋' },
      { id: 'analytics', label: 'Analytics', icon: '📊' }
    ]
  } else {
    return [ { id: 'form', label: 'Submit Feedback', icon: '📝' } ]
  }
})

const categories = [
  'Product Quality', 'Customer Support', 'Delivery & Shipping',
  'Pricing & Value', 'Website & App', 'Returns & Refunds', 'Other',
]

const ratingScale = ['Excellent', 'Good', 'Average', 'Poor', 'Terrible']

function goToPortal() {
  window.location.href = "/";
}

function logout() {
  localStorage.removeItem("crm_token");
  localStorage.removeItem("crm_user");
  window.location.href = "/";
}

async function checkApiHealth() {
  try {
    await checkHealth()
    apiStatus.value = 'healthy'
  } catch {
    apiStatus.value = 'offline'
  }
}

function onFeedbackSubmitted() {
  // Do not switch tabs. The FeedbackForm component displays its own success message.
}

onMounted(() => {
  const token = localStorage.getItem('crm_token')
  if (!token) {
    window.location.href = '/' // Redirect to portal
    return
  }
  checkApiHealth()
})
</script>
