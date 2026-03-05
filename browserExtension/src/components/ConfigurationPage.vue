<script lang="ts" setup>
import { ref } from 'vue';
import { refreshAll } from "@/utils/refresh";
import { Storges } from "@/stores/storages";

const PasswordServerUrl = Storges().PasswordServerUrl.vue

</script>

<template>
  <el-form label-position="top">
    <el-form-item label="服务器地址">
      <el-input v-model="PasswordServerUrl" placeholder="请输入服务器地址" clearable />
    </el-form-item>
    <el-form-item>
      <el-button type="info" @click="refreshAll"> 保存 </el-button>
      <el-button>取消</el-button>
    </el-form-item>
  </el-form>
</template>

<style scoped>
.el-form {
  padding: 10px;

}
</style>
