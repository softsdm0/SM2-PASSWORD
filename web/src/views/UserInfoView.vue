<template>
  <transition name="fade">
    <div class="user-info-container">
      <!-- 返回首页按钮 -->
      <div class="back-button-container">
        <el-button icon="ArrowLeft" @click="goBackHome" type="text" class="back-button">
          <span class="back-button-text">返回首页</span>
        </el-button>
      </div>
      <h1 class="title">用户信息</h1>
      <!-- 加载状态 -->
      <div v-if="userInfoStore.loading" class="loading">
        <el-loading v-bind="loadingOptions"></el-loading>
      </div>
      <!-- 错误状态 -->
      <div v-if="userInfoStore.error" class="error">
        <el-alert title="获取用户信息失败" type="error" :description="userInfoStore.error" closable></el-alert>
      </div>
      <!-- 正常显示用户信息 -->
      <div v-else-if="userInfoStore.userInfo" class="user-info">
        <el-card class="user-info-card">
          <template #header>
            <div class="card-header">
              <!-- 添加条件判断，如果没有头像，显示透明占位 -->
              <img 
                :src="userInfoStore.userInfo?.data.avatar || 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7'" 
                alt="Avatar" 
                class="avatar" 
              />
              <div class="user-name">{{ userInfoStore.userInfo?.data.displayName }}</div>
            </div>
          </template>
          <InfoItem label="ID" :value="userInfoStore.userInfo?.data.id" class="info-item" />
          <InfoItem label="姓名" :value="userInfoStore.userInfo?.data.name" class="info-item" />
          <InfoItem label="邮箱" :value="userInfoStore.userInfo?.data.email" class="info-item" />
          <InfoItem label="电话" :value="userInfoStore.userInfo?.data.phone" class="info-item" />
          <InfoItem label="地址" :value="userInfoStore.userInfo?.data.location" class="info-item" />
          <InfoItem label="性别" :value="userInfoStore.userInfo?.data.gender" class="info-item" />
          <InfoItem label="生日" :value="userInfoStore.userInfo?.data.birthday" class="info-item" />
          <!-- 新增安全信息 -->
          <InfoItem label="注册时间" :value="formatDate(userInfoStore.userInfo?.data.createdTime)" class="info-item" />
          <InfoItem label="上次登录时间" :value="formatDate(userInfoStore.userInfo?.data.lastLoginDate)" class="info-item" />
          <InfoItem label="在线状态" :value="userInfoStore.userInfo?.data.isOnline ? '在线' : '离线'" class="info-item" />
        </el-card>
      </div>
      <!-- 获取用户信息按钮 -->
      <!-- <div class="button-container">
        <el-button type="primary" @click="getcasdoorUserInfo" :loading="userInfoStore.loading" class="fetch-button">
          获取casdoor用户信息
        </el-button>
      </div> -->
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useUserInfoStore } from '../stores/userinfo';
import { ElLoading } from 'element-plus';
import { useRouter } from 'vue-router';
import InfoItem from '../components/InfoItem.vue'; 
import dayjs from 'dayjs'; // 引入 dayjs

const userInfoStore = useUserInfoStore();
const router = useRouter();

const loadingOptions = {
  lock: true,
  text: '加载中...',
  spinner: 'el-icon-loading',
  background: 'rgba(0, 0, 0, 0.7)'
};

// 定义返回首页的方法
const goBackHome = () => {
  router.push('/'); 
};

// 封装获取用户信息的方法
const fetchUserInfo = async () => {
  try {
    if (!userInfoStore.userInfo) {
      await userInfoStore.fetchUserInfo();
    }
  } catch (error) {
    console.error('获取用户信息失败:', error);
  }
};

// import { getCurrentInstance } from 'vue'
// const instance = getCurrentInstance()
// const getcasdoorUserInfo = () => {
//   window.location.href = instance.proxy.getMyProfileUrl()
// }

onMounted(() => {
  fetchUserInfo();
});

// 定义格式化日期的函数
const formatDate = (date: string | number | Date | null | undefined) => {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '暂无数据';
};
</script>

