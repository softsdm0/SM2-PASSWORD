import { saveAs } from 'file-saver'
import type { IPasswordInfo } from '@/types/password'

export const exportPasswords = (passwords: IPasswordInfo[]) => {
  const data = JSON.stringify(passwords, null, 2)
  const blob = new Blob([data], { type: 'application/json' })
  saveAs(blob, `passwords-${new Date().toISOString().split('T')[0]}.json`)
}

export const importPasswords = (file: File): Promise<IPasswordInfo[]> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const data = JSON.parse(e.target?.result as string)
        resolve(data)
      } catch (error) {
        reject(new Error('文件格式错误'))
      }
    }
    reader.onerror = () => reject(new Error('读取文件失败'))
    reader.readAsText(file)
  })
}
