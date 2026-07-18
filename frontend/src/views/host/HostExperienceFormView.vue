<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { categoryApi, hostExperienceApi, type ExperienceFormPayload } from '../../api/endpoints'
import { extractErrorMessage } from '../../api/client'
import type { Category, Experience } from '../../api/types'
import { useFormatters } from '../../composables/useFormatters'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import AppIcon from '../../components/AppIcon.vue'
import LocationMap from '../../components/LocationMap.vue'

const props = defineProps<{ id?: string }>()
const { t, locale } = useI18n()
const { formatDateTime } = useFormatters()
const router = useRouter()

const isEditing = computed(() => Boolean(props.id))
const isLoading = ref(true)
const isSaving = ref(false)
const saveError = ref('')
const saveSuccess = ref(false)

const categories = ref<Category[]>([])
const experience = ref<Experience | null>(null)

const form = reactive({
  title: '',
  summary: '',
  description: '',
  category_id: '' as string,
  priceToman: 0,
  duration_minutes: 120,
  max_guests: 6,
  min_guests: 1,
  city: '',
  address: '',
  latitude: 0,
  longitude: 0,
  status: 'draft' as 'draft' | 'published' | 'archived',
})

function fillFormFrom(exp: Experience) {
  form.title = exp.title
  form.summary = exp.summary || ''
  form.description = exp.description || ''
  form.category_id = exp.category_id || ''
  form.priceToman = Math.round(exp.price_amount / 10)
  form.duration_minutes = exp.duration_minutes
  form.max_guests = exp.max_guests
  form.min_guests = exp.min_guests
  form.city = exp.city
  form.address = exp.address || ''
  form.latitude = exp.latitude
  form.longitude = exp.longitude
  form.status = exp.status
}

async function load() {
  isLoading.value = true
  try {
    categories.value = await categoryApi.list()
    if (props.id) {
      experience.value = await hostExperienceApi.get(props.id)
      fillFormFrom(experience.value)
    }
  } finally {
    isLoading.value = false
  }
}

function onLocationPicked({ lat, lng }: { lat: number; lng: number }) {
  form.latitude = lat
  form.longitude = lng
}

function buildPayload(): ExperienceFormPayload {
  return {
    title: form.title,
    summary: form.summary || undefined,
    description: form.description || undefined,
    category_id: form.category_id || null,
    price_amount: Math.round(form.priceToman * 10),
    duration_minutes: form.duration_minutes,
    max_guests: form.max_guests,
    min_guests: form.min_guests,
    city: form.city,
    address: form.address || undefined,
    latitude: form.latitude,
    longitude: form.longitude,
    status: form.status,
  }
}

async function save() {
  saveError.value = ''
  saveSuccess.value = false
  if (!form.title || !form.city || !form.priceToman || !form.latitude || !form.longitude) {
    saveError.value = t('host.requiredFieldsHint')
    return
  }

  isSaving.value = true
  try {
    if (isEditing.value && props.id) {
      experience.value = await hostExperienceApi.update(props.id, buildPayload())
      saveSuccess.value = true
    } else {
      const created = await hostExperienceApi.create(buildPayload())
      router.replace({ name: 'host-experience-edit', params: { id: created.id } })
    }
  } catch (err) {
    saveError.value = extractErrorMessage(err, t('common.errorGeneric'))
  } finally {
    isSaving.value = false
  }
}

// --- Photos ---
const newPhotoUrl = ref('')

async function addPhoto() {
  if (!props.id || !newPhotoUrl.value) return
  await hostExperienceApi.addPhoto(props.id, { url: newPhotoUrl.value, is_cover: !experience.value?.photos?.length })
  newPhotoUrl.value = ''
  experience.value = await hostExperienceApi.get(props.id)
}

async function removePhoto(photoId: string) {
  if (!props.id) return
  await hostExperienceApi.deletePhoto(props.id, photoId)
  experience.value = await hostExperienceApi.get(props.id)
}

