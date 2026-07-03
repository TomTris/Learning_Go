<script setup lang="ts">
import { loadIncidents } from '@/api';
import IncidentListItem from '@/components/IncidentListItem.vue';
import type { Incident } from '@/types';
import { computed, onMounted, ref } from 'vue';

const incidents = ref<Incident[]>([])

onMounted(async () => {
  incidents.value = await loadIncidents()
})

const filterStatus = ref('')
const service = ref('')
const filteredIncidents = computed(() => {
  let result = incidents.value
  if (service.value.trim() !== '') {
    result = result.filter((inc) => inc.service.includes(service.value.trim()))
  }
  if (filterStatus.value !== '') {
    result = result.filter((inc) => inc.status === filterStatus.value)
  }
  return result
})
</script>

<template>
  <main>
    <div class="page">
      <div class="dash-head">
        <div>
          <p class="eyebrow">Active board</p>
          <h1 class="page-title">Incidents</h1>
        </div>
        <div class="spacer"></div>
        <RouterLink :to="{name: 'incidents-new'}" class="btn btn-primary">+ New Incident</RouterLink>
      </div>

      <div class="filters">
        <div class="filter">
          <label class="field-label">Status</label>
          <select class="select" v-model="filterStatus">
            <option value="">All</option>
            <option value="active">Active</option>
            <option value="triggered">Triggered</option>
            <option value="investigating">investigating</option>
            <option value="mitigated">Mitigated</option>
            <option value="resolved">Resolved</option>
          </select>
        </div>
        <div class="filter">
          <label class="field-label">Service</label>
          <input class="input" type="text" v-model="service" placeholder="e.g Payment">
        </div>
      </div>

      <ul class="incident-list">
        <li v-for="inc in filteredIncidents">
          <IncidentListItem :key="inc.id" :inc="inc"></IncidentListItem>
        </li>
      </ul>


    </div>
  </main>
</template>

<style scoped>
.dash-head {
  align-items: flex-end;
  display: flex;
  margin-bottom: 24px;
}

.filters {
  display: flex;
  gap: 16px;
  list-style: none;
  margin-bottom: 24px;
}

.filter {
  width: 220px;
}

.incident-list{
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}


@media (max-width: 768px) {
  .dash-head {
    align-items: stretch;
    flex-direction: column;
    gap: 14px;
  }

  .filters {
    flex-direction: column;
  }

  .filter {
    width: 100%;
  }

}
</style>