<style scoped>
/* 整体容器样式 */
.user-info-container {
  padding: 15px !important;
  margin: 0px !important;
  /* max-width: 800px; */
  /* margin: 2rem auto; */
  position: relative;
  animation: fadeIn 0.5s ease;
  font-family: 'Inter', sans-serif; 
  background-color: #f9f9f9; /* 调整背景颜色 */
  border-radius: 15px;
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1); /* 调整阴影 */
  /* border: 1px solid #e0e0e0; */
  transition: all 0.3s ease; /* 添加过渡效果 */
}

.user-info-container:hover {
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.15); /* 鼠标悬停时阴影加深 */
}

/* 返回按钮样式 */
.back-button-container {
  position: absolute;
  top: 2rem;
  left: 2rem;
}

.back-button {
  display: flex;
  align-items: center;
  color: #777;
  transition: all 0.3s ease;
  background-color: transparent;
  border: none;
  cursor: pointer;
}

.back-button:hover {
  color: #333;
  transform: scale(1.1);
}

.back-button-text {
  margin-left: 0.6rem;
  font-size: 1rem;
  font-weight: 500;
}


/* 标题样式 */
.title {
  text-align: center;
  margin-bottom: 2.5rem;
  font-size: 2.5rem;
  color: #2a2a2a;
  font-weight: 700;
  letter-spacing: -0.02em;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: color 0.3s ease; /* 添加颜色过渡效果 */
}

.title:hover {
  color: #007BFF; /* 鼠标悬停时标题颜色变化 */
}

/* 加载状态样式 */
.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 250px;
}

/* 错误状态样式 */
.error {
  margin-bottom: 2.5rem;
}

/* 用户信息卡片样式 */
.user-info {
  margin-bottom: 2.5rem;
}

.user-info-card {
  border-radius: 15px;
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.06);
  border: none;
  background-color: #ffffff; /* 调整卡片背景颜色 */
  transition: all 0.3s ease;
}

.user-info-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
}

.card-header {
  display: flex;
  align-items: center;
  padding: 2rem;
  border-bottom: 1px solid #eaeaea;
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  margin-right: 2rem;
  object-fit: cover;
  border: 4px solid #ffffff;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.avatar:hover {
  transform: scale(1.15);
}

.user-name {
  font-size: 1.8rem;
  font-weight: 600;
  color: #333;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 1.2rem;
  padding: 0 2rem;
}

.label {
  width: 150px;
  font-weight: 600;
  color: #666;
}

.value {
  flex: 1;
  color: #333;
  word-break: break-all; /* 让长文本自动换行 */
  font-size: 1.05rem;
}


/* 获取用户信息按钮样式 */
.button-container {
  text-align: center;
}

.fetch-button {
  padding: 1rem 2.2rem;
  font-size: 1.15rem;
  border-radius: 10px;
  background-color: #007BFF;
  color: #fff;
  border: none;
  transition: all 0.3s ease;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.fetch-button:hover {
  background-color: #0056b3;
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

/* 淡入动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-info-container {
    padding: 2rem;
    margin: 1.5rem auto;
  }

  .title {
    font-size: 2.2rem;
  }

  .avatar {
    width: 90px;
    height: 90px;
  }

  .user-name {
    font-size: 1.6rem;
  }

  .label {
    width: 120px;
  }

  .fetch-button {
    padding: 0.9rem 2rem;
  }
}

@media (max-width: 576px) {
  .user-info-container {
    padding: 1.5rem;
    margin: 1rem auto;
  }

  .title {
    font-size: 2rem;
  }

  .avatar {
    width: 70px;
    height: 70px;
  }

  .user-name {
    font-size: 1.4rem;
  }

  .label {
    width: 100px;
  }

  .fetch-button {
    padding: 0.8rem 1.8rem;
  }

  .info-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .label {
    margin-bottom: 0.5rem;
  }

  .back-button-container {
    top: 1rem; /* 减小顶部间距 */
    left: 1rem; /* 减小左侧间距 */
  }

  .back-button {
    transform: scale(0.8); /* 缩小按钮 */
    transform-origin: top left; /* 设置缩放原点为左上角 */
  }
}
</style>