<template>
  <ContextMenu>
    <ContextMenuTrigger asChild>
      <Button variant="ghost" size="icon" @click="toggleTheme">
        <span class="sr-only">Toggle Theme</span>
        <div class="flex flex-row items-center justify-center">
          <div class="">
            <!-- Display icon based on the selected mode (light, dark, or auto/system) -->
            <span v-if="store === 'light'" class="theme-icon">
              <SunIcon class="size-4" />
            </span>
            <span v-else-if="store === 'dark'" class="theme-icon">
              <MoonIcon class="size-4" />
            </span>
            <span v-else-if="store === 'auto'">
              <!-- 'auto' corresponds to 'system' -->
              <SunMoonIcon class="size-4" />
            </span>
            <span v-else>
              <CircleHelpIcon class="size-4" />
            </span>
          </div>
        </div>
      </Button>
    </ContextMenuTrigger>

    <ContextMenuContent>
      <ContextMenuGroup>
        <ContextMenuItem
          value="light"
          @click="store = 'light'"
          class="flex items-center justify-between"
        >
          <div class="flex items-center gap-2"><SunIcon /> Light</div>
          <CheckIcon v-if="store === 'light'" class="size-4" />
        </ContextMenuItem>
        <ContextMenuItem
          value="dark"
          @click="store = 'dark'"
          class="flex items-center justify-between"
        >
          <div class="flex items-center gap-2"><MoonIcon /> Dark</div>
          <CheckIcon v-if="store === 'dark'" class="size-4" />
        </ContextMenuItem>
        <ContextMenuItem
          value="system"
          @click="store = 'auto'"
          class="flex items-center justify-between"
        >
          <div class="flex items-center gap-2"><SunMoonIcon /> System</div>
          <CheckIcon v-if="store === 'auto'" class="size-4" />
        </ContextMenuItem>
      </ContextMenuGroup>
    </ContextMenuContent>
  </ContextMenu>
</template>

<script lang="ts" setup>
import { CheckIcon, SunIcon, MoonIcon, SunMoonIcon, CircleHelpIcon } from 'lucide-vue-next'
import { useColorMode } from '@vueuse/core'

const { system, store } = useColorMode({
  storageKey: 'theme',
  selector: 'html',
  attribute: 'class',
  initialValue: 'auto',
  modes: {
    dark: 'dark',
    light: '',
  },
})

const toggleTheme = () => {
  if (store.value === 'auto' && system.value === 'light') {
    store.value = 'dark'
  } else if (store.value === 'auto' && system.value === 'dark') {
    store.value = 'light'
  } else {
    store.value = store.value === 'light' ? 'dark' : 'light'
  }
}
</script>
