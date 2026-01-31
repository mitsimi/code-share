<template>
  <Button
    :variant="variant"
    @click.stop="authStore.isAuthenticated() && handleSave()"
    :class="[
      { 'pointer-events-none': !authStore.isAuthenticated() },
      isSaved
        ? 'text-primary border-primary hover:bg-primary/10'
        : 'text-muted-foreground hover:bg-secondary/10',
    ]"
  >
    <template v-if="isPending">
      <LoaderCircleIcon class="size-4 animate-spin" />
    </template>
    <template v-else>
      <BookmarkIcon
        class="size-4 transition-transform duration-200"
        :class="{ 'scale-110': isSaved }"
        :fill="isSaved ? 'currentColor' : 'none'"
        stroke-width="2"
      />
    </template>
  </Button>
</template>

<script setup lang="ts">
import { BookmarkIcon, LoaderCircleIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useSaveSnippet } from '@/composables/useSaveSnippet'
import type { ButtonProps } from '@/components/ui/button'
import { Button } from '@/components/ui/button'

const authStore = useAuthStore()
const { mutate: toggleSave, isPending } = useSaveSnippet()

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

const handleSave = () => {
  if (!props.snippetId || isPending.value) return

  const action = props.isSaved ? 'unsave' : 'save'
  toggleSave({
    snippetId: props.snippetId,
    action,
    currentIsSaved: props.isSaved,
  })
}
</script>
