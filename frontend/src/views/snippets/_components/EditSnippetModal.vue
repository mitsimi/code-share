<template>
  <Dialog :open="show" @update:open="$emit('close')">
    <DialogContent class="sm:max-w-[600px]">
      <DialogHeader>
        <DialogTitle>Edit Snippet</DialogTitle>
        <DialogDescription>
          Update your code snippet. Make sure to provide a clear title and description.
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

              <ComboboxGroup>
                <ComboboxItem
                  v-for="language in filteredLanguages"
                  :key="language.value"
                  :value="language"
                  @select="handleLanguageSelect"
                  class="flex items-center justify-between"
                >
                  <div class="flex items-center gap-2">
                    {{ language.label }}
                    <span class="text-muted-foreground text-xs">.{{ language.extension }}</span>
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
          <Button type="button" variant="outline" @click="$emit('close')"> Cancel </Button>
          <Button type="submit" variant="reverse" :disabled="isLoading">
            <LoaderCircleIcon v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
            Save Changes
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
import { getAvailableLanguages, getLanguageName } from '@/lib/languages'

const props = defineProps<{
  show: boolean
  isLoading: boolean
  snippet: Snippet
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'submit', data: { title: string; content: string; language: string }): void
}>()

const formData = ref({
  title: '',
  content: '',
  language: '',
})

// Transform languages to label/value format
const languageOptions = computed(() => {
  return getAvailableLanguages().map((lang) => ({
    label: lang.name,
    value: lang.name.toLowerCase(),
    extension: lang.extension,
    aliases: lang.aliases,
  }))
})

const selectedLanguage = ref<{ label: string; value: string; extension: string } | null>(null)

const filteredLanguages = computed(() => {
  return languageOptions.value
})

const handleLanguageSelect = (event: any) => {
  const language = event.detail.value
  selectedLanguage.value = language
  formData.value.language = language.value
}

// Update form data when snippet changes
watch(
  () => props.snippet,
  (newSnippet) => {
    if (newSnippet) {
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
    }
  },
  { immediate: true },
)

const handleSubmit = () => {
  emit('submit', {
    title: formData.value.title,
    content: formData.value.content,
    language: getLanguageName(formData.value.language) || formData.value.language,
  })
  emit('close')
}
</script>
