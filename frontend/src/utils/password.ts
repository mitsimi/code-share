import * as z from 'zod'

// Regex patterns
export const specialCharRegex = /[\p{P}\p{S}]/u
export const uppercaseRegex = /\p{Lu}/u
export const lowercaseRegex = /\p{Ll}/u
export const numberRegex = /\p{N}/u

// Validation functions
export const hasUpperCase = (str: string) => uppercaseRegex.test(str)
export const hasLowerCase = (str: string) => lowercaseRegex.test(str)
export const hasNumber = (str: string) => numberRegex.test(str)
export const hasSpecialChar = (str: string) => specialCharRegex.test(str)

// Password validation schema
export const passwordSchema = z
  .string()
  .min(8, 'Password must be at least 8 characters')
  .regex(uppercaseRegex, 'Password must contain at least one uppercase letter')
  .regex(lowercaseRegex, 'Password must contain at least one lowercase letter')
  .regex(numberRegex, 'Password must contain at least one number')
  .regex(specialCharRegex, 'Password must contain at least one special character')

// Password requirements for UI display
export const passwordRequirements = [
  {
    id: 'length',
    label: 'At least 8 characters',
    validator: (password: string) => password.length >= 8,
  },
  {
    id: 'uppercase',
    label: 'At least one uppercase letter',
    validator: hasUpperCase,
  },
  {
    id: 'lowercase',
    label: 'At least one lowercase letter',
    validator: hasLowerCase,
  },
  {
    id: 'number',
    label: 'At least one number',
    validator: hasNumber,
  },
  {
    id: 'special',
    label: 'At least one special character',
    validator: hasSpecialChar,
  },
]
