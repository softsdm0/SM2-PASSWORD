<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="logo-container">
        <img src="@/assets/passwordserver.svg" class="icon" alt="logo" />
      </div>

      <el-form :model="loginForm" @submit.prevent="handleLogin">
        <el-form-item>
          <el-input v-model="loginForm.identity" placeholder="请输入用户名/邮箱" prefix-icon="User" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" block> 登录 </el-button>
        </el-form-item>
      </el-form>
      <el-button type="primary" native-type="submit" @click="casdoorLogin()">
        casdoor 登录
      </el-button>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/api/request'
import { IApiResponse} from '@/types/password'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const loginForm = reactive({
  identity: '',
  password: '',
})

const handleLogin = async () => {
  try {
    loading.value = true
    const res:IApiResponse  = await request.post('/auth/login', loginForm)
    console.log(res);
    
    userStore.setToken(res.data.token)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}

// casdoor登录
import { onMounted, getCurrentInstance } from 'vue'
// import backend from '@/backend/backend'

const instance = getCurrentInstance()
const casdoorLogin = () => {
  window.location.href = instance.proxy.getSigninUrl()
}

onMounted(() => {
    //get account
    // backend.getAccount().then((res) => {
    //   if (res['status'] === 'ok') {
    //     console.log('success:', res)
    //   } else {
    //     console.log('fail:', res)
    //   }
    // })
  })
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
}

.login-card {
  width: 400px;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.logo-container {
  text-align: center;
  /* 让图片水平居中 */
  margin-bottom: 20px;
  /* 为图片下方添加间距 */
}

.icon {
  width: 150px;
  /* 设置图片宽度 */
  height: auto;
  /* 保持图片高度自适应 */
  transition: transform 0.3s ease;
  /* 添加过渡效果 */
}

.icon:hover {
  transform: scale(1.1);
  /* 悬停时放大图片 */
}
</style>
