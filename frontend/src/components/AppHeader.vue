<template>
  <header class="bg-muted sticky top-0 z-50 border-b-2">
    <div class="container mx-auto px-4 py-4">
      <!-- Desktop Navigation -->
      <div class="hidden sm:flex sm:items-center sm:justify-between">
        <!-- Logo and Title -->
        <div class="flex items-center gap-4">
          <div class="bg-background rounded-lg border-2 p-2 shadow">
            <img
              src="../assets/logo.svg"
              alt="CodeShare Logo"
              class="size-12"
              @error="handleLogoError"
            />
          </div>
          <div class="flex flex-col">
            <h1 class="text-accent-foreground font-mono text-2xl font-bold">CodeShare</h1>
            <p class="text-accent-foreground/80 text-sm">Share your code snippets</p>
          </div>
        </div>

        <!-- Desktop Navigation Links -->
        <nav class="flex items-center gap-2">
          <Button variant="link" @click="router.push('/snippets')"> Snippets </Button>

          <Button variant="link" @click="router.push('/about')"> About </Button>

          <template v-if="authStore.isAuthenticated()">
            <Button variant="destructive" @click="handleLogout"> Logout </Button>
          </template>
          <template v-else>
            <Button variant="boring" @click="router.push('/login')"> Login </Button>
            <Button @click="router.push('/signup')"> Signup </Button>
          </template>
        </nav>
      </div>

      <!-- Mobile Navigation -->
      <div class="sm:hidden">
        <div class="flex items-center justify-between">
          <!-- Logo and Title -->
          <div class="flex items-center gap-3">
            <div class="neobrutalism bg-background rounded-lg p-2">
              <img
                src="../assets/logo.svg"
                alt="CodeShare Logo"
                class="h-8 w-8"
                @error="handleLogoError"
              />
            </div>
            <h1 class="text-accent-foreground font-mono text-xl font-bold">CodeShare</h1>
          </div>

          <!-- Mobile Menu Button -->
          <Button
            variant="ghost"
            size="icon"
            @click="isMenuOpen = !isMenuOpen"
            class="relative"
            aria-label="Toggle menu"
          >
            <Menu v-if="!isMenuOpen" class="size-6" />
            <X v-else class="size-6" />
          </Button>
        </div>

        <!-- Mobile Menu -->
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="transform -translate-y-4 opacity-0"
          enter-to-class="transform translate-y-0 opacity-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="transform translate-y-0 opacity-100"
          leave-to-class="transform -translate-y-4 opacity-0"
        >
          <nav v-if="isMenuOpen" class="mt-4 flex flex-col gap-2">
            <router-link to="/snippets" @click="isMenuOpen = false"> Snippets </router-link>

            <router-link to="/about" @click="isMenuOpen = false"> About </router-link>

            <template v-if="authStore.isAuthenticated()">
              <Button variant="destructive" @click="handleLogout"> Logout </Button>
            </template>
            <template v-else>
              <router-link to="/login" @click="isMenuOpen = false"> Login </router-link>
              <router-link to="/signup" @click="isMenuOpen = false"> Sign Up </router-link>
            </template>
          </nav>
        </Transition>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Menu, X } from 'lucide-vue-next'
import { Button } from './ui/button'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const isMenuOpen = ref(false)

const handleLogoError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>
