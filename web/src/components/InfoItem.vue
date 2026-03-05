<template>
  <div class="info-item">
    <span class="label">{{ label }}</span>
    <span class="value">{{ value }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: String,
    default: ''
  }
});
</script>

<style scoped>
.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.label {
  width: 100px;
  font-weight: bold;
}

.value {
  flex: 1;
}
</style>
