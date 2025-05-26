<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import { useVModel } from '@vueuse/core'

const props = defineProps<{
  defaultValue?: string | number
  modelValue?: string | number
  class?: HTMLAttributes['class']
}>()

const emits = defineEmits<{
  (e: 'update:modelValue', payload: string | number): void
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
  defaultValue: props.defaultValue,
})
</script>

<template>
  <input
    v-model="modelValue"
    data-slot="input"
    :class="
      cn(
        'file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 bg-input flex h-10 w-full min-w-0 rounded-md border-4 border-black px-3 py-1 text-base shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] transition-all outline-none file:inline-flex file:h-8 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
        'focus-visible:border-ring focus-visible:translate-x-[2px] focus-visible:translate-y-[2px] focus-visible:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] focus-visible:ring-0',
        'aria-invalid:border-destructive aria-invalid:text-destructive aria-invalid:shadow-[4px_4px_0px_0px_rgba(220,38,38,0.5)]',
        'hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]',
        'active:translate-x-[4px] active:translate-y-[4px] active:shadow-none',
        props.class,
      )
    "
  />
</template>
