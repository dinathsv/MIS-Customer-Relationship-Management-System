<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-logo">
      <div class="logo-icon"><i class='bx bx-star' style='color: var(--color-primary);'></i></div>
      <span class="logo-text" v-show="!collapsed">SalesHub</span>
    </div>


    <nav class="sidebar-nav">
      <template v-for="item in navItems" :key="item.path">
        <a v-if="item.external" :href="item.path" target="_blank" class="nav-item">
          <span class="nav-icon"><i :class="item.icon"></i></span>
          <span class="nav-label" v-show="!collapsed">{{ item.label }}</span>
        </a>
        <router-link v-else :to="item.path" class="nav-item" :class="{ active: $route.path === item.path }">
          <span class="nav-icon"><i :class="item.icon"></i></span>
          <span class="nav-label" v-show="!collapsed">{{ item.label }}</span>
        </router-link>
      </template>
    </nav>

    <div class="sidebar-footer">
      <button class="nav-item" @click="collapsed = !collapsed">
        <span class="nav-icon"><i :class="collapsed ? 'bx bx-chevron-right' : 'bx bx-chevron-left'"></i></span>
        <span class="nav-label" v-show="!collapsed">Collapse</span>
      </button>
    </div>
  </aside>
</template>

<script>
export default {
  name: 'AppSidebar',
  data() {
    return {
      collapsed: false,
      navItems: [
        { path: '/', label: 'Dashboard', icon: 'bx bx-home-alt' },
        { path: '/sales', label: 'Sales', icon: 'bx bx-dollar-circle' },
        { path: '/reports', label: 'Reports', icon: 'bx bx-line-chart' }
      ]
    }
  }
}
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: var(--sidebar-width);
  height: 100vh;
  background: var(--bg-card);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: width var(--transition-slow);
}

.sidebar.collapsed { width: var(--sidebar-collapsed); }

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.logo-icon { font-size: 1.8rem; }

.logo-text {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--text-primary);
  white-space: nowrap;
}


.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  text-decoration: none;
  transition: all var(--transition-fast);
  cursor: pointer;
  border: none;
  background: none;
  font-family: inherit;
  font-size: 0.9rem;
  width: 100%;
  text-align: left;
}

.nav-item:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--color-info-bg);
  color: var(--color-info);
  border-left: 3px solid var(--color-info);
}

.nav-icon { font-size: 1.2rem; min-width: 24px; text-align: center; }
.nav-label { white-space: nowrap; }

.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--border-color);
}
</style>
