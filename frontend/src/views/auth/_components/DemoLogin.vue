<template>
  <Button
    type="button"
    variant="outline"
    class="relative w-full"
    :disabled="isLoading"
    @click="handleDemoLogin"
  >
    <span class="flex items-center justify-center gap-2">
      Try Demo Version
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger>
            <InfoIcon class="text-muted-foreground h-4 w-4" />
          </TooltipTrigger>
          <TooltipContent>
            <p class="max-w-xs">
              The demo version gives you access to a pre-configured account with sample data. All
              changes are temporary and will be reset periodically.
            </p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </span>
    <span
      v-if="isLoading"
      class="bg-background/80 absolute inset-0 flex items-center justify-center"
    >
      <span class="animate-spin"><LoaderCircleIcon /></span>
    </span>
  </Button>

  <div class="bg-muted rounded-lg p-3 text-sm">
    <p class="text-muted-foreground text-center text-sm">
      Demo credentials:<br />
      Username: demo<br />
      Password: password123
    </p>
  </div>
</template>

<script setup lang="ts">
import { InfoIcon, LoaderCircleIcon } from 'lucide-vue-next'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { ref } from 'vue'
import { authService } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import { toast } from 'vue-sonner'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isLoading = ref(false)

const handleDemoLogin = async () => {
  try {
    isLoading.value = true
    const response = await authService.login({
      username: 'demo',
      password: 'password123',
    })
    authStore.setAuth({
      user: response.user,
      token: response.token,
      refreshToken: response.refreshToken,
      expiresAt: response.expiresAt,
    })
    const redirectPath = route.query.redirect as string
    router.push(redirectPath || '/snippets')
  } catch (error) {
    console.log(error)
    toast.error(error instanceof Error ? error.message : 'Failed to login with demo credentials')
  } finally {
    isLoading.value = false
  }
}
</script>
