<template>
  <div>
    <NuxtRouteAnnouncer />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    
    <!-- Auth Modal -->
    <AuthModal ref="authModalRef" />
  </div>
</template>

<script setup>
import { useAuthStore } from "~/stores/auth"
import { useMoviesStore } from "~/stores/movies"

// Initialize stores
const authStore = useAuthStore()
const moviesStore = useMoviesStore()

// Auth modal reference
const authModalRef = ref()

// Provide auth modal to components
provide('authModal', {
  show: (mode) => authModalRef.value?.show(mode),
  hide: () => authModalRef.value?.hide()
})

onMounted(async () => {
  await authStore.initializeAuth()
  await moviesStore.fetchMovies()
  await moviesStore.fetchCinemas()
})
</script>