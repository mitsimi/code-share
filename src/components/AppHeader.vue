<template>
  <header class="bg-secondary p-4">
    <div
      class="bg-background mx-auto flex max-w-7xl -rotate-[0.5deg] items-center justify-between border-4 border-black px-8 py-4 shadow-[5px_5px_0_0_#000]"
    >
      <div class="flex items-center gap-2">
        <img
          src="../assets/logo.svg"
          @error="handleLogoError"
          alt="SnippetShare Logo"
          class="size-8 w-auto"
        />
        <span class="text-2xl font-black text-black">SnippetShare</span>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const handleLogoError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.src = 'https://img.logoipsum.com/284.svg'
}
</script>
