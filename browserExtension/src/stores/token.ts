import { LocalToken, Storges } from "./storages";

// 定义用户信息的 store
export const useTokenStore = () => {
  const getToken = async () => {
    const Token = await Storges().Token.value;
    return Token || "";
  };

  const setToken = async (token: string) => {
    await Storges().Token.setItem(token);
  };

  const removeToken = async () => {
    await Storges().Token.removeItem();
  };

  return {
    getToken,
    setToken,
    removeToken,
  };
};
