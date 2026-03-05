<template>
  <div class="edit-password">
    <el-card>
      <template #header>
        <div class="card-header">
          <h3>编辑密码</h3>
          <el-button type="primary" link @click="handleViewOriginalPassword">
            查看原密码
          </el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item label="应用名称" prop="AppName">
          <el-autocomplete v-model="formData.AppName" :fetch-suggestions="queryAppName" placeholder="请输入应用名称"
            @select="handleAppSelect" />
        </el-form-item>

        <el-form-item label="账户类型" prop="AccountType">
          <el-autocomplete v-model="formData.AccountType" :fetch-suggestions="queryAccountType" placeholder="请输入账户类型"
            @select="handleAppSelect" />
        </el-form-item>

        <el-form-item label="账户" prop="Account">
          <el-autocomplete v-if="formData.Account" v-model="formData.Account" :fetch-suggestions="queryAccount"
            placeholder="请输入账户" @select="handleAppSelect" />
          <el-autocomplete v-else v-model="formData.Account" :fetch-suggestions="queryAccount"
            placeholder="账户异常请求重新设置账号" @select="handleAppSelect" />
        </el-form-item>

        <el-form-item label="密码" prop="Password">
          <div class="password-input-group">
            <div class="password-input-wrapper">
              <el-input v-if="password" v-model="formData.Password" type="password" placeholder="请输入新密码（留空表示不修改）"
                show-password @input="checkStrength" />
              <el-input v-else v-model="formData.Password" type="password" placeholder="密码异常请重新设置密码" show-password
                @input="checkStrength" />
              <div v-if="formData.Password" class="password-strength">
                密码强度：<span :style="{ color: strengthInfo.color }">{{
                  strengthInfo.level
                }}</span>
              </div>
            </div>
            <el-button @click="showGenerateDialog"> 生成密码 </el-button>
          </div>
        </el-form-item>

        <el-form-item label="网站" prop="Url">
          <el-autocomplete v-model="formData.Url" :fetch-suggestions="queryUrl" placeholder="请输入 URL"
            @select="handleAppSelect" />
        </el-form-item>

        <el-form-item label="备注" prop="Notes">
          <el-input v-model="formData.Notes" type="textarea" placeholder="请输入备注" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading"> 保存 </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-dialog v-model="generateDialogVisible" title="生成密码" width="500px">
      <el-form :model="passwordOptions" label-width="120px">
        <el-form-item label="密码长度">
          <el-slider v-model="passwordOptions.length" :min="8" :max="32" />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="passwordOptions.includeUppercase">包含大写字母</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="passwordOptions.includeLowercase">包含小写字母</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="passwordOptions.includeNumbers">包含数字</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="passwordOptions.includeSymbols">包含特殊字符</el-checkbox>
        </el-form-item>
        <el-form-item v-if="previewPassword" label="预览">
          <div class="password-preview">
            <span>{{ previewPassword }}</span>
            <el-button type="primary" link :icon="DocumentCopy" @click="copyPassword(previewPassword)">
              复制
            </el-button>
          </div>
        </el-form-item>
        <el-button type="primary" @click="previewGeneratePassword"> 预览生成 </el-button>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="generateDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleGeneratePassword"> 生成并使用 </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { DocumentCopy } from '@element-plus/icons-vue'
