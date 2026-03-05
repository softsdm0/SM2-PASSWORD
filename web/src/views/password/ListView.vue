<template>
  <div class="password-list">
    <div class="statistics-cards">
      <el-card class="stat-card">
        <template #header>
          <div class="stat-header">
            <span class="stat-title">应用数量</span>
          </div>
        </template>
        <div class="stat-content">{{ statistics.appCount }}</div>
      </el-card>

      <el-card class="stat-card">
        <template #header>
          <div class="stat-header">
            <span class="stat-title">总密码数</span>
          </div>
        </template>
        <div class="stat-content">{{ statistics.total }}</div>
      </el-card>

      <el-card class="stat-card">
        <template #header>
          <div class="stat-header">
            <span class="stat-title">弱密码数量</span>
          </div>
        </template>
        <div class="stat-content warning">{{ statistics.weakPasswords }}</div>
      </el-card>

      <el-card class="stat-card" @click="handleViewRecordsDetail">
        <template #header>
          <div class="stat-header">
            <span class="stat-title">最近7天查看密码次数</span>
          </div>
        </template>
        <div class="stat-content">{{ recent7DaysViewCount }}</div>
      </el-card>
    </div>

    <div class="list-header">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="搜索应用名称、账户类型或账户" prefix-icon="Search" clearable
          @clear="handleSearchClear" class="search-input" />
        <el-select v-model="sortBy" placeholder="排序方式" class="sort-select">
          <el-option label="应用名称" value="AppName" />
          <el-option label="账户类型" value="AccountType" />
          <el-option label="创建时间" value="CreatedAt" />
          <el-option label="更新时间" value="UpdatedAt" />
        </el-select>
        <el-switch v-model="sortDesc" active-text="降序" inactive-text="升序" />
      </div>
    </div>
    <!-- <div class="action-area">
      <el-button type="success" :icon="Upload" @click="handleImport"> 导入 </el-button>
      <el-button type="success" :icon="Download" @click="handleExport"> 导出 </el-button>
      <el-button type="danger" :icon="Delete" :disabled="!selectedIds.length" @click="handleBatchDelete">
        批量删除
      </el-button>
    </div> -->

    <el-collapse v-model="activeGroups" v-loading="passwordStore.loading">
      <template v-if="filteredAndSortedPasswords.length">
        <el-collapse-item v-for="(group, appName) in passwordGroups" :key="appName"
          :title="`${appName} (${group.length})`" :name="appName">
          <template #title>
            <div class="group-title">
              <span>{{ appName }}</span>
              <el-tag size="small" round>{{ group.length }}</el-tag>
            </div>
          </template>

          <el-table :data="group" style="width: 100%" :row-class-name="getRowClassName"
            :header-cell-style="{ background: '#f5f7fa' }" table-layout="fixed">
            <el-table-column prop="AccountType" label="账户类型" min-width="80" />
            <el-table-column prop="Account" label="账户" min-width="100">
              <template #default="{ row }">
                <span v-if="row.Account" class="account-link"
                  @click="router.push({ name: 'passwordDetail', params: { id: row.ID } })">
                  {{ row.Account }}
                </span>
                <span v-else class="account-link"
                  @click="router.push({ name: 'passwordDetail', params: { id: row.ID } })">
                  请重新设置账号
                </span>
                <span v-if="row.PasswordStrength <= 2" class="weak-password-label">弱</span>
                <el-tooltip content="复制账号" placement="top">
                  <el-button type="success" :icon="CopyDocument" circle class="action-btn"
                    @click="handleCopy(row.Account)" />
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" align="center"> <!-- 增加列宽度以容纳新按钮 -->
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-tooltip content="复制密码" placement="top">
                    <el-button type="success" :icon="CopyDocument" circle class="action-btn"
                      @click="handleCopyPassword(row.ID)" />
                  </el-tooltip>
                  <el-tooltip content="查看密码" placement="top">
                    <el-button type="primary" :icon="View" circle class="action-btn"
                      @click="handleViewPassword(row.ID)" />
                  </el-tooltip>
                  <el-tooltip content="查看详情" placement="top">
                    <el-button type="info" :icon="InfoFilled" circle class="action-btn"
                      @click="handleViewDetail(row)" />
                  </el-tooltip>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-collapse-item>
      </template>
      <el-empty v-else :description="searchQuery ? '没有找到匹配的密码记录' : '暂无密码记录'">
        <template #extra>
          <el-button type="primary" @click="handleCreate"> 立即创建 </el-button>
        </template>
      </el-empty>
    </el-collapse>

    <!-- <input ref="fileInputRef" type="file" accept=".json" style="display: none" @change="onFileSelected" /> -->
    <el-button class="floating-create-btn" type="primary" :icon="Plus" circle @click="handleCreate" />
  </div>
