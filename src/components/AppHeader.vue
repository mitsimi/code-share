<template>
  <header class="border-b-4 border-black bg-white">
    <div class="mx-auto flex max-w-7xl items-center justify-between px-4 py-6">
      <div class="flex items-center gap-2">
        <img
          src="../assets/logo.svg"
          @error="handleLogoError"
          alt="SnippetShare Logo"
          class="h-8 w-8"
        />
        <h1 class="text-3xl font-bold">SnippetShare</h1>
      </div>
      <nav>
        <ul class="flex gap-4">
          <li>
            <a
              href="#"
              class="rounded-lg border-4 border-black bg-white px-4 py-2 font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
              >Home</a
            >
          </li>
          <li>
            <a
              href="#"
              class="rounded-lg border-4 border-black bg-white px-4 py-2 font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
              >About</a
            >
          </li>
        </ul>
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
