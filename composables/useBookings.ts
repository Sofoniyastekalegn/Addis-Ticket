import { useQuery, useMutation } from '@vue/apollo-composable'
import GET_USER_BOOKINGS from '~/graphql/queries/bookings.gql'
import GET_BOOKING_BY_ID from '~/graphql/queries/bookings.gql'
import CREATE_BOOKING from '~/graphql/mutations/bookings.gql'
import UPDATE_BOOKING_STATUS from '~/graphql/mutations/bookings.gql'
import CANCEL_BOOKING from '~/graphql/mutations/bookings.gql'

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

export const useBookings = () => {
  const currentBooking = ref<BookingDetails | null>(null)
  const seats = ref<Seat[]>([])
  const loading = ref(false)

  // Mutations
  const { mutate: createBookingMutation } = useMutation(CREATE_BOOKING)
  const { mutate: updateBookingStatusMutation } = useMutation(UPDATE_BOOKING_STATUS)
  const { mutate: cancelBookingMutation } = useMutation(CANCEL_BOOKING)

  // Computed
  const selectedSeats = computed(() => seats.value.filter(seat => seat.status === 'selected'))
  const totalAmount = computed(() => {
    return selectedSeats.value.reduce((total, seat) => total + seat.price, 0)
  })

  // Methods
  const generateSeats = (rows: number = 10, seatsPerRow: number = 12) => {
    const newSeats: Seat[] = []
    const rowLabels = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    
    for (let r = 0; r < rows; r++) {
      for (let s = 1; s <= seatsPerRow; s++) {
        const seatType = r < 3 ? 'premium' : r < 7 ? 'standard' : 'vip'
        const basePrice = seatType === 'premium' ? 180 : seatType === 'vip' ? 250 : 120
        
        newSeats.push({
          id: `${rowLabels[r]}${s}`,
          row: rowLabels[r],
          number: s,
          type: seatType,
          status: Math.random() > 0.7 ? 'occupied' : 'available',
          price: basePrice
        })
      }
    }
    
    seats.value = newSeats
  }

  const selectSeat = (seatId: string) => {
    const seat = seats.value.find(s => s.id === seatId)
    if (seat && seat.status === 'available') {
      seat.status = 'selected'
    }
  }

  const deselectSeat = (seatId: string) => {
    const seat = seats.value.find(s => s.id === seatId)
    if (seat && seat.status === 'selected') {
      seat.status = 'available'
    }
  }

  const initializeBooking = async (movieId: string, showtimeId: string) => {
    loading.value = true
    try {
      generateSeats()
      
      currentBooking.value = {
        movieId,
        showtimeId,
        selectedSeats: [],
        totalAmount: 0,
        customerInfo: {
          name: '',
          email: '',
          phone: ''
        },
        paymentMethod: 'chapa'
      }
    } finally {
      loading.value = false
    }
  }

  const processPayment = async (paymentData: any) => {
    loading.value = true
    try {
      // Create booking in database
      const bookingResult = await createBookingMutation({
        input: {
          user_id: paymentData.userId,
          showtime_id: currentBooking.value?.showtimeId,
          seats: JSON.stringify(selectedSeats.value),
          total_amount: totalAmount.value,
          payment_method: currentBooking.value?.paymentMethod,
          booking_status: 'pending',
          payment_status: 'pending',
          booking_reference: Math.random().toString(36).substr(2, 9).toUpperCase()
        }
      })

      if (bookingResult?.data?.insert_bookings_one) {
        const booking = bookingResult.data.insert_bookings_one
        
        // Simulate payment processing
        await new Promise(resolve => setTimeout(resolve, 2000))
        
        // Update booking status
        await updateBookingStatusMutation({
          id: booking.id,
          status: 'confirmed',
          paymentStatus: 'completed'
        })
        
        // Clear current booking
        currentBooking.value = null
        seats.value = []
        
        return { success: true, booking }
      }
      
      return { success: false, message: 'Booking creation failed' }
    } catch (error) {
      return { success: false, message: 'Payment failed' }
    } finally {
      loading.value = false
    }
  }

  const getUserBookings = (userId: string) => {
    const { result, loading: bookingsLoading, refetch } = useQuery(
      GET_USER_BOOKINGS,
      { userId }
    )
    
    return {
      bookings: computed(() => result.value?.bookings || []),
      loading: bookingsLoading,
      refetch
    }
  }

  const getBookingById = (id: string) => {
    const { result, loading: bookingLoading } = useQuery(
      GET_BOOKING_BY_ID,
      { id }
    )
    
    return {
      booking: computed(() => result.value?.bookings_by_pk || null),
      loading: bookingLoading
    }
  }

  const cancelBooking = async (bookingId: string) => {
    try {
      await cancelBookingMutation({ id: bookingId })
      return { success: true }
    } catch (error) {
      return { success: false, message: 'Failed to cancel booking' }
    }
  }

  const clearBooking = () => {
    currentBooking.value = null
    seats.value = []
  }

  return {
    // State
    currentBooking: readonly(currentBooking),
    seats: readonly(seats),
    selectedSeats,
    totalAmount,
    loading: readonly(loading),
    
    // Methods
    generateSeats,
    selectSeat,
    deselectSeat,
    initializeBooking,
    processPayment,
    getUserBookings,
    getBookingById,
    cancelBooking,
    clearBooking
  }
}