<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { extractErrorMessage } from '../api/client'
import AppIcon from '../components/AppIcon.vue'
import { useAuthStore } from '../stores/auth'
import { setLocale } from '../i18n'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

const fullName = ref(authStore.user?.full_name || '')
const isSaving = ref(false)
const savedMessage = ref('')

async function saveProfile() {
  isSaving.value = true
  savedMessage.value = ''
  try {
    await authStore.updateProfile({ full_name: fullName.value })
    savedMessage.value = t('profile.saved')
  } finally {
    isSaving.value = false
  }
}

function onLocaleChange(e: Event) {
  const value = (e.target as HTMLSelectElement).value
  setLocale(value)
}

const hostDisplayName = ref('')
const hostBio = ref('')
const hostCity = ref('')
const isBecomingHost = ref(false)
const becomeHostError = ref('')

async function becomeHost() {
  isBecomingHost.value = true
  becomeHostError.value = ''
  try {
    await authStore.becomeHost({
      display_name: hostDisplayName.value || undefined,
      bio: hostBio.value || undefined,
      city: hostCity.value || undefined,
    })
    router.push({ name: 'host-dashboard' })
  } catch (err) {
    becomeHostError.value = extractErrorMessage(err, t('common.errorGeneric'))
  } finally {
    isBecomingHost.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-2xl space-y-8 px-4 py-10 sm:px-6">
    <div class="surface-card animate-fade-up p-6 sm:p-7">
      <h1 class="text-xl font-extrabold text-gray-900">{{ t('profile.title') }}</h1>

      <div class="mt-5 space-y-4">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('profile.phoneNumber') }}</label>
          <input :value="authStore.user?.phone" dir="ltr" disabled class="input-fancy bg-gray-50 text-gray-500" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('profile.fullName') }}</label>
          <input v-model="fullName" type="text" class="input-fancy" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('profile.language') }}</label>
          <select :value="locale" class="input-fancy" @change="onLocaleChange">
            <option value="fa">فارسی</option>
            <option value="en">English</option>
          </select>
        </div>
        <p v-if="savedMessage" class="rounded-lg bg-green-50 px-3 py-2 text-sm text-green-600">{{ savedMessage }}</p>
        <button type="button" class="btn-primary" :disabled="isSaving" @click="saveProfile">
          {{ isSaving ? t('common.loading') : t('profile.save') }}
        </button>
      </div>
    </div>

    <div v-if="!authStore.isHost" class="surface-card animate-fade-up relative overflow-hidden p-6 sm:p-7">
      <div class="blob h-40 w-40 bg-brand-accent/15 -top-10 end-0" />
      <AppIcon icon="heroicons-outline:sparkles" size="2xl" icon-class="mb-2 text-brand-primary" />
      <h2 class="text-lg font-extrabold text-gray-900">{{ t('profile.becomeHostTitle') }}</h2>
      <p class="mt-1 text-sm text-gray-500">{{ t('profile.becomeHostBody') }}</p>

      <div class="relative mt-5 space-y-4">
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('profile.displayName') }}</label>
          <input v-model="hostDisplayName" type="text" class="input-fancy" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('profile.city') }}</label>
          <input v-model="hostCity" type="text" class="input-fancy" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('profile.bio') }}</label>
          <textarea v-model="hostBio" rows="3" class="input-fancy resize-none" />
        </div>
        <p v-if="becomeHostError" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{{ becomeHostError }}</p>
        <button type="button" class="btn-outline !border-2 !border-brand-primary !text-brand-primary hover:!bg-brand-primary hover:!text-white" :disabled="isBecomingHost" @click="becomeHost">
          {{ isBecomingHost ? t('common.loading') : t('profile.submit') }}
        </button>
      </div>
    </div>
  </div>
</template>
