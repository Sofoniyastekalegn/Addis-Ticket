<template>
    <div class="glass-effect rounded-lg p-6 w-80">
      <h3 class="text-lg font-heading font-bold text-white mb-4 text-center">Quick Book</h3>
      
      <div class="space-y-4">
        <!-- Choose a Venue -->
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Choose a Venue</label>
          <select 
            v-model="selectedCinema"
            class="w-full bg-dark-800 text-white border border-dark-700 rounded-lg px-3 py-2 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20"
          >
            <option value="">Select Cinema</option>
            <option v-for="cinema in moviesStore.cinemas" :key="cinema.id" :value="cinema.id">
              {{ cinema.name }}
            </option>
          </select>
        </div>
  
        <!-- Choose a Film or Event -->
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Choose a Film or Event</label>
          <select 
            v-model="selectedMovie"
            class="w-full bg-dark-800 text-white border border-dark-700 rounded-lg px-3 py-2 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20"
          >
            <option value="">Select Movie</option>
            <option v-for="movie in moviesStore.movies" :key="movie.id" :value="movie.id">
              {{ movie.title }}
            </option>
          </select>
        </div>
  
        <!-- Choose a Date -->
        <div>
          <label class="block text-sm font-medium text-gray-300 mb-2">Choose a Date</label>
          <input 
            v-model="selectedDate"
            type="date"
            :min="today"
            class="w-full bg-dark-800 text-white border border-dark-700 rounded-lg px-3 py-2 focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20"
          >
        </div>
  
        <!-- Find Times Button -->
        <button 
          @click="findTimes"
          :disabled="!canFindTimes"
          class="w-full btn-primary py-3 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Find Times
        </button>
      </div>
  
      <!-- Available Showtimes -->
      <div v-if="availableShowtimes.length > 0" class="mt-6 space-y-3">
        <h4 class="text-sm font-medium text-gray-300">Available Times:</h4>
        <div class="grid grid-cols-2 gap-2">
          <button 
            v-for="showtime in availableShowtimes"
            :key="showtime.id"
            @click="selectShowtime(showtime)"
            class="bg-dark-800 hover:bg-primary-500 text-white text-sm py-2 px-3 rounded transition-colors border border-dark-700 hover:border-primary-500"
          >
            {{ showtime.time }}
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { useMoviesStore } from '~/stores/movies'
  import { useAuthStore } from '~/stores/auth'
  
  const moviesStore = useMoviesStore()
  const authStore = useAuthStore()
  
  const selectedCinema = ref('')
  const selectedMovie = ref('')
  const selectedDate = ref('')
  const availableShowtimes = ref([])
  
  const today = computed(() => {
    return new Date().toISOString().split('T')[0]
  })
  
  const canFindTimes = computed(() => {
    return selectedCinema.value && selectedMovie.value && selectedDate.value
  })
  
  const findTimes = () => {
    if (!canFindTimes.value) return
  
    const movie = moviesStore.getMovieById(selectedMovie.value)
    if (movie) {
      // Filter showtimes by selected cinema and date
      availableShowtimes.value = movie.showtimes.filter(showtime => 
        showtime.cinemaId === selectedCinema.value &&
        showtime.date === selectedDate.value
      )
    }
  }
  
  const selectShowtime = (showtime) => {
    if (!authStore.isAuthenticated) {
      const authModal = inject('authModal')
      authModal?.show('login')
      return
    }
    
    navigateTo(`/booking/${selectedMovie.value}?showtime=${showtime.id}`)
  }
  
  // Initialize with today's date
  onMounted(() => {
    selectedDate.value = today.value
  })
  </script>