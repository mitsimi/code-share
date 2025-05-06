<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click="$emit('close')"
  >
    <div
      class="bg-background w-full max-w-2xl rounded-lg border-4 border-black p-6 shadow-[8px_8px_0_0_#000]"
      @click.stop
    >
      <div class="mb-6 flex items-center justify-between">
        <h2 class="text-2xl font-bold">Add New Snippet</h2>
        <button
          class="flex h-8 w-8 items-center justify-center rounded-lg border-4 border-black bg-white shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
          @click="$emit('close')"
        >
          <X class="h-4 w-4" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label for="title" class="mb-2 block font-bold">Title</label>
          <input
            id="title"
            v-model="title"
            type="text"
            required
            class="focus:ring-accent w-full rounded-lg border-4 border-black bg-white p-2 shadow-[4px_4px_0_0_#000] focus:ring-2 focus:outline-none"
          />
        </div>

        <div>
          <label for="code" class="mb-2 block font-bold">Code</label>
          <textarea
            id="code"
            v-model="code"
            required
            rows="8"
            class="focus:ring-accent w-full rounded-lg border-4 border-black bg-white p-2 font-mono text-sm shadow-[4px_4px_0_0_#000] focus:ring-2 focus:outline-none"
          ></textarea>
        </div>

        <div>
          <label for="author" class="mb-2 block font-bold">Author</label>
          <input
            id="author"
            v-model="author"
            type="text"
            required
            class="focus:ring-accent w-full rounded-lg border-4 border-black bg-white p-2 shadow-[4px_4px_0_0_#000] focus:ring-2 focus:outline-none"
          />
        </div>

        <div class="flex justify-end gap-4">
          <button
            type="button"
            class="rounded-lg border-4 border-black bg-white px-6 py-2 font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
            @click="$emit('close')"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="bg-accent rounded-lg border-4 border-black px-6 py-2 font-bold text-white shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
          >
            Submit
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { X } from 'lucide-vue-next'

defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', data: { title: string; code: string; author: string }): void
}>()

const title = ref('')
const code = ref('')
const author = ref('')

const handleSubmit = () => {
  emit('submit', {
    title: title.value,
    code: code.value,
    author: author.value,
  })
  title.value = ''
  code.value = ''
  author.value = ''
}
</script>
