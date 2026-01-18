<script setup>
const auth = useAuthStore()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await $fetch('/api/auth/login', {
      method: 'POST',
      body: { email: email.value, password: password.value }
    })
    await auth.fetchUser()
    navigateTo('/')
  } catch (err) {
    error.value = 'Invalid email or password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center p-4">
    <UCard class="w-full max-w-sm shadow-xl ring-1 ring-slate-200">
      <template #header>
        <div class="text-center">
          <h1 class="text-2xl font-bold text-slate-900 ">Welcome Back</h1>
          <p class="text-sm text-slate-500 mt-1">Sign in to your philosophical journal</p>
        </div>
      </template>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <UFormField label="Email">
          <UInput v-model="email" type="email" placeholder="you@example.com" icon="i-lucide-mail" class="w-full" />
        </UFormField>

        <UFormField label="Password">
          <UInput v-model="password" type="password" icon="i-lucide-lock" class="w-full" />
        </UFormField>

        <UButton type="submit" label="Sign In" block :loading="loading" color="neutral" variant="solid" class="mt-6" />

        <p v-if="error" class="text-center text-sm text-red-500 font-medium">{{ error }}</p>
      </form>

      <template #footer>
        <p class="text-center text-xs text-slate-400">
          Your entries are private and encrypted.
        </p>
      </template>
    </UCard>
  </div>
</template>