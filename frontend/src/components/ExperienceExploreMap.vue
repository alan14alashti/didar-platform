<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import L from 'leaflet'
import type { Experience } from '../api/types'
import type { GeoPoint } from '../composables/useGeolocation'
import { useFormatters } from '../composables/useFormatters'

const props = withDefaults(
  defineProps<{
    experiences: Experience[]
    userLocation?: GeoPoint | null
    activeSlug?: string | null
  }>(),
  { activeSlug: null },
)

const emit = defineEmits<{
  select: [slug: string]
}>()

const { formatPrice } = useFormatters()

const el = ref<HTMLDivElement | null>(null)
let map: L.Map | null = null
let markersLayer: L.LayerGroup | null = null
let userMarker: L.CircleMarker | null = null
let userPulse: L.Circle | null = null
const markerBySlug = new Map<string, L.Marker>()

const DEFAULT_CENTER: [number, number] = [32.4279, 53.6880]
const DEFAULT_ZOOM = 5

function markerHtml(exp: Experience, active: boolean) {
  const price = formatPrice(exp.price_amount).replace(/\s.+$/, '')
  return `<div class="exp-map-marker ${active ? 'exp-map-marker--active' : ''}">${price}</div>`
}

function syncMarkers() {
  if (!map || !markersLayer) return

  markersLayer.clearLayers()
  markerBySlug.clear()

  for (const exp of props.experiences) {
    if (!exp.latitude || !exp.longitude) continue
    const active = props.activeSlug === exp.slug
    const marker = L.marker([exp.latitude, exp.longitude], {
      icon: L.divIcon({
        className: 'exp-map-marker-wrap',
        html: markerHtml(exp, active),
        iconSize: [0, 0],
      }),
      zIndexOffset: active ? 1000 : 0,
    })
    marker.on('click', () => emit('select', exp.slug))
    marker.addTo(markersLayer)
    markerBySlug.set(exp.slug, marker)
  }

  syncUserMarker()
  fitVisibleBounds()
}

function syncUserMarker() {
  if (!map) return

  userMarker?.remove()
  userPulse?.remove()
  userMarker = null
  userPulse = null

  if (!props.userLocation) return

  const { lat, lng } = props.userLocation
  userPulse = L.circle([lat, lng], {
    radius: 1200,
    color: '#2563eb',
    weight: 1,
    opacity: 0.35,
    fillColor: '#3b82f6',
    fillOpacity: 0.12,
  }).addTo(map)

  userMarker = L.circleMarker([lat, lng], {
    radius: 8,
    color: '#ffffff',
    weight: 3,
    fillColor: '#2563eb',
    fillOpacity: 1,
  }).addTo(map)
}

function fitVisibleBounds() {
  if (!map) return

  const points: L.LatLngExpression[] = props.experiences
    .filter((e) => e.latitude && e.longitude)
    .map((e) => [e.latitude, e.longitude])

  if (props.userLocation) {
    points.push([props.userLocation.lat, props.userLocation.lng])
  }

  if (points.length === 0) {
    map.setView(DEFAULT_CENTER, DEFAULT_ZOOM)
    return
  }

  if (points.length === 1) {
    map.setView(points[0], 11)
    return
  }

  map.fitBounds(L.latLngBounds(points), { padding: [48, 48], maxZoom: 12 })
}

function panToSlug(slug: string | null) {
  if (!map || !slug) return
  const marker = markerBySlug.get(slug)
  if (!marker) return
  map.panTo(marker.getLatLng(), { animate: true })
}

onMounted(() => {
  if (!el.value) return

  map = L.map(el.value, { zoomControl: true }).setView(DEFAULT_CENTER, DEFAULT_ZOOM)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap',
    maxZoom: 19,
  }).addTo(map)

  markersLayer = L.layerGroup().addTo(map)
  syncMarkers()
})

watch(() => props.experiences, syncMarkers, { deep: true })
watch(() => props.userLocation, () => {
  syncUserMarker()
  fitVisibleBounds()
})
watch(() => props.activeSlug, (slug) => {
  syncMarkers()
  panToSlug(slug)
})

onBeforeUnmount(() => {
  map?.remove()
  map = null
  markersLayer = null
  markerBySlug.clear()
})
</script>

<template>
  <div ref="el" class="h-full w-full" />
</template>

<style>
.exp-map-marker-wrap {
  background: transparent;
  border: none;
}

.exp-map-marker {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
  border: 1.5px solid rgb(112 141 35 / 0.35);
  background: white;
  padding: 0.35rem 0.65rem;
  font-size: 0.72rem;
  font-weight: 800;
  color: #1f2937;
  box-shadow: 0 4px 14px rgb(15 23 42 / 0.12);
  transform: translate(-50%, -100%);
  white-space: nowrap;
  cursor: pointer;
  transition: transform 160ms ease, background-color 160ms ease, color 160ms ease;
}

.exp-map-marker:hover,
.exp-map-marker--active {
  background: var(--color-primary, #708d23);
  border-color: var(--color-primary, #708d23);
  color: white;
  transform: translate(-50%, -100%) scale(1.06);
  z-index: 1000;
}
</style>
