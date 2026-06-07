 <template>
  <div class="glass-card p-8 animate-fade-up">
    <!-- Header -->
    <div class="mb-8">
      <p class="section-label mb-2">Customer Experience</p>
      <h2 class="text-2xl font-bold font-display text-white">Submit Feedback</h2>
      <p class="text-sm mt-1" style="color: var(--text-secondary)">
        Your insight drives product excellence. All feedback is linked to your CRM profile.
      </p>
    </div>

    <!-- Success State -->
    <Transition name="fade">
      <div v-if="submitted" class="rounded-xl p-6 mb-6 text-center toast-success">
        <div class="text-3xl mb-2">🎉</div>
        <p class="font-semibold text-emerald-300">Feedback submitted successfully!</p>
        <p class="text-xs mt-1 text-emerald-400/70">Feedback ID: <span class="font-mono">{{ lastFeedbackId }}</span></p>
        <button @click="resetForm" class="btn-secondary mt-4 text-xs">Submit Another</button>
      </div>
    </Transition>

    <form v-if="!submitted" @submit.prevent="handleSubmit" novalidate>
      <!-- Email Address -->
      <div class="mb-5">
        <label class="block text-xs font-semibold uppercase tracking-wider mb-2" style="color: var(--text-secondary)">
          Email Address <span class="text-rose-400">*</span>
        </label>
        <div class="relative">
          <input
            type="email"
            v-model="form.customer_id"
            placeholder="Enter your email"
            class="glass-input w-full"
            :class="{ 'border-rose-500/60': errors.customer_id }"
            @input="validateField('customer_id')"
          />
        </div>
        <p v-if="errors.customer_id" class="text-xs text-rose-400 mt-1.5">{{ errors.customer_id }}</p>
      </div>

      <!-- Star Rating -->
      <div class="mb-5">
        <label class="block text-xs font-semibold uppercase tracking-wider mb-3" style="color: var(--text-secondary)">
          Overall Rating <span class="text-rose-400">*</span>
        </label>
        <div class="flex items-center gap-1">
          <span
            v-for="star in 5"
            :key="star"
            class="star-btn"
            :class="{ active: star <= (hoveredStar || form.rating) }"
            @click="form.rating = star"
            @mouseenter="hoveredStar = star"
            @mouseleave="hoveredStar = 0"
            role="button"
            :aria-label="`Rate ${star} star${star > 1 ? 's' : ''}`"
          >★</span>
          <span class="ml-3 text-sm font-medium" style="color: var(--text-secondary)">
            {{ ratingLabel }}
          </span>
        </div>
        <p v-if="errors.rating" class="text-xs text-rose-400 mt-1.5">{{ errors.rating }}</p>
      </div>

      <!-- Category -->
      <div class="mb-5">
        <label class="block text-xs font-semibold uppercase tracking-wider mb-2" style="color: var(--text-secondary)">
          Category <span class="text-rose-400">*</span>
        </label>
        <div class="relative">
          <select
            v-model="form.category"
            class="glass-select pr-10"
            :class="{ 'border-rose-500/60': errors.category }"
            @change="validateField('category')"
          >
            <option value="" disabled>Select feedback category...</option>
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
          <span class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-sm"
                style="color: var(--text-secondary)">▼</span>
        </div>
        <p v-if="errors.category" class="text-xs text-rose-400 mt-1.5">{{ errors.category }}</p>
      </div>

      <!-- Comments -->
      <div class="mb-6">
        <label class="block text-xs font-semibold uppercase tracking-wider mb-2" style="color: var(--text-secondary)">
          Comments <span style="color: var(--text-secondary)" class="normal-case font-normal">(optional)</span>
        </label>
        <textarea
          v-model="form.comments"
          rows="4"
          placeholder="Share your experience in detail..."
          class="glass-input resize-none"
          maxlength="1000"
        ></textarea>
        <p class="text-right text-xs mt-1" style="color: var(--text-secondary)">
          {{ form.comments.length }}/1000
        </p>
      </div>

      <!-- Error Banner -->
      <Transition name="fade">
        <div v-if="submitError" class="toast-error rounded-xl px-4 py-3 text-sm mb-4">
          ⚠️ {{ submitError }}
        </div>
      </Transition>

      <!-- Submit -->
      <button type="submit" class="btn-primary w-full py-3.5" :disabled="loading">
        <span v-if="!loading" class="flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
          </svg>
          Submit Feedback
        </span>
        <span v-else class="flex items-center justify-center gap-2">
          <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
          </svg>
          Submitting...
        </span>
      </button>
    </form>

    <!-- Recent Submissions (last 3) -->
    <div v-if="recentList.length && !submitted" class="mt-8 pt-6 border-t" style="border-color: var(--glass-border)">
      <p class="section-label mb-4">Recent Submissions</p>
      <div class="space-y-2">
        <div v-for="fb in recentList" :key="fb.feedback_id"
             class="flex items-center justify-between rounded-xl p-3"
             style="background: rgba(255,255,255,0.04); border: 1px solid var(--glass-border)">
          <div>
            <p class="text-xs font-mono" style="color: var(--text-secondary)">{{ fb.customer_id }}</p>
            <p class="text-xs text-white/70">{{ fb.category }}</p>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-amber-400 text-sm">{{ '★'.repeat(fb.rating) }}<span class="text-white/20">{{ '★'.repeat(5 - fb.rating) }}</span></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useApi } from '../composables/useApi.js'

