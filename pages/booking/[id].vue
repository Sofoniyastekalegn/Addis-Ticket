<template>
    <div class="pt-16 min-h-screen bg-dark-950">
      <div v-if="loading" class="flex items-center justify-center min-h-screen">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
      </div>
  
      <div v-else-if="movie" class="container mx-auto px-4 py-8">
        <!-- Movie Header -->
        <div class="bg-dark-900 rounded-xl p-6 mb-8 border border-dark-800">
          <div class="flex flex-col md:flex-row gap-6">
            <img 
              :src="movie.poster" 
              :alt="movie.title"
              class="w-32 h-48 object-cover rounded-lg flex-shrink-0"
            >
            <div class="flex-1">
              <h1 class="text-3xl font-heading font-bold text-white mb-2">{{ movie.title }}</h1>
              <div class="flex items-center space-x-4 text-gray-400 mb-4">
                <span class="flex items-center space-x-1">
                  <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                  <span>{{ movie.rating }}</span>
                </span>
                <span>{{ movie.duration }}</span>
                <span>{{ movie.year }}</span>
              </div>
              <p class="text-gray-300 mb-4">{{ movie.description }}</p>
              <p class="text-sm text-gray-400">
                Directed by <span class="text-white">{{ movie.director }}</span>
              </p>
            </div>
          </div>
        </div>
  
        <!-- Booking Steps -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <!-- Step 1: Select Showtime -->
          <div class="lg:col-span-2">
            <div class="bg-dark-900 rounded-xl p-6 border border-dark-800 mb-6">
              <h2 class="text-xl font-heading font-bold text-white mb-4">Select Showtime</h2>
              <div class="space-y-4">
                <div 
                  v-for="showtime in movie.showtimes" 
                  :key="showtime.id"
                  @click="selectShowtime(showtime)"
                  class="p-4 bg-dark-800 rounded-lg border border-dark-700 hover:border-primary-500 cursor-pointer transition-all"
                  :class="{ 'border-primary-500 bg-primary-500/10': selectedShowtime?.id === showtime.id }"
                >
                  <div class="flex items-center justify-between">
                    <div>
                      <h3 class="font-semibold text-white">{{ showtime.cinema.name }}</h3>
                      <p class="text-gray-400 text-sm">{{ showtime.cinema.location }}</p>
                    </div>
                    <div class="text-right">
                      <p class="text-white font-semibold">{{ showtime.time }}</p>
                      <p class="text-gray-400 text-sm">{{ showtime.screenType }}</p>
                    </div>
                    <div class="text-right">
                      <p class="text-primary-400 font-bold">{{ showtime.price }} ETB</p>
                      <p class="text-gray-400 text-sm">{{ showtime.availableSeats }} seats left</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
  
            <!-- Step 2: Select Seats -->
            <div v-if="selectedShowtime" class="bg-dark-900 rounded-xl p-6 border border-dark-800">
              <h2 class="text-xl font-heading font-bold text-white mb-4">Select Seats</h2>
              <SeatSelection 
                :seats="seats"
                @seat-selected="selectSeat"
                @seat-deselected="deselectSeat"
              />
            </div>
          </div>
  
          <!-- Step 3: Booking Summary -->
          <div class="lg:col-span-1">
            <div class="bg-dark-900 rounded-xl p-6 border border-dark-800 sticky top-24">
              <h2 class="text-xl font-heading font-bold text-white mb-4">Booking Summary</h2>
              
              <div v-if="selectedShowtime" class="space-y-4">
                <!-- Movie Info -->
                <div class="pb-4 border-b border-dark-800">
                  <h3 class="font-semibold text-white">{{ movie.title }}</h3>
                  <p class="text-gray-400 text-sm">{{ selectedShowtime.cinema.name }}</p>
                  <p class="text-gray-400 text-sm">{{ selectedShowtime.date }} at {{ selectedShowtime.time }}</p>
                </div>
  
                <!-- Selected Seats -->
                <div v-if="selectedSeats.length > 0" class="pb-4 border-b border-dark-800">
                  <h4 class="font-medium text-white mb-2">Selected Seats</h4>
                  <div class="flex flex-wrap gap-2">
                    <span 
                      v-for="seat in selectedSeats" 
                      :key="seat.id"
                      class="bg-primary-500/20 text-primary-400 px-2 py-1 rounded text-sm"
                    >
                      {{ seat.id }}
                    </span>
                  </div>
                </div>
  
                <!-- Price Breakdown -->
                <div v-if="selectedSeats.length > 0" class="pb-4 border-b border-dark-800">
                  <div class="space-y-2">
                    <div class="flex justify-between text-gray-400">
                      <span>Tickets ({{ selectedSeats.length }})</span>
                      <span>{{ totalAmount }} ETB</span>
                    </div>
                    <div class="flex justify-between text-gray-400">
                      <span>Service Fee</span>
                      <span>10 ETB</span>
                    </div>
                    <div class="flex justify-between font-bold text-white text-lg">
                      <span>Total</span>
                      <span>{{ totalAmount + 10 }} ETB</span>
                    </div>
                  </div>
                </div>
  
                <!-- Proceed Button -->
                <button 
                  v-if="selectedSeats.length > 0"
                  @click="proceedToPayment"
                  class="w-full btn-primary py-3"
                >
                  Proceed to Payment
                </button>
              </div>
  
              <div v-else class="text-center py-8">
                <svg class="w-12 h-12 text-gray-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"></path>
                </svg>
                <p class="text-gray-400">Select a showtime to continue</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { useMovies } from '~/composables/useMovies'
  import { useBookings } from '~/composables/useBookings'
  import { useAuthStore } from '~/stores/auth'
  
  const route = useRoute()
  const authStore = useAuthStore()
  
  const movieId = route.params.id as string
  const selectedShowtime = ref(null)
  const loading = ref(true)
  
  const { getMovieById } = useMovies()
  const { 
    seats, 
    selectedSeats, 
    totalAmount, 
    initializeBooking, 
    selectSeat, 
    deselectSeat 
  } = useBookings()
  
  const movie = ref(null)
  
  // Initialize booking
  onMounted(async () => {
    try {
      movie.value = await getMovieById(movieId)
      if (movie.value) {
        await initializeBooking(movieId, '')
      }
    } finally {
      loading.value = false
    }
  })
  
  const selectShowtime = (showtime) => {
    selectedShowtime.value = showtime
  }
  
  const proceedToPayment = () => {
    if (!authStore.isAuthenticated) {
      const authModal = inject('authModal')
      authModal?.show('login')
      return
    }
    
    navigateTo(`/payment/${movieId}?showtime=${selectedShowtime.value.id}`)
  }
  
  // Page meta
  useSeoMeta({
    title: `Book ${movie.value?.title || 'Movie'} - Tizetabox`,
    description: `Book tickets for ${movie.value?.title || 'this movie'} at Tizetabox cinemas.`
  })
  </script>
  </template>