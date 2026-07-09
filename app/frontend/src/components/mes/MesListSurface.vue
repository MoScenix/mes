<template>
  <section class="list-surface">
    <div class="list-toolbar">
      <div>
        <p v-if="eyebrow">{{ eyebrow }}</p>
        <h2>{{ title }}</h2>
      </div>
      <div class="toolbar-actions">
        <slot name="actions" />
      </div>
    </div>
    <div v-if="$slots.filters" class="list-filters">
      <slot name="filters" />
    </div>
    <slot />
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  eyebrow?: string
}>()
</script>

<style scoped>
.list-surface {
  display: grid;
  gap: 14px;
  border: 1px solid #e5e5ea;
  border-radius: 8px;
  padding: 16px;
  background: #fff;
}

.list-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.list-toolbar p {
  margin: 0 0 4px;
  color: #6e6e73;
  font-size: 12px;
}

.list-toolbar h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 650;
  letter-spacing: 0;
}

.toolbar-actions,
.list-filters {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 720px) {
  .list-toolbar {
    flex-direction: column;
  }

  .toolbar-actions,
  .list-filters {
    justify-content: flex-start;
    width: 100%;
  }
}
</style>