const { submitFeedback, listFeedback } = useApi()

const categories = [
  'Product Quality',
  'Customer Support',
  'Delivery & Shipping',
  'Pricing & Value',
  'Website & App',
  'Returns & Refunds',
  'Other',
]

const ratingLabels = ['', 'Terrible', 'Poor', 'Average', 'Good', 'Excellent']

const users = ref([])
const form = ref({ customer_id: '', rating: 0, category: '', comments: '' })
const errors = ref({})
const hoveredStar = ref(0)
const loading = ref(false)
const submitted = ref(false)
const submitError = ref('')
const lastFeedbackId = ref('')
const recentList = ref([])

const ratingLabel = computed(() => ratingLabels[hoveredStar.value || form.value.rating] || '')

function validateField(field) {
  errors.value[field] = ''
  if (field === 'customer_id') {
    if (!form.value.customer_id) {
      errors.value.customer_id = 'Please enter your email'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.customer_id)) {
      errors.value.customer_id = 'Please enter a valid email address'
    }
  }
  if (field === 'rating' && !form.value.rating) {
    errors.value.rating = 'Please select a rating'
  }
  if (field === 'category' && !form.value.category) {
    errors.value.category = 'Please select a category'
  }
}

function validate() {
  ;['customer_id', 'rating', 'category'].forEach(validateField)
  return !Object.values(errors.value).some(Boolean)
}

async function handleSubmit() {
  if (!validate()) return
  loading.value = true
  submitError.value = ''
  try {
    const result = await submitFeedback({
      customer_id: String(form.value.customer_id),
      rating:      form.value.rating,
      category:    form.value.category,
      comments:    form.value.comments.trim(),
    })
    lastFeedbackId.value = result.feedback?.feedback_id || 'N/A'
    submitted.value = true
    await refreshRecent()
    emit('submitted')
  } catch (err) {
    submitError.value = err?.response?.data?.error || 'Failed to submit. Please try again.'
  } finally {
    loading.value = false
  }
}

async function refreshRecent() {
  try {
    const data = await listFeedback({ limit: 3 })
    recentList.value = data.data || []
  } catch (_) {}
}

function resetForm() {
  form.value = { customer_id: '', rating: 0, category: '', comments: '' }
  errors.value = {}
  submitted.value = false
  submitError.value = ''
}

const emit = defineEmits(['submitted'])

// Auto-fill email if logged in
try {
  const crmUser = localStorage.getItem('crm_user');
  if (crmUser) {
    const user = JSON.parse(crmUser);
    if (user.email) {
      form.value.customer_id = user.email;
    }
  }
} catch (e) {}

// Load recent on mount
refreshRecent()
</script>
