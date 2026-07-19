<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import { bookingApi, experienceApi } from '../api/endpoints'
import { extractErrorMessage } from '../api/client'
import type { Experience, Photo, Slot } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import { useAuthStore } from '../stores/auth'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import EmptyState from '../components/EmptyState.vue'
import AppIcon from '../components/AppIcon.vue'
import LocationMap from '../components/LocationMap.vue'

const props = defineProps<{ slug: string }>()
const { t } = useI18n()
const { formatPrice, formatDateTime } = useFormatters()
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const experience = ref<Experience | null>(null)
const isLoading = ref(true)
const notFound = ref(false)

const selectedSlot = ref<Slot | null>(null)
const guestsCount = ref(1)
const guestNote = ref('')
const isBooking = ref(false)
const bookingError = ref('')
const bookingSuccess = ref(false)

const lightboxPhoto = ref<Photo | null>(null)

async function load() {
  isLoading.value = true
  notFound.value = false
  try {
    experience.value = await experienceApi.getBySlug(props.slug)
    if (experience.value.slots?.length) {
      selectedSlot.value = experience.value.slots[0]
    }
  } catch {
    notFound.value = true
  } finally {
    isLoading.value = false
  }
}

const totalPrice = computed(() => (experience.value ? experience.value.price_amount * guestsCount.value : 0))

const maxGuestsForSlot = computed(() => {
  if (!selectedSlot.value || !experience.value) return experience.value?.max_guests ?? 10
  return Math.min(selectedSlot.value.capacity - selectedSlot.value.booked_count, experience.value.max_guests)
})

const sortedPhotos = computed(() =>
  [...(experience.value?.photos ?? [])].sort((a, b) => a.sort_order - b.sort_order),
)

const coverPhoto = computed(() => {
  const photos = sortedPhotos.value
  if (!photos.length) return null
  return photos.find((photo) => photo.is_cover) ?? photos[0]
})

const galleryPhotos = computed(() => sortedPhotos.value.filter((photo) => !photo.is_cover))

async function submitBooking() {
  if (!authStore.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: route.fullPath } })
    return
  }
  if (!selectedSlot.value) return

  isBooking.value = true
  bookingError.value = ''
  try {
    await bookingApi.create({
      slot_id: selectedSlot.value.id,
      guests_count: guestsCount.value,
      guest_note: guestNote.value || undefined,
    })
    bookingSuccess.value = true
  } catch (err) {
    bookingError.value = extractErrorMessage(err, t('common.errorGeneric'))
  } finally {
    isBooking.value = false
  }
}

onMounted(load)
</script>

