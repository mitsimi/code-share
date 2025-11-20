<template>
  <pre class="highlighted-code" :class="{ 'highlighted-code--wrap': wrap }" v-bind="$attrs">
    <code ref="codeRef" :class="languageClass" class="px-2 py-1"></code>
  </pre>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue';
import Prism from 'prismjs';
import { languages } from '@/lib/languages';

const props = defineProps<{
  code: string
  language?: string | null
  wrap?: boolean
}>()

const codeRef = ref<HTMLElement | null>(null)

const aliasToLanguageId = languages.reduce<Record<string, string>>((map, lang) => {
  const canonical = lang.id
  const addAlias = (alias?: string) => {
    if (!alias) return
    map[alias.toLowerCase()] = canonical
  }

  addAlias(lang.id)
  addAlias(lang.extension)
  lang.aliases?.forEach(addAlias)

  return map
}, {})

const resolvedLanguage = computed(() => {
  if (!props.language) return 'plaintext'
  const normalized = props.language.toLowerCase()
  return aliasToLanguageId[normalized] || normalized
})

const languageClass = computed(() => `language-${resolvedLanguage.value}`)

const highlight = async () => {
  await nextTick()
  if (!codeRef.value) return
  codeRef.value.textContent = props.code ?? ''
  Prism.highlightElement(codeRef.value)
}

onMounted(() => {
  highlight()
})

watch(
  () => [props.code, resolvedLanguage.value],
  () => {
    highlight()
  },
)
</script>

<style scoped>
.highlighted-code {
  margin: 0;
  padding: 0;
  font-family: 'JetBrains Mono', 'Fira Code', Menlo, Consolas, monospace;
  font-size: 0.875rem;
  line-height: 1.6;
  color: var(--foreground, #f5f5f4);
  background: transparent;
  white-space: nowrap;
  overflow: auto;
}

.highlighted-code code {
  display: block;
}

.highlighted-code--wrap {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>

