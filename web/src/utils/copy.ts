import { ElMessage } from 'element-plus'

export const handleCopy = async (cp: any) => {
  try {
    if (!navigator.clipboard) {
      ElMessage.error('当前浏览器不支持复制功能，请使用支持该功能的浏览器。')
      return
    }
    await navigator.clipboard.writeText(cp)
    ElMessage.success('已复制')
  } catch (error) {
    ElMessage.error('复制密码失败')
    console.log(error)
  }
}
