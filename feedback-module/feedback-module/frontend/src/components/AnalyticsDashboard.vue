<template>
  <div class="space-y-6">
    <!-- Dashboard Header -->
    <div class="flex items-center justify-between flex-wrap gap-4">
      <div>
        <p class="section-label mb-1">CRM Analytics</p>
        <h2 class="text-2xl font-bold font-display text-white">Feedback Dashboard</h2>
      </div>
      <div class="flex items-center gap-3">
        <span v-if="lastRefreshed" class="text-xs" style="color: var(--text-secondary)">
          Updated {{ lastRefreshed }}
        </span>
        <button @click="loadData" :disabled="loading" class="btn-secondary flex items-center gap-2 text-xs">
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Skeleton Loader -->
    <template v-if="loading && !summary">
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div v-for="i in 4" :key="i" class="glass-card p-6 h-28">
          <div class="shimmer h-4 w-20 mb-3 rounded"></div>
          <div class="shimmer h-8 w-16 rounded"></div>
        </div>
      </div>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div v-for="i in 2" :key="i" class="glass-card p-6 h-72">
          <div class="shimmer h-4 w-32 mb-6 rounded"></div>
          <div class="shimmer h-full w-full rounded"></div>
        </div>
      </div>
    </template>

    <!-- Error State -->
    <div v-else-if="error" class="glass-card p-8 text-center toast-error">
      <p class="text-4xl mb-3">⚠️</p>
      <p class="font-semibold">Failed to load analytics</p>
      <p class="text-xs mt-1 opacity-70">{{ error }}</p>
      <button @click="loadData" class="btn-secondary mt-4 text-sm">Try Again</button>
    </div>

    <!-- Main Content -->
    <template v-else-if="summary">

      <!-- ── KPI Stat Cards ─────────────────────────────────── -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="stat-card animate-fade-up" style="animation-delay:0s">
          <p class="section-label mb-3">Total Feedback</p>
          <p class="text-4xl font-bold font-display gradient-text">{{ summary.total_feedback }}</p>
          <p class="text-xs mt-2" style="color: var(--text-secondary)">All time submissions</p>
        </div>

        <div class="stat-card animate-fade-up" style="animation-delay:0.05s">
          <p class="section-label mb-3">Avg. Rating</p>
          <div class="flex items-baseline gap-1">
            <p class="text-4xl font-bold font-display" :style="{ color: ratingColor(summary.average_rating) }">
              {{ summary.average_rating.toFixed(1) }}
            </p>
            <span class="text-amber-400 text-lg">★</span>
          </div>
          <p class="text-xs mt-2" style="color: var(--text-secondary)">Out of 5.0</p>
        </div>

        <div class="stat-card animate-fade-up" style="animation-delay:0.1s">
          <p class="section-label mb-3">Top Category</p>
          <p class="text-lg font-semibold text-white leading-tight">{{ topCategory?.category || '—' }}</p>
          <p class="text-xs mt-2" style="color: var(--text-secondary)">{{ topCategory?.count }} responses</p>
        </div>

        <div class="stat-card animate-fade-up" style="animation-delay:0.15s">
          <p class="section-label mb-3">5-Star Rate</p>
          <p class="text-4xl font-bold font-display" style="color: var(--accent-green)">{{ fiveStarRate }}%</p>
          <p class="text-xs mt-2" style="color: var(--text-secondary)">Excellent ratings</p>
        </div>
      </div>

      <!-- ── Charts Row ─────────────────────────────────────── -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

        <!-- Category Volume Bar Chart -->
        <div class="glass-card-static p-6 animate-fade-up" style="animation-delay:0.2s">
          <div class="flex items-center justify-between mb-5">
            <div>
              <p class="section-label mb-1">Volume by Category</p>
              <p class="text-sm font-semibold text-white">Feedback Distribution</p>
            </div>
            <span class="badge" style="background: rgba(124,58,237,0.2); color: #a78bfa; border: 1px solid rgba(124,58,237,0.3)">
              Bar Chart
            </span>
          </div>
          <div class="h-60">
            <Bar v-if="categoryChartData" :data="categoryChartData" :options="barOptions" />
          </div>
        </div>

        <!-- Rating Distribution Doughnut -->
        <div class="glass-card-static p-6 animate-fade-up" style="animation-delay:0.25s">
          <div class="flex items-center justify-between mb-5">
            <div>
              <p class="section-label mb-1">Rating Breakdown</p>
              <p class="text-sm font-semibold text-white">Star Distribution</p>
            </div>
            <span class="badge" style="background: rgba(0,229,255,0.15); color: #67e8f9; border: 1px solid rgba(0,229,255,0.25)">
              Doughnut
            </span>
          </div>
          <div class="flex items-center gap-6 h-60">
            <div class="flex-1 h-full">
              <Doughnut v-if="ratingChartData" :data="ratingChartData" :options="doughnutOptions" />
            </div>
            <!-- Legend -->
            <div class="space-y-2 shrink-0">
              <div v-for="(item, i) in ratingLegend" :key="i" class="flex items-center gap-2 text-xs">
                <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ background: item.color }"></span>
                <span style="color: var(--text-secondary)">{{ item.label }}</span>
                <span class="font-semibold text-white ml-auto pl-3">{{ item.count }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Daily Trend Line ───────────────────────────────── -->
      <div v-if="summary.recent_trend?.length" class="glass-card-static p-6 animate-fade-up" style="animation-delay:0.3s">
        <div class="flex items-center justify-between mb-5">
          <div>
            <p class="section-label mb-1">Activity Trend</p>
            <p class="text-sm font-semibold text-white">14-Day Submission Volume</p>
          </div>
          <span class="badge" style="background: rgba(16,185,129,0.15); color: #6ee7b7; border: 1px solid rgba(16,185,129,0.25)">
            Line Chart
          </span>
        </div>
        <div class="h-48">
          <Line v-if="trendChartData" :data="trendChartData" :options="lineOptions" />
        </div>
      </div>

      <!-- ── Category Avg Rating Table ─────────────────────── -->
      <div class="glass-card-static p-6 animate-fade-up" style="animation-delay:0.35s">
        <p class="section-label mb-1">Per-Category Performance</p>
        <p class="text-sm font-semibold text-white mb-5">Average Rating by Category</p>
        <div class="space-y-3">
          <div v-for="cv in summary.category_volumes" :key="cv.category"
               class="flex items-center gap-4">
            <div class="w-36 shrink-0">
              <p class="text-xs font-medium text-white truncate">{{ cv.category }}</p>
              <p class="text-xs" style="color: var(--text-secondary)">{{ cv.count }} responses</p>
            </div>
            <!-- Progress bar -->
            <div class="flex-1 h-2 rounded-full overflow-hidden" style="background: rgba(255,255,255,0.08)">
              <div class="h-full rounded-full transition-all duration-700"
                   :style="{ width: (cv.avg_rating / 5 * 100) + '%', background: progressGradient(cv.avg_rating) }">
              </div>
            </div>
            <div class="w-14 text-right shrink-0">
              <span class="text-sm font-bold" :style="{ color: ratingColor(cv.avg_rating) }">
                {{ cv.avg_rating.toFixed(1) }}
              </span>
              <span class="text-amber-400 text-xs"> ★</span>
            </div>
          </div>
        </div>
      </div>

    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  Chart as ChartJS, CategoryScale, LinearScale, BarElement,
  ArcElement, PointElement, LineElement, Tooltip, Legend, Filler
} from 'chart.js'
import { Bar, Doughnut, Line } from 'vue-chartjs'
import { useApi } from '../composables/useApi.js'

