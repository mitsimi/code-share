<template>
  <div class="bg-card rounded-lg border-2 p-6 shadow">
    <form @submit.prevent="updateProfile" class="space-y-6">
      <!-- Avatar Upload -->
      <div class="flex flex-col items-center gap-6 sm:flex-row">
        <UserAvatar
          :user="authStore.user"
          :avatar-class="'size-32 ring-4'"
          :username-class="'hidden'"
        />
        <div class="flex w-full flex-col gap-2">
          <form @submit.prevent="updateAvatar" class="flex gap-2">
            <FormField v-slot="{ componentField, errorMessage }" name="avatarUrl">
              <FormItem class="w-full">
                <FormLabel>Avatar URL</FormLabel>
                <div class="flex gap-2">
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="url"
                      placeholder="Enter avatar URL"
                      :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
                    />
                  </FormControl>
                  <Button type="submit" variant="outline" :disabled="isUpdatingAvatar">
                    {{ isUpdatingAvatar ? 'Setting...' : 'Set' }}
                  </Button>
                </div>
                <FormMessage />
              </FormItem>
            </FormField>
          </form>
          <p class="text-muted-foreground text-sm">Enter a URL to your profile picture</p>
        </div>
      </div>

      <!-- Username -->
      <FormField v-slot="{ componentField, errorMessage }" name="username">
        <FormItem>
          <FormLabel>Username</FormLabel>
          <FormControl>
            <Input
              v-bind="componentField"
              type="text"
              placeholder="Enter your username"
              :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              :model-value="values.username"
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <!-- Email -->
      <FormField v-slot="{ componentField, errorMessage }" name="email">
        <FormItem>
          <FormLabel>Email</FormLabel>
          <FormControl>
            <Input
              v-bind="componentField"
              type="email"
              placeholder="Enter your email"
              :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              :model-value="values.email"
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <!-- Password Change Section -->
      <div class="space-y-4">
        <h3 class="text-lg font-medium">Change Password</h3>
        <FormField v-slot="{ componentField, errorMessage }" name="currentPassword">
          <FormItem>
            <FormLabel>Current Password</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="password"
                placeholder="Enter your current password"
                :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField, errorMessage }" name="newPassword">
          <PasswordInput
            label="New Password"
            placeholder="Enter your new password"
            :component-field="componentField"
            :error-message="errorMessage"
          />
        </FormField>

        <FormField v-slot="{ componentField, errorMessage }" name="confirmPassword">
          <FormItem>
            <FormLabel>Confirm New Password</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="password"
                placeholder="Confirm your new password"
                :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <Button type="submit" class="w-full" :disabled="isUpdatingProfile">
        {{ isUpdatingProfile ? 'Saving changes...' : 'Save Changes' }}
      </Button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { toast } from 'vue-sonner'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { passwordSchema } from '@/lib/password'
import PasswordInput from '@/components/ui/password-input.vue'
import { useMutation } from '@tanstack/vue-query'
import { usersService } from '@/services/users'

const authStore = useAuthStore()
const avatarUrl = ref(authStore.user?.avatar || '')

// Form validation schema
const formSchema = toTypedSchema(
  z
    .object({
      username: z.string().min(2, 'Username must be at least 2 characters'),
      email: z.string().email('Please enter a valid email address'),
      currentPassword: z.string().optional(),
      newPassword: passwordSchema.optional(),
      confirmPassword: z.string().optional(),
    })
    .refine(
      (data) => {
        if (data.newPassword && !data.currentPassword) {
          return false
        }
        return true
      },
      {
        message: 'Current password is required to change password',
        path: ['currentPassword'],
      },
    )
    .refine(
      (data) => {
        if (data.newPassword && data.newPassword !== data.confirmPassword) {
          return false
        }
        return true
      },
      {
        message: "Passwords don't match",
        path: ['confirmPassword'],
      },
    ),
)

const { handleSubmit, values } = useForm({
  validationSchema: formSchema,
  initialValues: {
    username: authStore.user?.username || '',
    email: authStore.user?.email || '',
  },
})

// Avatar update mutation
const { mutate: updateAvatarMutation, isPending: isUpdatingAvatar } = useMutation({
  mutationKey: ['updateAvatar'],
  mutationFn: usersService.updateAvatar,
  onSuccess: (newAvatarUrl) => {
    authStore.setUser({
      ...authStore.user!,
      avatar: newAvatarUrl,
    })
    avatarUrl.value = newAvatarUrl
    toast.success('Avatar updated successfully')
  },
  onError: (error) => {
    toast.error(error instanceof Error ? error.message : 'Failed to update avatar')
  },
})

// Profile update mutation
const { mutate: updateProfileMutation, isPending: isUpdatingProfile } = useMutation({
  mutationKey: ['updateProfile'],
  mutationFn: async (values: {
    username: string
    email: string
    currentPassword?: string
    newPassword?: string
  }) => {
    const { currentPassword, newPassword, ...profileData } = values

    // Update profile
    const updatedUser = await usersService.updateProfile(profileData)

    // Update password if provided
    if (currentPassword && newPassword) {
      await usersService.updatePassword({ currentPassword, newPassword })
    }

    return updatedUser
  },
  onSuccess: (data) => {
    authStore.setUser({
      ...authStore.user!,
      ...data,
    })
    toast.success('Profile updated successfully')
  },
  onError: (error) => {
    toast.error(error instanceof Error ? error.message : 'Failed to update profile')
  },
})

// Avatar form
const avatarForm = useForm({
  validationSchema: toTypedSchema(
    z.object({
      avatarUrl: z.string().url('Please enter a valid URL').optional(),
    }),
  ),
})

const updateAvatar = avatarForm.handleSubmit((values) => {
  if (!values.avatarUrl) {
    toast.error('Please enter a valid avatar URL')
    return
  }
  updateAvatarMutation(values.avatarUrl)
})

const updateProfile = handleSubmit((values) => {
  updateProfileMutation(values)
})
</script>
