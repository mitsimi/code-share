<template>
  <div>
    <div v-if="isLoading" class="flex justify-center py-8">
      <LoaderCircleIcon class="size-8 animate-spin" />
    </div>
    <div v-else-if="snippets.length === 0" class="py-8 text-center">
      <p class="text-muted-foreground">{{ emptyMessage }}</p>
    </div>
    <div v-else class="grid gap-4">
      <SnippetCard
        v-for="snippet in snippets"
        :key="snippet.id"
        :snippet="snippet"
        @click="router.push(`/snippets/${snippet.id}`)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { LoaderCircleIcon } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import type { Snippet } from '@/types'

const router = useRouter()

defineProps<{
  snippets: Snippet[]
  isLoading: boolean
  emptyMessage: string
}>()
</script>
