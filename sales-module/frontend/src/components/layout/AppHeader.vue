<template>
  <header class="app-header-unified">
    <div class="unified-header-left">
      <div class="unified-logo">◈</div>
      <div class="unified-brand">
        <p class="unified-subtitle">CRM ENTERPRISE</p>
        <p class="unified-title">Sales Module</p>
      </div>
    </div>
    
    <div class="unified-header-center">
      <h2>{{ $route.name }}</h2>
    </div>

    <div class="unified-header-right">
      <button @click="goToPortal" class="unified-btn-portal" title="Back to Portal">
        🏠 Portal
      </button>
      <div class="unified-user-info">
        <span class="unified-role">{{ user?.role || 'User' }}</span>
        <span class="unified-username">{{ user?.username || 'Admin' }}</span>
      </div>
      <button @click="logout" class="unified-btn-logout">
        Logout
      </button>
    </div>
  </header>
</template>

<script>
import { useAuthStore } from '../../stores/auth'

export default {
  name: 'AppHeader',
  computed: {
    user() { return useAuthStore().user }
  },
  methods: {
    goToPortal() {
      window.location.href = '/'
    },
    logout() {
      useAuthStore().logout()
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.app-header-unified {
  height: 70px;
  background: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 50;
  color: #e2e8f0;
}

.unified-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.unified-logo {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  color: white;
  background: linear-gradient(135deg, #7c3aed, #4f46e5);
  box-shadow: 0 0 16px rgba(124, 58, 237, 0.5);
}

.unified-brand {
  display: flex;
  flex-direction: column;
}

.unified-subtitle {
  font-size: 0.65rem;
  font-weight: 600;
  color: #94a3b8;
  letter-spacing: 0.1em;
  margin: 0;
  line-height: 1.2;
}

.unified-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #f8fafc;
  margin: 0;
  line-height: 1.2;
}

.unified-header-center h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #f8fafc;
  margin: 0;
}

.unified-header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.unified-btn-portal {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  cursor: pointer;
  transition: all 0.2s;
}

.unified-btn-portal:hover {
  background: rgba(255, 255, 255, 0.1);
}

.unified-user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.05);
  padding: 4px 12px 4px 4px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.unified-role {
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
  background: rgba(124, 58, 237, 0.2);
  color: #c4b5fd;
  text-transform: uppercase;
}

.unified-username {
  font-size: 0.9rem;
  color: #e2e8f0;
  font-weight: 500;
}

.unified-btn-logout {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  color: #ef4444;
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.3);
  cursor: pointer;
  transition: all 0.2s;
}

.unified-btn-logout:hover {
  background: rgba(239, 68, 68, 0.1);
}
</style>
