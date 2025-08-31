import { useAuthStore } from '~/stores/auth'
import { useMoviesStore } from '../stores/movies'
  
  interface Movie {
    id: string
    title: string
    description: string
    poster: string
    rating: number
    duration: string
    year: string
    director: string
    genres: string[]
    showtimes: any[]
  }
  
  const props = defineProps<{
    movie: Movie
  }>();
  
  const authStore = useAuthStore()
  
  const isFavorite = computed(() => {
    return authStore.user?.favorites.includes(props.movie.id) ?? false
  })
  
  const toggleFavorite = () => {
    if (!authStore.isAuthenticated) {
      const authModal = inject('authModal')
      authModal?.show('login')
      return
    }
    
    authStore.toggleFavorite(props.movie.id)
  }
  
  const bookTickets = () => {
    if (!authStore.isAuthenticated) {
      const authModal = inject('authModal')
      authModal?.show('login')
      return
    }
    
    navigateTo(`/booking/${props.movie.id}`)
  }
  </script>
  
  <style>
  .line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  </style>
  </script>