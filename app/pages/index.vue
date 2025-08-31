<template>
  <div>
    <!-- Hero Section -->
    <HeroSection />

    <!-- Featured Movies Section -->
    <section class="py-20 bg-dark-950">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between mb-12">
          <div>
            <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Featured Movies</h2>
            <p class="text-gray-400 text-lg">Hand-picked favorites just for you</p>
          </div>
          <NuxtLink to="/movies" class="btn-secondary">
            View All Movies
          </NuxtLink>
        </div>

        <!-- Movies Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          <MovieCard 
            v-for="movie in featuredMovies" 
            :key="movie.id" 
            :movie="movie"
            class="animate-slide-up"
          />
        </div>
      </div>
    </section>

    <!-- Now Showing Section -->
    <section class="py-20 bg-dark-900">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between mb-12">
          <div>
            <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Now Showing</h2>
            <p class="text-gray-400 text-lg">Currently playing in theaters</p>
          </div>
          <div class="flex space-x-4">
            <button 
              @click="scrollMovies('left')"
              class="w-10 h-10 bg-dark-800 hover:bg-primary-500 text-white rounded-full flex items-center justify-center transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </button>
            <button 
              @click="scrollMovies('right')"
              class="w-10 h-10 bg-dark-800 hover:bg-primary-500 text-white rounded-full flex items-center justify-center transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </button>
          </div>
        </div>

        <!-- Horizontal Scroll Movies -->
        <div ref="scrollContainer" class="flex space-x-6 overflow-x-auto scrollbar-hide pb-4" style="scroll-behavior: smooth;">
          <div 
            v-for="movie in nowShowingMovies" 
            :key="movie.id"
            class="flex-none w-80"
          >
            <MovieCard :movie="movie" />
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-20 bg-dark-950">
      <div class="container mx-auto px-4">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Why Choose Tizetabox?</h2>
          <p class="text-gray-400 text-lg max-w-2xl mx-auto">Experience the best of Ethiopian cinema with our premium booking platform</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="text-center p-6 rounded-lg bg-dark-900 border border-dark-800 hover:border-primary-500/50 transition-all duration-300 group">
            <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:bg-primary-500/20 transition-colors">
              <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <h3 class="font-heading font-bold text-xl text-white mb-4">Secure Payments</h3>
            <p class="text-gray-400">Pay safely with Chapa, TeleBirr, or M-Pesa. Your transactions are fully secured and encrypted.</p>
          </div>

          <!-- Feature 2 -->
          <div class="text-center p-6 rounded-lg bg-dark-900 border border-dark-800 hover:border-primary-500/50 transition-all duration-300 group">
            <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:bg-primary-500/20 transition-colors">
              <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"></path>
              </svg>
            </div>
            <h3 class="font-heading font-bold text-xl text-white mb-4">Easy Booking</h3>
            <p class="text-gray-400">Book your tickets in just a few clicks. Select your seats, choose your showtime, and you're done!</p>
          </div>

          <!-- Feature 3 -->
          <div class="text-center p-6 rounded-lg bg-dark-900 border border-dark-800 hover:border-primary-500/50 transition-all duration-300 group">
            <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:bg-primary-500/20 transition-colors">
              <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
            </div>
            <h3 class="font-heading font-bold text-xl text-white mb-4">Multiple Locations</h3>
            <p class="text-gray-400">Choose from premium cinemas across Addis Ababa. Find the perfect location and showtime for you.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 bg-gradient-to-r from-primary-600 to-primary-500">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Ready for Your Next Movie Night?</h2>
        <p class="text-primary-100 text-lg mb-8 max-w-2xl mx-auto">Join thousands of movie lovers who trust Tizetabox for their cinema experience</p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <NuxtLink to="/movies" class="bg-white text-primary-600 hover:bg-gray-100 font-bold py-4 px-8 rounded-lg transition-colors">
            Browse Movies
          </NuxtLink>
          <button @click="showAuth('register')" class="bg-transparent border-2 border-white text-white hover:bg-white hover:text-primary-600 font-bold py-4 px-8 rounded-lg transition-all">
            Sign Up Today
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { useMoviesStore } from '~/stores/movies'

const moviesStore = useMoviesStore()
const scrollContainer = ref()

const featuredMovies = computed(() => moviesStore.movies.slice(0, 4))
const nowShowingMovies = computed(() => moviesStore.nowShowing)

const showAuth = (mode) => {
  const authModal = inject('authModal')
  authModal?.show(mode)
}

const scrollMovies = (direction) => {
  if (!scrollContainer.value) return
  
  const scrollAmount = 320 // Width of movie card + gap
  const currentScroll = scrollContainer.value.scrollLeft
  
  if (direction === 'left') {
    scrollContainer.value.scrollLeft = currentScroll - scrollAmount
  } else {
    scrollContainer.value.scrollLeft = currentScroll + scrollAmount
  }
}

// Page meta
useSeoMeta({
  title: 'Tizetabox - Ethiopia\'s Premier Cinema Booking Platform',
  ogTitle: 'Tizetabox - Ethiopia\'s Premier Cinema Booking Platform',
  description: 'Book movie tickets online in Ethiopia. Discover the latest movies, choose your seats, and enjoy premium cinema experiences across Addis Ababa.',
  ogDescription: 'Book movie tickets online in Ethiopia. Discover the latest movies, choose your seats, and enjoy premium cinema experiences across Addis Ababa.',
  ogImage: 'https://images.pexels.com/photos/7991579/pexels-photo-7991579.jpeg',
  twitterCard: 'summary_large_image'
})
</script>

<style>
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.animate-slide-up {
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>