export interface PasswordValidationResult {
  isValid: boolean
  errors: string[]
}

export const validatePassword = (password: string): PasswordValidationResult => {
  const errors: string[] = []

  if (password.length < 8) {
    errors.push('密码长度至少为8位')
  }

  if (!/[A-Z]/.test(password)) {
    errors.push('密码必须包含大写字母')
  }

  if (!/[a-z]/.test(password)) {
    errors.push('密码必须包含小写字母')
  }

  if (!/[0-9]/.test(password)) {
    errors.push('密码必须包含数字')
  }

  if (!/[^A-Za-z0-9]/.test(password)) {
    errors.push('密码必须包含特殊字符')
  }

  return {
    isValid: errors.length === 0,
    errors,
  }
}
