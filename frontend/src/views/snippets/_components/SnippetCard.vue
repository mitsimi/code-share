<template>
  <Card class="group relative min-h-full cursor-pointer overflow-hidden" @click="$emit('click')">
    <!-- Header with title and language tag -->
    <CardHeader class="pb-3">
      <div class="flex items-start justify-between gap-3">
        <CardTitle
          class="group-hover:text-primary line-clamp-2 text-lg leading-tight font-bold break-words transition-colors"
        >
          {{ snippet.title }}
        </CardTitle>
        <Badge variant="secondary" class="shrink-0 text-xs">
          <!-- TODO: Add language property to snippet type -->
          {{ hardcodedLanguage }}
        </Badge>
      </div>
    </CardHeader>

    <!-- Code content with syntax highlighting preview -->
    <CardContent class="flex-grow px-6 pb-4">
      <div class="bg-muted/30 relative overflow-hidden rounded-md border">
        <!-- Code header bar -->
        <div class="bg-muted/50 flex items-center justify-between border-b px-3 py-2">
          <div class="flex items-center gap-2">
            <span class="text-muted-foreground font-mono text-xs">
              <!-- TODO: Add file extension property to snippet type -->
              snippet.{{ hardcodedFileExtension }}
            </span>
          </div>
          <Button
            variant="ghost"
            size="sm"
            class="h-6 w-6 p-0 opacity-0 transition-opacity group-hover:opacity-100"
            @click.stop="copyToClipboard"
          >
            <CopyIcon class="h-3 w-3" />
          </Button>
        </div>

        <!-- Code content -->
        <div class="relative max-h-32 overflow-x-clip overflow-y-auto p-3">
          <pre
            class="text-foreground/90 font-mono text-xs leading-relaxed"
          ><code>{{ truncatedContent }}</code></pre>
          <!-- Fade overlay for long content -->
          <div
            v-if="snippet.content.length > 200"
            class="from-muted/30 absolute right-0 bottom-0 left-0 h-8 bg-gradient-to-t to-transparent"
          ></div>
        </div>
      </div>
    </CardContent>

    <!-- Footer with interactions -->
    <CardFooter class="px-6 pt-0">
      <div class="flex w-full items-center justify-between">
        <!-- Left side - Tags and stats -->
        <div class="text-muted-foreground flex items-center gap-4 text-xs">
          <!-- Author info with avatar -->
          <div class="mt-3 flex items-center gap-2">
            <div
              class="bg-primary text-primary-foreground flex h-8 w-8 items-center justify-center rounded-full text-xs font-medium"
            >
              <!-- TODO: Add author avatar or generate from author name -->
              {{ getAuthorInitials(snippet.author) }}
            </div>
            <div class="flex flex-col">
              <span class="text-foreground text-sm font-medium">{{ snippet.author }}</span>
              <span class="text-muted-foreground text-xs">
                <!-- TODO: Add createdAt property to snippet type -->
                {{ dayjs(snippet.createdAt).fromNow() }}
              </span>
            </div>
          </div>
        </div>

        <!-- Right side - Action buttons -->
        <div class="flex items-center gap-2">
          <!-- Save/Bookmark button -->
          <Button
            variant="ghost"
            size="sm"
            class="h-8 w-8 p-0"
            :class="{ 'pointer-events-none': !authStore.isAuthenticated() }"
            @click.stop="authStore.isAuthenticated() && toggleSave()"
          >
            <BookmarkIcon
              :class="[
                'h-4 w-4 transition-colors',
                hardcodedIsSaved
                  ? 'text-primary fill-current'
                  : 'text-muted-foreground hover:text-foreground',
              ]"
            />
          </Button>

          <!-- Like button - your existing component -->
          <div @click.stop>
            <LikeButton :likes="snippet.likes" :isLiked="snippet.isLiked" :snippetId="snippet.id" />
          </div>
        </div>
      </div>
    </CardFooter>

    <!-- Hover overlay for better interaction feedback -->
    <div
      class="bg-primary/5 pointer-events-none absolute inset-0 rounded-lg opacity-0 transition-opacity group-hover:opacity-100"
    ></div>
  </Card>
</template>

<script setup lang="ts">
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import type { Snippet } from '@/types'
import { CopyIcon, BookmarkIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import LikeButton from './LikeButton.vue'
import { toast } from 'vue-sonner'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

dayjs.extend(relativeTime)

defineEmits<{
  (e: 'click'): void
}>()

// TODO: Replace these hardcoded values with actual props when available
const hardcodedLanguage = 'JavaScript' // Add to snippet type: language: string
const hardcodedFileExtension = 'js' // Add to snippet type: filename?: string
const hardcodedIsSaved = false // Add to snippet type or user state: isSaved: boolean

const props = defineProps<{
  snippet: Snippet
}>()

// Computed properties
const truncatedContent = computed(() => {
  const maxLength = 200
  if (props.snippet.content.length <= maxLength) {
    return props.snippet.content
  }
  return props.snippet.content.substring(0, maxLength) + '...'
})

// Helper functions
const getAuthorInitials = (authorName: string): string => {
  return authorName
    .split(' ')
    .map((name) => name.charAt(0))
    .join('')
    .toUpperCase()
    .substring(0, 2)
}

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(props.snippet.content)
    toast.success('Code copied to clipboard!')
  } catch (err) {
    console.error('Failed to copy code:', err)
    // TODO: Add error toast
    // You can use: toast.error('Failed to copy code')
  }
}

const toggleSave = () => {
  // TODO: Implement save/bookmark functionality similar to LikeButton
  // You might want to create a SaveButton component following the same pattern
  console.log('Toggle save for snippet:', props.snippet.id)
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
