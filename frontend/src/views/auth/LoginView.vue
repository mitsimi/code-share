<template>
  <div class="bg-background flex min-h-[calc(100vh-4rem)] items-start justify-center p-4">
    <div class="w-full max-w-md space-y-6">
      <div class="space-y-2 text-center">
        <h1 class="text-3xl font-bold">Login</h1>
        <p class="text-muted-foreground">Enter your credentials to access your account</p>
      </div>

      <form @submit="onSubmit" class="space-y-4">
        <FormField v-slot="{ componentField }" name="username">
          <FormItem>
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input v-bind="componentField" type="text" placeholder="Enter your username" />
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
      </form>

      <div class="flex justify-center">
        <p class="text-muted-foreground text-sm">
          Don't have an account?
          <RouterLink to="/register" class="text-primary font-medium hover:underline">
            Register
          </RouterLink>
        </p>
      </div>
    </div>
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
import { authService } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isLoading = ref(false)

const formSchema = toTypedSchema(
  z.object({
    username: z.string(),
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
</script>
