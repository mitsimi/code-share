import fs from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import yaml from 'js-yaml'
import { bundledLanguages } from 'shiki'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const LINGUIST_URL =
  'https://raw.githubusercontent.com/github/linguist/master/lib/linguist/languages.yml'
const OUT_PATH = path.resolve(__dirname, '../src/assets/languages.json')

interface LinguistData {
  id: number

  type: string
  ace_mode: string
  extensions: string[]
  tm_scope: string

  color?: string
}

async function generate() {
  console.log('Fetching GitHub Linguist data...')
  const response = await fetch(LINGUIST_URL)
  const text = await response.text()
  const data = yaml.load(text) as Record<string, LinguistData>

  // Get the Set of valid Shiki IDs (keys of the bundledLanguages object)
  const validShikiIds = new Set(Object.keys(bundledLanguages))

  const languages: Array<LanguageDefinition> = []

  for (const [name, config] of Object.entries(data)) {
    if (!config.extensions || config.extensions.length === 0 || config.tm_scope == 'none') continue

    const tm_scope = config.tm_scope.split('.').pop()

    const entry = {
      displayName: name,
      name: name.toLowerCase(),

      color: config.color || '#888888', // Fallback gray if no color
      extensions: config.extensions,
    }

    if (validShikiIds.has(name.toLowerCase())) {
      entry.shikiId = name.toLowerCase()
    } else if (validShikiIds.has(tm_scope)) {
      entry.shikiId = tm_scope
    } else {
      // Optional: Log missing ones to see what we are dropping
      // console.warn(`Skipping ${name} (TM_Scope: ${tm_scope}) - Not in Shiki`)
      continue
    }

    languages.push(entry)
  }

  languages.push({
    displayName: 'Plain Text',
    name: 'text',
    shikiId: 'text',
    color: '#888888',
    extensions: '.txt',
  })

  languages.sort((a, b) => a.name.localeCompare(b.name))

  await fs.mkdir(path.dirname(OUT_PATH), { recursive: true })
  await fs.writeFile(OUT_PATH, JSON.stringify(languages, null, 2))

  console.log(`✅ Generated list of ${languages.length} languages.`)
}

generate()
