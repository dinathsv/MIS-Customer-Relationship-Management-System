<template>
  <div class="products-page">
    <div class="page-header" style="display:flex; justify-content:space-between; align-items:center;">
      <div>
        <h1>Products</h1>
        <p>Product inventory management</p>
      </div>
      <button class="btn btn-primary" @click="openModal('add')"><i class="bx bx-plus"></i> Add Product</button>
    </div>

    <div class="filters-bar" style="margin-bottom:20px; display:flex; gap:12px;">
      <div style="position:relative; display:flex; align-items:center; max-width:300px; width:100%;">
        <i class="bx bx-search" style="position:absolute; left:12px; color:var(--text-muted);"></i>
        <input v-model="search" type="text" class="form-control" placeholder="Search products..." style="padding-left:36px; width:100%;" />
      </div>
      
      <select v-model="selectedCategory" class="form-control" style="max-width:200px;">
        <option value="">All Categories</option>
        <option v-for="category in categories" :key="category" :value="category">
          {{ category }}
        </option>
      </select>
    </div>

    <LoadingSpinner v-if="loading" />

    <template v-else>
      <div class="card">
        <table class="data-table">
          <thead>
            <tr>
              <th>PRODUCT NAME</th>
              <th>SKU</th>
              <th>CATEGORY</th>
              <th>PRICE</th>
              <th>STOCK</th>
              <th>DESCRIPTION</th>
              <th>ACTIONS</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in filteredProducts" :key="p.id">
              <td class="text-primary font-medium">{{ p.name }}</td>
              <td class="text-muted">{{ p.sku || '-' }}</td>
              <td>{{ p.category || '-' }}</td>
              <td style="font-weight:600;">Rs. {{ fmt(p.price) }}</td>
              <td :class="{'text-danger': p.stock_qty <= 5}">{{ p.stock_qty || 0 }}</td>
              <td class="text-muted">{{ truncate(p.description, 30) || '-' }}</td>
              <td style="display:flex; gap:8px;">
                <button class="btn btn-ghost btn-sm" title="Edit"><i class="bx bx-edit"></i></button>
                <button class="btn btn-danger btn-sm" title="Delete" style="background:#fee2e2; color:#ef4444; border:none; padding:4px 8px; border-radius:4px; cursor:pointer;"><i class="bx bx-trash"></i></button>
              </td>
            </tr>
            <tr v-if="!filteredProducts.length">
                <td colspan="7" class="text-center text-muted" style="padding:40px;">No products found</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <!-- Product Modal -->
    <Modal v-if="showModal" :title="modalType === 'add' ? 'Add New Product' : 'Edit Product'" @close="closeModal" maxWidth="500px">
      <form @submit.prevent="saveProduct">
        <div class="form-group" style="margin-bottom:15px;">
          <label class="form-label">Product Name *</label>
          <input v-model="form.name" type="text" class="form-control" required style="width:100%;" />
        </div>
        <div class="form-group" style="margin-bottom:15px; display:flex; gap:15px;">
          <div style="flex:1">
            <label class="form-label">SKU *</label>
            <input v-model="form.sku" type="text" class="form-control" required style="width:100%;" />
          </div>
          <div style="flex:1">
            <label class="form-label">Category</label>
            <input v-model="form.category" type="text" class="form-control" style="width:100%;" />
          </div>
        </div>
        <div class="form-group" style="margin-bottom:15px; display:flex; gap:15px;">
          <div style="flex:1">
            <label class="form-label">Price *</label>
            <input v-model.number="form.price" type="number" step="0.01" class="form-control" required style="width:100%;" />
          </div>
          <div style="flex:1">
            <label class="form-label">Stock Qty</label>
            <input v-model.number="form.stock_qty" type="number" class="form-control" style="width:100%;" />
          </div>
        </div>
        <div class="form-group" style="margin-bottom:20px;">
          <label class="form-label">Description</label>
          <textarea v-model="form.description" class="form-control" rows="3" style="width:100%;"></textarea>
        </div>
        <div style="display:flex; justify-content:flex-end; gap:10px; border-top:1px solid var(--border-color); padding-top:15px;">
          <button type="button" class="btn" @click="closeModal">Cancel</button>
          <button type="submit" class="btn btn-primary" :disabled="saving">
            {{ saving ? 'Saving...' : 'Save Product' }}
          </button>
        </div>
      </form>
    </Modal>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia'
import { useProductsStore } from '../stores/products'
import LoadingSpinner from '../components/common/LoadingSpinner.vue'
import Modal from '../components/common/Modal.vue'

export default {
  name: 'ProductsView',
  components: { LoadingSpinner, Modal },
  data() {
    return {
      search: '',
      selectedCategory: '',
      showModal: false,
      modalType: 'add',
      saving: false,
      form: {
        name: '',
        sku: '',
        category: '',
        price: null,
        stock_qty: null,
        description: ''
      }
    }
  },
  computed: {
    ...mapState(useProductsStore, ['products', 'loading']),
    categories() {
      const cats = this.products.map(p => p.category).filter(Boolean)
      return [...new Set(cats)]
    },
    filteredProducts() {
      let filtered = this.products
      
      if (this.selectedCategory) {
        filtered = filtered.filter(p => p.category === this.selectedCategory)
      }
      
      if (this.search) {
        const q = this.search.toLowerCase()
        filtered = filtered.filter(p => 
          p.name.toLowerCase().includes(q) || 
          (p.sku || '').toLowerCase().includes(q)
        )
      }
      return filtered
    }
  },
  methods: {
    ...mapActions(useProductsStore, ['fetchProducts', 'createProduct']),
    fmt(n) { 
      return Number(n || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) 
    },
    truncate(str, length) {
      if (!str) return ''
      return str.length > length ? str.substring(0, length) + '...' : str
    },
    openModal(type) {
      this.modalType = type
      this.form = { name: '', sku: '', category: '', price: null, stock_qty: null, description: '' }
      this.showModal = true
    },
    closeModal() {
      this.showModal = false
    },
    async saveProduct() {
      this.saving = true
      try {
        if (this.modalType === 'add') {
          await this.createProduct(this.form)
          alert('Product added successfully!')
        }
        this.closeModal()
      } catch (err) {
        alert(err.message)
      } finally {
        this.saving = false
      }
    }
  },
  mounted() {
    this.fetchProducts()
  }
}
</script>

<style scoped>
.page-header {
  margin-bottom: 24px;
}
.page-header h1 {
  font-size: 1.5rem;
  margin: 0 0 4px 0;
}
.page-header p {
  color: var(--text-muted);
  margin: 0;
  font-size: 0.9rem;
}
</style>
