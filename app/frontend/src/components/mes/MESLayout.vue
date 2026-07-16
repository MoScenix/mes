<template>
  <a-layout class="mes-shell">
    <aside class="mes-sidebar">
      <div class="sidebar-title">
        <span class="title-mark"></span>
        <div>
          <strong>MES 工作台</strong>
        </div>
      </div>
      <nav class="side-nav" aria-label="MES 功能">
        <section v-for="group in visibleNavGroups" :key="group.name" class="nav-group">
          <div class="nav-group-title">{{ group.name }}</div>
          <button
            v-for="item in group.items"
            :key="item.key"
            class="side-nav-item"
            :class="{ active: isActive(item) }"
            type="button"
            @click="go(item)"
          >
            <component :is="item.icon" />
            <span>{{ item.label }}</span>
          </button>
        </section>
      </nav>
    </aside>

    <a-layout class="mes-main-layout">
      <main class="mes-content">
        <router-view />
      </main>
    </a-layout>

    <FloatingAssistant v-if="route.path !== '/mes/assistant'" />
  </a-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLoginUserStore } from '@/stores/loginUser'
import FloatingAssistant from '@/components/mes/FloatingAssistant.vue'
import { normalizeMesRole } from '@/utils/mesRole'
import {
  mesNavGroups,
  mesNavTargetFor,
  visibleMesNavItems,
  type MesNavItem,
} from '@/components/mes/mesNav'

const router = useRouter()
const route = useRoute()
const loginUserStore = useLoginUserStore()

const normalizedRole = computed(() => normalizeMesRole(loginUserStore.loginUser.userRole))

const visibleNavItems = computed(() => visibleMesNavItems(normalizedRole.value))
const visibleNavGroups = computed(() =>
  mesNavGroups
    .map((name) => ({ name, items: visibleNavItems.value.filter((item) => item.group === name) }))
    .filter((group) => group.items.length),
)

const targetFor = (item: MesNavItem) => mesNavTargetFor(item, normalizedRole.value)

const isActive = (item: MesNavItem) => {
  const target = targetFor(item)
  if (route.path !== target.path) return false
  if (target.scanMode) return String(route.query.mode || '') === target.scanMode
  return (
    String(route.query.panel || '') === (target.panel || '') &&
    String(route.query.view || '') === target.view &&
    String(route.query.businessType || '') === (target.businessType || '')
  )
}

const go = async (item: MesNavItem) => {
  const target = targetFor(item)
  const query = target.scanMode
    ? { mode: target.scanMode }
    : target.panel
      ? { panel: target.panel, view: target.view, businessType: target.businessType }
      : { view: target.view }
  const active = target.scanMode
    ? route.path === target.path && String(route.query.mode || '') === target.scanMode
    : route.path === target.path &&
      String(route.query.panel || '') === (target.panel || '') &&
      String(route.query.view || '') === target.view &&
      String(route.query.businessType || '') === (target.businessType || '')
  if (!active) {
    await router.push({ path: target.path, query })
  }
}
</script>

<style scoped>
.mes-shell {
  min-height: calc(100vh - 52px);
  display: flex;
  flex-direction: row;
  align-items: stretch;
  background: #fafafa;
  color: #1d1d1f;
  position: relative;
}

.mes-sidebar {
  width: 200px;
  flex: 0 0 200px;
  min-height: calc(100vh - 52px);
  padding: 20px 10px;
  background: #fff;
  border-right: 1px solid var(--border);
}

.sidebar-title {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 10px 18px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 12px;
}

.title-mark {
  width: 4px;
  height: 20px;
  border-radius: 2px;
  background: var(--primary);
}

.sidebar-title strong {
  display: block;
  font-size: 14px;
  line-height: 1.3;
  font-weight: 600;
  color: var(--foreground);
}

.side-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-group + .nav-group {
  margin-top: 14px;
}

.nav-group-title {
  padding: 0 10px 6px;
  color: #8c8c8c;
  font-size: 11px;
  font-weight: 600;
}

.side-nav-item {
  width: 100%;
  min-height: 36px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 0;
  border-radius: 6px;
  padding: 0 10px;
  background: transparent;
  color: var(--muted-foreground);
  font-size: 13px;
  line-height: 1.2;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;
}

.side-nav-item:hover {
  background: var(--muted);
  color: var(--foreground);
}

.side-nav-item.active {
  background: var(--muted);
  color: var(--primary);
  font-weight: 500;
}

.side-nav-item span {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.mes-main-layout {
  flex: 1 1 auto;
  min-width: 0;
  background: transparent;
}

.mes-content {
  min-width: 0;
  padding: 20px 24px 96px;
}

@media (max-width: 768px) {
  .mes-shell {
    min-height: calc(100vh - 48px);
  }

  .mes-sidebar {
    display: none;
  }

  .mes-content {
    padding: 14px 12px 80px;
  }
}
</style>
