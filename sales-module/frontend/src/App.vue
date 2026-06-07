<template>
  <AppLayout>
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </AppLayout>
</template>

<script>
import AppLayout from './components/layout/AppLayout.vue'

export default {
  name: 'App',
  components: { AppLayout }
}
</script>
