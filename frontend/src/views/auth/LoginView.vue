<template>
  <div class="bg-background flex min-h-[calc(100vh-4rem)] items-start justify-center p-4 pt-16">
    <Card
      class="w-full max-w-md border-4 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
    >
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
                  placeholder="Enter your password"
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
            {{ isLoading ? 'Logging in...' : 'Login' }}
          </Button>
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
      refreshToken: response.refresh_token,
      expiresAt: response.expires_at,
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