async function makeCover(photoId: string) {
  if (!props.id || !experience.value) return
  const photo = experience.value.photos?.find((p) => p.id === photoId)
  if (!photo) return
  await hostExperienceApi.addPhoto(props.id, { url: photo.url, alt: photo.alt || undefined, is_cover: true })
  await hostExperienceApi.deletePhoto(props.id, photoId)
  experience.value = await hostExperienceApi.get(props.id)
}

// --- Slots ---
const newSlotStart = ref('')
const newSlotEnd = ref('')
const newSlotCapacity = ref(6)

async function addSlot() {
  if (!props.id || !newSlotStart.value || !newSlotEnd.value) return
  await hostExperienceApi.addSlot(props.id, {
    start_at: new Date(newSlotStart.value).toISOString(),
    end_at: new Date(newSlotEnd.value).toISOString(),
    capacity: newSlotCapacity.value,
  })
  newSlotStart.value = ''
  newSlotEnd.value = ''
  experience.value = await hostExperienceApi.get(props.id)
}

async function removeSlot(slotId: string) {
  if (!props.id) return
  await hostExperienceApi.deleteSlot(props.id, slotId)
  experience.value = await hostExperienceApi.get(props.id)
}

async function archive() {
  if (!props.id) return
  await hostExperienceApi.archive(props.id)
  router.push({ name: 'host-experiences' })
}

onMounted(load)
</script>

