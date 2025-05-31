<template>
  <header
    class="bg-background/95 supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50 w-full border-b backdrop-blur"
  >
    <div class="container mx-auto flex h-14 items-center px-4">
      <!-- Logo/Brand -->
      <div class="mr-4 flex">
        <router-link to="/" class="mr-6 flex items-center space-x-2">
          <CodepenIcon />
          <span class="hidden font-bold sm:inline-block">CodeShare</span>
        </router-link>
      </div>

      <!-- Desktop Navigation -->
      <nav class="hidden items-center space-x-6 text-sm font-medium md:flex">
        <router-link
          v-for="link in navigationLinks"
          :key="link.href"
          :to="link.href"
          class="hover:text-foreground/80 text-foreground/60 transition-colors"
          active-class="text-primary"
        >
          {{ link.label }}
        </router-link>
      </nav>

      <!-- Right Side Actions -->
      <div class="flex flex-1 items-center justify-end space-x-2">
        <!-- Desktop Auth Buttons -->
        <div class="hidden items-center space-x-2 md:flex">
          <template v-if="!authStore.isAuthenticated()">
            <Button variant="ghost" size="sm" @click="handleLogin"> Login </Button>
            <Button variant="default" size="sm" @click="handleSignup"> Sign Up </Button>
          </template>
          <template v-else>
            <Button variant="ghost" size="sm" @click="handleLogout">
              <LogOut class="mr-2 h-4 w-4" />
              Logout
            </Button>
          </template>
        </div>

        <!-- Theme Toggle -->
        <Button variant="ghost" size="icon" @click="toggleTheme" class="h-9 w-9 px-0">
          <Sun class="h-4 w-4 scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90" />
          <Moon
            class="absolute h-4 w-4 scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0"
          />
          <span class="sr-only">Toggle theme</span>
        </Button>

        <!-- Mobile Menu Button -->
        <Button variant="ghost" size="sm" class="h-9 w-9 px-0 md:hidden" @click="toggleMobileMenu">
          <Menu v-if="!isMobileMenuOpen" class="h-4 w-4" />
          <X v-else class="h-4 w-4" />
          <span class="sr-only">Toggle menu</span>
        </Button>
      </div>
    </div>

    <!-- Mobile Navigation Menu -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform -translate-y-4 opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform translate-y-0 opacity-100"
      leave-to-class="transform -translate-y-4 opacity-0"
    >
      <div v-if="isMobileMenuOpen" class="bg-background border-t md:hidden">
        <div class="container mx-auto space-y-2 px-4 py-2">
          <!-- Mobile Navigation Links -->
          <nav class="flex flex-col space-y-2">
            <router-link
              v-for="link in navigationLinks"
              :key="link.href"
              :to="link.href"
              class="hover:text-foreground/80 text-foreground/60 hover:bg-accent block rounded-md px-2 py-2 text-sm font-medium transition-colors"
              active-class="text-foreground bg-accent"
              @click="closeMobileMenu"
            >
              {{ link.label }}
            </router-link>
          </nav>

          <!-- Mobile Auth Buttons -->
          <div class="flex flex-col space-y-2 border-t pt-2">
            <template v-if="!authStore.isAuthenticated()">
              <Button variant="ghost" class="justify-start" @click="handleLogin"> Login </Button>
              <Button variant="secondary" class="justify-start" @click="handleSignup">
                Sign Up
              </Button>
            </template>
            <template v-else>
              <Button variant="ghost" class="justify-start" @click="handleLogout">
                <LogOut class="mr-2 h-4 w-4" />
                Logout
              </Button>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Button } from '@/components/ui/button'
import { Sun, Moon, Menu, X, LogOut, CodepenIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const authStore = useAuthStore()

// Mobile menu state
const isMobileMenuOpen = ref(false)

// Navigation links - customize these based on your app
const navigationLinks = [
  { href: '/', label: 'Home' },
  { href: '/snippets', label: 'Explore' },
  { href: '/about', label: 'About' },
]

const handleLogin = async () => {
  router.push('/login')
}

const handleSignup = async () => {
  router.push('/signup')
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

// Theme handling
const toggleTheme = () => {
  const html = document.documentElement
  const isDark = html.classList.contains('dark')

  if (isDark) {
    html.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  } else {
    html.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  }
}

// Mobile menu handling
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

// Initialize theme on mount
onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches

  if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
    document.documentElement.classList.add('dark')
  }
})
</script>
