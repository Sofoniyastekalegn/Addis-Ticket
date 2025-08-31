import { defineStore } from 'pinia'

interface Movie {
  id: number
  title: string
  description: string
  poster: string
  backdrop: string
  rating: number
  duration: number
  year: string
  director: string
  genres: string[]
  cast: string[]
  trailer?: string
  schedules: Schedule[]
  isFeatured: boolean
}

interface Schedule {
  id: number
  movie_id: number
  cinema_hall_id: number
  start_time: string
  end_time: string
  price: number
  available_seats: number
  cinema_hall: {
    name: string
    cinema: {
      name: string
      city: string
    }
  }
}

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

interface MoviesState {
  movies: Movie[]
  cinemas: Cinema[]
  loading: boolean
  searchQuery: string
  selectedGenres: string[]
  featuredMovie: Movie | null
}

export const useMoviesStore = defineStore('movies', {
  state: (): MoviesState => ({
    movies: [],
    cinemas: [],
    loading: false,
    searchQuery: '',
    selectedGenres: [],
    featuredMovie: null
  }),

  getters: {
    filteredMovies: (state) => {
      let filtered = state.movies

      if (state.searchQuery) {
        filtered = filtered.filter(movie => 
          movie.title.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
          movie.description.toLowerCase().includes(state.searchQuery.toLowerCase())
        )
      }

      if (state.selectedGenres.length > 0) {
        filtered = filtered.filter(movie => 
          movie.genres.some(genre => state.selectedGenres.includes(genre))
        )
      }

      return filtered
    },

    nowShowing: (state) => state.movies.filter(movie => movie.schedules.length > 0),

    allGenres: (state) => {
      const genres = new Set<string>()
      state.movies.forEach(movie => {
        movie.genres.forEach(genre => genres.add(genre))
      })
      return Array.from(genres).sort()
    }
  },

  actions: {
    async fetchMovies() {
      this.loading = true
      try {
        // Mock data for now - replace with actual API call
        this.movies = [
          {
            id: 1,
            title: "The Lion King",
            description: "A young lion prince flees his kingdom only to learn the true meaning of responsibility and bravery.",
            poster: "/images/movie-placeholder.jpg",
            backdrop: "/images/movie-placeholder.jpg",
            rating: 4.5,
            duration: 118,
            year: "2019",
            director: "Jon Favreau",
            genres: ["Animation", "Adventure", "Drama"],
            cast: ["Donald Glover", "Beyoncé", "James Earl Jones"],
            schedules: [],
            isFeatured: true
          },
          {
            id: 2,
            title: "Avengers: Endgame",
            description: "After the devastating events of Avengers: Infinity War, the universe is in ruins.",
            poster: "/images/movie-placeholder.jpg",
            backdrop: "/images/movie-placeholder.jpg",
            rating: 4.8,
            duration: 181,
            year: "2019",
            director: "Anthony Russo",
            genres: ["Action", "Adventure", "Drama"],
            cast: ["Robert Downey Jr.", "Chris Evans", "Mark Ruffalo"],
            schedules: [],
            isFeatured: true
          },
          {
            id: 3,
            title: "Black Panther",
            description: "T'Challa, heir to the hidden but advanced kingdom of Wakanda, must step forward to lead his people into a new future.",
            poster: "/images/movie-placeholder.jpg",
            backdrop: "/images/movie-placeholder.jpg",
            rating: 4.7,
            duration: 134,
            year: "2018",
            director: "Ryan Coogler",
            genres: ["Action", "Adventure", "Sci-Fi"],
            cast: ["Chadwick Boseman", "Michael B. Jordan", "Lupita Nyong'o"],
            schedules: [],
            isFeatured: true
          },
          {
            id: 4,
            title: "Wonder Woman",
            description: "When a pilot crashes and tells of conflict in the outside world, Diana, an Amazonian warrior in training, leaves home to fight a war.",
            poster: "/images/movie-placeholder.jpg",
            backdrop: "/images/movie-placeholder.jpg",
            rating: 4.6,
            duration: 141,
            year: "2017",
            director: "Patty Jenkins",
            genres: ["Action", "Adventure", "Fantasy"],
            cast: ["Gal Gadot", "Chris Pine", "Robin Wright"],
            schedules: [],
            isFeatured: true
          }
        ]
        this.featuredMovie = this.movies[0]
      } catch (error) {
        console.error('Error fetching movies:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchCinemas() {
      try {
        // Mock data for now - replace with actual API call
        this.cinemas = [
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
      }
    },

    setSearchQuery(query: string) {
      this.searchQuery = query
    },

    toggleGenreFilter(genre: string) {
      const index = this.selectedGenres.indexOf(genre)
      if (index > -1) {
        this.selectedGenres.splice(index, 1)
      } else {
        this.selectedGenres.push(genre)
      }
    },

    clearFilters() {
      this.searchQuery = ''
      this.selectedGenres = []
    }
  }
})