ChartJS.register(
  CategoryScale, LinearScale, BarElement, ArcElement,
  PointElement, LineElement, Tooltip, Legend, Filler
)

const { getAnalyticsSummary } = useApi()

const summary     = ref(null)
const loading     = ref(false)
const error       = ref('')
const lastRefreshed = ref('')

// ── Derived KPIs ──────────────────────────────────────────────
const topCategory = computed(() =>
  summary.value?.category_volumes?.[0] ?? null
)
const fiveStarRate = computed(() => {
  if (!summary.value?.rating_distribution?.length) return 0
  const total = summary.value.total_feedback
  const fives = summary.value.rating_distribution.find(r => r.rating === 5)?.count || 0
  return total ? Math.round((fives / total) * 100) : 0
})

// ── Chart Palettes ────────────────────────────────────────────
const catColors = [
  'rgba(124,58,237,0.75)', 'rgba(0,229,255,0.75)', 'rgba(244,63,94,0.75)',
  'rgba(245,158,11,0.75)', 'rgba(16,185,129,0.75)', 'rgba(99,102,241,0.75)',
  'rgba(236,72,153,0.75)',
]
const ratingColors = [
  '#ef4444', '#f97316', '#eab308', '#22c55e', '#10b981'
]

// ── Chart Data ────────────────────────────────────────────────
const categoryChartData = computed(() => {
  const vols = summary.value?.category_volumes
  if (!vols?.length) return null
  return {
    labels: vols.map(v => v.category),
    datasets: [{
      label: 'Responses',
      data: vols.map(v => v.count),
      backgroundColor: catColors.slice(0, vols.length),
      borderColor: catColors.map(c => c.replace('0.75', '1')),
      borderWidth: 1,
      borderRadius: 6,
      borderSkipped: false,
    }],
  }
})

