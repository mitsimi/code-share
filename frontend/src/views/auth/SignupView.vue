<template>
  <div class="bg-background flex min-h-[calc(100vh-4rem)] items-start justify-center p-4 pt-16">
    <Card
      class="w-full max-w-md border-4 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
    >
      <CardHeader>
        <CardTitle class="text-3xl font-bold">Sign Up</CardTitle>
        <CardDescription>Create a new account to get started</CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit="onSubmit" class="space-y-4">
          <FormField v-slot="{ componentField }" name="username">
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="text"
                  placeholder="Enter your username"
                  class="border-2 border-black focus:ring-2 focus:ring-black"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="email">
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="email"
                  placeholder="Enter your email"
                  class="border-2 border-black focus:ring-2 focus:ring-black"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="password">
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="password"
                  placeholder="Create a password"
                  class="border-2 border-black focus:ring-2 focus:ring-black"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="confirmPassword">
            <FormItem>
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="password"
                  placeholder="Confirm your password"
                  class="border-2 border-black focus:ring-2 focus:ring-black"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <Button
            type="submit"
            class="bg-primary text-primary-foreground hover:bg-primary/90 w-full border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] active:translate-x-[3px] active:translate-y-[3px] active:shadow-none"
            :disabled="isLoading"
          >
            {{ isLoading ? 'Creating account...' : 'Sign Up' }}
          </Button>
        </form>
      </CardContent>
      <CardFooter class="flex justify-center">
        <p class="text-muted-foreground text-sm">
          Already have an account?
          <RouterLink to="/login" class="text-primary font-medium hover:underline">
            Login
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
import { useRouter } from 'vue-router'
import { toast } from 'sonner'
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
const authStore = useAuthStore()
const isLoading = ref(false)

const formSchema = toTypedSchema(
  z
    .object({
      username: z.string().min(2, 'Username must be at least 2 characters'),
      email: z.string().email('Please enter a valid email address'),
      password: z.string().min(8, 'Password must be at least 8 characters'),
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
    const { confirmPassword, ...signupData } = values
    const response = await authService.signup(signupData)
    authStore.setAuth({
      user: response.user,
      token: response.token,
      expiresAt: response.expires_at,
    })
    toast.success('Your account has been created successfully')
    router.push('/snippets')
  } catch (error) {
    toast.error(error instanceof Error ? error.message : 'Failed to create account')
  } finally {
    isLoading.value = false
  }
})
</script>
