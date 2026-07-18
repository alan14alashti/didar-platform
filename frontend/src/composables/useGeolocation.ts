import { onMounted, ref } from 'vue'

export interface GeoPoint {
  lat: number
  lng: number
}

export function useGeolocation() {
  const location = ref<GeoPoint | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  function requestLocation() {
    if (!navigator.geolocation) {
      error.value = 'unsupported'
      return
    }

    isLoading.value = true
    error.value = null

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        location.value = { lat: pos.coords.latitude, lng: pos.coords.longitude }
        isLoading.value = false
      },
      () => {
        error.value = 'denied'
        isLoading.value = false
      },
      { enableHighAccuracy: true, timeout: 12000, maximumAge: 300_000 },
    )
  }

  onMounted(requestLocation)

  return { location, isLoading, error, requestLocation }
}
