<script setup>
console.log('Test component loaded')
</script>

<template>
  <div class="test-container">
    <h2>测试页面</h2>
    <p>这是一个测试页面，用于验证Vue组件是否能正常渲染</p>
  </div>
</template>

<style scoped>
.test-container {
  max-width: 400px;
  margin: 50px auto;
  padding: 30px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.test-container h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

.test-container p {
  text-align: center;
  color: #666;
}
</style>