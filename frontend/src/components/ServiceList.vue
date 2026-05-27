<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <span>{{ t('sidebar.services', { count: displayServices.length }) }}</span>
      <div class="sidebar-actions">
        <button class="btn-sm btn-secondary" @click="$emit('new')">{{ t('sidebar.newBtn') }}</button>
        <button class="btn-sm btn-primary" @click="$emit('refresh')">{{ t('sidebar.refreshBtn') }}</button>
      </div>
    </div>
    <div class="service-list" v-if="displayServices.length > 0">
      <div
        v-for="svc in displayServices"
        :key="svc.name + '-' + svc.source"
        class="service-item"
        :class="{ active: selectedService === svc.name }"
        @click="$emit('select', svc)"
      >
        <div class="service-item-row">
          <div class="service-item-info">
            <div class="service-item-name">
              {{ svc.displayName || svc.name }}
              <span v-if="svc.source === 'file'" class="source-badge">File</span>
            </div>
            <div class="service-item-meta">
              <span class="status-badge" :class="statusClass(svc.status)">{{ svc.status }}</span>
              <span class="start-type">{{ svc.startType }}</span>
            </div>
          </div>
          <button class="btn-sm btn-copy" @click.stop="$emit('copy', svc)" :title="t('sidebar.copyTooltip')">
            &#x2398;
          </button>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <p>{{ t('sidebar.empty') }}</p>
      <p class="hint">{{ t('sidebar.emptyHint') }}</p>
    </div>
  </aside>
</template>

<script>
import { useI18n } from 'vue-i18n'

export default {
  name: 'ServiceList',
  props: {
    displayServices: { type: Array, required: true },
    selectedService: { type: String, default: '' },
  },
  emits: ['new', 'refresh', 'select', 'copy'],
  setup() {
    const { t } = useI18n()

    function statusClass(status) {
      switch (status) {
        case 'Running': return 'status-running'
        case 'Stopped': return 'status-stopped'
        case 'Not Installed': return 'status-file'
        default: return 'status-other'
      }
    }

    return { t, statusClass }
  },
}
</script>
