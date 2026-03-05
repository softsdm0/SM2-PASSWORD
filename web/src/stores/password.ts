import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { IApiResponse, IPassword, IPasswordInfo, IPasswordRecords } from '../types/password'
import request from '../api/request'
import { AxiosError, type AxiosResponse } from 'axios'

export const usePasswordStore = defineStore('password', () => {
  const passwordList = ref<IPasswordInfo[]>([])
  const loading = ref(false)

  async function fetchPasswordList(): Promise<void> {
    try {
      loading.value = true
      const resp: AxiosResponse<IApiResponse<IPasswordInfo[]>> = await request.get('/password')
      passwordList.value = resp.data.data
    } catch (error) {
      console.error('获取密码列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  async function deletePassword(id: string): Promise<void> {
    await request.delete(`/password/${id}`)
    await fetchPasswordList()
  }

  async function createPassword(
    data: Omit<IPasswordInfo, 'ID' | 'UserId' | 'CreatedAt' | 'UpdatedAt' | 'PasswordStrength'>
  ): Promise<void> {
    await request.post('/password', data)
    await fetchPasswordList()
  }

  async function updatePassword(id: string, data: Partial<IPasswordInfo>): Promise<void> {
    await request.post(`/password/${id}`, data)
    await fetchPasswordList()
  }

  async function getPassword(id: string): Promise<IPassword> {
    const resp: AxiosResponse<IApiResponse<IPassword>> = await request.get(`/password/${id}`)
    return resp.data.data
  }

  async function getPasswordRecords(): Promise<IPasswordRecords[]> {
    const resp: AxiosResponse<IApiResponse<IPasswordRecords[]>> = await request.get(
      `/password/record`
    )
    return resp.data.data
  }

  return {
    passwordList,
    loading,
    fetchPasswordList,
    deletePassword,
    createPassword,
    updatePassword,
    getPassword,
    getPasswordRecords,
  }
})
