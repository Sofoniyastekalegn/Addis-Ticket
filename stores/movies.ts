
import { defineStore } from "pinia"

import { useMovies } from '~/composables/useMovies'

interface Movie {
    id: string
    title: string
    description: string
    poster: string
    backdrop: string
    rating: number
    duration: string
    year: string
    director: string
    genres: string[]
    cast: string[]
    trailer?: string
    showtimes: Showtime[]
    isFeatured: boolean

}

interface Showtime  {
    id: string
    movieId: string
    cinemaId: string
    cinemaName: string
    date: string
    time: string
    price: number
    availableSeats: number
    totalSeats: number
    screenType: string
}


interface Cinema {
    id: string
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
    featuredMovie: Movie |  null

}

export const useMoviesStore = defineStore("movies", {
    state: (): MoviesState => ({
        movies: [],
        cinemas: [],
        loading: false,
        searchQuery: "",
        selectedGenres: [],
        featuredMovie: null
        
    }),


    getters: {
        filteredMovies: (state) => state.movies,
        nowShowing: (state) => state.movies,

        allGenres: (state) => state.movies,

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
            const { movies, featuredMovies, nowShowingMovies, featuredMovie } = useMovies()
            this.movies = movies.value
            this.featuredMovie = featuredMovie.value
        },

        async fetchCinemas() {
            const { cinemas } = useCinemas()

            this.cinemas = cinemas.value
        },

        setSearchQuery(query: string) {
            const { setSearchQuery } = useMovies()
            setSearchQuery(query)
            this.searchQuery = query
        },
        
    }
})