<template>
  <div>
    <router-link :to="{ name: 'host-experiences' }" class="mb-4 inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      <AppIcon icon="heroicons-outline:arrow-left" size="sm" />
      {{ t('host.backToList') }}
    </router-link>

    <h1 class="text-xl font-bold text-gray-900">
      {{ isEditing ? t('host.editExperience') : t('host.createExperience') }}
    </h1>

    <LoadingSpinner v-if="isLoading" :label="t('common.loading')" />

    <div v-else class="mt-6 grid gap-6 lg:grid-cols-3">
      <form class="surface-card space-y-4 p-6 lg:col-span-2" @submit.prevent="save">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.titleLabel') }}</label>
          <input v-model="form.title" type="text" class="input-fancy" required />
        </div>

        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.summaryLabel') }}</label>
          <input v-model="form.summary" type="text" class="input-fancy" />
        </div>

        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.descriptionLabel') }}</label>
          <textarea v-model="form.description" rows="4" class="input-fancy" />
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.categoryLabel') }}</label>
            <select v-model="form.category_id" class="input-fancy">
              <option value="">—</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ locale === 'fa' ? c.title_fa : c.title_en }}</option>
            </select>
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.statusLabel') }}</label>
            <select v-model="form.status" class="input-fancy">
              <option value="draft">{{ t('host.statusDraft') }}</option>
              <option value="published">{{ t('host.statusPublished') }}</option>
              <option value="archived">{{ t('host.statusArchived') }}</option>
            </select>
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.priceLabel') }}</label>
            <input v-model.number="form.priceToman" type="number" min="0" class="input-fancy" required />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.durationLabel') }}</label>
            <input v-model.number="form.duration_minutes" type="number" min="15" class="input-fancy" />
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.minGuestsLabel') }}</label>
            <input v-model.number="form.min_guests" type="number" min="1" class="input-fancy" />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.maxGuestsLabel') }}</label>
            <input v-model.number="form.max_guests" type="number" min="1" class="input-fancy" />
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.cityLabel') }}</label>
            <input v-model="form.city" type="text" class="input-fancy" required />
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.addressLabel') }}</label>
            <input v-model="form.address" type="text" class="input-fancy" />
          </div>
        </div>

        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700">{{ t('host.locationLabel') }}</label>
          <p class="mb-2 text-xs text-gray-500">{{ t('host.pickOnMapHint') }}</p>
          <div class="h-64 overflow-hidden rounded-xl border border-gray-200">
            <LocationMap :latitude="form.latitude" :longitude="form.longitude" interactive @update:location="onLocationPicked" />
          </div>
        </div>

        <p v-if="saveError" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{{ saveError }}</p>
        <p v-if="saveSuccess" class="rounded-lg bg-green-50 px-3 py-2 text-sm text-green-600">{{ t('host.savedSuccess') }}</p>

        <div class="flex items-center gap-3">
          <button type="submit" class="btn-primary" :disabled="isSaving">
            {{ isSaving ? t('common.loading') : isEditing ? t('host.save') : t('host.createAndContinue') }}
          </button>
          <button v-if="isEditing" type="button" class="text-sm font-medium text-red-600 hover:underline" @click="archive">
            {{ t('host.archiveExperience') }}
          </button>
        </div>
      </form>

      <div v-if="isEditing" class="space-y-6">
        <!-- Photos -->
        <div class="surface-card p-5">
          <h2 class="mb-3 text-sm font-bold text-gray-900">{{ t('host.photosLabel') }}</h2>
          <div class="space-y-2">
            <div v-for="photo in experience?.photos" :key="photo.id" class="flex items-center gap-2">
              <img :src="photo.url" class="h-12 w-16 rounded-lg object-cover" alt="" />
              <span v-if="photo.is_cover" class="rounded bg-brand-primary/10 px-1.5 py-0.5 text-[10px] font-medium text-brand-primary">
                {{ t('host.coverPhoto') }}
              </span>
              <button v-else type="button" class="text-[11px] text-gray-500 hover:text-brand-primary" @click="makeCover(photo.id)">
                {{ t('host.makeCover') }}
              </button>
              <button type="button" class="ms-auto text-xs text-red-500 hover:underline" @click="removePhoto(photo.id)">
                {{ t('common.delete') }}
              </button>
            </div>
          </div>
          <div class="mt-3 flex gap-2">
            <input v-model="newPhotoUrl" type="text" :placeholder="t('host.photoUrlPlaceholder')" dir="ltr" class="flex-1 rounded-lg border border-gray-300 px-2.5 py-2 text-xs focus:border-brand-primary focus:outline-none" />
            <button type="button" class="rounded-lg bg-gray-900 px-3 py-2 text-xs font-semibold text-white" @click="addPhoto">
              {{ t('host.addPhotoUrl') }}
            </button>
          </div>
        </div>

        <!-- Slots -->
        <div class="surface-card p-5">
          <h2 class="mb-3 text-sm font-bold text-gray-900">{{ t('host.slotsLabel') }}</h2>
          <div class="space-y-2">
            <div v-for="slot in experience?.slots" :key="slot.id" class="flex items-center justify-between rounded-lg bg-gray-50 px-3 py-2 text-xs">
              <span>{{ formatDateTime(slot.start_at) }}</span>
              <span>{{ slot.booked_count }}/{{ slot.capacity }}</span>
              <button type="button" class="text-red-500 hover:underline" @click="removeSlot(slot.id)">{{ t('common.delete') }}</button>
            </div>
          </div>
          <div class="mt-3 space-y-2">
            <input v-model="newSlotStart" type="datetime-local" class="w-full rounded-lg border border-gray-300 px-2.5 py-2 text-xs focus:border-brand-primary focus:outline-none" />
            <input v-model="newSlotEnd" type="datetime-local" class="w-full rounded-lg border border-gray-300 px-2.5 py-2 text-xs focus:border-brand-primary focus:outline-none" />
            <input v-model.number="newSlotCapacity" type="number" min="1" :placeholder="t('host.capacity')" class="w-full rounded-lg border border-gray-300 px-2.5 py-2 text-xs focus:border-brand-primary focus:outline-none" />
            <button type="button" class="w-full rounded-lg bg-gray-900 px-3 py-2 text-xs font-semibold text-white" @click="addSlot">
              {{ t('host.addSlot') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