import { checkPasswordStrength } from '../../utils/passwordStrength'
import { generatePassword, type PasswordOptions } from '../../utils/passwordGenerator'
import { ref, reactive, onMounted, watch, onBeforeUnmount, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { usePasswordStore } from '../../stores/password'
import { encryptPassword } from '../../utils/crypto'
import type { IPassword, IPasswordInfo } from '../../types/password'
import { ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()
const formRef = ref<FormInstance>()
const loading = ref(false)
const passwordStore = usePasswordStore()
let password = ref('')

let id: string;
if (typeof route.params.id === 'string') {
  id = route.params.id;
} else {
  // 处理无法转换的情况
  id = '';
}

const formData = reactive({
  ID: '',
  AppName: '',
  AccountType: '',
  Account: '',
  Password: '',
  Url: '',
  Notes: '',
})

// 加载密码数据
const loadPasswordData = async () => {
  try {
    // 查询密码

    const data: IPassword = await passwordStore.getPassword(id)

    // 记录一下密码
    password.value = data.Password

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
          Url: item.Url,
          Notes: item.Notes,
          Password: '', // 不显示原密码
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

onMounted(async () => {
  await loadPasswordData()
})

const rules = {
  AppName: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  AccountType: [{ required: true, message: '请输入账户类型', trigger: 'blur' }],
  Account: [{ required: true, message: '请输入账户', trigger: 'blur' }],
  Password: [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (!value) {
          callback()
          return
        }
        // const { isValid, errors } = validatePassword(value)
        // if (!isValid) {
        //   callback(new Error(errors.join('；')))
        // } else {
        //   callback()
        // }
      },
      trigger: 'blur',
    },
  ],
  // 添加 Url 验证规则
  Url: [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (!value) {
          // 允许为空
          callback();
          return;
        }
        const urlRegex = /^(http|https):\/\/.*$/;
        if (!urlRegex.test(value)) {
          callback(new Error('URL 必须以 http:// 或 https:// 开头'));
          return;
        }
        try {
          new URL(value);
          callback();
        } catch (error) {
          callback(new Error('请输入有效的 URL'));
        }
      },
      trigger: 'blur',
    },
  ],
}

// 添加表单是否被修改的标记
const isFormModified = ref(false)

// 监听表单数据变化
watch(
  formData,
  () => {
    isFormModified.value = true
  },
  { deep: true },
)

// 添加路由离开守卫
// onBeforeRouteLeave(async (to, from, next) => {
//   if (isFormModified.value) {
//     try {
//       await ElMessageBox.confirm(
//         '有未保存的修改，确定要离开吗？',
//         '提示',
//         {
//           confirmButtonText: '确定',
//           cancelButtonText: '取消',
//           type: 'warning'
//         }
//       )
//       next()
//     } catch {
//       next(false)
//     }
//   } else {
//     next()
//   }
// })

// 修改提交处理函数
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    if (formData.Url) {
      // 检查url是否符合格式
      if (checkUrl(formData.Url)) {
        ElMessage.error(checkUrl(formData.Url));
        return
      }
      // 处理 Url 结尾的 / 号
      formData.Url = formData.Url.replace(/\/+$/, '');
    }
    // await formRef.value.validate()
    loading.value = true
    const submitData = { ...formData }

    if (submitData.Password) {
      submitData.Password = encryptPassword(submitData.Password)
    } else {
      submitData.Password = password.value
    }

    await passwordStore.updatePassword(submitData.ID, submitData)
    isFormModified.value = false
    ElMessage.success('更新成功')
    router.push('/')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新失败:', error)
      ElMessage.error('更新失败：' + (error as Error).message)
    }
  } finally {
    loading.value = false
  }
}

const checkUrl = (url: string): string => {
  const urlRegex = /^(http|https):\/\/.*$/;
  if (!urlRegex.test(url)) {
    return 'URL 必须以 http:// 或 https:// 开头'
  }
  try {
    new URL(url);
  } catch (error) {
    return '请输入有效的 URL'
  }

  return
}

// 修改取消处理函数
const handleCancel = async () => {
  if (isFormModified.value) {
    try {
      await ElMessageBox.confirm('确定要放弃未保存的修改吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
      router.push('/')
    } catch {
      // 用户取消操作，不做任何处理
    }
  } else {
    router.push('/')
  }
}

const generateDialogVisible = ref(false)
const passwordOptions = reactive<PasswordOptions>({
  length: 16,
  includeUppercase: true,
  includeLowercase: true,
  includeNumbers: true,
  includeSymbols: true,
})

const showGenerateDialog = () => {
  generateDialogVisible.value = true
}

const previewPassword = ref('')

const previewGeneratePassword = () => {
  previewPassword.value = generatePassword(passwordOptions)
}

const copyPassword = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('密码已复制到剪贴板')
  } catch (err) {
    ElMessage.error('复制失败')
  }
}

