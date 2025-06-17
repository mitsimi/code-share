<template>
  <FormItem>
    <FormLabel>{{ label }}</FormLabel>
    <FormControl>
      <Input
        v-bind="componentField"
        type="password"
        :placeholder="placeholder"
        class="pr-8"
        :class="{ 'ring-2 ring-red-500 ring-offset-2': errorMessage }"
        @focus="showRequirements = true"
        @blur="showRequirements = false"
        v-model="password"
      />
    </FormControl>
    <template v-if="showRequirements">
      <div class="mt-2 text-sm">
        <p class="text-muted-foreground font-medium">Password must contain:</p>
        <ul class="mt-1 space-y-1">
          <li
            v-for="requirement in passwordRequirements"
            :key="requirement.id"
            class="flex items-center gap-2"
          >
            <template v-if="requirement.validator(password)">
              <CheckCircle2Icon class="size-4 text-green-500" />
            </template>
            <template v-else>
              <CircleIcon class="text-muted-foreground size-4" />
            </template>
            <span
              :class="{
                'text-green-500': requirement.validator(password),
                'text-muted-foreground': !requirement.validator(password),
              }"
            >
              {{ requirement.label }}
            </span>
          </li>
        </ul>
      </div>
    </template>
    <FormMessage />
  </FormItem>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CheckCircle2Icon, CircleIcon } from 'lucide-vue-next'
import { FormControl, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { passwordRequirements } from '@/utils/password'

defineProps<{
  label: string
  placeholder: string
  componentField: Record<string, any>
  errorMessage?: string
}>()

const password = ref('')
const showRequirements = ref(false)
</script>
