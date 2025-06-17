<template>
  <div class="container mx-auto max-w-4xl px-4">
    <Card>
      <CardHeader>
        <CardTitle class="text-3xl font-bold">Profile Settings</CardTitle>
        <CardDescription>Manage your account settings and preferences</CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs v-model="activeTab" class="w-full">
          <TabsList class="grid w-full grid-cols-4">
            <TabsTrigger value="settings">Settings</TabsTrigger>
            <TabsTrigger value="mine">My Snippets</TabsTrigger>
            <TabsTrigger value="liked">Liked Snippets</TabsTrigger>
            <TabsTrigger value="saved">Saved Snippets</TabsTrigger>
          </TabsList>

          <!-- Settings Tab -->
          <TabsContent value="settings" class="mt-6">
            <form @submit.prevent="updateProfile" class="space-y-6">
              <!-- Avatar Upload -->
              <div class="flex items-center gap-6">
                <Avatar class="size-24">
                  <AvatarImage :src="avatarUrl" />
                  <AvatarFallback class="text-4xl">{{
                    authStore.user?.username[0].toUpperCase()
                  }}</AvatarFallback>
                </Avatar>
                <div class="flex flex-col gap-2">
                  <form @submit.prevent="updateAvatar" class="flex gap-2">
                    <FormField v-slot="{ componentField, errorMessage }" name="avatarUrl">
                      <FormItem>
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
          </TabsContent>

          <!-- My Snippets Tab -->
          <TabsContent value="mine" class="mt-6">
            <div v-if="isLoadingMySnippets" class="flex justify-center py-8">
              <LoaderCircleIcon class="size-8 animate-spin" />
            </div>
            <div v-else-if="mySnippets.length === 0" class="py-8 text-center">
              <p class="text-muted-foreground">You haven't created any snippets yet.</p>
            </div>
            <div v-else class="grid gap-4">
              <SnippetCard
                v-for="snippet in mySnippets"
                :key="snippet.id"
                :snippet="snippet"
                @click="router.push(`/snippets/${snippet.id}`)"
              />
            </div>
          </TabsContent>

          <!-- Liked Snippets Tab -->
          <TabsContent value="liked" class="mt-6">
            <div v-if="isLoadingLiked" class="flex justify-center py-8">
              <LoaderCircleIcon class="size-8 animate-spin" />
            </div>
            <div v-else-if="likedSnippets.length === 0" class="py-8 text-center">
              <p class="text-muted-foreground">You haven't liked any snippets yet.</p>
            </div>
            <div v-else class="grid gap-4">
              <SnippetCard
                v-for="snippet in likedSnippets"
                :key="snippet.id"
                :snippet="snippet"
                @click="router.push(`/snippets/${snippet.id}`)"
              />
            </div>
          </TabsContent>

          <!-- Saved Snippets Tab -->
          <TabsContent value="saved" class="mt-6">
            <div v-if="isLoadingSaved" class="flex justify-center py-8">
              <LoaderCircleIcon class="size-8 animate-spin" />
            </div>
            <div v-else-if="savedSnippets.length === 0" class="py-8 text-center">
              <p class="text-muted-foreground">You haven't saved any snippets yet.</p>
            </div>
            <div v-else class="grid gap-4">
              <SnippetCard
                v-for="snippet in savedSnippets"
                :key="snippet.id"
                :snippet="snippet"
                @click="router.push(`/snippets/${snippet.id}`)"
              />
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { useRouter, useRoute } from 'vue-router'
import { toast } from 'vue-sonner'
import { useAuthStore } from '@/stores/auth'
import { useFetch } from '@/composables/useCustomFetch'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { LoaderCircleIcon } from 'lucide-vue-next'
import SnippetCard from '@/views/snippets/_components/SnippetCard.vue'
import { passwordSchema } from '@/utils/password'
import PasswordInput from '@/components/ui/password-input.vue'
import TabsTrigger from '@/components/ui/tabs/TabsTrigger.vue'
import type { Snippet } from '@/types'
import { useQuery, useMutation } from '@tanstack/vue-query'

const router = useRouter()
const authStore = useAuthStore()
const route = useRoute()
const activeTab = ref((route.query.tab as string) || 'settings')
const avatarUrl = ref(authStore.user?.avatar || '')

// Watch for tab changes to update URL
watch(activeTab, (newTab) => {
  router.replace({ query: { ...route.query, tab: newTab } })
})

// Watch for URL changes to update active tab
watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab && typeof newTab === 'string') {
      activeTab.value = newTab
    }
  },
  { immediate: true },
)

// Form validation schema
const formSchema = toTypedSchema(
  z
    .object({
      username: z.string().min(2, 'Username must be at least 2 characters'), // not optional because they are prefilled
      email: z.string().email('Please enter a valid email address'), // not optional because they are prefilled
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

// Fetch liked snippets
const { data: likedSnippets, isLoading: isLoadingLiked } = useQuery({
  queryKey: ['liked-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(`/users/${authStore.user?.id}/liked`).json()
    if (error.value) throw new Error('Failed to fetch liked snippets')
    return data.value || []
  },
})

// Fetch saved snippets
const { data: savedSnippets, isLoading: isLoadingSaved } = useQuery({
  queryKey: ['saved-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(`/users/${authStore.user?.id}/saved`).json()
    if (error.value) throw new Error('Failed to fetch saved snippets')
    return data.value || []
  },
})

// Fetch user's snippets
const { data: mySnippets, isLoading: isLoadingMySnippets } = useQuery({
  queryKey: ['my-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(
      `/users/${authStore.user?.id}/snippets`,
    ).json()
    if (error.value) throw new Error('Failed to fetch your snippets')
    return data.value || []
  },
})

// Avatar update mutation
const { mutate: updateAvatarMutation, isPending: isUpdatingAvatar } = useMutation({
  mutationKey: ['updateAvatar'],
  mutationFn: async (avatarUrl: string) => {
    const { data, error } = await useFetch(`/users/${authStore.user?.id}/avatar`, {
      method: 'PATCH',
      body: JSON.stringify({ avatarUrl }),
    }).json()

    if (error.value) throw new Error('Failed to update avatar')
    if (!data.value) throw new Error('No data received from server')

    return data.value.avatar
  },
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
    const { data, error } = await useFetch(`/users/${authStore.user?.id}`, {
      method: 'PATCH',
      body: JSON.stringify(profileData),
    }).json()

    if (error.value) throw new Error('Failed to update profile')
    if (!data.value) throw new Error('No data received from server')

    // Update password if provided
    if (currentPassword && newPassword) {
      const { error: passwordError } = await useFetch(`/users/${authStore.user?.id}/password`, {
        method: 'PATCH',
        body: JSON.stringify({
          currentPassword,
          newPassword,
        }),
      }).json()

      if (passwordError.value) throw new Error('Failed to update password')
    }

    return data.value
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
  initialValues: {
    avatarUrl: authStore.user?.avatar || '',
  },
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
