<template>
  <div class="bg-background min-h-screen">
    <div class="container mx-auto max-w-7xl px-8">
      <!-- Navigation -->
      <div class="mb-8">
        <Button variant="outline" @click="router.back()" class="gap-2">
          <ArrowLeftIcon class="h-4 w-4" />
          Back to snippets
        </Button>
      </div>

      <!-- Loading state -->
      <div v-if="isPending" class="space-y-8 md:grid md:grid-cols-[30%_70%] md:space-x-8">
        <!-- Left column skeleton -->
        <Card class="p-8">
          <div class="space-y-6">
            <!-- Title skeleton -->
            <div class="space-y-3">
              <div class="bg-muted h-8 w-3/4 animate-pulse rounded-lg"></div>
              <div class="bg-muted h-6 w-1/3 animate-pulse rounded-lg"></div>
            </div>

            <!-- Author skeleton -->
            <div class="flex items-center gap-4 border-t pt-4">
              <div class="bg-muted h-12 w-12 animate-pulse rounded-full"></div>
              <div class="space-y-2">
                <div class="bg-muted h-4 w-24 animate-pulse rounded"></div>
                <div class="bg-muted h-3 w-16 animate-pulse rounded"></div>
              </div>
            </div>

            <!-- Stats skeleton -->
            <div class="flex gap-6 pt-4">
              <div class="bg-muted h-4 w-16 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-16 animate-pulse rounded"></div>
            </div>

            <!-- Actions skeleton -->
            <div class="flex gap-2 pt-4">
              <div class="bg-muted h-10 w-16 animate-pulse rounded"></div>
              <div class="bg-muted h-10 w-16 animate-pulse rounded"></div>
            </div>
          </div>
        </Card>

        <!-- Right column skeleton -->
        <Card class="p-8">
          <div class="space-y-4">
            <div class="flex items-center justify-between border-b pb-4">
              <div class="bg-muted h-5 w-32 animate-pulse rounded"></div>
              <div class="bg-muted h-8 w-20 animate-pulse rounded"></div>
            </div>
            <div class="space-y-2">
              <div class="bg-muted h-4 w-full animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-5/6 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-4/5 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-3/4 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-2/3 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-3/4 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-5/6 animate-pulse rounded"></div>
              <div class="bg-muted h-4 w-1/2 animate-pulse rounded"></div>
            </div>
          </div>
        </Card>
      </div>

      <!-- Error state -->
      <Card v-else-if="isError" class="border-destructive bg-destructive/5 p-8">
        <div class="space-y-4 text-center">
          <div
            class="bg-destructive/10 mx-auto flex h-16 w-16 items-center justify-center rounded-full"
          >
            <AlertCircleIcon class="text-destructive h-8 w-8" />
          </div>
          <div>
            <h2 class="text-destructive mb-2 text-xl font-bold">Something went wrong</h2>
            <p class="text-destructive/80">
              {{ error?.message || 'An unexpected error occurred while loading the snippet.' }}
            </p>
          </div>
          <Button variant="outline" @click="router.back()" class="mt-4"> Go back </Button>
        </div>
      </Card>

      <!-- Snippet content -->
      <div v-else-if="snippet" class="space-y-8">
        <!-- Title - Full width -->
        <div class="text-start">
          <h1 class="text-foreground text-4xl font-bold tracking-tight break-words">
            {{ snippet.title }}
          </h1>
        </div>

        <!-- Two column layout -->
        <div class="space-y-8 md:grid md:grid-cols-[30%_70%] md:space-x-8">
          <!-- Left column: Metadata and actions -->
          <Card class="h-fit p-8">
            <div class="space-y-6">
              <!-- Author info -->
              <UserAvatar :user="snippet.author" :subtitle="dayjs(snippet.createdAt).fromNow()" />

              <Separator />

              <!-- Metadata -->
              <div class="space-y-4">
                <!-- Language -->
                <div v-if="snippet.language" class="text-muted-foreground text-sm">
                  <span class="text-foreground font-medium">Language:</span>
                  {{ getLanguageName(snippet.language) || 'Plain Text' }}
                </div>

                <!-- Creation date -->
                <div class="text-muted-foreground text-sm">
                  <span class="text-foreground font-medium">Created:</span>
                  {{ dayjs(snippet.createdAt).format('MMMM D, YYYY') }}
                </div>

                <!-- Stats -->
                <div
                  class="text-muted-foreground flex items-center gap-2 text-sm md:flex-col md:items-start lg:flex-row lg:gap-6"
                >
                  <div class="flex items-center gap-1">
                    <HeartIcon class="h-4 w-4" />
                    <span>{{ snippet.likes }} {{ snippet.likes === 1 ? 'like' : 'likes' }}</span>
                  </div>
                  <div class="flex items-center gap-1">
                    <EyeIcon class="h-4 w-4" />
                    <span>{{ Math.floor(Math.random() * 100) + 20 }} views</span>
                  </div>
                </div>
              </div>

              <Separator />

              <!-- Actions -->
              <div class="space-y-3">
                <!--h3 class="text-foreground text-sm font-medium">Actions</h3-->

                <!-- Save and Like buttons -->
                <Authenticated>
                  <div class="grid grid-cols-2 gap-2">
                    <SaveButton
                      variant="outline"
                      size="sm"
                      class="flex-1"
                      :isSaved="snippet.isSaved"
                      :snippetId="snippet.id"
                    />
                    <LikeButton
                      variant="outline"
                      size="sm"
                      class="flex-1"
                      :likes="snippet.likes"
                      :isLiked="snippet.isLiked"
                      :snippetId="snippet.id"
                      :hideCount="true"
                    />

                    <!-- Author-only actions -->

                    <IsAuthor :authorId="snippet.author.id">
                      <Button
                        variant="outline"
                        size="sm"
                        class="flex-1"
                        @click="showEditModal = true"
                      >
                        <EditIcon class="size-4 shrink-0" />
                        <span>Edit</span>
                      </Button>
                    </IsAuthor>
                    <IsAuthor :authorId="snippet.author.id">
                      <Button
                        variant="destructive"
                        size="sm"
                        class="flex-1"
                        :disabled="isDeleting"
                        @click="deleteSnippet()"
                      >
                        <Trash2Icon class="size-4 shrink-0" />
                        <span>Delete</span>
                      </Button>
                    </IsAuthor>
                  </div>
                </Authenticated>

                <!-- Share button -->
                <Button variant="outline" size="sm" @click="shareSnippet" class="w-full">
                  <ShareIcon class="size-4 shrink-0" />
                  <span class="ml-2 hidden sm:inline">Share</span>
                </Button>
              </div>
            </div>
          </Card>

          <!-- Right column: Code display -->
          <Card class="min-w-0 gap-0 overflow-hidden py-0">
            <div class="bg-muted/30 flex items-center justify-between border-b px-6 py-4">
              <!-- File header -->
              <div class="flex items-center gap-3">
                <div class="flex items-center gap-2">
                  <div class="h-3 w-3 rounded-full bg-red-500"></div>
                  <div class="h-3 w-3 rounded-full bg-yellow-500"></div>
                  <div class="h-3 w-3 rounded-full bg-green-500"></div>
                </div>
                <span class="text-muted-foreground font-mono text-sm">
                  snippet.{{ getLanguageExtension(snippet.language) || 'txt' }}
                </span>
              </div>

              <!-- Copy button -->
              <Button
                variant="ghost"
                size="sm"
                @click="copySnippetToClipboard"
                class="text-muted-foreground hover:text-foreground gap-2"
              >
                <CopyIcon class="h-4 w-4" />
                Copy
              </Button>
            </div>

            <!-- Code content -->
            <CardContent class="min-w-0 overflow-y-auto p-6">
              <pre
                class="text-foreground/90 overflow-x-auto font-mono text-sm leading-relaxed"
              ><code>{{ snippet.content }}</code></pre>
            </CardContent>
          </Card>
        </div>
      </div>

      <!-- Edit Modal -->
      <SnippetFormModal
        v-if="snippet"
        :show="showEditModal"
        :is-loading="isUpdating"
        :snippet="snippet"
        @close="showEditModal = false"
        @submit="updateSnippet"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeftIcon,
  HeartIcon,
  EyeIcon,
  CopyIcon,
  ShareIcon,
  AlertCircleIcon,
  Trash2Icon,
  EditIcon,
} from 'lucide-vue-next'
import UserAvatar from '@/components/UserAvatar.vue'
import { getLanguageName, getLanguageExtension } from '@/lib/languages'
import { copyToClipboard as copyTextToClipboard } from '@/lib/utils'
import { toast } from 'vue-sonner'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import CardContent from '@/components/ui/card/CardContent.vue'
import SnippetFormModal from './_components/SnippetFormModal.vue'

