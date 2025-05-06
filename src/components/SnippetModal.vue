<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click="$emit('close')"
  >
    <div
      class="bg-main-bg relative max-h-[90vh] w-[90%] max-w-2xl -rotate-1 overflow-y-auto border-4 border-black p-8 shadow-[8px_8px_0_0_#000]"
      @click.stop
    >
      <div class="mb-6 flex items-center justify-between border-b-3 border-black pb-4">
        <h2 class="m-0 text-2xl font-black uppercase">Submit a new snippet</h2>
        <button
          class="hover:bg-accent flex h-9 w-9 cursor-pointer items-center justify-center border-3 border-black bg-white text-2xl shadow-[3px_3px_0_0_#000] hover:text-white"
          @click="$emit('close')"
        >
          &times;
        </button>
      </div>

      <form
        class="grid grid-cols-[auto_1fr] items-start gap-4"
        @submit.prevent="$emit('submit', formData)"
      >
        <label for="title-modal" class="self-center text-base font-bold text-black">Title: </label>
        <input
          type="text"
          v-model="formData.title"
          id="title-modal"
          required
          class="box-border w-full resize-none border-3 border-black bg-white p-2 text-black shadow-[3px_3px_0_0_#000]"
        />

        <label for="snippet-modal" class="self-center text-base font-bold text-black"
          >Snippet:
        </label>
        <textarea
          v-model="formData.code"
          id="snippet-modal"
          required
          @keydown="handleTabKey"
          class="box-border h-[150px] w-full resize-none border-3 border-black bg-white p-2 text-black shadow-[3px_3px_0_0_#000]"
        ></textarea>

        <label for="name-modal" class="self-center text-base font-bold text-black"
          >Your Name:
        </label>
        <input
          type="text"
          v-model="formData.author"
          id="name-modal"
          required
          class="box-border w-full resize-none border-3 border-black bg-white p-2 text-black shadow-[3px_3px_0_0_#000]"
        />

        <button
          type="submit"
          class="bg-accent col-span-2 mt-2 cursor-pointer border-3 border-black p-3 text-base font-black text-white uppercase shadow-[3px_3px_0_0_#000] transition-all duration-200"
        >
          Submit
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface FormData {
  title: string
  code: string
  author: string
}

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', data: FormData): void
}>()

const formData = ref<FormData>({
  title: '',
  code: '',
  author: '',
})

const handleTabKey = (event: KeyboardEvent) => {
  if (event.key === 'Tab') {
    event.preventDefault()

    const textarea = event.target as HTMLTextAreaElement
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = textarea.value

    if (start !== end && value.substring(start, end).includes('\n')) {
      const selectedText = value.substring(start, end)
      const lines = selectedText.split('\n')

      const newLines = lines.map((line) => (event.shiftKey ? removeTab(line) : '    ' + line))
      const newText = newLines.join('\n')

      textarea.value = value.substring(0, start) + newText + value.substring(end)
      textarea.selectionStart = start
      textarea.selectionEnd = start + newText.length
    } else {
      const spaces = '    '
      textarea.value = value.substring(0, start) + spaces + value.substring(end)
      textarea.selectionStart = textarea.selectionEnd = start + spaces.length
    }

    formData.value.code = textarea.value
  }
}

const removeTab = (line: string) => {
  if (line.startsWith('\t')) {
    return line.substring(1)
  }
  if (line.startsWith('    ')) {
    return line.substring(4)
  }
  if (line.startsWith('  ')) {
    return line.substring(2)
  }
  return line
}
</script>

<style scoped>
.bg-main-bg {
  background-color: #fcf7de;
}
.bg-accent {
  background-color: #ff3e4d;
}
</style>
