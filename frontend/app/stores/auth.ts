import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const initialized = ref(false)
  const user = ref(null)
  const isLoggedIn = computed(() => !!user.value)

  async function fetchUser() {
    try {
      const data = await $fetch('/api/auth/me')
      user.value = data
    } catch (err) {
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  async function logout() {
    await $fetch('/api/auth/logout', { method: 'POST' })
    user.value = null
    navigateTo('/login')
  }

  return { user, isLoggedIn, initialized, fetchUser, logout }
})
