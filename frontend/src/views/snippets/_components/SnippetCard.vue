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
        <Badge v-if="snippet.language" variant="secondary" class="shrink-0 text-xs">
          {{ getLanguageName(snippet.language) || 'Text' }}
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
              snippet.{{ getLanguageExtension(snippet.language || '') || 'txt' }}
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
            <Avatar class="ring-0">
              <AvatarImage src="" :alt="snippet.author" />
              <AvatarFallback class="bg-primary text-primary-foreground">{{
                snippet.author.username[0].toUpperCase()
              }}</AvatarFallback>
            </Avatar>
            <div class="flex flex-col">
              <span class="text-foreground text-sm font-medium">{{ snippet.author.username }}</span>
              <span class="text-muted-foreground text-xs">
                {{ dayjs(snippet.createdAt).fromNow() }}
              </span>
            </div>
          </div>
        </div>

        <!-- Right side - Action buttons -->
        <div class="flex items-center gap-2">
          <!-- Save/Bookmark button -->
          <div @click.stop>
            <SaveButton :isSaved="snippet.isSaved" :snippetId="snippet.id" />
          </div>

          <!-- Like button -->
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
import { CopyIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import LikeButton from './LikeButton.vue'
import SaveButton from './SaveButton.vue'
import { toast } from 'vue-sonner'
import { getLanguageExtension, getLanguageName } from '@/utils/languages'

dayjs.extend(relativeTime)

defineEmits<{
  (e: 'click'): void
}>()

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

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(props.snippet.content)
    toast.success('Code copied to clipboard!')
  } catch (err) {
    console.error('Failed to copy code:', err)
  }
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
