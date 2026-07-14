<template>
  <main class="home-page">
    <section class="home-card">
      <h1>MES 单体系统</h1>
      <p>应用生成模块已下线，当前入口为 MES 工作台。</p>
      <a-button type="primary" @click="router.push('/mes')">进入工作台</a-button>
    </section>
  </main>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.home-page {
  min-height: calc(100vh - 120px);
  display: grid;
  place-items: center;
  padding: 24px;
  background: #f8fafc;
}

.home-card {
  width: min(520px, 100%);
  padding: 32px;
  border: 1px solid #e2e8f0;
  border-radius: 20px;
  background: white;
  text-align: center;
  box-shadow: 0 18px 50px rgba(15, 23, 42, 0.08);
}

.home-card h1 {
  margin-bottom: 12px;
  color: #0f172a;
  font-size: 28px;
}

.home-card p {
  margin-bottom: 24px;
  color: #64748b;
}
</style>
