import { StorageItemKey } from "wxt/utils/storage";
import { IPasswordInfo } from "@/types/password";
import { IConfig } from "@/types/config";
import { UserInfo, LoginStats } from "@/types/userInfo";

// PasswordServer后端地址
export const LocalPasswordServerUrl: StorageItemKey = "local:PasswordServerUrl";

// 密码信息列表
export const LocalPasswordInfoList: StorageItemKey = "local:PasswordInfoList";

// config
export const LocalConfig: StorageItemKey = "local:Config";

// sm2PublicKey
export const LocalSm2PublicKey: StorageItemKey = "local:Sm2PublicKey";

// sessionId
export const LocalSessionId: StorageItemKey = "local:SessionId";

// sm4Key
export const LocalSm4Key: StorageItemKey = "local:Sm4Key";

// token登录状态
export const LocalToken: StorageItemKey = "local:Token";

// 登录状态
export const LocalLoginStatus: StorageItemKey = "local:LoginStatus";

// 用户信息
export const LocalUserInfo: StorageItemKey = "local:UserInfo";

// 获取存储
export const Storges = () => {
  const PasswordServerUrl = CreateStorage<string>(LocalPasswordServerUrl, "");
  const PasswordInfoList = CreateStorage<IPasswordInfo[]>(
    LocalPasswordInfoList
  );
  const Config = CreateStorage<IConfig>(LocalConfig);
  const Sm2PublicKey = CreateStorage<string>(LocalSm2PublicKey, "");
  const SessionId = CreateStorage<string>(LocalSessionId, "");
  const Sm4Key = CreateStorage<string>(LocalSm4Key, "");
  const Token = CreateStorage<string>(LocalToken, "");
  const LoginStatus = CreateStorage<LoginStats>(LocalLoginStatus, "initial");
  const UserInfo = CreateStorage<UserInfo>(LocalUserInfo);

  return {
    PasswordServerUrl,
    PasswordInfoList,
    Config,
    Sm2PublicKey,
    SessionId,
    Sm4Key,
    Token,
    LoginStatus,
    UserInfo,
  };
};

// 分别生成vue使用的，和普通storage使用的
const CreateStorage = <T>(key: StorageItemKey, defaultValue?: T) => {
  const { state: vue } = useStoredValue<T>(key, defaultValue);
  const value = storage.getItem<T>(key, { defaultValue: defaultValue });
  const setItem = (v: T) => storage.setItem<T>(key, v);
  const removeItem = () => storage.removeItem(key);
  return {
    vue,
    value,
    setItem,
    removeItem,
  };
};
