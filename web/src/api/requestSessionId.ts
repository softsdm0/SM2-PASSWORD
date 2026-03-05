import axios, { AxiosError } from 'axios'
import {
  generateRandomSm4Key,
  encryptWithSM2,
  encryptWithSM4,
  decryptWithSM4,
} from '../utils/cryptoUtils'
import { useSessionIdStore } from '@/stores/sessionId'
import { ElMessage } from 'element-plus'
import type { IApiResponse } from '@/types/password'

// 定义请求体结构
interface RequestBody {
  Sm4Key: string
  Msg: string
}

// 定义响应数据结构
interface EncryptedResponse {
  data: string
  message: string
  status: number
}

// 定义解密后的数据结构
interface DecryptedResponse {
  SessionId: string
}

// 请求会话 ID
export const requestSessionId = async (): Promise<DecryptedResponse> => {
  // 获取sm4Key
  const sm4Key = generateRandomSm4Key()

  // 构造请求体
  const requestBody: RequestBody = {
    Sm4Key: sm4Key,
    Msg: encryptWithSM4('ok.', sm4Key), // 使用sm4加密ok.作为Msg
  }

  // 转为json格式
  const jsonBody = JSON.stringify(requestBody)
  // 使用sm2公钥加密
  const encryptedBody = encryptWithSM2(jsonBody)

  try {
    // 使用sm2加密后的密文请求到后端
    const response = await axios.post<EncryptedResponse>('/public/api/session/id', encryptedBody)
    if (response.data.status !== 0) {
      console.log('请求SessionId失败:', response.data.message, response.data.status)
      throw new Error(response.data.status.toString())
    }
    // 后端返回的是sm4加密后的密文
    const encryptedData = response.data.data
    // 使用sm4key解密
    const decryptedData = decryptWithSM4(encryptedData, sm4Key)
    // 清理解密后的数据，去除额外的字符
    const cleanedData = decryptedData.trim()
    // 解释为json格式，并且转为DecryptedResponse类型
    const decryptedResponse = JSON.parse(cleanedData) as DecryptedResponse
    // 设置缓存
    useSessionIdStore().setSessionIdAndSm4Key(decryptedResponse.SessionId, sm4Key)

    return decryptedResponse
  } catch (error: AxiosError | any) {
    ElMessage.error('请求SessionId接口失败')
    console.log('请求SessionId失败:', error)
    if (error instanceof AxiosError) {
      if (error.response?.data as EncryptedResponse) {
        const respData: EncryptedResponse = error.response?.data
        ElMessage.error('后端返回: ' + respData.message + ', 错误码: ' + respData.status)
      } else {
        ElMessage.error(
          '错误信息: ' + error.response?.statusText + ', HttpCode: ' + error.response?.status
        )
      }
    }
    return { SessionId: '' }
  }
}
