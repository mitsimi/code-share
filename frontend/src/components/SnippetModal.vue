<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click="$emit('close')"
  >
    <div class="bg-card w-full max-w-2xl rounded-lg border-2 p-6 shadow" @click.stop>
      <div class="mb-6 flex items-center justify-between">
        <h2 class="text-2xl font-bold">Add New Snippet</h2>
        <Button variant="outline" size="icon" @click="$emit('close')">
          <X class="size-4" />
        </Button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <Label for="title" class="mb-2 block font-bold">Title</Label>
          <Input id="title" v-model="title" type="text" required />
        </div>

        <div>
          <Label for="code" class="mb-2 block font-bold">Code</Label>
          <Textarea
            id="code"
            v-model="code"
            required
            :rows="8"
            class="resize-none font-mono"
            @keydown="handleTab"
          />
        </div>

        <div>
          <Label for="author" class="mb-2 block font-bold">Author</Label>
          <Input id="author" v-model="author" type="text" required />
        </div>

        <div class="flex justify-end gap-4">
          <Button variant="outline" type="button" @click="$emit('close')"> Cancel </Button>
          <Button variant="noShadow" type="submit"> Submit </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Textarea } from './ui/textarea'
import { Input } from './ui/input'
import { Label } from './ui/label'
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

const handleTab = (e: KeyboardEvent) => {
  if (e.key === 'Tab') {
    e.preventDefault()
    const textarea = e.target as HTMLTextAreaElement
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = code.value

    if (e.shiftKey) {
      // Handle unindent
      const lines = value.split('\n')
      const startLine = value.substring(0, start).split('\n').length - 1
      const endLine = value.substring(0, end).split('\n').length - 1

      const newLines = lines.map((line, i) => {
        if (i >= startLine && i <= endLine && line.startsWith('\t')) {
          return line.substring(1)
        }
        return line
      })

      code.value = newLines.join('\n')
      textarea.selectionStart = start - (startLine === endLine ? 1 : 0)
      textarea.selectionEnd = end - (endLine - startLine + 1)
    } else {
      // Handle indent
      if (start === end) {
        // Single line indent
        code.value = value.substring(0, start) + '\t' + value.substring(end)
        textarea.selectionStart = textarea.selectionEnd = start + 1
      } else {
        // Multi-line indent
        const lines = value.split('\n')
        const startLine = value.substring(0, start).split('\n').length - 1
        const endLine = value.substring(0, end).split('\n').length - 1

        const newLines = lines.map((line, i) => {
          if (i >= startLine && i <= endLine) {
            return '\t' + line
          }
          return line
        })

        code.value = newLines.join('\n')
        textarea.selectionStart = start + 1
        textarea.selectionEnd = end + (endLine - startLine + 1)
      }
    }
  }
}

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
