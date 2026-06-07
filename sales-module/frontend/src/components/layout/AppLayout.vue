<template>
  <div class="app-layout-wrapper">
    <AppHeader />
    <div class="app-layout">
      <AppSidebar />
      <div class="main-area">
        <main class="main-content">
          <slot />
        </main>
      </div>
    </div>
  </div>
</template>

<script>
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

export default {
  name: 'AppLayout',
  components: { AppSidebar, AppHeader }
}
</script>

<style scoped>
.app-layout-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.app-layout {
  display: flex;
  flex: 1;
  padding-top: 70px; /* Space for the fixed full-width header */
}

.main-area {
  flex: 1;
  margin-left: var(--sidebar-width);
  transition: margin-left var(--transition-slow);
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 70px);
}

.main-content {
  flex: 1;
  padding: 28px;
  animation: fadeIn 0.3s ease;
}
</style>
