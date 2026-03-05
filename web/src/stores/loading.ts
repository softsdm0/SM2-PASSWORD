import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const isLoading = ref(false)
  const loadingText = ref('加载中...')

  function startLoading(text?: string) {
    if (text) {
      loadingText.value = text
    }
    isLoading.value = true
  }

  function endLoading() {
    isLoading.value = false
    loadingText.value = '加载中...'
  }

  return {
    isLoading,
    loadingText,
    startLoading,
    endLoading,
  }
})
