<template>
    <section class="relative h-screen flex items-center justify-center overflow-hidden">
      <!-- Background Image with Overlay -->
      <div class="absolute inset-0">
        <div class="absolute inset-0 bg-gradient-to-r from-black via-black/70 to-transparent z-10"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-dark-950 via-transparent to-transparent z-10"></div>
        <img 
          v-if="featuredMovie"
          :src="featuredMovie.backdrop" 
          :alt="featuredMovie.title"
          class="w-full h-full object-cover"
        >
      </div>
  
      <!-- Content -->
      <div class="relative z-20 container mx-auto px-4 flex items-center min-h-screen">
        <div class="max-w-2xl">
          <!-- Genre Tags -->
          <div v-if="featuredMovie" class="flex flex-wrap gap-2 mb-4">
            <span 
              v-for="genre in featuredMovie.genres.slice(0, 3)" 
              :key="genre"
              class="bg-primary-500/20 text-primary-400 px-3 py-1 rounded-full text-sm font-medium border border-primary-500/30"
            >
              {{ genre }}
            </span>
          </div>
  
          <!-- Title -->
          <h1 v-if="featuredMovie" class="text-5xl md:text-7xl font-heading font-bold text-white mb-6 animate-slide-up">
            {{ featuredMovie.title }}
          </h1>
  
          <!-- Movie Meta -->
          <div v-if="featuredMovie" class="flex items-center space-x-6 mb-6 text-gray-300">
            <div class="flex items-center space-x-2">
              <svg class="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              <span class="font-semibold">{{ featuredMovie.rating }}</span>
            </div>
            <div class="flex items-center space-x-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>{{ featuredMovie.duration }}</span>
            </div>
            <div class="flex items-center space-x-2">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
              <span>{{ featuredMovie.year }}</span>
            </div>
          </div>
  
          <!-- Director -->
          <p v-if="featuredMovie" class="text-gray-300 mb-6">
            Directed by <span class="text-white font-medium">{{ featuredMovie.director }}</span>
          </p>
  
          <!-- Description -->
          <p v-if="featuredMovie" class="text-gray-300 text-lg mb-8 max-w-xl leading-relaxed">
            {{ featuredMovie.description }}
          </p>
  
          <!-- Action Buttons -->
          <div class="flex flex-col sm:flex-row gap-4">
            <button 
              v-if="featuredMovie"
              @click="bookNow"
              class="btn-primary text-lg px-8 py-4 flex items-center justify-center space-x-2 group"
            >
              <svg class="w-5 h-5 group-hover:scale-110 transition-transform" fill="currentColor" viewBox="0 0 24 24">
                <path d="M8 5v14l11-7z"/>
              </svg>
              <span>Book Tickets</span>
            </button>
            
            <button 
              v-if="featuredMovie"
              @click="watchTrailer"
              class="btn-secondary text-lg px-8 py-4 flex items-center justify-center space-x-2 group"
            >
              <svg class="w-5 h-5 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.586a1 1 0 01.707.293l2.414 2.414a1 1 0 00.707.293H15M9 10V9a2 2 0 012-2h2a2 2 0 012 2v1m-4 1v1m0 0v1m0-1h1m-1 0H9"/>
              </svg>
              <span>Watch Trailer</span>
            </button>
          </div>
        </div>
      </div>
  
      <!-- Scroll Indicator -->
      <div class="absolute bottom-8 left-1/2 transform -translate-x-1/2 z-20">
        <div class="animate-bounce">
          <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
          </svg>
        </div>
      </div>
  
      <!-- Quick Book Widget -->
      <div class="absolute bottom-8 right-8 hidden lg:block z-20">
        <QuickBookWidget />
      </div>
    </section>
  </template>
  
  <script setup>
  import { useMoviesStore } from '~/stores/movies'
  import { useAuthStore } from '~/stores/auth'
  
  const moviesStore = useMoviesStore()
  const authStore = useAuthStore()
  
  const featuredMovie = computed(() => moviesStore.featuredMovie)
  
  const bookNow = () => {
    if (!authStore.isAuthenticated) {
      const authModal = inject('authModal')
      authModal?.show('login')
      return
    }
    
    if (featuredMovie.value) {
      navigateTo(`/booking/${featuredMovie.value.id}`)
    }
  }
  
  const watchTrailer = () => {
    // Implement trailer functionality
    console.log('Watch trailer for:', featuredMovie.value?.title)
  }
  </script>