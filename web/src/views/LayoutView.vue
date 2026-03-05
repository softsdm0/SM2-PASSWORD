<template>
  <el-container class="layout-container">
    <el-header height="80px">
      <el-row class="header-content" :gutter="24">
        <el-col :span="4" class="header-left">
          <img src="@/assets/passwordserver.svg" class="icon" alt="logo" style="height: 40px;" />
        </el-col>
        <el-col :span="16" class="header-middle">
          <span @click="goToUserInfo" class="user-name">{{ userInfoStore.userInfo?.data.name }}</span>
        </el-col>
        <el-col :span="4" class="header-right">
          <el-button type="danger" @click="handleLogout">
            <el-icon>
              <SwitchButton />
            </el-icon>
            <span class="header-right-button-span">
              退出登录
            </span>
          </el-button>
        </el-col>
      </el-row>
    </el-header>

    <el-main>
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { SwitchButton } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '../stores/user'
import { useUserInfoStore } from '../stores/userinfo'
import { useSessionIdStore } from "../stores/sessionId";
import { onMounted } from 'vue'
import Cookies from 'js-cookie';

const router = useRouter();
const userStore = useUserStore();
const userInfoStore = useUserInfoStore();

const handleLogout = () => {
  ElMessageBox.confirm('确认退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      userStore.clearToken();
      useSessionIdStore().clearSessionIdAndSm4Key();
      Cookies.remove('casdoorAuthToken');
      Cookies.remove('authToken');
      router.push('/login');
    })
    .catch(() => { })
}

const goToUserInfo = () => {
  router.push('/userinfo')
}

onMounted(() => {
  userInfoStore.fetchUserInfo()
})
</script>

<style scoped>
/* 整体容器样式 */
.layout-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f9fafb;
}


/* 头部内容样式 */
.header-content {
  display: flex;
  align-items: center;
  height: 100%;
  width: 100%;
  margin: 0px !important;
  text-align: center;
}

/* 头部样式 */
.el-header {
  padding: 0px;
  background-color: #ffffff;
  border-bottom: 2px solid #e5e7eb;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  position: relative;
  z-index: 1;
  transition: all 0.3s ease;
}

/* 头部悬停效果 */
.el-header:hover {
  box-shadow: 0 6px 8px -1px rgba(0, 0, 0, 0.15), 0 3px 6px -1px rgba(0, 0, 0, 0.1);
}

/* 头部渐变背景 */
.el-header::before {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(to bottom, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.8));
  z-index: -1;
}

.icon {
  height: 40px;
  transition: transform 0.3s ease;

  &:hover {
    transform: scale(1.1);
  }
}


/* 主内容区域样式 */
.el-main {
  background-color: #ffffff;
  flex: 1;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  padding: 20px 0px;
}

/* 路由切换过渡效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 用户名称样式 */
.user-name {
  cursor: pointer;
  color: #4b5563;
  font-weight: 600;
  font-size: 16px;
  transition: color 0.3s ease;
}

.user-name:hover {
  text-decoration: underline;
  color: #1d4ed8;
}

/* 按钮样式 */
.el-button {
  transition: all 0.3s ease;
  border-radius: 6px;
}

.el-button:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* 超小屏幕（手机竖屏）响应式样式 */
@media (max-width: 375px) {
  .user-name {
    font-size: 12px;
  }

  .el-button {
    font-size: 10px;
    padding: 6px 10px;
  }

  .header-right-button-span {
    display: none;
  }
}

/* 小屏幕（手机横屏、小平板）响应式样式 */
@media (min-width: 376px) and (max-width: 576px) {

  .user-name {
    font-size: 14px;
  }

  .el-button {
    font-size: 12px;
    padding: 8px 12px;
  }

  .header-right-button-span {
    display: none;
  }
}

/* 中等屏幕（平板）响应式样式 */
@media (min-width: 577px) and (max-width: 768px) {


  .header-right-button-span {
    display: none;
  }
}

/* 大屏幕（桌面）响应式样式 */
@media (min-width: 769px) and (max-width: 1200px) {}

/* 超大屏幕响应式样式 */
@media (min-width: 1201px) {}
</style>
