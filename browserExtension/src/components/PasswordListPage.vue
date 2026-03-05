<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import type { IPasswordInfo } from '@/types/password';
import { Storges } from "@/stores/storages";

const expandedApps = ref<Record<string, boolean>>({});
const passwordList = Storges().PasswordInfoList.vue;

// 加载密码列表
onMounted(async () => {
});

// 按照AppName分组
const groupedPasswords = computed(() => {
  const groups: Record<string, IPasswordInfo[]> = {};
  passwordList.value?.forEach(password => {
    // console.log(password);

    if (!groups[password.AppName]) {
      groups[password.AppName] = [];
    }
    groups[password.AppName].push(password);
  });
  return groups;
});

const toggleApp = (appName: string) => {
  expandedApps.value[appName] = !expandedApps.value[appName];
};

// 新增方法
const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text);
};

const openUrl = (url: string) => {
  window.open(url, '_blank');
};
</script>

<template>
  <div class="password-list-container">
    <el-scrollbar>
      <el-collapse v-model="expandedApps" accordion>
        <el-collapse-item v-for="(passwords, appName) in groupedPasswords" :key="appName" :title="appName"
          :name="appName">
          <el-table :data="passwords" style="width: 100%">
            <el-table-column prop="AccountType" label="账号类型" width="120" />
            <el-table-column prop="Account" label="账号" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="primary" link @click="copyToClipboard(row.Account)">
                  复制
                </el-button>
                <el-button type="primary" link @click="openUrl(row.Url)">
                  打开
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-collapse-item>
      </el-collapse>
    </el-scrollbar>
  </div>
</template>

<style scoped>
.password-list-container {
  /* padding: 10px; */
  padding-left: 10px;
  height: 100%;
}

/* 自定义滚动条样式 */
:deep(.el-scrollbar__view) {
  padding-right: 10px;
}

:deep(.el-scrollbar__bar.is-vertical) {
  right: 2px;
}
</style>