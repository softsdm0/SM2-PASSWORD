<script lang="ts" setup>
import { Storges } from '@/stores/storages';
import { useUserInfoStore } from "@/stores/userinfo";
import { onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { refreshAll } from '@/utils/refresh'


const status = Storges().LoginStatus.vue;
const userInfo = Storges().UserInfo.vue;

onMounted(async () => {
  refreshAll()
})
</script>

<template>
  <el-row justify="center">
    <el-col justify="center" align="middle" :span="2">
      <div class="status-indicator-container">
        <div class="status-indicator" :class="{
      'initial': status === 'initial',
      'logged': status === 'logged',
      'unlogged': status === 'unlogged',
      'error': status === 'error'
    }">
        </div>
      </div>
    </el-col>
    <el-col justify="center" align="middle" :span="7">
      <div class="text-container">
        <span v-if="status === 'logged' && userInfo" class="username">
          {{ userInfo.data.name }}
        </span>
        <span v-else class="error-message">
          未登录
        </span>
      </div>
    </el-col>
    <el-col :span="13">
    </el-col>
    <el-col :span="2">
      <el-button @click="refreshAll" :icon="Refresh" circle size="small" />
    </el-col>
  </el-row>
</template>

<style scoped>
.status-container {
  display: flex;
  align-items: center;
}

.username {
  font-size: 14px;
  color: #909399;
  /* 修改为Element UI的浅灰色 */
}

.status-indicator-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.status-indicator {
  width: 16px;
  height: 16px;
  border-radius: 50%;
}

.status-indicator.initial {
  background-color: #9E9E9E;
  /* 灰色 */
}

.status-indicator.logged {
  background-color: #4CAF50;
  /* 绿色 */
}

.status-indicator.unlogged {
  background-color: #FFC107;
  /* 黄色 */
}

.status-indicator.error {
  background-color: #F44336;
  /* 红色 */
}

.error-message {
  font-size: 12px;
  color: #F44336;
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.text-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 100%;
}

.username,
.error-message {
  text-align: center;
}
</style>
