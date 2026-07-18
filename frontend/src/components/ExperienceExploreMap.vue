<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import maplibregl from 'maplibre-gl'
import type { Map as MlMap, Marker as MlMarker } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'

import type { Experience } from '../api/types'
import type { GeoPoint } from '../composables/useGeolocation'
import { useFormatters } from '../composables/useFormatters'
import ExperienceMapPopup from './ExperienceMapPopup.vue'

const props = withDefaults(
  defineProps<{
    experiences: Experience[]
    userLocation?: GeoPoint | null
    selectedSlug?: string | null
  }>(),
  { selectedSlug: null },
)

const emit = defineEmits<{
  select: [slug: string]
  openDetail: [slug: string]
}>()

const { formatPrice } = useFormatters()

const el = ref<HTMLDivElement | null>(null)
const popupStyle = ref<Record<string, string>>({ display: 'none' })

let map: MlMap | null = null
const markerBySlug = new globalThis.Map<string, MlMarker>()
let userMarker: MlMarker | null = null

const DEFAULT_CENTER: [number, number] = [53.688, 32.428]
const DEFAULT_ZOOM = 5

const selectedExperience = computed(() =>
  props.experiences.find((e) => e.slug === props.selectedSlug) ?? null,
)

function priceLabel(amountRial: number) {
  return formatPrice(amountRial).replace(/\s.+$/, '')
}

function createMarkerElement(exp: Experience, active: boolean) {
  const btn = document.createElement('button')
  btn.type = 'button'
  btn.className = `exp-map-marker${active ? ' exp-map-marker--active' : ''}`
  btn.textContent = priceLabel(exp.price_amount)
  btn.addEventListener('click', (event) => {
    event.stopPropagation()
    emit('select', exp.slug)
  })
  return btn
}

function updatePopupPosition() {
  if (!map || !selectedExperience.value) {
    popupStyle.value = { display: 'none' }
    return
  }

  const { longitude, latitude } = selectedExperience.value
  const point = map.project([longitude, latitude])
  popupStyle.value = {
    display: 'block',
    left: `${point.x}px`,
    top: `${point.y - 12}px`,
  }
}

function syncMarkers() {
  if (!map) return

  for (const marker of markerBySlug.values()) marker.remove()
  markerBySlug.clear()

  for (const exp of props.experiences) {
    if (!exp.latitude || !exp.longitude) continue
    const element = createMarkerElement(exp, props.selectedSlug === exp.slug)
    const marker = new maplibregl.Marker({ element, anchor: 'bottom' })
      .setLngLat([exp.longitude, exp.latitude])
      .addTo(map)
    markerBySlug.set(exp.slug, marker)
  }

  syncUserMarker()
  fitVisibleBounds()
  updatePopupPosition()
}

function syncUserMarker() {
  userMarker?.remove()
  userMarker = null
  if (!map || !props.userLocation) return

  const dot = document.createElement('div')
  dot.className = 'exp-map-user-dot'
  dot.title = 'You'

  userMarker = new maplibregl.Marker({ element: dot, anchor: 'center' })
    .setLngLat([props.userLocation.lng, props.userLocation.lat])
    .addTo(map)
}

function fitVisibleBounds() {
  if (!map) return

  const coords: [number, number][] = props.experiences
    .filter((e) => e.latitude && e.longitude)
    .map((e) => [e.longitude, e.latitude])

  if (props.userLocation) {
    coords.push([props.userLocation.lng, props.userLocation.lat])
  }

  if (coords.length === 0) {
    map.easeTo({ center: DEFAULT_CENTER, zoom: DEFAULT_ZOOM })
    return
  }

  if (coords.length === 1) {
    map.easeTo({ center: coords[0], zoom: 11 })
    return
  }

  const bounds = coords.reduce(
    (b, coord) => b.extend(coord),
    new maplibregl.LngLatBounds(coords[0], coords[0]),
  )
  map.fitBounds(bounds, { padding: 60, maxZoom: 12, duration: 500 })
}

function panToSlug(slug: string | null) {
  if (!map || !slug) return
  const exp = props.experiences.find((e) => e.slug === slug)
  if (!exp) return
  map.easeTo({ center: [exp.longitude, exp.latitude], zoom: Math.max(map.getZoom(), 10), duration: 450 })
}

onMounted(() => {
  if (!el.value) return

  map = new maplibregl.Map({
    container: el.value,
    style: 'https://tiles.openfreemap.org/styles/liberty',
    center: DEFAULT_CENTER,
    zoom: DEFAULT_ZOOM,
    attributionControl: { compact: true },
  })

  map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'top-left')
  map.on('load', syncMarkers)
  map.on('move', updatePopupPosition)
  map.on('zoom', updatePopupPosition)
  map.on('click', () => emit('select', ''))
})

watch(() => props.experiences, syncMarkers, { deep: true })
watch(() => props.userLocation, () => {
  syncUserMarker()
  fitVisibleBounds()
})
watch(
  () => props.selectedSlug,
  (slug) => {
    syncMarkers()
    if (slug) panToSlug(slug)
    updatePopupPosition()
  },
)

onBeforeUnmount(() => {
  for (const marker of markerBySlug.values()) marker.remove()
  markerBySlug.clear()
  userMarker?.remove()
  map?.remove()
  map = null
})
</script>

<template>
  <div class="relative h-full w-full">
    <div ref="el" class="h-full w-full" />

    <ExperienceMapPopup
      v-if="selectedExperience"
      :experience="selectedExperience"
      :style="popupStyle"
      @close="emit('select', '')"
      @open="emit('openDetail', selectedExperience.slug)"
    />
  </div>
</template>

<style>
.maplibregl-ctrl-attrib {
  font-size: 10px;
}

.exp-map-marker {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
  border: 1.5px solid rgb(112 141 35 / 0.35);
  background: white;
  padding: 0.4rem 0.7rem;
  font-size: 0.72rem;
  font-weight: 800;
  color: #1f2937;
  box-shadow: 0 4px 14px rgb(15 23 42 / 0.12);
  white-space: nowrap;
  cursor: pointer;
  transition: transform 160ms ease, background-color 160ms ease, color 160ms ease;
}

.exp-map-marker:hover,
.exp-map-marker--active {
  background: var(--color-primary, #708d23);
  border-color: var(--color-primary, #708d23);
  color: white;
  transform: scale(1.06);
}

.exp-map-user-dot {
  width: 14px;
  height: 14px;
  border-radius: 9999px;
  background: #2563eb;
  border: 3px solid white;
  box-shadow: 0 0 0 6px rgb(37 99 235 / 0.2);
}
</style>
