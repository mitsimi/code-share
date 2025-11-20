import { ref, type Ref } from 'vue'
import { codeToHtml } from 'shiki'
import { getShikiLanguageId } from '@/lib/languages'
import { escapeHtml } from '@/lib/utils'

export async function highlightCode(code: string, language?: string | null): Promise<string> {
  try {
    const langId = getShikiLanguageId(language || '')
    return await codeToHtml(code, {
      lang: langId,
      themes: {
        light: 'one-light',
        dark: 'one-dark-pro',
      },
    })
  } catch (err) {
    console.error('Error highlighting code:', err)

    // 3. SAFE FALLBACK:
    // Return valid HTML structure that matches Shiki's output
    // so the UI doesn't jump or break.
    return `<pre class="shiki fallback"><code>${escapeHtml(code)}</code></pre>`
  }
}

export function useShiki() {
  const isReady = ref(true)

  return {
    isReady: isReady as Readonly<Ref<boolean>>,
    highlight: highlightCode,
  }
}
