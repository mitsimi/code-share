export interface Language {
  name: string
  extension: string
  aliases?: string[]
}

export const languages: Language[] = [
  { name: 'JavaScript', extension: 'js', aliases: ['javascript', 'js'] },
  { name: 'TypeScript', extension: 'ts', aliases: ['typescript', 'ts'] },
  { name: 'JSX', extension: 'jsx', aliases: ['jsx', 'react'] },
  { name: 'TSX', extension: 'tsx', aliases: ['tsx', 'react-ts'] },
  { name: 'Vue', extension: 'vue', aliases: ['vue', 'vuejs'] },
  { name: 'Python', extension: 'py', aliases: ['python', 'py'] },
  { name: 'Java', extension: 'java', aliases: ['java'] },
  { name: 'Kotlin', extension: 'kt', aliases: ['kotlin', 'kt'] },
  { name: 'C', extension: 'c', aliases: ['c'] },
  { name: 'C++', extension: 'cpp', aliases: ['cpp', 'c++', 'cplusplus'] },
  { name: 'C#', extension: 'cs', aliases: ['csharp', 'cs', 'c#'] },
  { name: 'Go', extension: 'go', aliases: ['go', 'golang'] },
  { name: 'Rust', extension: 'rs', aliases: ['rust', 'rs'] },
  { name: 'Ruby', extension: 'rb', aliases: ['ruby', 'rb'] },
  { name: 'PHP', extension: 'php', aliases: ['php'] },
  { name: 'Swift', extension: 'swift', aliases: ['swift'] },
  { name: 'HTML', extension: 'html', aliases: ['html', 'htm'] },
  { name: 'CSS', extension: 'css', aliases: ['css'] },
  { name: 'SQL', extension: 'sql', aliases: ['sql'] },
  { name: 'Shell', extension: 'sh', aliases: ['shell', 'bash', 'sh'] },
  { name: 'YAML', extension: 'yml', aliases: ['yaml', 'yml'] },
  { name: 'TOML', extension: 'toml', aliases: ['toml'] },
  { name: 'JSON', extension: 'json', aliases: ['json'] },
  { name: 'Markdown', extension: 'md', aliases: ['markdown', 'md'] },
  { name: 'GraphQL', extension: 'graphql', aliases: ['graphql', 'gql'] },
  { name: 'Elixir', extension: 'ex', aliases: ['elixir', 'ex'] },
  { name: 'Haskell', extension: 'hs', aliases: ['haskell', 'hs'] },
  { name: 'Scala', extension: 'scala', aliases: ['scala'] },
]

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
