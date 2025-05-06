<template>
  <button
    class="bg-accent fixed right-8 bottom-8 z-99 flex h-[60px] w-[60px] cursor-pointer items-center justify-center border-4 border-black text-white shadow-[5px_5px_0_0_#000] transition-all duration-200 hover:-translate-y-1 hover:shadow-[7px_7px_0_0_#000] active:translate-x-0.5 active:translate-y-0.5 active:shadow-[3px_3px_0_0_#000]"
    @click="$emit('click')"
  >
    <LucidePlus class="size-8" />
  </button>
</template>

<script setup lang="ts">
import { LucidePlus } from 'lucide-vue-next'

defineEmits<{
  (e: 'click'): void
}>()
</script>

<style scoped>
.bg-accent {
  background-color: #ff3e4d;
}
</style>
