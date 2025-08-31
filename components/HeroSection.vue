<template>
  <section class="relative h-screen flex items-center justify-center overflow-hidden">
    <!-- Background Image -->
    <div class="absolute inset-0">
      <img 
        src="/images/hero-bg.jpg" 
        alt="Cinema Background"
        class="w-full h-full object-cover"
      />
      <div class="absolute inset-0 bg-black bg-opacity-60"></div>
    </div>
    
    <!-- Hero Content -->
    <div class="relative z-10 text-center text-white px-4 max-w-4xl mx-auto">
      <h1 class="text-5xl md:text-7xl font-heading font-bold mb-6 animate-fade-in">
        Welcome to <span class="text-primary-400">Tizetabox</span>
      </h1>
      <p class="text-xl md:text-2xl text-gray-200 mb-8 animate-fade-in-delay">
        Ethiopia's Premier Cinema Booking Platform
      </p>
      <p class="text-lg text-gray-300 mb-12 max-w-2xl mx-auto animate-fade-in-delay-2">
        Discover the latest movies, book your tickets online, and enjoy premium cinema experiences across Addis Ababa
      </p>
      
      <!-- CTA Buttons -->
      <div class="flex flex-col sm:flex-row gap-4 justify-center animate-fade-in-delay-3">
        <NuxtLink 
          to="/movies" 
          class="bg-primary-600 hover:bg-primary-700 text-white font-bold py-4 px-8 rounded-lg transition-colors text-lg"
        >
          Browse Movies
        </NuxtLink>
        <button 
          @click="showAuth('register')" 
          class="bg-transparent border-2 border-white text-white hover:bg-white hover:text-gray-900 font-bold py-4 px-8 rounded-lg transition-all text-lg"
        >
          Sign Up Today
        </button>
      </div>
      
      <!-- Features Preview -->
      <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8 animate-fade-in-delay-4">
        <div class="text-center">
          <div class="w-16 h-16 bg-primary-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <h3 class="font-bold text-lg mb-2">Secure Payments</h3>
          <p class="text-gray-300">Pay safely with Chapa, TeleBirr, or M-Pesa</p>
        </div>
        
        <div class="text-center">
          <div class="w-16 h-16 bg-primary-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"></path>
            </svg>
          </div>
          <h3 class="font-bold text-lg mb-2">Easy Booking</h3>
          <p class="text-gray-300">Book tickets in just a few clicks</p>
        </div>
        
        <div class="text-center">
          <div class="w-16 h-16 bg-primary-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
            </svg>
          </div>
          <h3 class="font-bold text-lg mb-2">Multiple Locations</h3>
          <p class="text-gray-300">Choose from premium cinemas across Addis Ababa</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const showAuth = (mode) => {
  const authModal = inject('authModal')
  authModal?.show(mode)
}
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 1s ease-out;
}

.animate-fade-in-delay {
  animation: fadeIn 1s ease-out 0.2s both;
}

.animate-fade-in-delay-2 {
  animation: fadeIn 1s ease-out 0.4s both;
}

.animate-fade-in-delay-3 {
  animation: fadeIn 1s ease-out 0.6s both;
}

.animate-fade-in-delay-4 {
  animation: fadeIn 1s ease-out 0.8s both;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