</template>

<script setup lang="ts">
import { Upload, Download } from '@element-plus/icons-vue'
import { exportPasswords, importPasswords } from '../../utils/importExport'
import { Plus } from '@element-plus/icons-vue'
import { Edit } from '@element-plus/icons-vue'
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { View, Delete, CopyDocument, InfoFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { usePasswordStore } from '../../stores/password'
import { decryptPassword } from '../../utils/crypto'
import type { IPasswordInfo, IPassword } from '../../types/password'
import { checkPasswordStrength } from '../../utils/passwordStrength'
import { handleCopyPassword } from "@/utils/copyPassword";
import { handleCopy } from '@/utils/copy'

const router = useRouter()
const activeGroups = ref<string[]>([])
const passwordStore = usePasswordStore()
const recent7DaysViewCount = ref(0)

const searchQuery = ref('')
const sortBy = ref('AppName')
const sortDesc = ref(true)

const filteredAndSortedPasswords = computed(() => {
  // 检查 passwordStore.passwordList 是否为数组，如果不是则初始化为空数组
  let result = Array.isArray(passwordStore.passwordList) ? [...passwordStore.passwordList] : []

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(
      (item) =>
        item.AppName.toLowerCase().includes(query) ||
        item.AccountType.toLowerCase().includes(query) ||
        item.Account.toLowerCase().includes(query),
    )
  }

  // 排序
  result.sort((a, b) => {
    const aValue = a[sortBy.value]
    const bValue = b[sortBy.value]
    if (sortDesc.value) {
      return bValue.localeCompare(aValue)
    }
    return aValue.localeCompare(bValue)
  })

  return result
})

const passwordGroups = computed(() => {
  const filteredList = filteredAndSortedPasswords.value || []
  return filteredList.reduce((groups: Record<string, IPasswordInfo[]>, item) => {
    if (!groups[item.AppName]) {
      groups[item.AppName] = []
    }
    groups[item.AppName].push(item)
    return groups
  }, {})
})

const handleSearchClear = () => {
  searchQuery.value = ''
}

const handleViewPassword = async (id: string) => {
  try {
    const data: IPassword = await passwordStore.getPassword(id)
    let decryptedPassword = decryptPassword(data.Password)
    if (!decryptedPassword) {
      decryptedPassword = "密码异常请重新设置密码"
    }
    // 修改部分，使用 ElMessageBox 显示密码
    ElMessageBox.alert(
      `<span class="password-message">${decryptedPassword}</span>`,
      '',
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '确认',
        callback: () => {
          // 可以在这里添加确认后的回调逻辑
        }
      }
    )
  } catch (error) {
    ElMessage.error('获取密码失败')
    console.log(error)
  }
}

