<template>
  <div class="min-h-screen bg-dark-950">
    <!-- Header -->
    <div class="bg-dark-900 py-8">
      <div class="container mx-auto px-4">
        <h1 class="text-4xl font-heading font-bold text-white mb-4">All Movies</h1>
        <p class="text-gray-400 text-lg">Discover the latest and greatest films</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-dark-900 border-b border-dark-800 py-6">
      <div class="container mx-auto px-4">
        <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
          <!-- Search -->
          <div class="relative flex-1 max-w-md">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search movies..."
              class="w-full bg-dark-800 text-white px-4 py-3 pl-10 rounded-lg border border-dark-700 focus:border-primary-500 focus:outline-none"
            />
            <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>

          <!-- Genre Filters -->
          <div class="flex flex-wrap gap-2">
            <button
              v-for="genre in allGenres"
              :key="genre"
              @click="toggleGenreFilter(genre)"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
                selectedGenres.includes(genre)
                  ? 'bg-primary-500 text-white'
                  : 'bg-dark-800 text-gray-300 hover:bg-dark-700'
              ]"
            >
              {{ genre }}
            </button>
          </div>

          <!-- Clear Filters -->
          <button
            v-if="searchQuery || selectedGenres.length > 0"
            @click="clearFilters"
            class="text-primary-400 hover:text-primary-300 text-sm font-medium"
          >
            Clear Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Movies Grid -->
    <div class="container mx-auto px-4 py-12">
      <div v-if="loading" class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
      </div>

      <div v-else-if="filteredMovies.length === 0" class="text-center py-20">
        <div class="text-gray-400 text-lg mb-4">No movies found</div>
        <button @click="clearFilters" class="text-primary-400 hover:text-primary-300">
          Clear all filters
        </button>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
        <MovieCard
          v-for="movie in filteredMovies"
          :key="movie.id"
          :movie="movie"
          @book="bookMovie"
          @details="viewMovieDetails"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useMoviesStore } from '~/stores/movies'

const moviesStore = useMoviesStore()

const searchQuery = computed({
  get: () => moviesStore.searchQuery,
  set: (value) => moviesStore.setSearchQuery(value)
})

const selectedGenres = computed(() => moviesStore.selectedGenres)
const allGenres = computed(() => moviesStore.allGenres)
const filteredMovies = computed(() => moviesStore.filteredMovies)
const loading = computed(() => moviesStore.loading)

const toggleGenreFilter = (genre: string) => {
  moviesStore.toggleGenreFilter(genre)
}

const clearFilters = () => {
  moviesStore.clearFilters()
}

const bookMovie = (movie: any) => {
  navigateTo(`/booking/${movie.id}`)
}

const viewMovieDetails = (movie: any) => {
  navigateTo(`/movies/${movie.id}`)
}

// Fetch movies on page load
onMounted(() => {
  moviesStore.fetchMovies()
})

// Page meta
useSeoMeta({
  title: 'All Movies - Tizetabox',
  description: 'Browse all movies available for booking on Tizetabox. Find your next favorite film and book tickets online.',
})
</script>