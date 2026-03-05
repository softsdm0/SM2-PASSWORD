export interface IPasswordInfo {
  ID: string
  UserId: string
  AppName: string
  AccountType: string
  Account: string
  Password?: string
  Url: string
  Notes?: string
  PasswordStrength: number
  CreatedAt: string
  UpdatedAt: string
}

export interface ICreatePasswordForm {
  AppName: string
  AccountType: string
  Account: string
  Password: string
  Notes?: string
}

export interface IPassword {
  ID: string
  Password: string
}

export interface IApiResponse<T = any> {
  data: T
  message: string
  status: number
}

// 定义请求状态枚举类型
export enum GetPasswordRecordsStatus {
  Unknown = 0,
  Success = 1,
  Refused = 2,
  NotFound = 3,
}

export interface IPasswordRecords {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  // 请求者id
  UserId: string
  // 密码记录id
  PasswordId: string
  PasswordInfo: IPasswordInfo
  // 请求状态
  Status: GetPasswordRecordsStatus
  // ip地址
  IP: string
}
