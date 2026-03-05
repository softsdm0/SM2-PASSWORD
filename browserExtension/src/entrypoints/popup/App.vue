<script lang="ts" setup>
import { onMounted } from 'vue';
import { Storges } from '@/stores/storages';
import { useRouter } from 'vue-router';
import Status from "@/components/Status.vue";

const router = useRouter();

onMounted(async () => {
  const url = await Storges().PasswordServerUrl.value
  if (!url || url === "") {
    router.push('/config');
  }
});
</script>

<template>
  <el-container>
    <el-header height="66px">
      <el-row style="padding-top: 5px">
        <el-col>
          <Status></Status>
        </el-col>
      </el-row>
      <el-row justify="center" style="padding-top: 5px">
        <el-col align="center" :span="8">
          <router-link to="/" active-class="active-link" exact-active-class="exact-active-link">
            <el-button type="info" text>列表</el-button>
          </router-link>
        </el-col>
        <el-col align="center" :span="8">
          <router-link to="/config" active-class="active-link" exact-active-class="exact-active-link">
            <el-button type="info" text>配置</el-button>
          </router-link>
        </el-col>
        <el-col align="center" :span="8">
          <router-link to="/info" active-class="active-link" exact-active-class="exact-active-link">
            <el-button type="info" text>信息</el-button>
          </router-link>
        </el-col>
      </el-row>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
    <el-footer height="20px">
      Powered by Password
    </el-footer>
  </el-container>
</template>

<style scoped>
/* 整体容器样式 */
.el-container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f9fafb;
}

/* 头部样式 */
.el-header {
  padding-right: 5px;
  padding-left: 5px;
}

.el-header .el-button {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  text-decoration: none;
  /* 确保按钮也没有下划线 */
}

a {
  text-decoration: none;
}

.active-link .el-button {
  color: #606266 !important;
  /* 修改为深灰色 */
  font-weight: bold;
}

.exact-active-link .el-button {
  color: #606266 !important;
  /* 修改为深灰色 */
  font-weight: bold;
  border-bottom: 2px solid #606266;
  /* 下划线也改为深灰色 */
}

/* 主内容区域样式 */
.el-main {
  background-color: #ffffff;
  /* flex: 1; */
  padding: 0px;
}

.el-footer {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  color: #909399;
}
</style>
