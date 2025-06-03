<template>
  <div class="bg-background flex min-h-[calc(100vh-4rem)] items-start justify-center p-4 pt-16">
    <Card class="w-full max-w-md">
      <CardHeader>
        <CardTitle class="text-3xl font-bold">Login</CardTitle>
        <CardDescription>Enter your credentials to access your account</CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit="onSubmit" class="space-y-4">
          <FormField v-slot="{ componentField }" name="email">
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input v-bind="componentField" type="email" placeholder="Enter your email" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="password">
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input v-bind="componentField" type="password" placeholder="Enter your password" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <Button type="submit" class="w-full" :disabled="isLoading">
            {{ isLoading ? 'Logging in...' : 'Login' }}
          </Button>

          <div class="relative">
            <div class="absolute inset-0 flex items-center">
              <span class="w-full border-t" />
            </div>
            <div class="relative flex justify-center text-xs uppercase">
              <span class="bg-background text-muted-foreground px-2">Or</span>
            </div>
          </div>

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
                      The demo version gives you access to a pre-configured account with sample
                      data. All changes are temporary and will be reset periodically.
                    </p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </span>
            <span
              v-if="isLoading"
              class="bg-background/80 absolute inset-0 flex items-center justify-center"
            >
              <span class="animate-spin">⌛</span>
            </span>
          </Button>
          <p class="text-muted-foreground text-center text-sm">
            Demo credentials:<br />
            Email: demo@example.com<br />
            Password: password123
          </p>
          <div class="bg-muted rounded-lg p-3 text-sm">
            <p class="text-muted-foreground font-medium">Demo Account Information:</p>
            <ul class="text-muted-foreground mt-2 list-inside list-disc space-y-1">
              <li>Access to sample data and features</li>
              <li>Changes are temporary and will reset</li>
              <li>Limited functionality compared to full accounts</li>
              <li>Perfect for exploring the platform</li>
            </ul>
          </div>
        </form>
      </CardContent>
      <CardFooter class="flex justify-center">
        <p class="text-muted-foreground text-sm">
          Don't have an account?
          <RouterLink to="/signup" class="text-primary font-medium hover:underline">
            Sign up
          </RouterLink>
        </p>
      </CardFooter>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { authService } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isLoading = ref(false)

const formSchema = toTypedSchema(
  z.object({
    email: z.string().email('Please enter a valid email address'),
    password: z.string(),
  }),
)

const { handleSubmit } = useForm({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
  try {
    isLoading.value = true
    const response = await authService.login(values)
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
    toast.error(error instanceof Error ? error.message : 'Failed to login')
  } finally {
    isLoading.value = false
  }
})

const handleDemoLogin = async () => {
  try {
    isLoading.value = true
    const response = await authService.login({
      email: 'demo@example.com',
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
    toast.error(error instanceof Error ? error.message : 'Failed to login to demo account')
  } finally {
    isLoading.value = false
  }
}
</script>
