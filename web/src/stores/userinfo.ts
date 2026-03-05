import { defineStore } from 'pinia'
import axios from 'axios'
import { ref } from 'vue'
import type { UserInfo } from '@/types/userInfo'

// 定义用户信息的 store
export const useUserInfoStore = defineStore('userInfo', () => {
  const userInfo = ref<UserInfo>()
  const loading = ref(false)
  const error = ref<string | null>(null) // 添加 error 属性

  const fetchUserInfo = async () => {
    loading.value = true
    error.value = null // 开始请求时清空错误信息
    try {
      // 这里需要替换为实际的 API 地址
      const response = await axios.get<UserInfo>('/public/api/userinfo')
      userInfo.value = response.data
    } catch (err) {
      error.value =
        '获取用户信息失败，请稍后重试，错误码:' +
        userInfo.value?.status +
        '，错误信息:' +
        userInfo.value?.message // 捕获错误并赋值给 error 属性
    } finally {
      loading.value = false
    }
  }

  return {
    userInfo,
    loading,
    error, // 返回 error 属性
    fetchUserInfo,
  }
})
