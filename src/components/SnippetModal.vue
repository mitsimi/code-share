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
        <Button variant="outline" size="icon" @click="$emit('close')">
          <X class="size-4" />
        </Button>
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
          <Button variant="secondary" type="button" @click="$emit('close')"> Cancel </Button>
          <Button type="submit"> Submit </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Button } from './ui/button'
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
