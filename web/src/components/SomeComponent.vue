<template>
  <div>
    <button @click="fetchSessionId">获取会话 ID</button>
    <button @click="useSessionIdStore().clearSessionIdAndSm4Key()">清理会话 ID</button>
    <div v-if="sessionId">会话 ID: {{ sessionId }}</div>
    <div v-if="sm4Key">sm4Key: {{ sm4Key }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useSessionIdStore } from '../stores/sessionId';

const sessionId = ref<string>('');
const sm4Key = ref<string>('');

const fetchSessionId = async () => {
  try {
    sessionId.value = await useSessionIdStore().getSessionId()
    sm4Key.value = await useSessionIdStore().getSm4Key()
    console.log(sessionId.value,sm4Key.value);
  } catch (error) {
    console.error('获取会话 ID 时发生错误:', error);
  }
};
</script>

