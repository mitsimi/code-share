import languagesData from '@/assets/languages.json'

export interface LanguageDefinition {
  displayName: string
  name: string
  shikiId: string
  color: string
  extensions: string[]
}

// THE LIST (For Dropdowns)
export const allLanguages = languagesData as LanguageDefinition[]

// THE MAPS (For Fast Lookup)
const nameMap = new Map<string, LanguageDefinition>()
const idMap = new Map<string, LanguageDefinition>()
const extensionMap = new Map<string, LanguageDefinition>()

for (const lang of allLanguages) {
  nameMap.set(lang.name, lang)
  idMap.set(lang.shikiId, lang)

  for (let ext of lang.extensions) {
    ext = ext.substring(1)
    if (!extensionMap.has(ext)) extensionMap.set(ext, lang)
  }
}

/**
 * Get full language details from an extension, or ID.
 */
export function getLanguage(input: string | null | undefined): LanguageDefinition | null {
  if (!input) return null
  const query = input.toLowerCase().trim()

  if (nameMap.has(query)) {
    return nameMap.get(query) || null
  }

  if (extensionMap.has(query)) {
    return extensionMap.get(query) || null
  }

  return idMap.get(query) || null
}

/**
 * Helper for Shiki to just get the ID
 */
export function getShikiLanguageId(input: string | null | undefined): string {
  const lang = getLanguage(input)
  return lang ? lang.shikiId : 'text'
}
