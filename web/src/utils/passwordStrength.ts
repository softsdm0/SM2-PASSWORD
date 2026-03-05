export const checkPasswordStrength = (
  password: string,
): {
  score: number
  level: '弱' | '中' | '强'
  color: string
} => {
  let score = 0

  if (password.length >= 8) score += 1
  if (password.length >= 12) score += 1
  if (/[A-Z]/.test(password)) score += 1
  if (/[a-z]/.test(password)) score += 1
  if (/[0-9]/.test(password)) score += 1
  if (/[^A-Za-z0-9]/.test(password)) score += 1

  if (score <= 2) return { score, level: '弱', color: '#F56C6C' }
  if (score <= 4) return { score, level: '中', color: '#E6A23C' }
  return { score, level: '强', color: '#67C23A' }
}
