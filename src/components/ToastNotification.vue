<template>
  <div class="fixed bottom-8 left-8 z-1000">
    <div
      class="flex max-w-[300px] min-w-[200px] transform items-center border-4 border-black bg-[#4caf50] p-4 text-white shadow-[5px_5px_0_0_#000] transition-transform duration-300"
      :class="{ 'translate-y-0 rotate-2': show, 'translate-y-[200%] rotate-3': !show }"
    >
      <div class="flex items-center gap-3">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="shrink-0"
        >
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
        <span>{{ message }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  show: boolean
  message: string
}>()
</script>
