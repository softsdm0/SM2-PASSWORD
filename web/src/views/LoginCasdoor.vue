<template>
</template>

<script setup lang="ts">
// casdoor登录
import { getCurrentInstance } from 'vue'

const instance = getCurrentInstance()
const casdoorLogin = () => {
  window.location.href = instance.proxy.getSigninUrl()
}
casdoorLogin()
</script>

<style scoped>
</style>
