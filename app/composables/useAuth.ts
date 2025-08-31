import { ref, computed } from 'vue'
import { useCookie } from 'nuxt/app'

export const useAuth = () => {
  const user = ref(null)
  const token = useCookie('auth-token')
  const isAuthenticated = computed(() => !!token.value)

  const login = async (email: string, password: string) => {
    try {
      // Implement login logic here
      const response = await $fetch('/api/auth/login', {
        method: 'POST',
        body: { email, password }
      })
      
      token.value = response.token
      user.value = response.user
      return { success: true }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  const register = async (userData: any) => {
    try {
      const response = await $fetch('/api/auth/register', {
        method: 'POST',
        body: userData
      })
      
      token.value = response.token
      user.value = response.user
      return { success: true }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    register,
    logout
  }
}
