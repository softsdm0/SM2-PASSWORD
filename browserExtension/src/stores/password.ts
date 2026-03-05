import { defineStore } from "pinia";
import { ref } from "vue";
import type {
  IApiResponse,
  IPassword,
  IPasswordInfo,
  IPasswordRecords,
} from "@/types/password";
import request from "@/api/request";
import { AxiosError, type AxiosResponse } from "axios";
import { LocalPasswordInfoList, Storges } from "./storages";

export const usePasswordStore = () => {
  async function fetchPasswordList(): Promise<void> {
    try {
      const resp: AxiosResponse<IApiResponse<IPasswordInfo[]>> =
        await request.get("/password");
      await storage.setItem(LocalPasswordInfoList, resp.data.data);
    } catch (error) {
      console.error("获取密码列表失败:", error);
    } finally {
    }
  }

  async function getPasswordList() {
    const passwordList = await Storges().PasswordInfoList.value;
    if (!passwordList && passwordList === null) {
      await fetchPasswordList();
      // 重新加载一次缓存
      return await Storges().PasswordInfoList.value;
    }
    return passwordList;
  }

  async function deletePassword(id: number): Promise<void> {
    await request.delete(`/password/${id}`);
    await fetchPasswordList();
  }

  async function createPassword(
    data: Omit<
      IPasswordInfo,
      "ID" | "UserId" | "CreatedAt" | "UpdatedAt" | "PasswordStrength"
    >
  ): Promise<void> {
    await request.post("/password", data);
    await fetchPasswordList();
  }

  async function updatePassword(
    id: number,
    data: Partial<IPasswordInfo>
  ): Promise<void> {
    await request.post(`/password/${id}`, data);
    await fetchPasswordList();
  }

  async function getPassword(id: number): Promise<IPassword> {
    const resp: AxiosResponse<IApiResponse<IPassword>> = await request.get(
      `/password/${id}`
    );
    return resp.data.data;
  }

  async function getPasswordRecords(): Promise<IPasswordRecords[]> {
    const resp: AxiosResponse<IApiResponse<IPasswordRecords[]>> =
      await request.get(`/password/record`);
    return resp.data.data;
  }

  return {
    fetchPasswordList,
    getPasswordList,
    deletePassword,
    createPassword,
    updatePassword,
    getPassword,
    getPasswordRecords,
  };
};
