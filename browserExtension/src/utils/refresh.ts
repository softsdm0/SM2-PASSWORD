import { requestSessionId } from "@/api/requestSessionId";
import { usePasswordStore } from "@/stores/password";
import { useUserInfoStore } from "@/stores/userinfo";

export const refreshAll = async () => {
  // 刷新sessionId
  await requestSessionId().then(async() => {
    // 刷新用户信息
    await useUserInfoStore().fetchUserInfo();
    // 刷新密码列表
    await usePasswordStore().fetchPasswordList();

  });
};
