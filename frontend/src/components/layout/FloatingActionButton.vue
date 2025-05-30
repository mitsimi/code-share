<template>
  <Button class="fixed right-8 bottom-8 z-10 flex size-16" @click="$emit('click')">
    <Plus class="size-8" />
  </Button>
</template>

<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

defineEmits<{
  (e: 'click'): void
}>()
</script>
