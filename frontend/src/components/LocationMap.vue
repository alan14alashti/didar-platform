<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import maplibregl from 'maplibre-gl'
import type { Map as MlMap, Marker as MlMarker } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'

import { osmMapStyle } from '../lib/exploreMapStyle'

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

const DEFAULT_LAT = 35.6892
const DEFAULT_LNG = 51.389

const el = ref<HTMLDivElement | null>(null)

let map: MlMap | null = null
let marker: MlMarker | null = null
let mapReady = false
let bootAttempts = 0
let resizeObserver: ResizeObserver | null = null

function isValidCoord(lat: number, lng: number) {
  return Number.isFinite(lat) && Number.isFinite(lng) && Math.abs(lat) <= 90 && Math.abs(lng) <= 180 && (lat !== 0 || lng !== 0)
}

function currentCenter(): [number, number] {
  if (isValidCoord(props.latitude, props.longitude)) {
    return [props.longitude, props.latitude]
  }
  return [DEFAULT_LNG, DEFAULT_LAT]
}

function containerReady() {
  return Boolean(el.value && el.value.clientWidth > 0 && el.value.clientHeight > 0)
}

function ensureMapSize() {
  if (!map || !containerReady()) return
  map.resize()
}

function createMarkerElement() {
  const pin = document.createElement('div')
  pin.className = 'location-map-pin'
  return pin
}

function syncMarker(centerMap = false) {
  if (!map || !mapReady) return

  const [lng, lat] = currentCenter()

  if (!marker) {
    marker = new maplibregl.Marker({
      element: createMarkerElement(),
      draggable: props.interactive,
      anchor: 'bottom',
    })
      .setLngLat([lng, lat])
      .addTo(map)

    if (props.interactive) {
      marker.on('dragend', () => {
        const pos = marker!.getLngLat()
        emit('update:location', { lat: pos.lat, lng: pos.lng })
      })
    }
  } else {
    marker.setLngLat([lng, lat])
  }

  if (centerMap) {
    map.jumpTo({ center: [lng, lat], zoom: props.zoom })
  }
}

function initMap() {
  if (map || !el.value || !containerReady()) return false

  map = new maplibregl.Map({
    container: el.value,
    style: osmMapStyle,
    center: currentCenter(),
    zoom: props.zoom,
    attributionControl: { compact: true },
  })

  if (props.interactive) {
    map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'top-right')
    map.on('click', (event) => {
      marker?.setLngLat(event.lngLat)
      emit('update:location', { lat: event.lngLat.lat, lng: event.lngLat.lng })
    })
  }

  map.on('load', () => {
    mapReady = true
    syncMarker(true)
    ensureMapSize()
  })

  map.on('idle', ensureMapSize)

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

watch(
  () => [props.latitude, props.longitude],
  () => {
    if (!isValidCoord(props.latitude, props.longitude)) return
    syncMarker(true)
  },
)

onBeforeUnmount(() => {
  mapReady = false
  resizeObserver?.disconnect()
  resizeObserver = null
  marker?.remove()
  marker = null
  map?.remove()
  map = null
})
</script>

<template>
  <div class="location-map-shell h-full w-full">
    <div ref="el" class="location-map-canvas h-full w-full" />
  </div>
</template>

<style>
.location-map-shell,
.location-map-canvas,
.location-map-canvas .maplibregl-map {
  height: 100%;
  width: 100%;
}

.location-map-pin {
  width: 14px;
  height: 14px;
  border-radius: 9999px;
  background: var(--color-primary, #708d23);
  border: 3px solid white;
  box-shadow: 0 2px 6px rgb(15 23 42 / 0.25);
}

.maplibregl-ctrl-attrib {
  font-size: 10px;
}
</style>
