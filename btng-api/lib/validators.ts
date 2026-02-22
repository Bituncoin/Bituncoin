/**
 * Validation Utilities Module
 * Centralized validation functions for BTNG platform
 */

export interface ValidationResult {
  valid: boolean
  error?: string
}

/**
 * Validate email format
 */
export function validateEmail(email: string): ValidationResult {
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailPattern.test(email)) {
    return { valid: false, error: 'Invalid email format' }
  }
  return { valid: true }
}

/**
 * Validate phone number (international format)
 */
export function validatePhone(phone: string): ValidationResult {
  // Remove spaces, dashes, parentheses
  const cleaned = phone.replace(/[\s\-()]/g, '')
  
  // Check if it starts with + and has 10-15 digits
  const phonePattern = /^\+?[1-9]\d{9,14}$/
  
  if (!phonePattern.test(cleaned)) {
    return { valid: false, error: 'Invalid phone number format' }
  }
  return { valid: true }
}

/**
 * Validate Gold Card number format
 */
export function validateGoldCard(cardNumber: string): ValidationResult {
  const pattern = /^BTNG-[A-Z0-9]{4}-[A-Z0-9]{4}$/
  if (!pattern.test(cardNumber)) {
    return { valid: false, error: 'Invalid Gold Card format' }
  }
  return { valid: true }
}

/**
 * Validate wallet ID format
 */
export function validateWallet(walletId: string): ValidationResult {
  const pattern = /^BTNG-[A-Z0-9]+-[A-Z0-9]+$/
  if (!pattern.test(walletId)) {
    return { valid: false, error: 'Invalid wallet ID format' }
  }
  return { valid: true }
}

/**
 * Validate required string field
 */
export function validateRequired(value: string, fieldName: string): ValidationResult {
  if (!value || value.trim().length === 0) {
    return { valid: false, error: `${fieldName} is required` }
  }
  return { valid: true }
}

/**
 * Validate minimum length
 */
export function validateMinLength(
  value: string,
  minLength: number,
  fieldName: string
): ValidationResult {
  if (value.length < minLength) {
    return { valid: false, error: `${fieldName} must be at least ${minLength} characters` }
  }
  return { valid: true }
}

/**
 * Validate maximum length
 */
export function validateMaxLength(
  value: string,
  maxLength: number,
  fieldName: string
): ValidationResult {
  if (value.length > maxLength) {
    return { valid: false, error: `${fieldName} must be at most ${maxLength} characters` }
  }
  return { valid: true }
}

/**
 * Validate numeric value range
 */
export function validateRange(
  value: number,
  min: number,
  max: number,
  fieldName: string
): ValidationResult {
  if (value < min || value > max) {
    return { valid: false, error: `${fieldName} must be between ${min} and ${max}` }
  }
  return { valid: true }
}

/**
 * Validate business registration number
 */
export function validateBusinessRegistration(regNumber: string): ValidationResult {
  // Basic validation - non-empty and alphanumeric with dashes/slashes
  const pattern = /^[A-Z0-9\-/]+$/i
  if (!pattern.test(regNumber)) {
    return { valid: false, error: 'Invalid business registration format' }
  }
  return { valid: true }
}

/**
 * Validate country code
 */
export function validateCountryCode(code: string, validCodes: string[]): ValidationResult {
  if (!validCodes.includes(code)) {
    return { valid: false, error: 'Invalid country selection' }
  }
  return { valid: true }
}

/**
 * Batch validation - returns first error found
 */
export function validateAll(...validations: ValidationResult[]): ValidationResult {
  for (const validation of validations) {
    if (!validation.valid) {
      return validation
    }
  }
  return { valid: true }
}

/**
 * Validate proof-of-value type
 */
export function validatePOVType(type: string): ValidationResult {
  const validTypes = ['work', 'trade', 'trust', 'contribution']
  if (!validTypes.includes(type)) {
    return { valid: false, error: 'Invalid proof-of-value type' }
  }
  return { valid: true }
}

/**
 * Validate amount (positive number)
 */
export function validateAmount(amount: number): ValidationResult {
  if (amount <= 0) {
    return { valid: false, error: 'Amount must be greater than zero' }
  }
  if (!Number.isFinite(amount)) {
    return { valid: false, error: 'Invalid amount' }
  }
  return { valid: true }
}