const handleDelete = async (id: string) => {
  try {
    await ElMessageBox.confirm('确认删除该密码记录吗？', '提示', {
      type: 'warning',
    })
    await passwordStore.deletePassword(id)
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleViewDetail = (row: IPasswordInfo) => {
  router.push({
    name: 'passwordDetail',
    params: {
      id: row.ID,
    },
  })
}

const handleEdit = (id: string) => {
  router.push({
    name: 'passwordEdit',
    params: { id },
  })
}

const selectedIds = ref<string[]>([])

const handleCreate = () => {
  router.push({ name: 'passwordCreate' })
}

const handleSelectionChange = (selection: IPasswordInfo[]) => {
  selectedIds.value = selection.map((item) => item.ID)
}

const handleBatchDelete = async () => {
  if (!selectedIds.value.length) return

  try {
    await ElMessageBox.confirm(
      `确认删除选中的 ${selectedIds.value.length} 条密码记录吗？`,
      '批量删除',
      {
        type: 'warning',
      },
    )

    await Promise.all(selectedIds.value.map((id) => passwordStore.deletePassword(id)))
    ElMessage.success('批量删除成功')
    selectedIds.value = []
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

const fetchRecent7DaysViewCount = async () => {
  try {
    const records = await passwordStore.getPasswordRecords()
    const now = new Date()
    const sevenDaysAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)

    const recentRecords = records.filter(record => {
      const recordDate = new Date(record.CreatedAt)
      return recordDate >= sevenDaysAgo && recordDate <= now
    })

    recent7DaysViewCount.value = recentRecords.length
  } catch (error) {
    console.error('获取最近7天密码查看量失败:', error)
  }
}

onMounted(async () => {
  await passwordStore.fetchPasswordList()
  await fetchRecent7DaysViewCount()
})

const fileInputRef = ref<HTMLInputElement>()

const handleImport = () => {
  fileInputRef.value?.click()
}

const onFileSelected = async (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  try {
    const reader = new FileReader()
    const result = await new Promise((resolve, reject) => {
      reader.onload = () => resolve(reader.result)
      reader.onerror = () => reject(new Error('读取文件失败'))
      reader.readAsText(file)
    })

    const passwords = JSON.parse(result as string)
    await Promise.all(
      passwords.map((pwd) => {
        const { ID, CreatedAt, UpdatedAt, ...data } = pwd
        return passwordStore.createPassword(data)
      }),
    )

    ElMessage.success('导入成功')
    await passwordStore.fetchPasswordList()
  } catch (error) {
    ElMessage.error('导入失败：' + (error as Error).message)
  } finally {
    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
  }
}

const handleExport = () => {
  const exportData = passwordStore.passwordList.map(({ Password, ...rest }) => rest)
  const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)

  const link = document.createElement('a')
  link.href = url
  link.download = `passwords-${new Date().toISOString().split('T')[0]}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  ElMessage.success('导出成功')
}

const statistics = computed(() => {
  // 检查 passwordStore.passwordList 是否为 undefined 或 null
  const passwordList = passwordStore.passwordList || []
  const total = passwordList.length
  const appCount = new Set(passwordList.map((p) => p.AppName)).size
  const weakPasswords = passwordList.filter((p) => {
    return p.PasswordStrength <= 2
  }).length

  return {
    total,
    appCount,
    weakPasswords,
  }
})

const getRowClassName = ({ row }: { row: IPasswordInfo }) => {
  const password = row.Password
  const decryptedPassword = password ? decryptPassword(password) : undefined
  const strength = decryptedPassword ? checkPasswordStrength(decryptedPassword) : null
  return strength && strength.score <= 2 ? 'weak-password-row' : ''
}

const handleViewRecordsDetail = () => {
  router.push("/password/records/detail")
}
</script>


<style scoped>
.password-list {
  max-width: 100%;
  padding: 0 16px;
  /* 添加左右内边距 */
}

.list-header {
  display: flex;
  flex-wrap: wrap;
  /* 允许换行 */
  justify-content: space-between;
  margin-bottom: 10px;
  gap: 16px;
  /* 添加间距 */
}

.search-area {
  width: 100%;
  display: flex;
  flex-wrap: wrap;
  /* 允许换行 */
  gap: 24px;
  align-items: center;
  padding: 12px 0;
}

.action-area {
  display: flex;
  gap: 16px;
}

.search-input {
  width: 100%;
  /* 宽度自适应 */
  max-width: 300px;
  /* 最大宽度 */
}

.sort-select {
  width: 100%;
  /* 宽度自适应 */
  max-width: 150px;
  /* 最大宽度 */
}

:global(.password-message) {
  font-family: monospace;
  user-select: all;
  font-size: 24px;
  /* 增大字体大小，可根据需求调整 */
}

.statistics-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
}

.stat-card {
  border: none;
  /* 移除边框 */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  /* 添加阴影效果 */
  transition: box-shadow 0.3s ease;
  /* 添加过渡效果 */

  &:hover {
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
    /* 鼠标悬停时增强阴影效果 */
  }

  :deep(.el-card__header) {
    padding: 12px 16px;
    background-color: #f5f7fa;
    /* 设置头部背景颜色 */
    border-bottom: 1px solid #ebeef5;
    /* 添加底部边框 */
  }
}

.stat-header {
  display: flex;
  justify-content: center;
  /* 居中对齐标题 */
}

.stat-title {
  font-size: 14px;
  color: #606266;
  /* 调整标题颜色 */
  font-weight: 500;
  /* 增加标题字体粗细 */
}

.stat-content {
  font-size: 28px;
  padding: 16px 0;
  /* 增加内容内边距 */
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
  color: #303133;
  /* 调整内容颜色 */
}

.stat-content.warning {
  color: #e6a23c;
}

:deep(.el-empty) {
  padding: 40px 0;
}

:deep(.el-loading-mask) {
  background-color: rgba(255, 255, 255, 0.8);
}

.group-title {
  margin-left: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-collapse-item__header) {
  font-size: 16px;
  font-weight: 500;
}

:deep(.el-tag) {
  background-color: #409eff20;
  color: #409eff;
  border: none;
}

:deep(.weak-password-row) {
  background-color: #fdf6ec50;
}

:deep(.el-table) {
  width: 100% !important;

  &__body-wrapper,
  &__header-wrapper {
    overflow-x: auto;
    /* 允许水平滚动 */
  }
}

:deep(.el-collapse-item__content) {
  padding-bottom: 0;
}

:deep(.el-table) {
  --el-table-border-color: #ebeef5;
  --el-table-header-bg-color: #f5f7fa;
  border-radius: 4px;
  overflow: hidden;
}

:deep(.el-button.is-circle) {
  transition: transform 0.2s ease;
}

:deep(.el-button.is-circle:hover) {
  transform: scale(1.1);
}

:deep(.el-collapse-item__wrap) {
  padding: 12px;
  background-color: #f5f7fa20;
}

.floating-create-btn {
  position: fixed;
  bottom: 40px;
  left: 0;
  right: 0;
  margin: 0 auto;
  /* 水平居中 */
  width: 56px;
  height: 56px;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 2000;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    transform: scale(1.1);
    /* 只进行放大操作 */
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  }
}

:deep(svg) {
  width: 1.2em;
  height: 1.2em;
}

.account-link {
  cursor: pointer;
  color: var(--el-color-primary);
  transition: color 0.2s;
  display: inline-block;
  padding: 4px 0;
  line-height: 1.5;

  &:hover {
    color: var(--el-color-primary-light-3);
    text-decoration: underline;
  }
}

.action-buttons {
  display: flex;
  gap: 2px;
  justify-content: center;
}

.action-btn {
  transform: scale(0.9);
  transition: all 0.2s ease;

  &:hover {
    transform: scale(1);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  }
}

:deep(.el-table td) {
  padding: 12px 0;
}

:deep(.el-table th) {
  padding: 16px 0;
}

.weak-password-label {
  color: red;
  margin-left: 10px;
}

/* 小屏幕设备 */
@media (max-width: 768px) {
  /* 不想使用按列排列 */
  /* .list-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .search-area {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  } */

  .search-input {
    max-width: 100%;
  }

  .sort-select {
    max-width: 200px;
  }

  .statistics-cards {
    grid-template-columns: repeat(2, 1fr);
    /* 一行显示两个卡片 */
    gap: 6px;
    /* 进一步缩小卡片之间的间距 */
  }

  .stat-card {

    /* 调整卡片内边距 */
    :deep(.el-card__header) {
      padding: 4px 8px;
      /* 适当增加头部内边距 */
    }

    .stat-content {
      /* 调整内容字体大小 */
      font-size: 20px;
      /* 增大内容字体大小 */
      padding: 1px 0;
      /* 适当增加内容内边距 */
      min-height: auto;
      /* 去除最小高度限制 */
    }

    .stat-title {
      font-size: 14px;
      /* 增大标题字体大小 */
    }
  }

  /* .action-buttons {
    gap: 1px; 
    width: 30px;
  } */

  .action-btn {
    /* transform: scale(0.8); 进一步缩小按钮尺寸 */
    margin: 0px;
    border: 0px;
    padding: 0px;
  }
}
</style>
