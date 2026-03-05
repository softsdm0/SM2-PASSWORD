import { defineStore } from 'pinia'
import axios from 'axios'
import type { IConfig } from '@/types/config'

// 定义用户信息的 store
export const useConfigStore = defineStore('config', () => {
  let config: IConfig

  const getConfig = async (): Promise<IConfig> => {
    if (!config) {
      const resp = await axios.get<IConfig>('/config.json')
      config = resp.data

      // docker-compose 一键部署时，SM2 公钥和 Casdoor 内置应用配置由后端动态提供。
      try {
        const bootstrapResp = await axios.get<{
          data?: {
            sm2PublicKey?: string
            casdoorClientId?: string
            casdoorOrganization?: string
            casdoorApplication?: string
          }
        }>('/public/api/session/bootstrapConfig')
        const data = bootstrapResp.data?.data
        if (typeof data?.sm2PublicKey === 'string' && data.sm2PublicKey.length > 0) {
          config = {
            ...config,
            sm2_public_key_hex: data.sm2PublicKey,
          }
        }
        if (config.casdoor) {
          if (typeof data?.casdoorClientId === 'string' && data.casdoorClientId.length > 0) {
            config.casdoor.clientId = data.casdoorClientId
          }
          if (typeof data?.casdoorOrganization === 'string' && data.casdoorOrganization.length > 0) {
            config.casdoor.organizationName = data.casdoorOrganization
          }
          if (typeof data?.casdoorApplication === 'string' && data.casdoorApplication.length > 0) {
            config.casdoor.appName = data.casdoorApplication
          }
        }
      } catch (err) {
        console.warn('load bootstrap config failed, fallback to /config.json', err)
      }
    }
    return config
  }

  return {
    config,
    getConfig,
  }
})
