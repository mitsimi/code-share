<template>
  <Button
    :variant="variant"
    @click.stop="authStore.isAuthenticated() && toggleSave()"
    :class="[
      { 'pointer-events-none': !authStore.isAuthenticated() },
      isSaved
        ? 'text-primary border-primary hover:bg-primary/10'
        : 'text-muted-foreground hover:bg-secondary/10',
    ]"
  >
    <template v-if="isLoading">
      <LoaderCircleIcon class="size-4 animate-spin" />
    </template>
    <template v-else>
      <BookmarkIcon
        class="size-4 transition-transform duration-200"
        :class="{ 'scale-110': isSaved }"
        :fill="isSaved ? 'currentColor' : 'none'"
        :stroke="isSaved ? 'currentColor' : 'currentColor'"
        stroke-width="2"
      />
    </template>
  </Button>
</template>

<script setup lang="ts">
import { BookmarkIcon, LoaderCircleIcon } from 'lucide-vue-next'
import { type Snippet } from '@/types'
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { toast } from 'vue-sonner'
import type { ButtonProps } from '@/components/ui/button'
import { useSnippetsStore } from '@/stores/snippets'

const authStore = useAuthStore()
const snippetsStore = useSnippetsStore()
const queryClient = useQueryClient()

const props = withDefaults(
  defineProps<{
    snippetId: string
    isSaved: boolean
    variant?: ButtonProps['variant']
  }>(),
  {
    variant: 'ghost',
  },
)

const isLoading = ref(false)

const { mutate: updateSave } = useMutation<
  Snippet,
  Error,
  { snippetId: string; action: 'save' | 'unsave' }
>({
  mutationKey: ['saveMutation', props.snippetId],
  mutationFn: async ({ snippetId, action }) => {
    isLoading.value = true
    try {
      return await snippetsService.toggleSave(snippetId, action)
    } finally {
      isLoading.value = false
    }
  },
  onSuccess: (updatedSnippet) => {
    queryClient.setQueryData(['snippet', updatedSnippet.id], updatedSnippet)

    snippetsStore.updateSnippet(updatedSnippet.id, {
      isSaved: updatedSnippet.isSaved,
    })

    queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
  },
  onError: (error, { snippetId, action }) => {
    const revertAction = action === 'save' ? 'unsave' : 'save'
    snippetsStore.handleUserAction({
      action: revertAction,
      snippet_id: snippetId,
      value: true,
    })

    console.error('Like mutation failed:', error)
    toast.error(error.message || 'Please try again')
  },
})

const toggleSave = () => {
  if (!props.snippetId) {
    console.error('Cannot toggle save: snippetId is missing')
    return
  }

  if (isLoading.value) {
    return
  }

  const action = props.isSaved ? 'unsave' : 'save'
  updateSave({ snippetId: props.snippetId, action })
}
</script>
