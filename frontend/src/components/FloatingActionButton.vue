<template>
  <Button class="fixed right-8 bottom-8 flex size-16" @click="$emit('click')">
    <Plus class="size-8" />
  </Button>
</template>

<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import Button from './ui/button/Button.vue'

defineEmits<{
  (e: 'click'): void
}>()
</script>
