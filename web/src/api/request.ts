import axios, { AxiosError } from 'axios'
import type {
  AxiosInstance,
  AxiosRequestConfig,
  InternalAxiosRequestConfig,
  AxiosResponse,
} from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'
import { useLoadingStore } from '@/stores/loading'
import { useSessionIdStore } from '@/stores/sessionId'
import { decryptWithSM4, encryptWithSM4 } from '@/utils/cryptoUtils'
import type { IApiResponse } from '@/types/password'

const service: AxiosInstance = axios.create({
  baseURL: '/public/api',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求计数器
let requestCount = 0

const loadingStore = useLoadingStore()

const startLoading = () => {
  if (requestCount === 0) {
    loadingStore.startLoading()
  }
  requestCount++
}

const endLoading = () => {
  requestCount--
  if (requestCount === 0) {
    loadingStore.endLoading()
  }
}

service.interceptors.request.use(
  async (config: InternalAxiosRequestConfig): Promise<InternalAxiosRequestConfig> => {
    startLoading()
    // 不放这个了请求头太长了
    // const token = localStorage.getItem('token')
    // if (token && config.headers) {
    //   config.headers.Authorization = `Bearer ${token}`
    // }
    const sessionId = await useSessionIdStore().getSessionId()
    if (sessionId && config.headers) {
      config.headers.set('password-session-id', sessionId)
    }
    if (config.data) {
      const jsonData = JSON.stringify(config.data)
      // body加密
      const encryptedData = encryptWithSM4(jsonData, useSessionIdStore().Sm4Key)
      // 塞回去body
      config.data = encryptedData
    }

    return config
  },
  (error) => {
    console.log(error)
    endLoading()
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  (response: AxiosResponse) => {
    endLoading()
    const { data } = response

    if (data) {
      // data不为空解密
      const decryptData = decryptHandler(data)
      response.data = decryptData
      // console.log(decryptData)
      return response
    }

    return response
  },
  async (error: AxiosError) => {
    const prefix = 'RespErrorHandler:'
    endLoading()
    // 如果返回的是明文结构体
    if (error.response as AxiosResponse<IApiResponse>) {
      const resp = error.response as AxiosResponse<IApiResponse>

      // 登录过期
      if (resp.status === 401) {
        console.log(prefix, '登录已过期，请重新登录')
        ElMessage.error('登录已过期，请重新登录')
        localStorage.removeItem('token')
        router.push('/login')
      }

      // sessionid和sm4key过期
      if (error.response?.status === 400 && resp.data.status === 30001) {
        console.log(prefix, 'sessionid已经过期请重新发起请求')
        return await requestSessionIdAndRequestApi(error)
      }
    }

    // 如果返回的是密文string，就拿去解密
    if (error.response?.data !== '') {
      const decryptData = decryptHandler(error.response?.data as string) as IApiResponse
      console.log(prefix, decryptData)

      if (decryptData as IApiResponse) {
        ElMessage.error('后端返回: ' + decryptData.message + ', 错误码: ' + decryptData.status)
      } else {
        ElMessage.error(
          '错误信息: ' + error.response?.statusText + ', HttpCode: ' + error.response?.status
        )
      }
      return Promise.reject(error)
    }

    ElMessage.error(
      '错误信息: ' + error.response?.statusText + ', HttpCode: ' + error.response?.status
    )
    console.log(prefix, error)
    return Promise.reject(error)
  }
)

// 重新获取sessionid和sm4key,然后重新发起请求
async function requestSessionIdAndRequestApi(error: AxiosError): Promise<AxiosResponse> {
  const prefix = 'SessionIdHandler:'

  // 如果body里面有内容, 就先解密body:w
  if (error.config.data) {
    const data = error.config.data.replace(/"/g, '')
    // 把body拿出来解密后再丢回去body里面，不然会出现二次加密的问题。
    const decryptData = decryptWithSM4(data, useSessionIdStore().Sm4Key)
    error.config.data = JSON.parse(decryptData)
  }

  // 清理原有的sessionid和sm4key
  useSessionIdStore().clearSessionIdAndSm4Key()
  try {
    // 重新获取sessionid
    await useSessionIdStore().getSessionId()
    console.log(prefix, 'sessionid重新获取成功')
  } catch (error) {
    console.log(prefix, 'sessionid重新获取失败')
    return Promise.reject(error)
  }
  try {
    console.log(prefix, (error.config?.baseURL as string) + error.config?.url, '重新发起请求')
    // 重新发起请求
    const ar = await service.request(error.config as InternalAxiosRequestConfig)
    console.log(prefix, (error.config?.baseURL as string) + error.config?.url, '重新发起请求成功')
    return ar
  } catch (error: AxiosError | any) {
    console.log(prefix, (error.config?.baseURL as string) + error.config?.url, '重新发起请求失败')
    return Promise.reject(error)
  }
}

function decryptHandler(data: string): IApiResponse {
  // 如果数据是空的，就返回空
  if (data == '') {
    return {
      data: '',
      status: -90001,
      message: '输入的数据为空',
    }
  }

  const prefix = 'DecryptHandler:'
  try {
    // console.log(prefix, '开始解密')
    // 解密
    const decryptedData = decryptWithSM4(data, useSessionIdStore().Sm4Key)
    // console.log(prefix, '解密成功')

    // 清理解密后的数据，去除额外的字符
    const cleanedData = decryptedData.trim()
    // 解释为json格式，并且转为IPassword类型
    const decryptedResponse = JSON.parse(cleanedData) as IApiResponse

    return decryptedResponse
  } catch (error) {
    console.log(prefix, '解密失败, ', error)
    // 清空sessionid和sm4key
    useSessionIdStore().clearSessionIdAndSm4Key()
  }
  return {
    data: '',
    status: -90002,
    message: '解密失败',
  }
}

export default service
