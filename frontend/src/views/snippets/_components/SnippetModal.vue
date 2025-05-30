<template>
  <Dialog :open="show" @update:open="$emit('close')">
    <DialogContent class="sm:max-w-[600px]">
      <DialogHeader>
        <DialogTitle>Create New Snippet</DialogTitle>
        <DialogDescription>
          Share your code snippet with the community. Make sure to provide a clear title and
          description.
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div class="space-y-2">
          <Label for="title">Title</Label>
          <Input
            id="title"
            v-model="formData.title"
            placeholder="Enter a descriptive title"
            required
          />
        </div>

        <div class="space-y-2">
          <Label for="code">Code</Label>
          <Textarea
            id="code"
            v-model="formData.code"
            placeholder="Paste your code here"
            class="font-mono"
            rows="10"
            required
          />
        </div>

        <div class="flex justify-end gap-3">
          <Button type="button" variant="outline" @click="$emit('close')"> Cancel </Button>
          <Button type="submit" :disabled="isLoading">
            <LoaderCircleIcon v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
            Create Snippet
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { LoaderCircleIcon } from 'lucide-vue-next'
import { ref } from 'vue'

const props = defineProps<{
  show: boolean
  isLoading: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', data: { title: string; code: string }): void
}>()

const formData = ref({
  title: '',
  code: '',
})

const handleSubmit = () => {
  emit('submit', {
    title: formData.value.title,
    code: formData.value.code,
  })
}
</script>
