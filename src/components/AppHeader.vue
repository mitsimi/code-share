<template>
  <header class="sticky top-0 z-50 border-b-4 border-black bg-white">
    <div class="container mx-auto flex items-center justify-between px-4 py-4">
      <div class="flex items-center gap-4">
        <img
          src="../assets/logo.svg"
          alt="CodeShare Logo"
          class="h-12 w-12"
          @error="handleLogoError"
        />
        <h1 class="font-mono text-2xl font-bold">CodeShare</h1>
      </div>
      <nav class="flex items-center gap-4">
        <router-link
          to="/snippets"
          class="rounded-lg border-4 border-black px-4 py-2 font-mono font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
          active-class="bg-primary"
        >
          Snippets
        </router-link>
        <router-link
          to="/about"
          class="rounded-lg border-4 border-black px-4 py-2 font-mono font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
          active-class="bg-primary"
        >
          About
        </router-link>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
const handleLogoError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
}
</script>
