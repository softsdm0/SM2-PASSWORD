<template>
  <!-- 添加一个类名 content-container 用于设置边距 -->
  <div class="content-container">
    <div class="header-container">
      <el-button @click="goBack">返回</el-button>
      <h1 class="centered-title">所有查看密码记录详情</h1>
    </div>
    <el-table :data="recentRecords" style="width: 100%">
      <el-table-column prop="ID" label="ID" />
      <el-table-column label="创建时间">
        <template #default="{ row }">
          {{ formatDate(row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column prop="UserId" label="请求者ID" />
      <el-table-column prop="PasswordInfo.AppName" label="应用" />
      <el-table-column prop="PasswordInfo.AccountType" label="账户类型" />
      <el-table-column prop="PasswordInfo.Account" label="账户" />

      <el-table-column prop="Status" label="请求状态">
        <template #default="{ row }">
          {{ getStatusText(row.Status) }}
        </template>
      </el-table-column>
      <el-table-column prop="IP" label="IP地址" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { usePasswordStore } from '../../stores/password';
import { GetPasswordRecordsStatus } from '../../types/password';
import { useRouter } from 'vue-router';
import dayjs from 'dayjs';

const passwordStore = usePasswordStore();
const recentRecords = ref([]);
const router = useRouter();

const getStatusText = (status: GetPasswordRecordsStatus) => {
  switch (status) {
    case GetPasswordRecordsStatus.Unknown:
      return '未知错误';
    case GetPasswordRecordsStatus.Success:
      return '处理成功';
    case GetPasswordRecordsStatus.Refused:
      return '已拒绝';
    case GetPasswordRecordsStatus.NotFound:
      return '记录不存在';
    default:
      return '未知状态';
  }
};

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

const fetchRecentRecords = async () => {
  try {
    const records = await passwordStore.getPasswordRecords();
    // 对记录按照创建时间进行倒序排序
    const sortedRecords = records.sort((a, b) => {
      return new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime();
    });
    recentRecords.value = sortedRecords; // 赋值排序后的记录
    console.log(sortedRecords);

  } catch (error) {
    console.error('获取所有密码查看记录失败:', error);
  }
};

const goBack = () => {
  router.push('/')
};

onMounted(() => {
  fetchRecentRecords();
});
</script>

<style scoped>
.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.centered-title {
  text-align: center;
  flex-grow: 1;
}

/* 新增样式，设置左右边距 */
.content-container {
  padding-left: 20px;
  padding-right: 20px;
}
</style>
