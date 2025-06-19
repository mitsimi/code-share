<template>
  <header
    class="bg-background/95 supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50 w-full border-b backdrop-blur"
  >
    <div class="mx-auto flex h-14 max-w-7xl items-center px-4">
      <!-- Logo/Brand -->
      <div class="mx-4 flex">
        <router-link to="/" class="mr-6 flex items-center space-x-2">
          <span class="font-bold">CodeShare</span>
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
          <Unauthenticated>
            <Button variant="ghost" size="sm" @click="handleLogin"> Login </Button>
            <Button variant="reverse" size="sm" @click="handleRegister"> Register </Button>
          </Unauthenticated>
          <Authenticated>
            <div class="flex items-center space-x-2">
              <UserMenu @logout="handleLogout" @close="closeMobileMenu" />
            </div>
          </Authenticated>
        </div>

        <!-- Theme Toggle -->
        <ThemeSwitch />

        <!-- Mobile Menu Button -->
        <Button variant="ghost" size="sm" class="h-9 w-9 px-0 md:hidden" @click="toggleMobileMenu">
          <MenuIcon v-if="!isMobileMenuOpen" class="h-4 w-4" />
          <XIcon v-else class="h-4 w-4" />
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

          <Separator />

          <!-- Mobile Auth Buttons -->
          <div class="flex flex-col space-y-2">
            <Unauthenticated>
              <Button variant="ghost" class="justify-start" @click="handleLogin"> Login </Button>
              <Button variant="secondary" class="justify-start" @click="handleRegister">
                Sign Up
              </Button>
            </Unauthenticated>
            <Authenticated>
              <UserInfo class="px-2 py-2" />
              <MenuItems isMobile @logout="handleLogout" @close="closeMobileMenu" />
            </Authenticated>
          </div>
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'
import ThemeSwitch from './ThemeSwitch.vue'
import Separator from '../ui/separator/Separator.vue'
import UserInfo from '@/components/UserMenu/UserInfo.vue'
import MenuItems from '@/components/UserMenu/MenuItems.vue'
import { MenuIcon, XIcon } from 'lucide-vue-next'
import Unauthenticated from '../access/Unauthenticated.vue'

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
  isMobileMenuOpen.value = false
  router.push('/login')
}

const handleRegister = async () => {
  isMobileMenuOpen.value = false
  router.push('/register')
}

const handleLogout = async () => {
  await authStore.logout()
  isMobileMenuOpen.value = false
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
