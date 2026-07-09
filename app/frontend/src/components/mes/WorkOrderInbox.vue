<template>
  <a-dropdown trigger="click">
    <a-badge :count="unreadCount" :overflow-count="99" size="small">
      <a-button type="text" class="icon-button" aria-label="工单">
        <MailOutlined />
      </a-button>
    </a-badge>
    <template #overlay>
      <a-menu>
        <a-menu-item @click="router.push('/mes/workorders')">
          <MailOutlined />
          工单列表
        </a-menu-item>
        <a-menu-item @click="router.push({ path: '/mes/workorders', query: { mode: 'sent' } })">
          <SendOutlined />
          发工单
        </a-menu-item>
        <a-menu-item @click="router.push({ path: '/mes/create', query: { type: 'workOrder' } })">
          <FormOutlined />
          新建工单
        </a-menu-item>
      </a-menu>
    </template>
  </a-dropdown>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { FormOutlined, MailOutlined, SendOutlined } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { listWorkOrder } from '@/api/mesController'

const router = useRouter()
const unreadCount = ref(0)

const loadUnreadCount = async () => {
  const res = await listWorkOrder({ isTo: true, isUnread: true, limit: 30 })
  if (res.data.code === 0) {
    unreadCount.value = res.data.data?.totalRow ?? res.data.data?.records?.length ?? 0
  }
}

onMounted(loadUnreadCount)
</script>

<style scoped>
.icon-button {
  width: 42px;
  height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #1d1d1f;
  border: 0;
  background: transparent;
  box-shadow: none;
}

.icon-button :deep(svg) {
  width: 24px;
  height: 24px;
}

.icon-button:hover {
  color: #0066cc;
  background: transparent;
}
</style>
