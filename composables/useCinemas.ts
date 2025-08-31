import { ref, computed } from 'vue'

interface Cinema {
  id: number
  name: string
  location: string
  address: string
  phone: string
  facilities: string[]
  screens: number
  image: string
}

export const useCinemas = () => {
  const cinemas = ref<Cinema[]>([])
  const loading = ref(false)

  const fetchCinemas = async () => {
    loading.value = true
    try {
      // Mock data for now - replace with actual API call
      cinemas.value = [
        {
          id: 1,
          name: "CineMax Downtown",
          location: "Downtown",
          address: "123 Main Street, Downtown",
          phone: "+1-555-0101",
          facilities: ["Dolby Atmos", "Recliner Seats", "Food Service"],
          screens: 8,
          image: "/images/cinema-placeholder.jpg"
        }
      ]
    } catch (error) {
      console.error('Error fetching cinemas:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    cinemas: readonly(cinemas),
    loading: readonly(loading),
    fetchCinemas
  }
}
