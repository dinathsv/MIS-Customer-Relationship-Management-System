<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <p class="section-label mb-1">Feedback Management</p>
        <h2 class="text-2xl font-bold font-display text-white">Review Feedback</h2>
      </div>
      <div class="flex items-center gap-3">
        <button @click="loadFeedbacks" :disabled="loading" class="btn-secondary flex items-center gap-2 text-xs">
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="glass-card p-8 text-center toast-error">
      <p class="text-4xl mb-3">⚠️</p>
      <p class="font-semibold">Failed to load feedback</p>
      <p class="text-xs mt-1 opacity-70">{{ error }}</p>
      <button @click="loadFeedbacks" class="btn-secondary mt-4 text-sm">Try Again</button>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading && !feedbacks.length" class="space-y-4">
      <div v-for="i in 3" :key="i" class="glass-card p-6 h-32">
        <div class="shimmer h-4 w-1/3 mb-4 rounded"></div>
        <div class="shimmer h-3 w-1/4 mb-2 rounded"></div>
        <div class="shimmer h-3 w-full rounded"></div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!feedbacks.length" class="glass-card p-12 text-center">
      <p class="text-4xl mb-4">📭</p>
      <p class="text-lg font-medium text-white mb-2">No feedback found</p>
      <p style="color: var(--text-secondary)">There are no customer feedback submissions available yet.</p>
    </div>

    <!-- Feedback List -->
    <div v-else class="space-y-4">
      <div v-for="item in feedbacks" :key="item.feedback_id" class="glass-card-static p-6 animate-fade-up">
        <div class="flex flex-col md:flex-row md:items-start gap-4 justify-between mb-4">
          <div>
            <div class="flex items-center gap-3 mb-2">
              <span class="badge" style="background: rgba(124,58,237,0.15); color: #a78bfa; border: 1px solid rgba(124,58,237,0.3)">
                {{ item.category }}
              </span>
              <span class="text-xs" style="color: var(--text-secondary)">
                {{ formatDate(item.created_at) }}
              </span>
            </div>
            <p class="text-sm font-medium text-white">
              Customer ID: <span style="color: var(--accent-blue)">{{ item.customer_id }}</span>
            </p>
          </div>
          
          <div class="flex items-center gap-1 shrink-0">
            <span class="text-2xl font-bold font-display mr-1" :style="{ color: ratingColor(item.rating) }">
              {{ item.rating }}
            </span>
            <div class="flex text-lg">
              <span v-for="n in 5" :key="n" :style="{ color: n <= item.rating ? ratingColor(item.rating) : 'rgba(255,255,255,0.1)' }">
                ★
              </span>
            </div>
          </div>
        </div>
        
        <div class="p-4 rounded-lg" style="background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.05)">
          <p class="text-sm leading-relaxed" style="color: var(--text-secondary)">
            {{ item.comments || 'No written comment provided.' }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '../composables/useApi.js'

const { listFeedback } = useApi()

const feedbacks = ref([])
const loading = ref(false)
const error = ref('')

function ratingColor(r) {
  if (r >= 4.5) return '#10b981'
  if (r >= 3.5) return '#22c55e'
  if (r >= 2.5) return '#eab308'
  if (r >= 1.5) return '#f97316'
  return '#ef4444'
}

function formatDate(dateStr) {
  if (!dateStr) return 'Unknown date'
  const date = new Date(dateStr)
  return date.toLocaleString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric',
    hour: 'numeric', minute: '2-digit'
  })
}

async function loadFeedbacks() {
  loading.value = true
  error.value = ''
  try {
    const response = await listFeedback()
    const rawData = response?.data || []
    // Sort by created_at descending
    feedbacks.value = [...rawData].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  } catch (err) {
    error.value = err?.response?.data?.error || err.message || 'Unknown error'
  } finally {
    loading.value = false
  }
}

onMounted(loadFeedbacks)
</script>
