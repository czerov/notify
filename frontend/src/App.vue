<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()

// 应用启动时恢复认证状态
onMounted(() => {
  authStore.restoreAuth()
})
</script>

<template>
  <v-app>
    <router-view />
  </v-app>
</template>

<style></style>
