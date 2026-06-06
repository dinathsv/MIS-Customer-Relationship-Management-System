import { defineStore } from 'pinia'
import api from '../services/api'

export const useProductsStore = defineStore('products', {
  state: () => ({
    products: [],
    loading: false,
    error: null
  }),
  actions: {
    async fetchProducts() {
      this.loading = true
      this.error = null
      try {
        const { data } = await api.get('/products')
        this.products = data || []
      } catch (err) {
        this.error = err.response?.data?.error || err.message
        console.error('Failed to fetch products:', err)
      } finally {
        this.loading = false
      }
    },
    async createProduct(productData) {
      try {
        const { data } = await api.post('/products', productData)
        this.products.push(data)
        return data
      } catch (err) {
        throw new Error(err.response?.data?.error || 'Failed to create product')
      }
    }
  }
})