// 新增的补全部分
const strengthInfo = ref({
  level: '弱',
  color: 'red',
})

const checkStrength = () => {
  const password = formData.Password
  const { level, color } = checkPasswordStrength(password)
  strengthInfo.value = { level, color }
}

const handleGeneratePassword = () => {
  const password = previewPassword.value || generatePassword(passwordOptions)

  // const { isValid, errors } = validatePassword(password)

  // if (!isValid) {
  //   ElMessage.warning({
  //     message: '生成的密码不符合安全要求，请重新生成',
  //     duration: 3000,
  //   })
  //   return
  // }

  formData.Password = password
  generateDialogVisible.value = false
  previewPassword.value = ''
  checkStrength()
}

// 添加查看原密码功能
const handleViewOriginalPassword = async () => {
  try {
    const data = await passwordStore.getPassword(id)
    const decryptedPassword = data.Password
    ElMessage.success({
      message: `原密码: ${decryptedPassword}`,
      duration: 3000,
      customClass: 'password-message',
    })
  } catch (error) {
    ElMessage.error('获取原密码失败')
  }
}

// 在script部分添加以下代码（放在已有代码之后）
// 添加计算属性
const uniqueAppNames = computed(() => {
  return Array.from(new Set(passwordStore.passwordList.map((p) => p.AppName))).map((value) => ({
    value,
  }))
})

const uniqueAccountTypes = computed(() => {
  return Array.from(new Set(passwordStore.passwordList.map((p) => p.AccountType))).map((value) => ({
    value,
  }))
})

const uniqueAccounts = computed(() => {
  return Array.from(new Set(passwordStore.passwordList.map((p) => p.Account))).map((value) => ({
    value,
  }))
})

const uniqueUrls = computed(() => {
  return Array.from(new Set(passwordStore.passwordList.map((p) => p.Url))).map((value) => ({
    value,
  }))
})

// 添加查询方法
const queryAppName = (queryString: string, cb: Function) => {
  cb(
    queryString
      ? uniqueAppNames.value.filter((item) =>
        item.value.toLowerCase().includes(queryString.toLowerCase()),
      )
      : uniqueAppNames.value,
  )
}

const queryAccountType = (queryString: string, cb: Function) => {
  cb(
    queryString
      ? uniqueAccountTypes.value.filter((item) =>
        item.value.toLowerCase().includes(queryString.toLowerCase()),
      )
      : uniqueAccountTypes.value,
  )
}

const queryAccount = (queryString: string, cb: Function) => {
  cb(
    queryString
      ? uniqueAccounts.value.filter((item) =>
        item.value.toLowerCase().includes(queryString.toLowerCase()),
      )
      : uniqueAccounts.value,
  )
}

const queryUrl = (queryString: string, cb: Function) => {
  cb(
    queryString
      ? uniqueUrls.value.filter((item) =>
        item.value.toLowerCase().includes(queryString.toLowerCase()),
      )
      : uniqueUrls.value,
  )
}

// 添加选择处理方法
const handleAppSelect = (item: { value: string }) => {
  // 可以根据需要添加自动填充逻辑
}
</script>

<style scoped>
.create-password {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.password-input-group {
  display: flex;
  gap: 10px;
}

.password-input-wrapper {
  flex: 1;
}

.password-strength {
  font-size: 12px;
  margin-top: 4px;
  color: #606266;
}

.password-preview {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: monospace;
  background-color: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
}

:global(.password-message) {
  font-family: monospace;
  user-select: all;
}

.shortcut-tips {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

:deep(.el-tag) {
  cursor: help;
}

.edit-password {
  max-width: 800px;
  margin: 20px auto;
  padding: 0 20px;
}

@media (max-width: 768px) {
  .edit-password {
    padding: 0 10px;
  }

  .shortcut-tips {
    flex-wrap: wrap;
  }
}
</style>
