import { useQuery } from '@vue/apollo-composable'
import GET_CINEMAS from '~/graphql/queries/cinemas.gql'
import GET_CINEMA_BY_ID from '~/graphql/queries/cinemas.gql'

export const useCinemas = () => {
  // Queries
  const { result: cinemasResult, loading: cinemasLoading, refetch: refetchCinemas } = useQuery(GET_CINEMAS)

  // Computed properties
  const cinemas = computed(() => cinemasResult.value?.cinemas || [])

  // Methods
  const getCinemaById = async (id: string) => {
    const { result } = await useQuery(GET_CINEMA_BY_ID, { id })
    return result.value?.cinemas_by_pk || null
  }

  return {
    // Data
    cinemas,
    
    // Loading states
    cinemasLoading,
    
    // Methods
    getCinemaById,
    refetchCinemas
  }
}