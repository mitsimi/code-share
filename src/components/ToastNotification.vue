<template>
  <div class="fixed bottom-8 left-8 z-[1000] flex flex-col gap-2">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="bg-background flex max-w-[300px] min-w-[200px] transform items-center rounded-lg border-4 border-black p-4 shadow-[4px_4px_0_0_#000]"
        :class="[toastClasses[toast.type], 'translate-y-0 rotate-2']"
        @mouseenter="pauseToast(toast.id)"
        @mouseleave="resumeToast(toast.id)"
      >
        <div class="flex items-center gap-3">
          <component :is="icons[toast.type]" class="h-6 w-6 flex-shrink-0" />
          <span>{{ toast.message }}</span>
        </div>
        <Button variant="ghost" size="icon" @click="removeToast(toast.id)">
          <X class="size-6" />
        </Button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '../composables/useToast'
import { Info, CheckCircle, AlertCircle, X } from 'lucide-vue-next'
import { Button } from './ui/button'
const { toasts, removeToast, pauseToast, resumeToast } = useToast()

const toastClasses = {
  info: 'bg-primary text-primary-foreground',
  success: 'bg-secondary text-secondary-foreground',
  error: 'bg-destructive text-destructive-foreground',
}

const icons = {
  info: Info,
  success: CheckCircle,
  error: AlertCircle,
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(100%) rotate(3deg);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-100%) rotate(-3deg);
}
</style>
