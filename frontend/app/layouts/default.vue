<script setup>
const auth = useAuthStore()
const entries = ref([])
const loading = ref(true)

async function fetchHistory() {
  try {
    const response = await $fetch('/api/entries?page=1')
    entries.value = response.data
  } catch (err) {
    console.error('Failed to fetch history', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchHistory()
})

provide('refreshHistory', fetchHistory)
</script>

<template>
  <div class="flex h-screen bg-slate-50">
    <aside class="w-72 border-r border-slate-200 bg-white flex flex-col">
      <div class="p-6 border-b border-slate-100">
        <h2 class="text-sm font-bold text-slate-900 uppercase tracking-widest">Journal</h2>
      </div>

      <nav class="flex-1 overflow-y-auto p-4 space-y-2">
        <div v-if="loading" class="space-y-3 p-2">
          <div v-for="i in 3" :key="i" class="h-12 bg-slate-50 animate-pulse rounded-lg"></div>
        </div>

        <template v-else>
          <button v-for="entry in entries" :key="entry.id"
            class="w-full text-left p-3 rounded-lg hover:bg-slate-50 transition-colors group">
            <p class="text-xs text-slate-400 mb-1">
              {{ new Date(entry.created_at).toLocaleDateString() }}
            </p>
            <p class="text-sm text-slate-700 truncate font-medium">
              {{ entry.content }}
            </p>
          </button>
        </template>
      </nav>

      <div class="p-4 border-t border-slate-100">
        <div class="flex items-center gap-3 px-2">
          <div class="w-8 h-8 rounded-full bg-slate-900 flex items-center justify-center text-white text-xs">
            {{ auth.user?.name?.charAt(0) }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-slate-900 truncate">{{ auth.user?.name }}</p>
            <p class="text-xs text-slate-400 truncate">{{ auth.user?.is_pro ? 'Pro Plan' : 'Trial' }}</p>
          </div>
          <UButton icon="i-lucide-log-out" variant="ghost" color="neutral" size="xs" @click="auth.logout" />
        </div>
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto">
      <slot />
    </main>
  </div>
</template>


<style></style>