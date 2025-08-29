import { defineStore } from 'pinia'
import { useAuth } from '~/composables/useAuth'

interface User {
  id: string
  name: string
  email: string
  phone?: string
  avatar?: string
  favorites: string[]
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  loading: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    isAuthenticated: false,
    loading: false
  }),

  getters: {
    currentUser: (state) => state.user,
    isLoggedIn: (state) => state.isAuthenticated
  },

  actions: {
    async login(credentials: { email: string; password: string }) {
      const { login } = useAuth()
      const result = await login(credentials)
      
      if (result.success) {
        await this.initializeAuth()
      }
      
      return result
    },

    async register(userData: { name: string; email: string; password: string; phone?: string }) {
      const { register } = useAuth()
      const result = await register(userData)
      
      if (result.success) {
        await this.initializeAuth()
      }
      
      return result
    },

    logout() {
      const { logout } = useAuth()
      logout()
      
      this.user = null
      this.token = null
      this.isAuthenticated = false
    },

    async initializeAuth() {
      const { initializeAuth } = useAuth()
      await initializeAuth()
      
      if (process.client) {
        const token = localStorage.getItem('auth-token')
        const userData = localStorage.getItem('user')
        
        if (token && userData) {
          this.token = token
          this.user = JSON.parse(userData)
          this.isAuthenticated = true
        }
      }
    },

    toggleFavorite(movieId: string) {
      const { toggleFavorite } = useAuth()
      
      if (this.user) {
        toggleFavorite(movieId)
        
        const index = this.user.favorites.indexOf(movieId)
        if (index > -1) {
          this.user.favorites.splice(index, 1)
        } else {
          this.user.favorites.push(movieId)
        }
        
        if (process.client) {
          localStorage.setItem('user', JSON.stringify(this.user))
        }
      }
    }
  }
})