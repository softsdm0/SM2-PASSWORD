import { requestSessionId } from '@/api/requestSessionId'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSessionIdStore = defineStore('sessionId', () => {
  const SessionId = ref(localStorage.getItem('sessionId') || '')
  const Sm4Key = ref(localStorage.getItem('sm4Key') || '')

  async function getSessionId() {
    if (SessionId.value === '') {
      await requestSessionId()
    }
    return SessionId.value
  }

  async function getSm4Key() {
    if (Sm4Key.value === '') {
      await requestSessionId()
    }
    return Sm4Key.value
  }

  function setSessionIdAndSm4Key(sessionId: string, sm4Key: string) {
    SessionId.value = sessionId
    Sm4Key.value = sm4Key
    localStorage.setItem('sessionId', sessionId)
    localStorage.setItem('sm4Key', sm4Key)
  }

  function clearSessionIdAndSm4Key() {
    SessionId.value = ''
    Sm4Key.value = ''
    localStorage.removeItem('sessionId')
    localStorage.removeItem('sm4Key')
  }

  return {
    SessionId,
    Sm4Key,
    getSessionId,
    getSm4Key,
    setSessionIdAndSm4Key,
    clearSessionIdAndSm4Key,
  }
})
