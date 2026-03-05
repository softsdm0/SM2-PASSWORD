import axios from "axios";
import { Storges } from "@/stores/storages";

interface Sm2PublicKeyResponse {
  data: string;
  message: string;
  status: number;
}
// 获取sm2公钥
export const requestSm2PublicKey = async () => {
  const PasswordServerUrl = await Storges().PasswordServerUrl.value;
  console.log("PasswordServerUrl:", PasswordServerUrl, "");

  if (!PasswordServerUrl || PasswordServerUrl === "") {
    console.log("未设置PasswordServerUrl");
    return;
  }

  const response = await axios.get<Sm2PublicKeyResponse>(
    `${PasswordServerUrl}/public/api/session/sm2PublicKey`
  );
  if (response.data.status !== 0) {
    console.log(
      "请求sm2PublickKey失败:",
      response.data.message,
      response.data.status
    );
    throw new Error(response.data.status.toString());
  }

  // 设置sm2PublickKey
  await Storges().Sm2PublicKey.setItem(response.data.data);

  return;
};
