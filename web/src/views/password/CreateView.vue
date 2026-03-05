<template>
  <div class="create-password">
    <el-card>
      <template #header>
        <div class="card-header">
          <h3>新建密码</h3>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px">
        <el-form-item label="应用名称" prop="AppName">
          <el-autocomplete
            v-model="formData.AppName"
            :fetch-suggestions="queryAppName"
            placeholder="请输入应用名称"
            @select="handleAppSelect"
          />
        </el-form-item>

        <el-form-item label="账户类型" prop="AccountType">
          <el-autocomplete
            v-model="formData.AccountType"
            :fetch-suggestions="queryAccountType"
            placeholder="请输入账户类型"
            @select="handleAppSelect"
          />
        </el-form-item>

        <el-form-item label="账户" prop="Account">
          <el-autocomplete
            v-model="formData.Account"
            :fetch-suggestions="queryAccount"
            placeholder="请输入账户"
            @select="handleAppSelect"
          />
        </el-form-item>

        <el-form-item label="密码" prop="Password">
          <div class="password-input-group">
            <div class="password-input-wrapper">
              <el-input
                v-model="formData.Password"
                type="password"
                placeholder="请输入密码"
                show-password
                @input="checkStrength"
              />
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
          <el-autocomplete
            v-model="formData.Url"
            :fetch-suggestions="queryUrl"
            placeholder="请输入 URL"
            @select="handleAppSelect"
          />
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
            <el-button
              type="primary"
              link
              :icon="DocumentCopy"
              @click="copyPassword(previewPassword)"
            >
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
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { usePasswordStore } from '../../stores/password'
import { encryptPassword } from '../../utils/crypto'
import { validatePassword } from '../../utils/passwordValidator'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const passwordStore = usePasswordStore()

const formData = reactive({
  AppName: '',
  AccountType: '',
  Account: '',
  Password: '',
  Url: '',
  Notes: '',
})

const rules = {
  AppName: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  AccountType: [{ required: true, message: '请输入账户类型', trigger: 'blur' }],
  Account: [{ required: true, message: '请输入账户', trigger: 'blur' }],
  Password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (!value) {
          callback()
          return
        }
        // 检查密码强度
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

const handleSubmit = async () => {
  if (!formRef.value) return;

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
    loading.value = true;
    const submitData = { ...formData };
    submitData.Password = encryptPassword(formData.Password);
    await passwordStore.createPassword(submitData);
    ElMessage.success('创建成功');
    router.push('/');
  } catch (error) {
    console.log(error);
    // 校验失败会进入这里
    if (typeof error === 'object' && error !== null && 'message' in error) {
      ElMessage.error('表单校验失败，请检查输入内容');
    } else if (error !== 'cancel') {
      ElMessage.error('创建密码记录失败');
    }
  } finally {
    loading.value = false;
  }
}

const checkUrl = (url: string):string =>  {
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

const handleCancel = () => {
  router.push('/')
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
    console.log(err)
  }
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

const strengthInfo = reactive({
  score: 0,
  level: '弱',
  color: '#F56C6C',
})

const checkStrength = () => {
  if (!formData.Password) return
  const result = checkPasswordStrength(formData.Password)
  Object.assign(strengthInfo, result)
}

// 在现有响应式对象后添加计算属性
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
  // 可以选择自动填充其他字段，这里保持简单处理
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
</style>
