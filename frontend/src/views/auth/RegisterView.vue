<template>
  <div class="bg-background flex min-h-[calc(100vh-4rem)] items-start justify-center p-4 pt-16">
    <div class="w-full max-w-md space-y-6">
      <div class="space-y-2 text-center">
        <h1 class="text-3xl font-bold">Register</h1>
        <p class="text-muted-foreground">Create a new account to get started</p>
      </div>

      <form @submit="onSubmit" class="space-y-4">
        <FormField v-slot="{ componentField, errorMessage }" name="username">
          <FormItem>
            <FormLabel>Username</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="text"
                placeholder="Enter your username"
                :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField, errorMessage }" name="email">
          <FormItem>
            <FormLabel>Email</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="email"
                placeholder="Enter your email"
                :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField, errorMessage }" name="password">
          <PasswordInput
            label="Password"
            placeholder="Create a password"
            :component-field="componentField"
            :error-message="errorMessage"
          />
        </FormField>

        <FormField v-slot="{ componentField, errorMessage }" name="confirmPassword">
          <FormItem>
            <FormLabel>Confirm Password</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="password"
                placeholder="Confirm your password"
                :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <Button type="submit" class="w-full" :disabled="isLoading">
          {{ isLoading ? 'Creating account...' : 'Register' }}
        </Button>
      </form>

      <div class="flex justify-center">
        <p class="text-muted-foreground text-sm">
          Already have an account?
          <RouterLink to="/login" class="text-primary font-medium hover:underline">
            Login
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
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { authService } from '@/services/auth'
import { useAuthStore } from '@/stores/auth'
import { passwordSchema } from '@/lib/password'
import PasswordInput from '@/components/ui/password-input.vue'

const router = useRouter()
const authStore = useAuthStore()
const isLoading = ref(false)

const formSchema = toTypedSchema(
  z
    .object({
      username: z.string().min(2, 'Username must be at least 2 characters'),
      email: z.string().email('Please enter a valid email address'),
      password: passwordSchema,
      confirmPassword: z.string(),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: "Passwords don't match",
      path: ['confirmPassword'],
    }),
)

const { handleSubmit } = useForm({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
  try {
    isLoading.value = true
    const { confirmPassword, ...registerData } = values
    const response = await authService.register(registerData)
    authStore.setAuth({
      user: response.user,
      token: response.token,
      refreshToken: response.refreshToken,
      expiresAt: response.expiresAt,
    })
    router.push('/snippets')
  } catch (error) {
    console.log(error)
    toast.error(error instanceof Error ? error.message : 'Failed to create account')
  } finally {
    isLoading.value = false
  }
})
</script>
