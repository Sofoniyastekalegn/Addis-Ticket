<template>
  <div class="min-h-screen bg-dark-950">
    <!-- Hero Section -->
    <section class="relative h-screen flex items-center justify-center overflow-hidden">
      <!-- Background -->
      <div class="absolute inset-0 bg-gradient-to-br from-primary-900 to-dark-900"></div>
      
      <!-- Hero Content -->
      <div class="relative z-10 text-center text-white px-4 max-w-4xl mx-auto">
        <h1 class="text-5xl md:text-7xl font-bold mb-6">
          Welcome to <span class="text-primary-400">Tizetabox</span>
        </h1>
        <p class="text-xl md:text-2xl text-gray-200 mb-8">
          Ethiopia's Premier Cinema Booking Platform
        </p>
        <p class="text-lg text-gray-300 mb-12 max-w-2xl mx-auto">
          Discover the latest movies, book your tickets online, and enjoy premium cinema experiences across Addis Ababa
        </p>
        
        <!-- CTA Buttons -->
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <button class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-4 px-8 rounded-lg transition-colors text-lg">
            Browse Movies
          </button>
          <button class="bg-transparent border-2 border-white text-white hover:bg-white hover:text-gray-900 font-bold py-4 px-8 rounded-lg transition-all text-lg">
            Sign Up Today
          </button>
        </div>
      </div>
    </section>

    <!-- Featured Movies Section -->
    <section class="py-20 bg-dark-950">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between mb-12">
          <div>
            <h2 class="text-3xl md:text-4xl font-bold text-white mb-4">Featured Movies</h2>
            <p class="text-gray-400 text-lg">Hand-picked favorites just for you</p>
          </div>
        </div>

        <!-- Movies Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          <div 
            v-for="movie in movies" 
            :key="movie.id" 
            class="bg-dark-900 rounded-lg overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300"
          >
            <div class="h-64 bg-gradient-to-br from-primary-600 to-primary-800 flex items-center justify-center">
              <span class="text-white text-lg font-bold">{{ movie.title }}</span>
            </div>
            <div class="p-4">
              <h3 class="font-bold text-lg text-white mb-2">{{ movie.title }}</h3>
              <p class="text-gray-400 text-sm mb-3">{{ movie.description }}</p>
              <div class="flex items-center justify-between">
                <span class="text-yellow-400 text-sm font-semibold">{{ movie.rating }}/5</span>
                <span class="text-gray-300 text-sm">{{ movie.year }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
const movies = ref([
  {
    id: 1,
    title: "The Lion King",
    description: "A young lion prince flees his kingdom only to learn the true meaning of responsibility and bravery.",
    rating: 4.5,
    year: "2019"
  },
  {
    id: 2,
    title: "Avengers: Endgame",
    description: "After the devastating events of Avengers: Infinity War, the universe is in ruins.",
    rating: 4.8,
    year: "2019"
  },
  {
    id: 3,
    title: "Black Panther",
    description: "T'Challa, heir to the hidden but advanced kingdom of Wakanda, must step forward to lead his people into a new future.",
    rating: 4.7,
    year: "2018"
  },
  {
    id: 4,
    title: "Wonder Woman",
    description: "When a pilot crashes and tells of conflict in the outside world, Diana, an Amazonian warrior in training, leaves home to fight a war.",
    rating: 4.6,
    year: "2017"
  }
])
</script>
