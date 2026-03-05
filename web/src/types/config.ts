export interface IConfig {
  sm2_public_key_hex: string
  casdoor: {
    clientId: string
    organizationName?: string
    appName?: string
    serverUrl?: string
    redirectPath?: string
  }
}
