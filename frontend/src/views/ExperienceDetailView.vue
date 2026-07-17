<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import { bookingApi, experienceApi } from '../api/endpoints'
import { extractErrorMessage } from '../api/client'
import type { Experience, Slot } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import { useAuthStore } from '../stores/auth'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import EmptyState from '../components/EmptyState.vue'
import RatingStars from '../components/RatingStars.vue'
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

const activePhotoIndex = ref(0)

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

  <div v-else-if="notFound" class="mx-auto max-w-3xl px-4 py-16">
    <EmptyState :title="t('experience.notFound')" />
  </div>

  <div v-else-if="experience" class="mx-auto max-w-6xl px-4 py-8 sm:px-6">
    <!-- Gallery -->
    <div class="grid gap-2 overflow-hidden rounded-3xl sm:grid-cols-4 sm:grid-rows-2">
      <div class="aspect-[4/3] overflow-hidden bg-gray-100 sm:col-span-2 sm:row-span-2 sm:aspect-auto">
        <img
          v-if="experience.photos?.length"
          :src="experience.photos[activePhotoIndex]?.url"
          :alt="experience.title"
          class="h-full w-full object-cover transition-transform duration-500 hover:scale-105"
        />
      </div>
      <button
        v-for="(photo, idx) in experience.photos?.slice(1, 5)"
        :key="photo.id"
        type="button"
        class="hidden aspect-[4/3] overflow-hidden bg-gray-100 sm:block"
        @click="activePhotoIndex = idx + 1"
      >
        <img :src="photo.url" :alt="experience.title" class="h-full w-full object-cover transition-transform duration-500 hover:scale-105" />
      </button>
    </div>

    <div class="mt-8 grid gap-10 lg:grid-cols-3">
      <!-- Main content -->
      <div class="space-y-8 lg:col-span-2">
        <div>
          <div class="flex items-center gap-2 text-sm text-gray-500">
            <span class="flex items-center gap-1 font-medium text-gray-600">📍 {{ experience.city }}</span>
            <span>·</span>
            <RatingStars :rating="experience.rating_avg" />
            <span v-if="experience.rating_count">({{ experience.rating_count }})</span>
          </div>
          <h1 class="mt-2 text-2xl font-extrabold text-gray-900 sm:text-3xl">{{ experience.title }}</h1>
          <p class="mt-2 text-gray-600">{{ experience.summary }}</p>
        </div>

        <div class="flex flex-wrap gap-3 border-y border-gray-100 py-4 text-sm">
          <span class="flex items-center gap-1.5 rounded-full bg-gray-50 px-3 py-1.5 font-medium text-gray-700">
            ⏱️ {{ t('experience.durationMinutes', { n: experience.duration_minutes }) }}
          </span>
          <span class="flex items-center gap-1.5 rounded-full bg-gray-50 px-3 py-1.5 font-medium text-gray-700">
            👥 {{ t('experience.guests', { n: experience.min_guests, m: experience.max_guests }) }}
          </span>
        </div>

        <div v-if="experience.description" class="whitespace-pre-line leading-7 text-gray-700">
          {{ experience.description }}
        </div>

        <div v-if="experience.host" class="surface-card p-5">
          <h2 class="mb-3 text-lg font-bold text-gray-900">{{ t('experience.aboutHost') }}</h2>
          <div class="flex items-start gap-4">
            <img
              v-if="experience.host.avatar_url"
              :src="experience.host.avatar_url"
              class="h-14 w-14 rounded-full object-cover ring-2 ring-brand-primary/20"
              alt=""
            />
            <div v-else class="grid h-14 w-14 shrink-0 place-items-center rounded-full brand-gradient text-lg font-bold text-white">
              {{ (experience.host.display_name || experience.host.full_name || '?')[0] }}
            </div>
            <div>
              <p class="flex items-center gap-1.5 font-semibold text-gray-900">
                {{ experience.host.display_name || experience.host.full_name }}
                <span v-if="experience.host.is_verified" class="text-brand-accent" title="verified">✓</span>
              </p>
              <p class="text-sm text-gray-500">
                {{ experience.host.rating_count ? t('experience.rating', { rating: experience.host.rating_avg.toFixed(1), count: experience.host.rating_count }) : t('experience.noRatingYet') }}
              </p>
              <p v-if="experience.host.bio" class="mt-2 text-sm leading-6 text-gray-600">{{ experience.host.bio }}</p>
            </div>
          </div>
        </div>

        <div>
          <h2 class="mb-3 text-lg font-bold text-gray-900">{{ t('experience.location') }}</h2>
          <p class="mb-3 text-sm text-gray-500">{{ experience.address || experience.city }}</p>
          <div class="h-64 overflow-hidden rounded-2xl border border-gray-100 shadow-soft">
            <LocationMap :latitude="experience.latitude" :longitude="experience.longitude" />
          </div>
        </div>
      </div>

      <!-- Booking card -->
      <aside class="lg:col-span-1">
        <div class="surface-card sticky top-24 space-y-4 p-6">
          <p class="text-2xl font-extrabold text-gray-900 tabular-nums">
            {{ formatPrice(experience.price_amount) }}
            <span class="text-sm font-normal text-gray-500">/ {{ t('experience.perPerson') }}</span>
          </p>

          <div v-if="bookingSuccess" class="rounded-2xl bg-green-50 p-5 text-center">
            <p class="text-2xl">🎉</p>
            <p class="mt-1 font-semibold text-green-700">{{ t('booking.bookingSuccessTitle') }}</p>
            <p class="mt-1 text-sm text-green-600">{{ t('booking.bookingSuccessBody') }}</p>
            <router-link to="/trips" class="btn-primary mt-4 w-full !bg-green-600 shadow-none">
              {{ t('booking.goToMyTrips') }}
            </router-link>
          </div>

          <template v-else>
            <div v-if="!experience.slots?.length">
              <EmptyState :title="t('experience.noSlots')" />
            </div>

            <template v-else>
              <div>
                <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('experience.selectTime') }}</label>
                <div class="max-h-48 space-y-2 overflow-y-auto scrollbar-thin">
                  <button
                    v-for="slot in experience.slots"
                    :key="slot.id"
                    type="button"
                    class="w-full rounded-xl border px-3 py-2.5 text-start text-sm transition-all duration-200"
                    :class="selectedSlot?.id === slot.id ? 'border-brand-primary bg-brand-primary/5 shadow-soft' : 'border-gray-200 hover:border-gray-300'"
                    :disabled="slot.capacity - slot.booked_count <= 0"
                    @click="selectedSlot = slot"
                  >
                    <div class="flex items-center justify-between">
                      <span class="font-medium text-gray-800">{{ formatDateTime(slot.start_at) }}</span>
                      <span
                        class="rounded-full px-2 py-0.5 text-xs font-semibold"
                        :class="slot.capacity - slot.booked_count > 0 ? 'bg-green-50 text-green-600' : 'bg-red-50 text-red-500'"
                      >
                        {{ slot.capacity - slot.booked_count > 0 ? t('experience.seatsLeft', { n: slot.capacity - slot.booked_count }) : t('experience.fullyBooked') }}
                      </span>
                    </div>
                  </button>
                </div>
              </div>

              <div>
                <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('booking.guestsCount') }}</label>
                <input v-model.number="guestsCount" type="number" min="1" :max="maxGuestsForSlot" class="input-fancy" />
              </div>

              <div>
                <label class="mb-1.5 block text-sm font-semibold text-gray-700">
                  {{ t('booking.note') }} <span class="text-gray-400">({{ t('common.optional') }})</span>
                </label>
                <textarea v-model="guestNote" rows="2" :placeholder="t('booking.notePlaceholder')" class="input-fancy resize-none" />
              </div>

              <div class="flex items-center justify-between border-t border-gray-100 pt-3 text-sm font-semibold text-gray-800">
                <span>{{ t('booking.totalPrice') }}</span>
                <span class="text-lg tabular-nums text-brand-primary">{{ formatPrice(totalPrice) }}</span>
              </div>

              <p v-if="bookingError" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{{ bookingError }}</p>

              <button type="button" class="btn-primary w-full" :disabled="!selectedSlot || isBooking" @click="submitBooking">
                {{ isBooking ? t('common.loading') : t('booking.confirmBooking') }}
              </button>
            </template>
          </template>
        </div>
      </aside>
    </div>
  </div>
</template>
