<template>
  <div class="shiki-wrapper" v-html="highlightedCode"></div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { highlightCode } from '@/composables/useShiki'
import { escapeHtml } from '@/lib/utils'

const props = defineProps<{
  code: string
  language?: string | null
}>()

// Initial render: Plain text wrapped in pre/code to prevent CLS (Layout Shift)
// while Shiki loads lazily
const highlightedCode = ref<string>(
  `<pre class="shiki"><code>${escapeHtml(props.code ?? '')}</code></pre>`,
)

const performHighlight = async () => {
  highlightedCode.value = await highlightCode(props.code ?? '', props.language)
}

onMounted(performHighlight)

watch(() => [props.code, props.language], performHighlight)
</script>

<style>
.shiki-wrapper pre.shiki {
  font-family: 'JetBrains Mono', 'Fira Code', Menlo, Consolas, monospace;
  font-size: 0.875rem;
  line-height: 1.7;
  background-color: transparent !important;
}

html.dark .shiki,
html.dark .shiki span {
  color: var(--shiki-dark) !important;
  font-style: var(--shiki-dark-font-style) !important;
  font-weight: var(--shiki-dark-font-weight) !important;
  text-decoration: var(--shiki-dark-text-decoration) !important;
}
</style>
