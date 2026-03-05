import axios, { AxiosError } from "axios";
import type {
  AxiosInstance,
  AxiosRequestConfig,
  InternalAxiosRequestConfig,
  AxiosResponse,
} from "axios";
import { useSessionIdStore } from "@/stores/sessionId";
import { decryptWithSM4, encryptWithSM4 } from "@/utils/cryptoUtils";
import type { IApiResponse } from "../../../web/src/types/password";
import { Storges } from "@/stores/storages";
import { useTokenStore } from "@/stores/token";
import { useUserInfoStore } from "@/stores/userinfo";

const service: AxiosInstance = axios.create({
  // baseURL: `${(await Storges()).PasswordServerUrl}/public/api`,
  timeout: 5000,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求计数器
let requestCount = 0;

service.interceptors.request.use(
  async (
    config: InternalAxiosRequestConfig
  ): Promise<InternalAxiosRequestConfig> => {
    const url = await Storges().PasswordServerUrl.value;
    if (!url || url === "") {
      throw new Error("未设置PasswordServerUrl, 无法请求后端api。");
    }
    // 设置baseurl
    config.baseURL = `${url}/public/api`;
    // 不放这个了请求头太长了
    // const token = localStorage.getItem('token')
    // if (token && config.headers) {
    //   config.headers.Authorization = `Bearer ${token}`
    // }
    const sessionId = await useSessionIdStore().getSessionId();
    if (sessionId && config.headers) {
      config.headers.set("password-session-id", sessionId);
    }
    if (config.data) {
      const jsonData = JSON.stringify(config.data);
      const sm4k = await useSessionIdStore().getSm4Key();
      // body加密
      const encryptedData = encryptWithSM4(jsonData, sm4k);
      // 塞回去body
      config.data = encryptedData;
    }

    return config;
  },
  (error: any) => {
    console.log(error);
    return Promise.reject(error);
  }
);

service.interceptors.response.use(
  async (response: AxiosResponse) => {
    const { data } = response;

    if (data) {
      // data不为空解密
      const decryptData = await decryptHandler(data);
      response.data = decryptData;
      // console.log(decryptData)
      return response;
    }

    return response;
  },
  async (error: AxiosError | any) => {
    const prefix = "RespErrorHandler:";
    if (error === "未设置PasswordServerUrl, 无法请求后端api。") {
      console.log(prefix, error);
      return;
    }
    // 如果返回的是明文结构体
    if (error.response as AxiosResponse<IApiResponse>) {
      const resp = error.response as AxiosResponse<IApiResponse>;

      // 登录过期
      if (resp.status === 401) {
        console.log(prefix, "登录已过期，请重新登录");
        // 清理token
        (await useTokenStore()).removeToken();
        // 更新登录状态
        (await useUserInfoStore()).fetchUserInfo();
        return;
      }

      // sessionid和sm4key过期
      if (error.response?.status === 400 && resp.data.status === 30001) {
        console.log(prefix, "sessionid已经过期请重新发起请求");
        return await requestSessionIdAndRequestApi(error);
      }
    }

    // 如果返回的是密文string，就拿去解密
    if (error.response?.data !== "") {
      const decryptData = (await decryptHandler(
        error.response?.data as string
      )) as IApiResponse;
      console.log(prefix, decryptData);

      if (decryptData as IApiResponse) {
        console.log(
          "后端返回: " + decryptData.message + ", 错误码: " + decryptData.status
        );
      } else {
        console.log(
          "错误信息: " +
            error.response?.statusText +
            ", HttpCode: " +
            error.response?.status
        );
      }
      return Promise.reject(error);
    }

    console.log(
      "错误信息: " +
        error.response?.statusText +
        ", HttpCode: " +
        error.response?.status
    );
    console.log(prefix, error);
    return Promise.reject(error);
  }
);

// 重新获取sessionid和sm4key,然后重新发起请求
async function requestSessionIdAndRequestApi(
  error: AxiosError
): Promise<AxiosResponse> {
  const prefix = "SessionIdHandler:";

  // 如果body里面有内容, 就先解密body:w
  if (error.config?.data) {
    const data = error.config.data.replace(/"/g, "");
    // 把body拿出来解密后再丢回去body里面，不然会出现二次加密的问题。
    const decryptData = decryptWithSM4(
      data,
      await useSessionIdStore().getSm4Key()
    );
    error.config.data = JSON.parse(decryptData);
  }

  // 清理原有的sessionid和sm4key
  await useSessionIdStore().clearSessionIdAndSm4Key();
  try {
    // 重新获取sessionid
    await useSessionIdStore().getSessionId();
    console.log(prefix, "sessionid重新获取成功");
  } catch (error) {
    console.log(prefix, "sessionid重新获取失败");
    return Promise.reject(error);
  }
  try {
    console.log(
      prefix,
      (error.config?.baseURL as string) + error.config?.url,
      "重新发起请求"
    );
    // 重新发起请求
    const ar = await service.request(
      error.config as InternalAxiosRequestConfig
    );
    console.log(
      prefix,
      (error.config?.baseURL as string) + error.config?.url,
      "重新发起请求成功"
    );
    return ar;
  } catch (error: AxiosError | any) {
    console.log(
      prefix,
      (error.config?.baseURL as string) + error.config?.url,
      "重新发起请求失败"
    );
    return Promise.reject(error);
  }
}

async function decryptHandler(data: string): Promise<IApiResponse> {
  // 如果数据是空的，就返回空
  if (data === "") {
    return {
      data: "",
      status: -90001,
      message: "输入的数据为空",
    };
  }

  const prefix = "DecryptHandler:";
  try {
    const sm4key = await useSessionIdStore().getSm4Key();
    // 解密
    const decryptedData = decryptWithSM4(data, sm4key);

    // 清理解密后的数据，去除额外的字符
    const cleanedData = decryptedData.trim();
    // 解释为json格式，并且转为IPassword类型
    const decryptedResponse = JSON.parse(cleanedData) as IApiResponse;

    return decryptedResponse;
  } catch (error) {
    console.log(prefix, "解密失败, ", error);
    // 清空sessionid和sm4key
    await useSessionIdStore().clearSessionIdAndSm4Key();
  }
  return {
    data: "",
    status: -90002,
    message: "解密失败",
  };
}

export default service;
