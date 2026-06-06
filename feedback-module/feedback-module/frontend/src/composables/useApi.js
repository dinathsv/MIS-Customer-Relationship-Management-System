import axios from 'axios'

// Resolve API base — works both in Docker (Nginx proxy) and local dev (Vite proxy)
const BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

const api = axios.create({
  baseURL: BASE_URL,
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

export function useApi() {

  async function submitFeedback(payload) {
    const { data } = await api.post('/api/feedback', payload)
    return data
  }

  async function listFeedback(params = {}) {
    const { data } = await api.get('/api/feedback', { params })
    return data
  }

  async function getAnalyticsSummary() {
    const { data } = await api.get('/api/analytics/summary')
    return data
  }

  async function getFeedbackByCustomer(customerId) {
    const { data } = await api.get(`/api/feedback/customer/${customerId}`)
    return data
  }

  async function checkHealth() {
    const { data } = await api.get('/health')
    return data
  }

  return { submitFeedback, listFeedback, getAnalyticsSummary, getFeedbackByCustomer, checkHealth }
}
