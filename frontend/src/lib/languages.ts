import prismComponents from 'prismjs/components.json'

export interface Language {
  id: string
  name: string
  extension: string
  aliases?: string[]
}

interface PrismLanguageConfig {
  alias?: string | string[]
  title?: string
}

// Helper to normalize aliases to a string array
const getAliases = (langConfig: PrismLanguageConfig): string[] => {
  if (!langConfig.alias) return []
  return Array.isArray(langConfig.alias) ? langConfig.alias : [langConfig.alias]
}

// Generate the list automatically from Prism metadata
const prismLanguages = prismComponents.languages as Record<string, PrismLanguageConfig>

const autoLoadedLanguages: Language[] = Object.entries(prismLanguages)
  .filter(([id]) => id !== 'meta') // Remove metadata entry
  .map(([id, config]) => {
    const aliases = getAliases(config)

    // Heuristic: Use the shortest alias as the "extension", or fallback to the ID
    // e.g. id: 'javascript', aliases: ['js'] -> extension: 'js'
    const extension = aliases.sort((a, b) => a.length - b.length)[0] || id

    return {
      id,
      name: config.title || id, // Use the display title (e.g., "C#") or ID
      extension: extension,
      aliases: [id, ...aliases], // Include the ID in aliases for searchability
    }
  })

// Sort alphabetically for the dropdown
export const languages: Language[] = autoLoadedLanguages.sort((a, b) =>
  a.name.localeCompare(b.name),
)

// ... keep your existing helper functions (getLanguageExtension, etc.) ...

/**
 * Get the file extension for a given programming language
 * @param language The programming language name or alias
 * @returns The file extension for the language, or undefined if not found
 */
export function getLanguageExtension(language: string): string | undefined {
  if (!language) return undefined

  const normalizedLanguage = language.toLowerCase()
  const lang = languages.find(
    (l) => l.name.toLowerCase() === normalizedLanguage || l.aliases?.includes(normalizedLanguage),
  )
  return lang?.extension
}

/**
 * Get the display name for a given programming language
 * @param language The programming language name or alias
 * @returns The display name for the language, or undefined if not found
 */
export function getLanguageName(language: string): string | undefined {
  if (!language) return undefined

  const normalizedLanguage = language.toLowerCase()
  const lang = languages.find(
    (l) => l.name.toLowerCase() === normalizedLanguage || l.aliases?.includes(normalizedLanguage),
  )

  return lang?.name
}

/**
 * Get all available languages
 * @returns Array of language objects with name and extension
 */
export function getAvailableLanguages(): Language[] {
  return languages
}