<template>
  <LoadingSpinner v-if="isLoading" :label="t('common.loading')" />

  <div v-else-if="notFound" class="mx-auto max-w-3xl px-4 py-12">
    <EmptyState :title="t('experience.notFound')" />
  </div>

  <div v-else-if="experience" class="mx-auto max-w-5xl px-4 py-5 sm:px-6 sm:py-6">
    <div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_18rem] lg:gap-8 xl:grid-cols-[minmax(0,1fr)_20rem]">
      <!-- Main -->
      <div class="min-w-0 space-y-5">
        <div class="overflow-hidden rounded-xl bg-[#f7f7f7]">
          <div class="relative aspect-[16/10] w-full sm:aspect-[5/3]">
            <img
              v-if="coverPhoto"
              :src="coverPhoto.url"
              :alt="coverPhoto.alt || experience.title"
              class="h-full w-full object-cover"
            />
            <span
              v-if="coverPhoto"
              class="absolute start-3 top-3 rounded-md bg-black/55 px-2 py-1 text-[11px] font-medium text-white"
            >
              {{ t('experience.coverPhoto') }}
            </span>
          </div>
          <div
            v-if="galleryPhotos.length"
            class="grid grid-cols-4 gap-1.5 border-t border-gray-100 p-2 sm:grid-cols-4"
          >
            <button
              v-for="photo in galleryPhotos"
              :key="photo.id"
              type="button"
              class="group relative aspect-[4/3] overflow-hidden rounded-md bg-gray-100"
              :aria-label="t('experience.viewPhoto')"
              @click="lightboxPhoto = photo"
            >
              <img :src="photo.url" :alt="photo.alt || experience.title" class="h-full w-full object-cover transition group-hover:scale-105" />
            </button>
          </div>
        </div>

        <div
          v-if="lightboxPhoto"
          class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4"
          @click.self="lightboxPhoto = null"
        >
          <button
            type="button"
            class="absolute end-4 top-4 grid h-9 w-9 place-items-center rounded-full bg-white/15 text-white"
            @click="lightboxPhoto = null"
          >
            <AppIcon icon="heroicons-outline:x-mark" size="md" />
          </button>
          <img
            :src="lightboxPhoto.url"
            :alt="lightboxPhoto.alt || experience.title"
            class="max-h-[85vh] max-w-full rounded-lg object-contain"
          />
        </div>

        <div class="space-y-2">
          <div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-sm text-[#717171]">
            <span>{{ experience.city }}</span>
            <span v-if="experience.rating_count > 0" class="inline-flex items-center gap-0.5 text-[#222222]">
              ·
              <AppIcon icon="heroicons-outline:star" size="xs" />
              {{ experience.rating_avg.toFixed(2) }}
              <span class="text-[#717171]">({{ experience.rating_count }})</span>
            </span>
          </div>

          <h1 class="text-xl font-semibold leading-snug text-[#222222] sm:text-2xl">{{ experience.title }}</h1>
          <p class="text-sm leading-6 text-[#717171]">{{ experience.summary }}</p>
        </div>

        <div class="flex flex-wrap gap-x-4 gap-y-1 border-y border-gray-100 py-3 text-sm text-[#717171]">
          <span class="inline-flex items-center gap-1">
            <AppIcon icon="heroicons-outline:clock" size="xs" />
            {{ t('experience.durationMinutes', { n: experience.duration_minutes }) }}
          </span>
          <span class="inline-flex items-center gap-1">
            <AppIcon icon="heroicons-outline:users" size="xs" />
            {{ t('experience.guests', { n: experience.min_guests, m: experience.max_guests }) }}
          </span>
        </div>

        <section v-if="experience.description" class="space-y-2">
          <h2 class="text-sm font-semibold text-[#222222]">{{ t('experience.about') }}</h2>
          <p class="whitespace-pre-line text-sm leading-6 text-[#717171]">{{ experience.description }}</p>
        </section>

        <section v-if="experience.host" class="space-y-2 border-t border-gray-100 pt-4">
          <h2 class="text-sm font-semibold text-[#222222]">{{ t('experience.aboutHost') }}</h2>
          <div class="flex items-start gap-3">
            <img
              v-if="experience.host.avatar_url"
              :src="experience.host.avatar_url"
              class="h-10 w-10 rounded-full object-cover"
              alt=""
            />
            <div v-else class="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-[#222222] text-sm font-semibold text-white">
              {{ (experience.host.display_name || experience.host.full_name || '?')[0] }}
            </div>
            <div class="min-w-0">
              <p class="flex items-center gap-1 text-sm font-medium text-[#222222]">
                {{ experience.host.display_name || experience.host.full_name }}
                <AppIcon
                  v-if="experience.host.is_verified"
                  icon="heroicons-outline:check-badge"
                  size="xs"
                  icon-class="text-brand-primary"
                />
              </p>
              <p v-if="experience.host.bio" class="mt-1 text-sm leading-5 text-[#717171]">{{ experience.host.bio }}</p>
            </div>
          </div>
        </section>

        <section class="space-y-2 border-t border-gray-100 pt-4">
          <h2 class="text-sm font-semibold text-[#222222]">{{ t('experience.location') }}</h2>
          <p class="text-sm text-[#717171]">{{ experience.address || experience.city }}</p>
          <div class="h-44 overflow-hidden rounded-xl border border-gray-100">
            <LocationMap :latitude="experience.latitude" :longitude="experience.longitude" />
          </div>
        </section>
      </div>

      <!-- Booking -->
      <aside class="lg:col-span-1">
        <div class="rounded-xl border border-gray-200 bg-white p-4 shadow-[0_2px_8px_rgba(0,0,0,0.06)] lg:sticky lg:top-20">
          <p class="flex flex-wrap items-baseline gap-x-1 text-lg font-semibold text-[#222222] tabular-nums">
            <span>{{ formatPrice(experience.price_amount) }}</span>
            <span class="text-sm font-normal text-[#717171]">{{ t('experience.perPerson') }}</span>
          </p>

          <div v-if="bookingSuccess" class="mt-4 rounded-lg bg-green-50 p-4 text-center">
            <AppIcon icon="heroicons-outline:check-circle" size="lg" icon-class="mx-auto text-green-600" />
            <p class="mt-2 text-sm font-semibold text-green-700">{{ t('booking.bookingSuccessTitle') }}</p>
            <p class="mt-1 text-xs text-green-600">{{ t('booking.bookingSuccessBody') }}</p>
            <router-link to="/trips" class="btn-primary mt-3 w-full !py-2 text-sm !shadow-none">
              {{ t('booking.goToMyTrips') }}
            </router-link>
          </div>

          <template v-else>
            <div v-if="!experience.slots?.length" class="mt-4">
              <EmptyState :title="t('experience.noSlots')" />
            </div>

            <div v-else class="mt-4 space-y-3">
              <div>
                <label class="mb-1 block text-xs font-medium text-[#717171]">{{ t('experience.selectTime') }}</label>
                <div class="max-h-36 space-y-1.5 overflow-y-auto scrollbar-thin">
                  <button
                    v-for="slot in experience.slots"
                    :key="slot.id"
                    type="button"
                    class="w-full rounded-lg border px-2.5 py-2 text-start text-xs transition"
                    :class="selectedSlot?.id === slot.id ? 'border-[#222222] bg-[#f7f7f7]' : 'border-gray-200 hover:border-gray-300'"
                    :disabled="slot.capacity - slot.booked_count <= 0"
                    @click="selectedSlot = slot"
                  >
                    <div class="flex items-center justify-between gap-2">
                      <span class="font-medium text-[#222222]">{{ formatDateTime(slot.start_at) }}</span>
                      <span
                        class="shrink-0 text-[11px] font-medium"
                        :class="slot.capacity - slot.booked_count > 0 ? 'text-green-700' : 'text-red-500'"
                      >
                        {{ slot.capacity - slot.booked_count > 0 ? t('experience.seatsLeft', { n: slot.capacity - slot.booked_count }) : t('experience.fullyBooked') }}
                      </span>
                    </div>
                  </button>
                </div>
              </div>

              <div>
                <label class="mb-1 block text-xs font-medium text-[#717171]">{{ t('booking.guestsCount') }}</label>
                <input v-model.number="guestsCount" type="number" min="1" :max="maxGuestsForSlot" class="input-fancy !py-2 text-sm" />
              </div>

              <div>
                <label class="mb-1 block text-xs font-medium text-[#717171]">
                  {{ t('booking.note') }}
                  <span class="text-gray-400">({{ t('common.optional') }})</span>
                </label>
                <textarea v-model="guestNote" rows="2" :placeholder="t('booking.notePlaceholder')" class="input-fancy resize-none !py-2 text-sm" />
              </div>

              <div class="flex items-center justify-between border-t border-gray-100 pt-2 text-sm">
                <span class="text-[#717171]">{{ t('booking.totalPrice') }}</span>
                <span class="font-semibold tabular-nums text-[#222222]">{{ formatPrice(totalPrice) }}</span>
              </div>

              <p v-if="bookingError" class="rounded-lg bg-red-50 px-2.5 py-2 text-xs text-red-600">{{ bookingError }}</p>

              <button
                type="button"
                class="btn-primary w-full !py-2.5 text-sm"
                :disabled="!selectedSlot || isBooking"
                @click="submitBooking"
              >
                {{ isBooking ? t('common.loading') : t('booking.confirmBooking') }}
              </button>
            </div>
          </template>
        </div>
      </aside>
    </div>
  </div>
</template>
