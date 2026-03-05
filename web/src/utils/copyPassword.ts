import { IPassword } from '@/types/password'
import { usePasswordStore } from '../stores/password'
import { decryptPassword } from './crypto'
import { ElMessage } from 'element-plus'

export const handleCopyPassword = async (id: number) => {
  try {
    const data: IPassword = await usePasswordStore().getPassword(id)
    let decryptedPassword = decryptPassword(data.Password)
    if (!decryptedPassword) {
      decryptedPassword = '密码异常请重新设置密码'
    }

    // Check if the Clipboard API is supported
    if (!navigator.clipboard) {
      ElMessage.error('当前浏览器不支持复制功能，请使用支持该功能的浏览器。')
      return
    }
    await navigator.clipboard.writeText(decryptedPassword)
    ElMessage.success('已复制')
  } catch (error) {
    ElMessage.error('复制密码失败')
    console.log(error)
  }
}
