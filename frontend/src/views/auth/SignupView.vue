<template>
  <div class="bg-background flex min-h-[calc(100vh-4rem)] items-start justify-center p-4 pt-16">
    <Card class="w-full max-w-md border-2 shadow">
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
                <Input v-bind="componentField" type="text" placeholder="Enter your username" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

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
              <FormLabel> Password </FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                  type="password"
                  placeholder="Create a password"
                  class="pr-8"
                  @focus="showRequirements = true"
                  @blur="showRequirements = false"
                  v-model="password"
                />
              </FormControl>
              <div v-show="showRequirements" class="mt-2 text-sm">
                <p class="text-muted-foreground font-medium">Password must contain:</p>
                <ul class="mt-1 space-y-1">
                  <li class="flex items-center gap-2">
                    <template v-if="password.length >= 8">
                      <CheckCircle2Icon class="size-4 text-green-500" />
                    </template>
                    <template v-else>
                      <CircleIcon class="text-muted-foreground size-4" />
                    </template>
                    <span
                      :class="{
                        'text-green-500': password.length >= 8,
                        'text-muted-foreground': password.length < 8,
                      }"
                    >
                      At least 8 characters
                    </span>
                  </li>
                  <li class="flex items-center gap-2">
                    <template v-if="hasUpperCase(password)">
                      <CheckCircle2Icon class="size-4 text-green-500" />
                    </template>
                    <template v-else>
                      <CircleIcon class="text-muted-foreground size-4" />
                    </template>
                    <span
                      :class="{
                        'text-green-500': hasUpperCase(password),
                        'text-muted-foreground': !hasUpperCase(password),
                      }"
                    >
                      At least one uppercase letter
                    </span>
                  </li>
                  <li class="flex items-center gap-2">
                    <template v-if="hasLowerCase(password)">
                      <CheckCircle2Icon class="size-4 text-green-500" />
                    </template>
                    <template v-else>
                      <CircleIcon class="text-muted-foreground size-4" />
                    </template>
                    <span
                      :class="{
                        'text-green-500': hasLowerCase(password),
                        'text-muted-foreground': !hasLowerCase(password),
                      }"
                    >
                      At least one lowercase letter
                    </span>
                  </li>
                  <li class="flex items-center gap-2">
                    <template v-if="hasNumber(password)">
                      <CheckCircle2Icon class="size-4 text-green-500" />
                    </template>
                    <template v-else>
                      <CircleIcon class="text-muted-foreground size-4" />
                    </template>
                    <span
                      :class="{
                        'text-green-500': hasNumber(password),
                        'text-muted-foreground': !hasNumber(password),
                      }"
                    >
                      At least one number
                    </span>
                  </li>
                  <li class="flex items-center gap-2">
                    <template v-if="specialCharRegex.test(password)">
                      <CheckCircle2Icon class="size-4 text-green-500" />
                    </template>
                    <template v-else>
                      <CircleIcon class="text-muted-foreground size-4" />
                    </template>
                    <span
                      :class="{
                        'text-green-500': specialCharRegex.test(password),
                        'text-muted-foreground': !specialCharRegex.test(password),
                      }"
                    >
                      At least one special character
                    </span>
                  </li>
                </ul>
              </div>
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
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <Button type="submit" class="w-full" :disabled="isLoading">
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
import { CheckCircle2Icon, CircleIcon } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const isLoading = ref(false)
const password = ref('')
const showRequirements = ref(false)

const specialCharRegex = /[\p{P}\p{S}]/u
const uppercaseRegex = /\p{Lu}/u
const lowercaseRegex = /\p{Ll}/u
const numberRegex = /\p{N}/u
const hasUpperCase = (str: string) => uppercaseRegex.test(str)
const hasLowerCase = (str: string) => lowercaseRegex.test(str)
const hasNumber = (str: string) => numberRegex.test(str)

const formSchema = toTypedSchema(
  z
    .object({
      username: z.string().min(2, 'Username must be at least 2 characters'),
      email: z.string().email('Please enter a valid email address'),
      password: z
        .string()
        .min(8, 'Password must be at least 8 characters')
        .regex(uppercaseRegex, 'Password must contain at least one uppercase letter')
        .regex(lowercaseRegex, 'Password must contain at least one lowercase letter')
        .regex(numberRegex, 'Password must contain at least one number')
        .regex(specialCharRegex, 'Password must contain at least one special character'),
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
      refreshToken: response.refresh_token,
      expiresAt: response.expires_at,
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
