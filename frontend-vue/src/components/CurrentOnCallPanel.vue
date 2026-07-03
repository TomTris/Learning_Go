<script setup lang="ts">
import type { OnCallShiftEntry } from '@/types';

defineProps<{
  onCalls: OnCallShiftEntry[];
  errorMsg: string;
}>()
</script>

<template>
  <div class="panel oncall-panel">
    <h3 class="panel-title">Currently on call</h3>
    <p v-if="errorMsg" class="error">{{ errorMsg }}</p>
    <p v-else-if="!onCalls.length" class="dim oncall-empty">No one is currently on call.</p>
    <table v-else class="oncall-table">
      <thead>
        <tr>
          <th>Service</th>
          <th>On call</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="entry in onCalls" :key="entry.id">
          <td class="mono">{{ entry.service }}</td>
          <td class="mono">{{ entry.username }}</td>
        </tr>
        <p class="dim hint">Ask your admin to create more oncall shifts</p>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.panel-title {
  display: flex;
  align-items: center;
  flex-direction: column;
}

.oncall-table {
  width: 280px;
  font-size: 13px;
}

.oncall-table th {
  color: var(--color-text-dim);
  font-size: 11px;
  letter-spacing: 1px;
  text-transform: uppercase;
  text-align: left;
  padding: 6px 8px;
  border-bottom: 1px solid var(--color-border);
}

.oncall-table td {
  color: var(--color-text);
  padding: 8px;
  border-bottom: 1px solid var(--color-border);
}

.oncall-table tr:last-child td {
  border-bottom: none;
}

.oncall-empty {
  font-size: 13px;
}

@media (max-width: 768px) {
  .oncall-panel {
    width: 100%;
  }
}
</style>