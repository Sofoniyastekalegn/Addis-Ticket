<template>
  <div class="group relative bg-dark-900 rounded-lg overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1">
    <!-- Movie Poster -->
    <div class="relative overflow-hidden">
      <img 
        :src="movie.poster || '/images/movie-placeholder.jpg'" 
        :alt="movie.title"
        class="w-full h-64 object-cover group-hover:scale-105 transition-transform duration-300"
      />
      
      <!-- Overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300">
        <div class="absolute bottom-4 left-4 right-4">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center space-x-2">
              <span class="text-yellow-400 text-sm font-semibold">{{ movie.rating }}/5</span>
              <span class="text-gray-300 text-sm">{{ movie.year }}</span>
            </div>
            <span class="text-gray-300 text-sm">{{ formatDuration(movie.duration) }}</span>
          </div>
          
          <div class="flex space-x-2">
            <button 
              @click="$emit('book', movie)"
              class="flex-1 bg-primary-500 hover:bg-primary-600 text-white text-sm font-semibold py-2 px-3 rounded transition-colors"
            >
              Book Now
            </button>
            <button 
              @click="$emit('details', movie)"
              class="bg-dark-800 hover:bg-dark-700 text-white text-sm font-semibold py-2 px-3 rounded transition-colors"
            >
              Details
            </button>
          </div>
        </div>
      </div>
      
      <!-- Rating Badge -->
      <div class="absolute top-3 right-3 bg-primary-500 text-white text-xs font-bold px-2 py-1 rounded">
        {{ movie.rating }}
      </div>
    </div>
    
    <!-- Movie Info -->
    <div class="p-4">
      <h3 class="font-heading font-bold text-lg text-white mb-2 line-clamp-2 group-hover:text-primary-400 transition-colors">
        {{ movie.title }}
      </h3>
      
      <p class="text-gray-400 text-sm mb-3 line-clamp-2">
        {{ movie.description }}
      </p>
      
      <!-- Genres -->
      <div class="flex flex-wrap gap-1 mb-3">
        <span 
          v-for="genre in movie.genres.slice(0, 2)" 
          :key="genre"
          class="text-xs bg-dark-800 text-gray-300 px-2 py-1 rounded"
        >
          {{ genre }}
        </span>
        <span v-if="movie.genres.length > 2" class="text-xs text-gray-500">
          +{{ movie.genres.length - 2 }} more
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
interface Movie {
  id: number
  title: string
  description: string
  poster: string
  backdrop: string
  rating: number
  duration: number
  year: string
  genres: string[]
  schedules?: any[]
}

interface Props {
  movie: Movie
}

defineProps<Props>()
defineEmits(['book', 'details'])

const formatDuration = (minutes: number) => {
  const hours = Math.floor(minutes / 60)
  const mins = minutes % 60
  return `${hours}h ${mins}m`
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
