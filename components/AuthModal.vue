<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-8 max-w-md w-full mx-4">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900">
          {{ isLogin ? 'Sign In' : 'Sign Up' }}
        </h2>
        <button @click="hide" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div v-if="!isLogin">
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            v-model="form.email"
            type="email"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            v-model="form.password"
            type="password"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div v-if="error" class="text-red-600 text-sm">{{ error }}</div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 disabled:opacity-50"
        >
          {{ loading ? 'Loading...' : (isLogin ? 'Sign In' : 'Sign Up') }}
        </button>
      </form>

      <div class="mt-4 text-center">
        <button @click="toggleMode" class="text-primary-600 hover:text-primary-700">
          {{ isLogin ? "Don't have an account? Sign Up" : 'Already have an account? Sign In' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()

const isOpen = ref(false)
const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  email: '',
  password: ''
})

const show = (mode = 'login') => {
  isLogin.value = mode === 'login'
  isOpen.value = true
  error.value = ''
  resetForm()
}

const hide = () => {
  isOpen.value = false
  resetForm()
}

const toggleMode = () => {
  isLogin.value = !isLogin.value
  error.value = ''
  resetForm()
}

const resetForm = () => {
  form.name = ''
  form.email = ''
  form.password = ''
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    if (isLogin.value) {
      const result = await authStore.login(form.email, form.password)
      if (result.success) {
        hide()
      } else {
        error.value = result.error || 'Login failed'
      }
    } else {
      const result = await authStore.register(form.email, form.password, form.name)
      if (result.success) {
        // Auto-login after successful registration
        const loginResult = await authStore.login(form.email, form.password)
        if (loginResult.success) {
          hide()
        } else {
          error.value = 'Registration successful but login failed'
        }
      } else {
        error.value = result.error || 'Registration failed'
      }
    }
  } catch (err) {
    error.value = 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}

defineExpose({
  show,
  hide
})
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
```

```

