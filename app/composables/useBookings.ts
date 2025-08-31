import { ref, computed } from 'vue'
import { useQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const GET_USER_BOOKINGS = gql`
  query GetUserBookings($userId: Int!) {
    tickets(where: { user_id: { _eq: $userId } }, order_by: { created_at: desc }) {
      id
      user_id
      schedule_id
      seat_number
      price
      status
      payment_reference
      created_at
      updated_at
      schedule {
        id
        start_time
        end_time
        price
        movie {
          id
          title
          featured_image_url
        }
        cinema_hall {
          name
          cinema {
            name
            city
          }
        }
      }
    }
  }
`

const CREATE_BOOKING = gql`
  mutation CreateBooking($scheduleId: Int!, $seatNumber: String!, $price: numeric!) {
    insert_tickets_one(object: {
      schedule_id: $scheduleId,
      seat_number: $seatNumber,
      price: $price,
      status: "pending"
    }) {
      id
      schedule_id
      seat_number
      price
      status
    }
  }
`

export const useBookings = () => {
  const { result, loading, error, refetch } = useQuery(
    GET_USER_BOOKINGS,
    { userId: 1 } // This should come from auth context
  )

  const { mutate: createBooking, loading: creatingBooking } = useMutation(CREATE_BOOKING)

  const bookings = computed(() => {
    if (!result.value?.tickets) return []
    return result.value.tickets.map(ticket => ({
      ...ticket,
      movie: ticket.schedule.movie,
      cinema: ticket.schedule.cinema_hall.cinema,
      showtime: ticket.schedule.start_time
    }))
  })

  const bookTicket = async (scheduleId: number, seatNumber: string, price: number) => {
    try {
      const result = await createBooking({
        scheduleId,
        seatNumber,
        price
      })
      
      await refetch()
      return { success: true, booking: result.data.insert_tickets_one }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  return {
    bookings,
    loading,
    error,
    bookTicket,
    creatingBooking
  }
}
```

```

