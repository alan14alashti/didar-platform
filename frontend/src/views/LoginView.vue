<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import { extractErrorMessage } from '../api/client'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const step = ref<'email' | 'otp'>('email')
const email = ref('')
const code = ref('')
const isSubmitting = ref(false)
const errorMessage = ref('')
const cooldown = ref(0)
let cooldownTimer: ReturnType<typeof setInterval> | undefined

function startCooldown() {
  cooldown.value = 30
  clearInterval(cooldownTimer)
  cooldownTimer = setInterval(() => {
    cooldown.value -= 1
    if (cooldown.value <= 0) clearInterval(cooldownTimer)
  }, 1000)
}

async function sendCode() {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    await authStore.requestOtp(email.value)
    step.value = 'otp'
    startCooldown()
  } catch (err) {
    errorMessage.value = extractErrorMessage(err, t('auth.invalidEmail'))
  } finally {
    isSubmitting.value = false
  }
}

async function verify() {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    await authStore.verifyOtp(email.value, code.value)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err) {
    errorMessage.value = extractErrorMessage(err, t('common.errorGeneric'))
  } finally {
    isSubmitting.value = false
  }
}

function changeEmail() {
  step.value = 'email'
  code.value = ''
  errorMessage.value = ''
}
</script>

<template>
  <div class="relative flex min-h-[80vh] flex-col items-center justify-center overflow-hidden px-4 py-12">
    <div class="blob animate-blob h-72 w-72 bg-brand-primary/25 -top-10 -start-16" />
    <div class="blob animate-blob animate-blob-delay-2 h-80 w-80 bg-brand-accent/20 bottom-0 end-0" />

    <div class="relative w-full max-w-md">
      <div class="surface-card animate-fade-up overflow-hidden p-8">
        <div class="mx-auto mb-5 grid h-14 w-14 place-items-center rounded-2xl brand-gradient text-2xl shadow-glow">
          🔐
        </div>
        <h1 class="text-center text-2xl font-extrabold text-gray-900">{{ t('auth.loginTitle') }}</h1>
        <p class="mt-1.5 text-center text-sm text-gray-500">{{ t('auth.loginSubtitle') }}</p>

        <transition
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="opacity-0 translate-x-3"
          enter-to-class="opacity-100 translate-x-0"
          mode="out-in"
        >
          <form v-if="step === 'email'" key="email" class="mt-7 space-y-4" @submit.prevent="sendCode">
            <div>
              <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('auth.emailLabel') }}</label>
              <input
                v-model="email"
                type="email"
                inputmode="email"
                autocomplete="email"
                dir="ltr"
                :placeholder="t('auth.emailPlaceholder')"
                class="input-fancy text-center text-lg tracking-wide"
                required
              />
            </div>
            <p v-if="errorMessage" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{{ errorMessage }}</p>
            <button type="submit" class="btn-primary w-full" :disabled="isSubmitting">
              {{ isSubmitting ? t('common.loading') : t('auth.sendCode') }}
            </button>
          </form>

          <form v-else key="otp" class="mt-7 space-y-4" @submit.prevent="verify">
            <p class="rounded-xl bg-brand-primary/5 px-3 py-2.5 text-center text-sm text-gray-600" dir="ltr">
              {{ t('auth.otpSentTo', { email }) }}
            </p>
            <div>
              <label class="mb-1.5 block text-sm font-semibold text-gray-700">{{ t('auth.otpLabel') }}</label>
              <input
                v-model="code"
                type="text"
                inputmode="numeric"
                dir="ltr"
                :placeholder="t('auth.otpPlaceholder')"
                class="input-fancy text-center text-2xl font-bold tracking-[0.5em]"
                required
                autofocus
              />
            </div>
            <p v-if="errorMessage" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{{ errorMessage }}</p>
            <button type="submit" class="btn-primary w-full" :disabled="isSubmitting">
              {{ isSubmitting ? t('common.loading') : t('auth.verify') }}
            </button>
            <div class="flex items-center justify-between pt-1 text-sm">
              <button type="button" class="font-medium text-gray-500 transition hover:text-gray-700" @click="changeEmail">
                {{ t('auth.changeEmail') }}
              </button>
              <button
                type="button"
                class="font-semibold text-brand-primary transition hover:opacity-80 disabled:text-gray-400"
                :disabled="cooldown > 0"
                @click="sendCode"
              >
                {{ t('auth.resendCode') }}{{ cooldown > 0 ? ` (${cooldown})` : '' }}
              </button>
            </div>
          </form>
        </transition>
      </div>
    </div>
  </div>
</template>
