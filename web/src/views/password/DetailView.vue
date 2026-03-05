<template>
  <div class="password-detail">
    <el-card>
      <template #header>
        <div class="card-header">
          <h3>密码详情</h3>
          <div class="action-buttons">
            <el-tooltip content="编辑" placement="top">
              <el-button type="warning" :icon="Edit" circle @click="handleEdit" />
            </el-tooltip>
            <el-tooltip content="删除" placement="top">
              <el-button type="danger" :icon="Delete" circle @click="handleDelete" />
            </el-tooltip>
            <el-button @click="handleBack">返回列表</el-button>
          </div>
        </div>
      </template>

      <el-descriptions :column="1" border>
        <el-descriptions-item label="应用名称">
          {{ formData.AppName }}
        </el-descriptions-item>

        <el-descriptions-item label="账户类型">
          {{ formData.AccountType }}
        </el-descriptions-item>

        <el-descriptions-item v-if="formData.Account" label="账户">
          {{ formData.Account }}
        </el-descriptions-item>
        <el-descriptions-item v-else label="账户">
          后端解密账号失败, 请重新设置账号
        </el-descriptions-item>


        <el-descriptions-item label="密码" class="password-item">
          <div class="password-text">
            <!-- 显示密码或占位符 -->
            <span v-if="showPassword && formData.Password">{{ formData.Password }}</span>
            <span v-else-if="showPassword && !formData.Password">后端解密密码失败, 请重新设置密码</span>
            <span v-else>******</span>
          </div>
          <!-- 操作按钮 -->
          <div class="password-button-wrapper">
            <el-button @click="togglePasswordVisibility">
              {{ showPassword ? '隐藏' : '显示' }}
            </el-button>
            <el-button @click="copyPassword">复制</el-button>
          </div>

        </el-descriptions-item>
        <!-- 新增显示 URL 字段 -->
        <el-descriptions-item label="网站">
          <template v-if="formData.Url">
            <a :href="formData.Url" target="_blank" rel="noopener noreferrer" class="custom-link">
              {{ formData.Url }}
            </a>
          </template>
          <template v-else>
            无
          </template>

        </el-descriptions-item>
        <!-- 新增显示密码强度 -->
        <el-descriptions-item label="密码强度">
          {{ getPasswordStrengthText(formData.PasswordStrength) }}
        </el-descriptions-item>

        <el-descriptions-item label="备注">
          {{ formData.Notes }}
        </el-descriptions-item>

        <el-descriptions-item label="创建时间">
          {{ formatDate(formData.CreatedAt) }}
        </el-descriptions-item>

        <el-descriptions-item label="更新时间">
          {{ formatDate(formData.UpdatedAt) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import dayjs from 'dayjs'
import type { IPasswordInfo } from '../../types/password'
import { usePasswordStore } from '../../stores/password'
import { ElMessage } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const passwordStore = usePasswordStore()

const formData = reactive({
  ID: '',
  AppName: '',
  AccountType: '',
  Account: '',
  Password: '',
  PasswordStrength: 0, // 新增字段用于存储密码强度
  Url: '',
  Notes: '',
  CreatedAt: '',
  UpdatedAt: '',
})

const showPassword = ref(false)

const togglePasswordVisibility = async () => {
  // 如果当前显示的是不是密码，从后端获取密码
  if (showPassword.value == false) {
    try {
      const data = await passwordStore.getPassword(formData.ID)
      formData.Password = data.Password
    } catch (error) {
      ElMessage.error('显示密码失败')
    }
  }
  // 改变是否显示密码
  showPassword.value = !showPassword.value
}

const copyPassword = async () => {
  try {
    if (!showPassword.value) {
      await togglePasswordVisibility()
    }
    if ('clipboard' in navigator) {
      await navigator.clipboard.writeText(formData.Password)
      ElMessage.success('密码已复制到剪贴板')
    } else {
      ElMessage.error('当前浏览器不支持复制功能')
    }
  } catch (error) {
    ElMessage.error('复制密码失败')
  }
}

// 加载密码数据
const loadPasswordData = async () => {
  try {
    // 查询密码
    let id: string;
    if (typeof route.params.id === 'string') {
      id = route.params.id;
    } else {
      // 处理无法转换的情况
      id = '';
      return
    }

    // 匹配密码记录信息
    const result = Array.isArray(passwordStore.passwordList) ? [...passwordStore.passwordList] : []

    // 是否匹配到
    let isExist = false
    // 匹配
    result.filter((item) => {
      if (item.ID === id) {
        // 赋值显示
        Object.assign(formData, {
          ID: item.ID,
          AppName: item.AppName,
          AccountType: item.AccountType,
          Account: item.Account,
          Notes: item.Notes,
          Password: '', // 不显示原密码
          Url: item.Url,
          PasswordStrength: item.PasswordStrength, // 从后端数据中获取密码强度
          CreatedAt: item.CreatedAt,
          UpdatedAt: item.UpdatedAt,
        })
        // 匹配到
        isExist = true
        return
      }
    })

    // 如果不存在
    if (!isExist) {
      ElMessage.error('密码记录不存在')
      router.push('/')
    }
  } catch (error) {
    ElMessage.error('加载数据失败')
    router.push('/')
  }
}

onMounted(() => {
  loadPasswordData()
})

const formatDate = (date: string) => {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-'
}

const handleBack = () => {
  router.push('/')
}

// 添加删除处理逻辑
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm('确认删除该密码记录吗？', '警告', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await passwordStore.deletePassword(formData.ID)
    ElMessage.success('删除成功')
    router.push('/')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 添加编辑跳转逻辑
const handleEdit = () => {
  router.push({ name: 'passwordEdit', params: { id: formData.ID } })
}

// 定义 PasswordStrength 类型和对应的常量
enum PasswordStrength {
  VeryWeak = 1,
  Weak,
  Medium,
  Strong
}

// 根据强度等级返回对应的文本
const getPasswordStrengthText = (strength: number) => {
  switch (strength) {
    case PasswordStrength.VeryWeak:
      return '非常弱';
    case PasswordStrength.Weak:
      return '弱';
    case PasswordStrength.Medium:
      return '中等';
    case PasswordStrength.Strong:
      return '强';
    default:
      return '未知';
  }
}
</script>

<style scoped>
.password-detail {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
}

.password-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.password-button-wrapper {
  margin-left: auto;
}

.password-text {

  margin: 10px 0px;
}

.custom-link {
  text-decoration: none;
  /* 去掉下划线 */
  color: blue;
  /* 设置文字颜色为蓝色 */
}

/* 可选：鼠标悬停时也不显示下划线 */
.custom-link:hover {
  text-decoration: none;
  color: blue;
  /* 鼠标悬停时保持蓝色 */
}
</style>
