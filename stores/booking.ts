import { defineStore } from 'pinia'
import { useBookings } from '~/composables/useBookings'

interface Seat {
  id: string
  row: string
  number: number
  type: 'standard' | 'premium' | 'vip'
  status: 'available' | 'occupied' | 'selected'
  price: number
}

interface BookingDetails {
  movieId: string
  showtimeId: string
  selectedSeats: Seat[]
  totalAmount: number
  customerInfo: {
    name: string
    email: string
    phone: string
  }
  paymentMethod: 'chapa' | 'telebirr' | 'mpesa'
}

interface BookingState {
  currentBooking: BookingDetails | null
  seats: Seat[]
  loading: boolean
  bookingHistory: any[]
}

export const useBookingStore = defineStore('booking', {
  state: (): BookingState => ({
    currentBooking: null,
    seats: [],
    loading: false,
    bookingHistory: []
  }),

  getters: {
    selectedSeats: (state) => {
      const { selectedSeats } = useBookings()
      return selectedSeats.value
    },
    totalAmount: (state) => {
      const { totalAmount } = useBookings()
      return totalAmount.value
    }
  },

  actions: {
    generateSeats(rows: number = 10, seatsPerRow: number = 12) {
      const { generateSeats, seats } = useBookings()
      generateSeats(rows, seatsPerRow)
      this.seats = seats.value
    },

    selectSeat(seatId: string) {
      const { selectSeat } = useBookings()
      selectSeat(seatId)
    },

    deselectSeat(seatId: string) {
      const { deselectSeat } = useBookings()
      deselectSeat(seatId)
    },

    async initializeBooking(movieId: string, showtimeId: string) {
      const { initializeBooking, currentBooking } = useBookings()
      await initializeBooking(movieId, showtimeId)
      this.currentBooking = currentBooking.value
    },

    async processPayment(paymentData: any) {
      const { processPayment } = useBookings()
      const result = await processPayment(paymentData)
      
      if (result.success) {
        this.bookingHistory.push(result.booking)
        this.currentBooking = null
        this.seats = []
      }
      
      return result
    },

    clearBooking() {
      const { clearBooking } = useBookings()
      clearBooking()
      this.currentBooking = null
      this.seats = []
    }
  }
})