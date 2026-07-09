<template>
  <main class="mes-home-loading">
    <a-spin />
  </main>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLoginUserStore } from '@/stores/loginUser'
import { mesRoleHomePath } from '@/utils/mesRole'

const router = useRouter()
const loginUserStore = useLoginUserStore()

onMounted(async () => {
  if (!loginUserStore.loginUser.id) {
    await loginUserStore.fetchLoginUser()
  }
  await router.replace(mesRoleHomePath(loginUserStore.loginUser.userRole))
})
</script>

<style scoped>
.mes-home-loading {
  min-height: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
