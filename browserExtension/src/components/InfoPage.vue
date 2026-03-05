<script lang="ts" setup>
import { Storges, LocalLoginStatus } from '@/stores/storages';
import { useConfigStore } from "@/stores/config";
import { LoginStats } from "@/types/userInfo";
import { useUserInfoStore } from '@/stores/userinfo'
import { refreshAll } from "@/utils/refresh";
import { version } from '../../package.json'


const handleLogin = async () => {
  const PasswordServerUrl = await Storges().PasswordServerUrl.value;
  const Config = await useConfigStore().getConfig();

  const loginUrl = `${PasswordServerUrl}/login/oauth/authorize?client_id=${Config?.casdoor.clientId}&response_type=code&redirect_uri=${PasswordServerUrl}%2Fcallback&scope=read&state=hy87g8dtznm`

  browser.windows.create({
    url: loginUrl,
    type: 'popup',
    width: 500,
    height: 600
  }, (newWindow) => {
    if (newWindow) {
      const tabId = newWindow.tabs?.[0]?.id;
      if (!tabId) return;

      const tabUpdateListener = (updatedTabId: number, changeInfo: { url?: string }) => {
        if (updatedTabId === tabId && changeInfo.url && (changeInfo.url === `${PasswordServerUrl}/` || changeInfo.url === `${PasswordServerUrl}`)) {
          console.log('登录成功');
          // 关闭登录页面
          browser.tabs.remove(tabId);
          // 刷新状态
          refreshAll();
          // 关闭监听事件
          browser.tabs.onUpdated.removeListener(tabUpdateListener);
        }
      };

      browser.tabs.onUpdated.addListener(tabUpdateListener);
    }
  });
};

const { state: LoginStatus } = useStoredValue<LoginStats>(LocalLoginStatus, "initial");
const userInfo = Storges().UserInfo.vue;
const PasswordServerUrl = Storges().PasswordServerUrl.vue;


</script>

<template>
  <div class="content-container">
    <el-space direction="vertical" alignment="flex-start" class="login-container">
      <el-text type="info" size="large">
        用户名称: {{ userInfo?.data.name || '未登录' }}
      </el-text>
      <el-text type="info" size="large">
        插件版本: v{{ version }}
      </el-text>
      <el-text type="info" size="large">
        服务端版本: v{{ version }}
      </el-text>
      <el-text type="info" size="large">
        服务端地址: {{ PasswordServerUrl || '未配置' }}
      </el-text>
      <el-button v-if="LoginStatus !== 'logged'" @click="handleLogin" type="info" class="login-button">
        登录
      </el-button>
      <el-button v-else @click="useUserInfoStore().logout" type="info" class="logout-button">
        退出登录
      </el-button>
    </el-space>
  </div>
</template>

<style scoped>
.content-container {
  /* display: flex;
  justify-content: center;
  align-items: center; */
  padding: 10px;
}

.login-container {
  width: 100%;
  height: 100%;
  text-align: left;
}
</style>
