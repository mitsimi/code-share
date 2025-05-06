<template>
  <button
    class="bg-accent fixed right-8 bottom-8 flex h-16 w-16 items-center justify-center rounded-lg border-4 border-black text-white shadow-[8px_8px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
    @click="$emit('click')"
  >
    <Plus class="h-8 w-8" />
  </button>
</template>

<script setup lang="ts">
import { Plus } from 'lucide-vue-next'

defineEmits<{
  (e: 'click'): void
}>()
</script>

<style scoped>
.bg-accent {
  background-color: #ff3e4d;
}
</style>
