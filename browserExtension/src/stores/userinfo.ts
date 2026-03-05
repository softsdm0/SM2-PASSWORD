import axios, { AxiosError } from "axios";
import { type UserInfo } from "@/types/userInfo";
import { Storges } from "@/stores/storages";
import { useSessionIdStore } from "@/stores/sessionId";

// 定义用户信息的 store
export const useUserInfoStore = () => {
  const fetchUserInfo = async () => {
    try {
      const passwordServerUrl = await Storges().PasswordServerUrl.value;
      if (!passwordServerUrl || passwordServerUrl === "") {
        throw `未配置密码服务器地址, 请在设置中配置`;
      }
      // 这里需要替换为实际的 API 地址
      const response = await axios.get<UserInfo>(
        `${passwordServerUrl}/public/api/userinfo`
      );
      await Storges().UserInfo.setItem(response.data);
      await Storges().LoginStatus.setItem("logged");
      return response.data;
    } catch (err) {
      console.log("fetchUserInfo error:", err);
      if (err instanceof AxiosError) {
        if (err.response?.status === 401) {
          await Storges().LoginStatus.setItem("unlogged");
          throw "未登录";
        } else {
          await Storges().LoginStatus.setItem("error");
          throw `连接服务器失败, code: ${err.response?.status || "未知"}`;
        }
      }
      await Storges().LoginStatus.setItem("error");
      throw err;
    }
  };

  const logout = async () => {
    await Storges().UserInfo.removeItem();
    await Storges().LoginStatus.removeItem(); 
    await useSessionIdStore().clearSessionIdAndSm4Key()
    const PasswordServerUrl = await Storges().PasswordServerUrl.value;
    browser.cookies.remove({ url: PasswordServerUrl||'', name: 'casdoorAuthToken' });
    browser.cookies.remove({ url: PasswordServerUrl||'', name: 'authToken' });
    // 确认已经退出
    await fetchUserInfo();
  }

  return {
    fetchUserInfo,
    logout,
  };
};
