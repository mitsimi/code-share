<template>
  <Dialog :open="show" @update:open="$emit('close')">
    <DialogContent class="sm:max-w-[600px]">
      <DialogHeader>
        <DialogTitle>{{ isEditMode ? 'Edit Snippet' : 'Create New Snippet' }}</DialogTitle>
        <DialogDescription>
          {{
            isEditMode
              ? 'Update your code snippet. Make sure to provide a clear title and description.'
              : 'Share your code snippet with the community. Make sure to provide a clear title and description.'
          }}
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
          <Label>Language</Label>
          <Combobox v-model="selectedLanguage" by="value">
            <ComboboxAnchor>
              <div class="relative w-full items-center">
                <ComboboxInput
                  :display-value="(lang) => lang?.label ?? ''"
                  placeholder="Select language..."
                  class="w-full"
                />
                <ComboboxTrigger
                  class="absolute inset-y-0 end-0 flex items-center justify-center px-3"
                >
                  <ChevronsUpDownIcon class="text-muted-foreground size-4" />
                </ComboboxTrigger>
              </div>
            </ComboboxAnchor>

            <ComboboxList>
              <ComboboxEmpty> No language found. </ComboboxEmpty>

              <ComboboxGroup
                class="overflow-y-auto"
                :style="{
                  maxHeight: 'min(400px, calc(100vh - 200px))',
                }"
                ref="comboboxGroupRef"
              >
                <ComboboxItem
                  v-for="language in filteredLanguages"
                  :key="language.value"
                  :value="language"
                  @select="handleLanguageSelect"
                  class="flex items-center justify-between"
                >
                  <div class="flex items-center gap-2">
                    {{ language.label }}
                    <span class="text-muted-foreground text-xs">{{ language.extension }}</span>
                  </div>

                  <ComboboxItemIndicator>
                    <CheckIcon class="h-4 w-4" />
                  </ComboboxItemIndicator>
                </ComboboxItem>
              </ComboboxGroup>
            </ComboboxList>
          </Combobox>
        </div>

        <div class="space-y-2">
          <Label for="code">Code</Label>
          <Textarea
            id="code"
            v-model="formData.content"
            placeholder="Paste your code here"
            class="font-mono"
            rows="10"
            required
          />
        </div>

        <div class="flex justify-end gap-3">
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" variant="reverse" :disabled="isLoading">
            <LoaderCircleIcon v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
            {{ isEditMode ? 'Save Changes' : 'Create Snippet' }}
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { LoaderCircleIcon, ChevronsUpDownIcon, CheckIcon } from 'lucide-vue-next'
import { ref, watch, computed } from 'vue'
import type { Snippet } from '@/types'
import { Button } from '@/components/ui/button'
import Dialog from '@/components/ui/dialog/Dialog.vue'
import DialogContent from '@/components/ui/dialog/DialogContent.vue'
import DialogHeader from '@/components/ui/dialog/DialogHeader.vue'
import DialogTitle from '@/components/ui/dialog/DialogTitle.vue'
import DialogDescription from '@/components/ui/dialog/DialogDescription.vue'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Combobox,
  ComboboxAnchor,
  ComboboxInput,
  ComboboxTrigger,
  ComboboxList,
  ComboboxItem,
  ComboboxItemIndicator,
  ComboboxEmpty,
  ComboboxGroup,
} from '@/components/ui/combobox'
import { allLanguages } from '@/lib/languages'
import type { ListboxItemSelectEvent, AcceptableValue } from 'reka-ui'

const props = defineProps<{
  show: boolean
  isLoading: boolean
  snippet?: Snippet // Optional - if provided, we're in edit mode
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', data: { title: string; content: string; language?: string }): void
}>()

// Determine if we're in edit mode based on whether a snippet is provided
const isEditMode = computed(() => !!props.snippet)

const formData = ref({
  title: '',
  content: '',
  language: '',
})

// Transform languages to label/value format
const languageOptions = computed(() => {
  return allLanguages.map((lang) => ({
    label: lang.displayName,
    value: lang.name.toLowerCase(),
    extension: lang.extensions[0],
  }))
})

const selectedLanguage = ref<{ label: string; value: string; extension: string } | null>(null)

const filteredLanguages = computed(() => {
  return languageOptions.value
})

type LanguageOption = { label: string; value: string; extension: string }

const handleLanguageSelect = (event: ListboxItemSelectEvent<AcceptableValue>) => {
  const language = event.detail?.value as LanguageOption | undefined
  if (language) {
    selectedLanguage.value = language
    formData.value.language = language.value
  }
}

// Reset form data when modal closes
const resetForm = () => {
  formData.value = {
    title: '',
    content: '',
    language: '',
  }
  selectedLanguage.value = null
}

// Update form data when snippet changes (edit mode) or when modal opens (create mode)
watch(
  [() => props.snippet, () => props.show],
  ([newSnippet, showModal]) => {
    if (showModal) {
      if (newSnippet && isEditMode.value) {
        // Edit mode: populate with existing snippet data
        formData.value = {
          title: newSnippet.title,
          content: newSnippet.content,
          language: newSnippet.language.toLowerCase(),
        }

        // Set the selected language for the combobox
        const matchedLanguage = languageOptions.value.find(
          (lang) => lang.value === newSnippet.language.toLowerCase(),
        )
        selectedLanguage.value = matchedLanguage || null
      } else {
        // Create mode: reset form
        resetForm()
      }
    }
  },
  { immediate: true },
)

const handleSubmit = () => {
  const submitData: { title: string; content: string; language?: string } = {
    title: formData.value.title,
    content: formData.value.content,
  }

  // Only include language if it's set or we're in edit mode
  if (formData.value.language || isEditMode.value) {
    submitData.language = formData.value.language
  }

  emit('submit', submitData)
  emit('close')
}

const handleCancel = () => {
  resetForm()
  emit('close')
}
</script>
