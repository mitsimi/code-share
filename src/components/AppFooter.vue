<template>
  <footer class="bg-accent mt-12 border-t-4 border-black">
    <div class="mx-auto max-w-7xl px-4 py-12">
      <div class="grid gap-8 md:grid-cols-2">
        <!-- Left Section -->
        <div class="card bg-background rounded-lg p-6">
          <h2 class="mb-4 border-b-4 border-black pb-2 text-xl font-bold">About CodeShare</h2>
          <p class="font-mono">
            A platform for sharing and discovering code snippets, created for the course
            <a
              href="https://ssw.jku.at/Teaching/Lectures/IntroWebDev/2025/"
              class="text-accent font-bold hover:underline"
            >
              Special Topic: Introduction to Full Stack Web Development
            </a>
            at Johannes Kepler University.
          </p>
        </div>

        <!-- Right Section -->
        <div class="card bg-background rounded-lg p-6">
          <h2 class="mb-4 border-b-4 border-black pb-2 text-xl font-bold">Quick Links</h2>
          <div class="flex flex-col gap-2">
            <router-link
              to="/snippets"
              class="neobrutalism bg-background hover:bg-accent hover:text-accent-foreground rounded-lg px-4 py-2 font-mono font-bold"
            >
              Browse Snippets
            </router-link>
            <router-link
              to="/about"
              class="neobrutalism bg-background hover:bg-accent hover:text-accent-foreground rounded-lg px-4 py-2 font-mono font-bold"
            >
              About Page
            </router-link>
          </div>
        </div>
      </div>

      <!-- Bottom Section -->
      <div class="text-secondary-foreground/80 mt-8 text-center font-mono text-sm">
        © 2025 CodeShare. All rights reserved.
      </div>
    </div>
  </footer>
</template>
