<script setup>
const toast = useToast()
const content = ref('')
const isSaving = ref(false)
const auth = useAuthStore()
const showInsight = ref(false)
const lastInsight = ref(null)
const refreshHistory = inject('refreshHistory')

const route = useRoute()

watch(() => route.query.id, async (id) => {
  if (!id) {
    content.value = ''
    lastInsight.value = null
    return
  }

  try {
    const response = await $fetch(`/api/entries/${id}`)

    content.value = response.entry.content

    lastInsight.value = response.insight
    showInsight.value = true
  } catch (err) {
    console.error('Failed to load entry', err)
  }
}, { immediate: true })

function startNewEntry() {
  navigateTo('/')
}

async function saveEntry() {
  if (!content.value.trim()) return

  isSaving.value = true
  try {
    const response = await $fetch('/api/entries', {
      method: 'POST',
      body: { content: content.value }
    })

    if (response.insight) {
      lastInsight.value = response.insight
      showInsight.value = true
      content.value = ''
      if (refreshHistory) refreshHistory()
    } else {
      alert('Entrada salva, mas o terapeuta está em silêncio (sem insight).')
    }
  } catch (_) {
    toast.add({
      title: 'Failed to save entry',
      description: 'Internal server error... please try again later',
      icon: 'i-lucide-calendar-days'
    })
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto py-12 px-6">
    <header class="mb-12 flex justify-between items-end">
      <div>
        <p class="text-sm font-medium text-slate-400 uppercase tracking-widest mb-1">
          {{ route.query.id ? 'Historical Reflection' : 'New Reflection' }}
        </p>
        <h1 class="text-4xl font-serif text-slate-900">
          {{ route.query.id ? 'Past Thought' : 'Deep Reflection' }}
        </h1>
      </div>

      <div class="flex gap-3">
        <UButton v-if="route.query.id" label="New Entry" variant="ghost" color="neutral" icon="i-lucide-plus"
          class="cursor-pointer" @click="startNewEntry" />

        <UButton v-if="!route.query.id" label="Save Entry" color="black" :loading="isSaving" @click="saveEntry"
          class="cursor-pointer" />
      </div>
    </header>
    <main>
      <textarea v-model="content" :readonly="!!route.query.id" placeholder="Start typing your thoughts..."
        class="w-full h-[60vh] text-xl leading-relaxed text-slate-700 bg-transparent border-none focus:ring-0 resize-none placeholder:text-slate-300 font-serif" />
    </main>
    <USlideover v-model:open="showInsight" title="Reflection">
      <template #content>
        <div v-if="lastInsight" class="p-8 space-y-8 h-full bg-white overflow-y-auto">
          <header>
            <h2 class="text-2xl font-serif text-slate-900">
              Therapist's Reflection
            </h2>
          </header>

          <section>
            <h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-[0.2em] mb-3">
              Summary
            </h3>
            <p class="text-slate-600 leading-relaxed text-sm">
              {{ lastInsight.summary }}
            </p>
          </section>

          <section>
            <h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-[0.2em] mb-3">
              Feelings
            </h3>
            <div class="flex flex-wrap gap-2">
              <UBadge v-for="f in lastInsight.feelings" :key="f" variant="subtle" color="neutral">
                {{ f }}
              </UBadge>
            </div>
          </section>

          <section class="pt-8 border-t border-slate-100">
            <h3 class="text-[10px] font-bold text-slate-400 uppercase tracking-[0.2em] mb-4">
              Philosophical
              Mirror
            </h3>
            <p class="text-xl font-serif italic text-slate-800 leading-relaxed">
              "{{ lastInsight.reflection }}"
            </p>
          </section>

          <UButton label="Close" class="cursor-pointer mt-10" block variant="ghost" color="neutral"
            @click="showInsight = false" />
        </div>

        <div v-else class="p-12 text-center text-slate-400">
          <p>Waiting for reflection...</p>
        </div>
      </template>
    </USlideover>
  </div>
</template>
