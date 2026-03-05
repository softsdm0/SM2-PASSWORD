<template>
  <div class="callback-container">
    <h1 class="title">登录校验中</h1>
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading">
      <p>正在登录，请稍候...</p>
    </div>
    <!-- 错误信息 -->
    <div v-if="errorMessage" class="error">
      <p>{{ errorMessage }}</p>
    </div>
  </div>

</template>

<script setup lang="ts">
import { onMounted, getCurrentInstance, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const serverUrl = '/public'

// 加载状态
const isLoading = ref(false)
// 错误信息
const errorMessage = ref('')

function login() {
  isLoading.value = true
  const instance = getCurrentInstance()
  instance.proxy.signin(serverUrl).then((res) => {
    isLoading.value = false
    if (res.status === 0) {
      // 使用 ElMessage 提示登录成功
      ElMessage.success('Login success')
      userStore.setToken(res.data)
      router.push('/')
    } else {
      errorMessage.value = `Login failed: ${res.msg}`
      // 使用 ElMessage 提示登录失败
      ElMessage.error(`Login failed: ${res.msg}`)
      router.push('/login')
    }
  }).catch((err) => {
    isLoading.value = false
    errorMessage.value = `Login failed: ${err.message}`
    ElMessage.error(`Login failed: ${err.message}`)
    router.push('/login')
  })
}

onMounted(() => {
  login()
})
</script>

<style scoped>
.callback-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh; /* 让容器占满整个视口高度 */
  background-color: #f4f7fa; /* 背景颜色 */
}

.title {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
}

.loading {
  font-size: 16px;
  color: #666;
  animation: fadeIn 0.5s ease-in-out;
}

.error {
  font-size: 16px;
  color: #ff4d4f; /* 错误信息颜色 */
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>