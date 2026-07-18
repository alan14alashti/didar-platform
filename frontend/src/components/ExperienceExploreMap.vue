<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import maplibregl from 'maplibre-gl'
import type { Map as MlMap, Marker as MlMarker } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'

import type { Experience } from '../api/types'
import type { GeoPoint } from '../composables/useGeolocation'
import { osmMapStyle } from '../lib/exploreMapStyle'
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
let mapReady = false
let bootAttempts = 0
let resizeObserver: ResizeObserver | null = null
const markerBySlug = new globalThis.Map<string, MlMarker>()
let userMarker: MlMarker | null = null

const DEFAULT_CENTER: [number, number] = [53.688, 32.428]
const DEFAULT_ZOOM = 5

function isValidCoord(lat: number, lng: number) {
  return Number.isFinite(lat) && Number.isFinite(lng) && Math.abs(lat) <= 90 && Math.abs(lng) <= 180
}

function containerReady() {
  return Boolean(el.value && el.value.clientWidth > 0 && el.value.clientHeight > 0)
}

function ensureMapSize() {
  if (!map || !containerReady()) return
  map.resize()
}

function initMap() {
  if (map || !el.value || !containerReady()) return false

  map = new maplibregl.Map({
    container: el.value,
    style: osmMapStyle,
    center: DEFAULT_CENTER,
    zoom: DEFAULT_ZOOM,
    attributionControl: { compact: true },
  })

  map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'top-left')

  map.on('load', () => {
    mapReady = true
    ensureMapSize()
    syncMarkers(true)
  })

  map.on('idle', ensureMapSize)
  map.on('move', updatePopupPosition)
  map.on('zoom', updatePopupPosition)
  map.on('click', () => emit('select', ''))
  map.on('error', (event) => {
    console.error('MapLibre error:', event.error)
  })

  return true
}

function bootMap() {
  if (initMap()) {
    ensureMapSize()
    return
  }

  if (bootAttempts >= 30 || !el.value) return
  bootAttempts += 1
  requestAnimationFrame(bootMap)
}

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
  if (!isValidCoord(latitude, longitude)) {
    popupStyle.value = { display: 'none' }
    return
  }

  const point = map.project([longitude, latitude])
  popupStyle.value = {
    display: 'block',
    left: `${point.x}px`,
    top: `${point.y - 12}px`,
  }
}

function collectCoords(): [number, number][] {
  const coords: [number, number][] = props.experiences
    .filter((e) => isValidCoord(e.latitude, e.longitude))
    .map((e) => [e.longitude, e.latitude])

  if (props.userLocation && isValidCoord(props.userLocation.lat, props.userLocation.lng)) {
    coords.push([props.userLocation.lng, props.userLocation.lat])
  }

  return coords
}

function syncMarkers(fitBounds = false) {
  if (!map || !mapReady) return

  for (const marker of markerBySlug.values()) marker.remove()
  markerBySlug.clear()

  for (const exp of props.experiences) {
    if (!isValidCoord(exp.latitude, exp.longitude)) continue
    const element = createMarkerElement(exp, props.selectedSlug === exp.slug)
    const marker = new maplibregl.Marker({ element, anchor: 'bottom' })
      .setLngLat([exp.longitude, exp.latitude])
      .addTo(map)
    markerBySlug.set(exp.slug, marker)
  }

  syncUserMarker()
  if (fitBounds) fitVisibleBounds()
  updatePopupPosition()
}

function syncUserMarker() {
  userMarker?.remove()
  userMarker = null
  if (!map || !props.userLocation) return
  if (!isValidCoord(props.userLocation.lat, props.userLocation.lng)) return

  const dot = document.createElement('div')
  dot.className = 'exp-map-user-dot'
  dot.title = 'You'

  userMarker = new maplibregl.Marker({ element: dot, anchor: 'center' })
    .setLngLat([props.userLocation.lng, props.userLocation.lat])
    .addTo(map)
}

function fitVisibleBounds() {
  if (!map || !mapReady) return

  const coords = collectCoords()

  if (coords.length === 0) {
    map.jumpTo({ center: DEFAULT_CENTER, zoom: DEFAULT_ZOOM })
    return
  }

  if (coords.length === 1) {
    map.jumpTo({ center: coords[0], zoom: 11 })
    return
  }

  const bounds = coords.reduce(
    (b, coord) => b.extend(coord),
    new maplibregl.LngLatBounds(coords[0], coords[0]),
  )
  map.fitBounds(bounds, { padding: 60, maxZoom: 12, duration: 0 })
}

function panToSlug(slug: string | null) {
  if (!map || !mapReady || !slug) return
  const exp = props.experiences.find((e) => e.slug === slug)
  if (!exp || !isValidCoord(exp.latitude, exp.longitude)) return
  map.easeTo({ center: [exp.longitude, exp.latitude], zoom: Math.max(map.getZoom(), 10), duration: 450 })
}

async function refreshMap() {
  bootAttempts = 0
  bootMap()
  await nextTick()
  for (const delay of [0, 100, 300]) {
    await new Promise((resolve) => setTimeout(resolve, delay))
    ensureMapSize()
  }
  if (mapReady) syncMarkers(true)
}

defineExpose({ refreshMap })

onMounted(async () => {
  await nextTick()
  if (!el.value) return

  resizeObserver = new ResizeObserver(() => {
    if (!map) bootMap()
    ensureMapSize()
  })
  resizeObserver.observe(el.value)
  bootMap()
})

watch(() => props.experiences, () => syncMarkers(true), { deep: true })
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
  mapReady = false
  resizeObserver?.disconnect()
  resizeObserver = null
  for (const marker of markerBySlug.values()) marker.remove()
  markerBySlug.clear()
  userMarker?.remove()
  map?.remove()
  map = null
})
</script>

<template>
  <div class="exp-map-shell relative h-full w-full min-h-[20rem]">
    <div ref="el" class="exp-map-canvas absolute inset-0 h-full w-full" />

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
.exp-map-shell,
.exp-map-canvas,
.exp-map-canvas .maplibregl-map {
  height: 100%;
  width: 100%;
}

.maplibregl-ctrl-attrib {
  font-size: 10px;
}

.exp-map-marker {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
  border: 1px solid rgb(221 221 221);
  background: white;
  padding: 0.35rem 0.65rem;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #222222;
  box-shadow: 0 2px 4px rgb(0 0 0 / 0.18);
  white-space: nowrap;
  cursor: pointer;
  transition: transform 160ms ease, background-color 160ms ease, color 160ms ease, border-color 160ms ease;
}

.exp-map-marker:hover,
.exp-map-marker--active {
  background: #222222;
  border-color: #222222;
  color: white;
  transform: scale(1.04);
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
