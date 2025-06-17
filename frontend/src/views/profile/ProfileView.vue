<template>
  <div class="container mx-auto max-w-4xl px-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold">Profile Settings</h1>
      <p class="text-muted-foreground mt-2">Manage your account settings and preferences</p>
    </div>

    <Tabs v-model="activeTab" class="w-full">
      <TabsList class="flex w-full justify-start overflow-x-auto sm:grid sm:grid-cols-4">
        <TabsTrigger value="settings" class="flex-shrink-0">Settings</TabsTrigger>
        <TabsTrigger value="mine" class="flex-shrink-0">My Snippets</TabsTrigger>
        <TabsTrigger value="liked" class="flex-shrink-0">Liked Snippets</TabsTrigger>
        <TabsTrigger value="saved" class="flex-shrink-0">Saved Snippets</TabsTrigger>
      </TabsList>

      <!-- Settings Tab -->
      <TabsContent value="settings" class="mt-6">
        <SettingsTab />
      </TabsContent>

      <!-- My Snippets Tab -->
      <TabsContent value="mine" class="mt-6">
        <MySnippetsTab />
      </TabsContent>

      <!-- Liked Snippets Tab -->
      <TabsContent value="liked" class="mt-6">
        <LikedSnippetsTab />
      </TabsContent>

      <!-- Saved Snippets Tab -->
      <TabsContent value="saved" class="mt-6">
        <SavedSnippetsTab />
      </TabsContent>
    </Tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Tabs, TabsContent, TabsList } from '@/components/ui/tabs'
import TabsTrigger from '@/components/ui/tabs/TabsTrigger.vue'
import SettingsTab from './_components/SettingsTab.vue'
import MySnippetsTab from './_components/MySnippetsTab.vue'
import LikedSnippetsTab from './_components/LikedSnippetsTab.vue'
import SavedSnippetsTab from './_components/SavedSnippetsTab.vue'

const router = useRouter()
const route = useRoute()
const activeTab = ref((route.query.tab as string) || 'settings')

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
</script>