dayjs.extend(relativeTime)

const route = useRoute()
const router = useRouter()

const queryClient = useQueryClient()

const {
  data: snippet,
  isPending,
  isError,
  error,
} = useQuery({
  queryKey: ['snippet', route.params.snippetId],
  queryFn: () => snippetsService.getSnippet(route.params.snippetId as string),
})

// Delete mutation
const { mutate: deleteSnippet, isPending: isDeleting } = useMutation({
  mutationFn: () => snippetsService.deleteSnippet(route.params.snippetId as string),
  onSuccess: () => {
    toast.success('Snippet deleted successfully')
    // Invalidate and remove queries
    queryClient.removeQueries({ queryKey: ['snippet', route.params.snippetId] })
    queryClient.invalidateQueries({ queryKey: ['snippets'] })
    queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
    // Navigate back
    router.go(-1)
  },
  onError: (error) => {
    toast.error(error.message || 'Failed to delete snippet')
  },
})

// Update mutation
const { mutate: updateSnippet, isPending: isUpdating } = useMutation({
  mutationFn: (formData: { title: string; content: string; language?: string }) =>
    snippetsService.updateSnippet(route.params.snippetId as string, {
      title: formData.title,
      content: formData.content,
      language: formData.language || snippet.value?.language || '',
    }),
  onSuccess: (updatedSnippet) => {
    toast.success('Snippet updated successfully')
    // Update queries
    queryClient.setQueryData(['snippet', updatedSnippet.id], updatedSnippet)
    queryClient.invalidateQueries({ queryKey: ['snippets'] })
    queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
  },
  onError: (error) => {
    toast.error(error.message || 'Failed to update snippet')
  },
})

const copySnippetToClipboard = async () => {
  if (!snippet.value) return

  try {
    await copyTextToClipboard(snippet.value.content)
    toast.success('Code copied to clipboard!')
  } catch (err) {
    console.error('Failed to copy code:', err)
    toast.error('Failed to copy code to clipboard')
  }
}

const shareSnippet = async () => {
  const url = window.location.href

  if (navigator.share) {
    try {
      await navigator.share({
        title: snippet.value?.title,
        text: `Check out this code snippet: ${snippet.value?.title}`,
        url: url,
      })
      return // Successfully shared, don't fall back
    } catch (err: any) {
      // Only fall back to clipboard if sharing is not supported
      // Don't fall back if user just cancelled (AbortError)
      if (err.name === 'AbortError') {
        return // User cancelled, do nothing
      }
      console.log('Native share failed, falling back to clipboard')
    }
  }

  // Fallback: copy to clipboard
  try {
    await copyTextToClipboard(url)
    toast.success('Link copied to clipboard!')
  } catch (err) {
    console.error('Failed to copy link:', err)
    toast.error('Failed to copy link')
  }
}

const showEditModal = ref(false)

onMounted(() => {
  window.scrollTo({ top: 0, behavior: 'auto' })
})
</script>
