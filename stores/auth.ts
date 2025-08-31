import { defineStore } from 'pinia'

interface User {
  id: number
  email: string
  name: string
  role: string
}

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    loading: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin'
  },

  actions: {
    async initializeAuth() {
      // Check for existing token in localStorage
      if (process.client) {
        const token = localStorage.getItem('auth-token')
        if (token) {
          this.token = token
          // Verify token and get user info
          await this.verifyToken()
        }
      }
    },

    async login(email: string, password: string) {
      this.loading = true
      try {
        // This would be replaced with actual API call
        const response = await $fetch('/api/auth/login', {
          method: 'POST',
          body: { email, password }
        })
        
        this.token = response.token
        this.user = response.user
        if (process.client) {
          localStorage.setItem('auth-token', response.token)
        }
        
        return { success: true, user: response.user }
      } catch (error) {
        console.error('Login error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async register(email: string, password: string, name: string) {
      this.loading = true
      try {
        const response = await $fetch('/api/auth/register', {
          method: 'POST',
          body: { email, password, name }
        })
        
        return { success: true, message: response.message }
      } catch (error) {
        console.error('Register error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async verifyToken() {
      try {
        const response = await $fetch('/api/auth/verify', {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        })
        this.user = response.user
      } catch (error) {
        this.logout()
      }
    },

    logout() {
      this.user = null
      this.token = null
      if (process.client) {
        localStorage.removeItem('auth-token')
      }
    }
  }
})