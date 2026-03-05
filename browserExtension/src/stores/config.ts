import axios from "axios";
import type { IConfig } from "@/types/config";
import { LocalConfig, Storges } from "./storages";

// 定义用户信息的 store
export const useConfigStore = () => {

  const getConfig = async (): Promise<IConfig | null> => {
    const Config = await Storges().Config.value;
    const PasswordServerUrl = await Storges().PasswordServerUrl.value;
    
    if (!Config || Config === null) {
      const resp = await axios.get<IConfig>(`${PasswordServerUrl}/config.json`);

      await setConfig(resp.data);
      // 重新获取一次缓存
      const Config = await Storges().Config.value;
      return Config;
    }

    return Config;
  };

  const setConfig = async (config: IConfig) => {
    await Storges().Config.setItem(config);
  };

  const removeConfig = async () => {
    await Storges().Config.removeItem();
  };

  return {
    getConfig,
    setConfig,
    removeConfig,
  };
};
