import { requestSessionId } from "@/api/requestSessionId";
import { LocalSessionId, LocalSm4Key, Storges } from "@/stores/storages";

export const useSessionIdStore = () => {
  async function getSessionId() {
    const SessionId = await Storges().SessionId.value;
    if (!SessionId || SessionId === "") {
      const decryptedResponse = await requestSessionId();
      return decryptedResponse.SessionId;
    }
    return SessionId || "";
  }

  async function getSm4Key() {
    const Sm4Key = await Storges().Sm4Key.value;
    if (!Sm4Key || Sm4Key === "") {
      await requestSessionId();
    }
    return Sm4Key || "";
  }

  async function setSessionIdAndSm4Key(sessionId: string, sm4Key: string) {
    await Storges().SessionId.setItem(sessionId);
    await Storges().Sm4Key.setItem(sm4Key);
  }

  async function clearSessionIdAndSm4Key() {
    await Storges().SessionId.removeItem();
    await Storges().Sm4Key.removeItem();
  }

  return {
    getSessionId,
    getSm4Key,
    setSessionIdAndSm4Key,
    clearSessionIdAndSm4Key,
  };
};
