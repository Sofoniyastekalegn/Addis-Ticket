import { ref, computed } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

// GraphQL queries
const GET_MOVIES = gql`
  query GetMovies($limit: Int, $offset: Int) {
    movies(limit: $limit, offset: $offset, order_by: { created_at: desc }) {
      id
      title
      description
      duration
      director_id
      featured_image_url
      trailer_url
      release_date
      rating
      total_ratings
      created_at
      updated_at
      movie_images {
        id
        image_url
        is_featured
      }
      movie_genres {
        genre {
          id
          name
          description
        }
      }
      movie_stars {
        star {
          id
          name
          bio
          image_url
        }
      }
    }
  }
`

const GET_FEATURED_MOVIES = gql`
  query GetFeaturedMovies {
    movies(where: { movie_images: { is_featured: { _eq: true } } }, limit: 4) {
      id
      title
      description
      duration
      featured_image_url
      rating
      total_ratings
      release_date
      movie_images(where: { is_featured: { _eq: true } }, limit: 1) {
        image_url
      }
    }
  }
`

const GET_NOW_SHOWING = gql`
  query GetNowShowingMovies {
    movies(
      where: { 
        schedules: { 
          start_time: { _gte: "now()" } 
        } 
      }
      limit: 10
    ) {
      id
      title
      description
      duration
      featured_image_url
      rating
      total_ratings
      release_date
      schedules(where: { start_time: { _gte: "now()" } }, limit: 3) {
        id
        start_time
        end_time
        price
        available_seats
        cinema_hall {
          name
          cinema {
            name
            city
          }
        }
      }
      movie_images(limit: 1) {
        image_url
      }
    }
  }
`

export const useMovies = () => {
  const searchQuery = ref('')
  const selectedGenres = ref<string[]>([])

  // Fetch all movies
  const { result: moviesResult, loading: moviesLoading, error: moviesError } = useQuery(
    GET_MOVIES,
    { limit: 50, offset: 0 }
  )

  // Fetch featured movies
  const { result: featuredResult, loading: featuredLoading } = useQuery(GET_FEATURED_MOVIES)

  // Fetch now showing movies
  const { result: nowShowingResult, loading: nowShowingLoading } = useQuery(GET_NOW_SHOWING)

  // Computed properties
  const movies = computed(() => {
    if (!moviesResult.value?.movies) return []
    return moviesResult.value.movies.map(movie => ({
      ...movie,
      poster: movie.movie_images?.[0]?.image_url || movie.featured_image_url,
      backdrop: movie.featured_image_url,
      genres: movie.movie_genres?.map(mg => mg.genre.name) || [],
      cast: movie.movie_stars?.map(ms => ms.star.name) || [],
      year: movie.release_date ? new Date(movie.release_date).getFullYear().toString() : '',
      isFeatured: movie.movie_images?.some(mi => mi.is_featured) || false
    }))
  })

  const featuredMovies = computed(() => {
    if (!featuredResult.value?.movies) return []
    return featuredResult.value.movies.map(movie => ({
      ...movie,
      poster: movie.movie_images?.[0]?.image_url || movie.featured_image_url,
      backdrop: movie.featured_image_url,
      year: movie.release_date ? new Date(movie.release_date).getFullYear().toString() : '',
      isFeatured: true
    }))
  })

  const nowShowingMovies = computed(() => {
    if (!nowShowingResult.value?.movies) return []
    return nowShowingResult.value.movies.map(movie => ({
      ...movie,
      poster: movie.movie_images?.[0]?.image_url || movie.featured_image_url,
      backdrop: movie.featured_image_url,
      year: movie.release_date ? new Date(movie.release_date).getFullYear().toString() : '',
      isFeatured: false
    }))
  })

  const filteredMovies = computed(() => {
    let filtered = movies.value

    if (searchQuery.value) {
      filtered = filtered.filter(movie => 
        movie.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        movie.description.toLowerCase().includes(searchQuery.value.toLowerCase())
      )
    }

    if (selectedGenres.value.length > 0) {
      filtered = filtered.filter(movie => 
        movie.genres.some(genre => selectedGenres.value.includes(genre))
      )
    }

    return filtered
  })

  const allGenres = computed(() => {
    const genres = new Set<string>()
    movies.value.forEach(movie => {
      movie.genres.forEach(genre => genres.add(genre))
    })
    return Array.from(genres).sort()
  })

  const setSearchQuery = (query: string) => {
    searchQuery.value = query
  }

  const toggleGenreFilter = (genre: string) => {
    const index = selectedGenres.value.indexOf(genre)
    if (index > -1) {
      selectedGenres.value.splice(index, 1)
    } else {
      selectedGenres.value.push(genre)
    }
  }

  const clearFilters = () => {
    searchQuery.value = ''
    selectedGenres.value = []
  }

  return {
    movies,
    featuredMovies,
    nowShowingMovies,
    filteredMovies,
    allGenres,
    searchQuery,
    selectedGenres,
    moviesLoading,
    featuredLoading,
    nowShowingLoading,
    moviesError,
    setSearchQuery,
    toggleGenreFilter,
    clearFilters
  }
}
```

