import { useQuery, useMutation } from '@vue/apollo-composable'
import GET_MOVIES from '~/graphql/queries/movies.gql'
import GET_MOVIE_BY_ID from '~/graphql/queries/movies.gql'
import GET_FEATURED_MOVIES from '~/graphql/queries/movies.gql'
import GET_NOW_SHOWING from '~/graphql/queries/movies.gql'

export const useMovies = () => {
  const searchQuery = ref('')
  const selectedGenres = ref<string[]>([])
  const currentPage = ref(1)
  const itemsPerPage = 12

  // Computed variables for GraphQL
  const moviesWhere = computed(() => {
    const where: any = {}
    
    if (searchQuery.value) {
      where._or = [
        { title: { _ilike: `%${searchQuery.value}%` } },
        { director: { _ilike: `%${searchQuery.value}%` } },
        { genres: { _contains: searchQuery.value } }
      ]
    }
    
    if (selectedGenres.value.length > 0) {
      where.genres = { _contains: selectedGenres.value }
    }
    
    return where
  })

  // Queries
  const { result: moviesResult, loading: moviesLoading, refetch: refetchMovies } = useQuery(
    GET_MOVIES,
    () => ({
      limit: itemsPerPage,
      offset: (currentPage.value - 1) * itemsPerPage,
      where: moviesWhere.value
    })
  )

  const { result: featuredResult, loading: featuredLoading } = useQuery(GET_FEATURED_MOVIES)
  const { result: nowShowingResult, loading: nowShowingLoading } = useQuery(GET_NOW_SHOWING)

  // Computed properties
  const movies = computed(() => moviesResult.value?.movies || [])
  const featuredMovies = computed(() => featuredResult.value?.movies || [])
  const nowShowingMovies = computed(() => nowShowingResult.value?.movies || [])
  const featuredMovie = computed(() => featuredMovies.value[0] || null)

  // Methods
  const setSearchQuery = (query: string) => {
    searchQuery.value = query
    currentPage.value = 1
    refetchMovies()
  }

  const toggleGenreFilter = (genre: string) => {
    const index = selectedGenres.value.indexOf(genre)
    if (index > -1) {
      selectedGenres.value.splice(index, 1)
    } else {
      selectedGenres.value.push(genre)
    }
    currentPage.value = 1
    refetchMovies()
  }

  const clearFilters = () => {
    searchQuery.value = ''
    selectedGenres.value = []
    currentPage.value = 1
    refetchMovies()
  }

  const getMovieById = async (id: string) => {
    const { result } = await useQuery(GET_MOVIE_BY_ID, { id })
    return result.value?.movies_by_pk || null
  }

  return {
    // Data
    movies,
    featuredMovies,
    nowShowingMovies,
    featuredMovie,
    searchQuery: readonly(searchQuery),
    selectedGenres: readonly(selectedGenres),
    
    // Loading states
    moviesLoading,
    featuredLoading,
    nowShowingLoading,
    
    // Methods
    setSearchQuery,
    toggleGenreFilter,
    clearFilters,
    getMovieById,
    refetchMovies
  }
}