const ratingChartData = computed(() => {
  const dist = summary.value?.rating_distribution
  if (!dist?.length) return null
  return {
    labels: dist.map(d => `${d.rating} Star${d.rating > 1 ? 's' : ''}`),
    datasets: [{
      data: dist.map(d => d.count),
      backgroundColor: ratingColors,
      borderColor: 'rgba(255,255,255,0.05)',
      borderWidth: 2,
      hoverOffset: 8,
    }],
  }
})

const ratingLegend = computed(() => {
  const dist = summary.value?.rating_distribution || []
  return dist.map((d, i) => ({
    label: `${d.rating} Star${d.rating > 1 ? 's' : ''}`,
    color: ratingColors[i],
    count: d.count,
  }))
})

const trendChartData = computed(() => {
  const trend = summary.value?.recent_trend
  if (!trend?.length) return null
  return {
    labels: trend.map(t => new Date(t.date).toLocaleDateString('en-GB', { month: 'short', day: 'numeric' })),
    datasets: [
      {
        label: 'Submissions',
        data: trend.map(t => t.count),
        borderColor: '#00e5ff',
        backgroundColor: 'rgba(0,229,255,0.08)',
        borderWidth: 2,
        pointBackgroundColor: '#00e5ff',
        pointRadius: 4,
        pointHoverRadius: 6,
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Avg Rating',
        data: trend.map(t => t.avg_rating),
        borderColor: '#f59e0b',
        backgroundColor: 'transparent',
        borderWidth: 2,
        borderDash: [4, 3],
        pointBackgroundColor: '#f59e0b',
        pointRadius: 3,
        tension: 0.4,
        yAxisID: 'y1',
      },
    ],
  }
})

// ── Chart Options ─────────────────────────────────────────────
const commonTooltip = {
  backgroundColor: 'rgba(10,10,26,0.9)',
  titleColor: 'rgba(255,255,255,0.9)',
  bodyColor:  'rgba(255,255,255,0.6)',
  borderColor: 'rgba(255,255,255,0.1)',
  borderWidth: 1,
  padding: 12,
  cornerRadius: 8,
}

const barOptions = {
  responsive: true, maintainAspectRatio: false,
  plugins: { legend: { display: false }, tooltip: commonTooltip },
  scales: {
    x: {
      ticks: { color: 'rgba(255,255,255,0.45)', font: { size: 10 }, maxRotation: 30 },
      grid: { color: 'rgba(255,255,255,0.05)' },
    },
    y: {
      ticks: { color: 'rgba(255,255,255,0.45)', font: { size: 11 } },
      grid: { color: 'rgba(255,255,255,0.07)' },
      beginAtZero: true,
    },
  },
}

const doughnutOptions = {
  responsive: true, maintainAspectRatio: false,
  cutout: '68%',
  plugins: {
    legend: { display: false },
    tooltip: commonTooltip,
  },
}

const lineOptions = {
  responsive: true, maintainAspectRatio: false,
  plugins: {
    legend: {
      labels: { color: 'rgba(255,255,255,0.5)', font: { size: 11 }, usePointStyle: true },
    },
    tooltip: commonTooltip,
  },
  scales: {
    x: {
      ticks: { color: 'rgba(255,255,255,0.45)', font: { size: 10 } },
      grid: { color: 'rgba(255,255,255,0.05)' },
    },
    y: {
      ticks: { color: 'rgba(255,255,255,0.45)', font: { size: 11 } },
      grid: { color: 'rgba(255,255,255,0.07)' },
      beginAtZero: true,
    },
    y1: {
      position: 'right',
      min: 0, max: 5,
      ticks: { color: 'rgba(245,158,11,0.6)', font: { size: 11 } },
      grid: { display: false },
    },
  },
}

// ── Helpers ───────────────────────────────────────────────────
function ratingColor(r) {
  if (r >= 4.5) return '#10b981'
  if (r >= 3.5) return '#22c55e'
  if (r >= 2.5) return '#eab308'
  if (r >= 1.5) return '#f97316'
  return '#ef4444'
}

function progressGradient(r) {
  if (r >= 4)   return 'linear-gradient(90deg, #10b981, #6ee7b7)'
  if (r >= 3)   return 'linear-gradient(90deg, #eab308, #fde047)'
  return               'linear-gradient(90deg, #ef4444, #fca5a5)'
}

// ── Data Fetching ─────────────────────────────────────────────
async function loadData() {
  loading.value = true
  error.value = ''
  try {
    summary.value = await getAnalyticsSummary()
    lastRefreshed.value = new Date().toLocaleTimeString()
  } catch (err) {
    error.value = err?.response?.data?.error || err.message || 'Unknown error'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

// Expose for parent refresh calls
defineExpose({ loadData })
</script>
