import { sm2, sm4 } from 'sm-crypto'
import { useConfigStore } from '@/stores/config'

let sm2PublicKey = ''
// 获取 SM2 公钥
useConfigStore()
  .getConfig()
  .then((c) => {
    sm2PublicKey = c.sm2_public_key_hex
  })

// 使用 SM2 公钥加密数据
export const encryptWithSM2 = (data: string): string => {
  return sm2.doEncrypt(data, sm2PublicKey, 1)
}

// 生成随机的 SM4 密钥
export const generateRandomSm4Key = (): string => {
  const buffer = new Uint8Array(16)
  // 使用 crypto.getRandomValues 生成安全的随机字节
  crypto.getRandomValues(buffer)
  // 将字节数组转换为十六进制字符串
  return Array.from(buffer, (byte) => byte.toString(16).padStart(2, '0')).join('')
}

// 使用 SM4 密钥加密数据
export const encryptWithSM4 = (data: string, sm4Key: string): string => {
  return sm4.encrypt(data, sm4Key, {
    mode: 'ecb',
    padding: 'pkcs#7',
  })
}

// 使用 SM4 密钥解密数据
export const decryptWithSM4 = (encryptedData: string, sm4Key: string): string => {
  return sm4.decrypt(encryptedData, sm4Key, {
    mode: 'ecb',
    padding: 'pkcs#7',
  })
}
