import { describe, it, expect } from 'vitest'
import {
    hasUpperCase,
    hasLowerCase,
    hasNumber,
    hasSpecialChar,
    passwordSchema,
} from './password'

describe('Password Validation', () => {
  describe('Helper Functions', () => {
    it('hasUpperCase detects uppercase letters', () => {
      expect(hasUpperCase('abc')).toBe(false)
      expect(hasUpperCase('abC')).toBe(true)
    })

    it('hasLowerCase detects lowercase letters', () => {
      expect(hasLowerCase('ABC')).toBe(false)
      expect(hasLowerCase('ABc')).toBe(true)
    })

    it('hasNumber detects numbers', () => {
      expect(hasNumber('abc')).toBe(false)
      expect(hasNumber('ab1')).toBe(true)
    })

    it('hasSpecialChar detects special characters', () => {
      expect(hasSpecialChar('abc')).toBe(false)
      expect(hasSpecialChar('abc!')).toBe(true)
      expect(hasSpecialChar('abc@')).toBe(true)
    })
  })

  describe('Password Schema', () => {
    it('validates a correct password', () => {
      const validPassword = 'Password123!@#LongEnough'
      const result = passwordSchema.safeParse(validPassword)
      expect(result.success).toBe(true)
    })

    it('rejects short passwords', () => {
      const shortPassword = 'Pass1!'
      const result = passwordSchema.safeParse(shortPassword)
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('at least 15 characters')
      }
    })

    it('rejects passwords missing uppercase', () => {
      const noUpper = 'password123!@#longenough'
      const result = passwordSchema.safeParse(noUpper)
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('uppercase')
      }
    })

    it('rejects passwords missing lowercase', () => {
      const noLower = 'PASSWORD123!@#LONGENOUGH'
      const result = passwordSchema.safeParse(noLower)
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('lowercase')
      }
    })

    it('rejects passwords missing numbers', () => {
      const noNumber = 'Password!@#LongEnough'
      const result = passwordSchema.safeParse(noNumber)
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('number')
      }
    })

    it('rejects passwords missing special characters', () => {
      const noSpecial = 'Password123LongEnough'
      const result = passwordSchema.safeParse(noSpecial)
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('special character')
      }
    })
  })
})
