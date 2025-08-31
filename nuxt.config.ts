// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt'
  ],
  css: ['~/assets/css/main.css'],
  components: true, // Auto-import components
  app: {
    head: {
      title: 'Tizetabox - Cinema Tickets & Movie Schedules',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Book cinema tickets and browse movie schedules at Tizetabox' }
      ]
    }
  }
})
