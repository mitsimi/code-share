<template>
  <main class="mx-auto my-12 max-w-3xl px-4">
    <div class="rounded-lg border-4 border-black bg-white p-8 shadow-[8px_8px_0_0_#000]">
      <h1 class="mb-6 text-3xl font-bold">About CodeShare</h1>

      <div class="space-y-6">
        <div>
          <p>
            CodeShare is a platform for developers to share and discover useful code snippets.
            Whether you're looking for a quick solution or want to share your coding knowledge with
            others, CodeShare makes it easy and fun.
          </p>
        </div>

        <div>
          <h2 class="mb-4 border-b-4 border-black pb-2 text-xl font-bold">Course Information</h2>
          <p class="font-mono">
            This platform was created for the course
            <a
              href="https://ssw.jku.at/Teaching/Lectures/IntroWebDev/2025/"
              class="text-accent font-bold hover:underline"
            >
              Special Topic: Introduction to Full Stack Web Development
            </a>
            at Johannes Kepler University.
          </p>
        </div>

        <div>
          <h2 class="mb-4 border-b-4 border-black pb-2 text-xl font-bold">Features</h2>
          <ul class="list-inside list-disc space-y-2 pl-4">
            <li>Share code snippets with the community</li>
            <li>Like and save your favorite snippets</li>
            <li>Clean, modern interface with a playful design</li>
            <li>Instant feedback with toast notifications</li>
          </ul>
        </div>

        <div>
          <h2 class="mb-4 border-b-4 border-black pb-2 text-xl font-bold">Technology Stack</h2>
          <p>
            Built with Vue 3, TypeScript, and Tailwind CSS. The retro-inspired design features crisp
            black borders and playful shadows to create an engaging user experience.
          </p>
        </div>

        <!--div>
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
              About This Project
            </router-link>
          </div>
        </div-->
      </div>
    </div>
  </main>
</template>
