<script setup>
useHead({
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1' }
  ],
  link: [
    { rel: 'icon', href: '/favicon.ico' }
  ],
  htmlAttrs: {
    lang: 'en'
  }
})

const title = 'AI Journal'
const description = 'AI Journal'

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
  ogImage: 'https://ui.nuxt.com/assets/templates/nuxt/starter-light.png',
  twitterImage: 'https://ui.nuxt.com/assets/templates/nuxt/starter-light.png',
  twitterCard: 'summary_large_image'
})

const auth = useAuthStore()

onMounted(() => { auth.fetchUser() })

</script>

<template>
  <UApp>
    <NuxtLoadingIndicator color="#0f172a" />

    <div v-if="!auth.initialized" class="fixed inset-0 flex flex-col items-center justify-center bg-slate-50">
      <div class="w-12 h-12 border-4 border-slate-200 border-t-slate-900 rounded-full animate-spin mb-4"></div>
      <p class="text-slate-500 font-medium animate-pulse">Reflecting...</p>
    </div>

    <div v-else class="min-h-screen bg-slate-50 ">
      <div class="min-h-screen bg-slate-50 ">
        <NuxtLayout>
          <NuxtPage />
        </NuxtLayout>
      </div>
    </div>
  </UApp>
</template>
