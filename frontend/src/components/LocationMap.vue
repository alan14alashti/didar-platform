<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import L from 'leaflet'
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png'
import markerIcon from 'leaflet/dist/images/marker-icon.png'
import markerShadow from 'leaflet/dist/images/marker-shadow.png'

// Fix Leaflet's default marker icon paths, which break under bundlers
// that fingerprint asset URLs (Vite included).
delete (L.Icon.Default.prototype as unknown as { _getIconUrl?: unknown })._getIconUrl
L.Icon.Default.mergeOptions({ iconRetinaUrl: markerIcon2x, iconUrl: markerIcon, shadowUrl: markerShadow })

const props = withDefaults(
  defineProps<{
    latitude: number
    longitude: number
    interactive?: boolean
    zoom?: number
  }>(),
  { interactive: false, zoom: 13 },
)

const emit = defineEmits<{ 'update:location': [{ lat: number; lng: number }] }>()

// Default to Tehran when no coordinates are set yet (host creating a new
// experience), so the picker starts somewhere sensible for Iranian users.
const DEFAULT_LAT = 35.6892
const DEFAULT_LNG = 51.389

const el = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null
let marker: L.Marker | null = null

function currentLatLng(): [number, number] {
  const lat = props.latitude || DEFAULT_LAT
  const lng = props.longitude || DEFAULT_LNG
  return [lat, lng]
}

onMounted(() => {
  if (!el.value) return
  const [lat, lng] = currentLatLng()

  map = L.map(el.value, { zoomControl: props.interactive }).setView([lat, lng], props.zoom)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map)

  marker = L.marker([lat, lng], { draggable: props.interactive }).addTo(map)

  if (props.interactive) {
    map.on('click', (e: L.LeafletMouseEvent) => {
      marker?.setLatLng(e.latlng)
      emit('update:location', { lat: e.latlng.lat, lng: e.latlng.lng })
    })
    marker.on('dragend', () => {
      const pos = marker!.getLatLng()
      emit('update:location', { lat: pos.lat, lng: pos.lng })
    })
  }
})

watch(
  () => [props.latitude, props.longitude],
  ([lat, lng]) => {
    if (!map || !marker || !lat || !lng) return
    marker.setLatLng([lat, lng])
    map.setView([lat, lng])
  },
)

onBeforeUnmount(() => {
  map?.remove()
  map = null
})
</script>

<template>
  <div ref="el" class="h-full w-full rounded-xl" />
</template>
