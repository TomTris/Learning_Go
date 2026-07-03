<script setup lang="ts">
import { createIncident, currentOnCallAll } from '@/api';
import CurrentOnCallPanel from '@/components/CurrentOnCallPanel.vue';
import type { CreateIncidentRequest, OnCallShiftEntry, Severity } from '@/types';
import { onMounted, ref } from 'vue';

const title = ref('')
const service = ref('')
const severity = ref<Severity>('SEV1')
const errors = ref<{ title?: string; service?: string }>({})
const fetchError = ref('')
const success_msg = ref('')
const severities = [
  { value: "SEV1", desc: "Critical — full outage, page everyone" },
  { value: "SEV2", desc: "Major — degraded, customers affected" },
  { value: "SEV3", desc: "Minor — contained, low impact" },
]

const currentOnCalls = ref<OnCallShiftEntry[]>([])
const onCallError = ref('')

onMounted(async () => {
  try {
    currentOnCalls.value = await currentOnCallAll()
  } catch (e) {
    onCallError.value = 'Could not load current on-call.'
  }
})

async function handleCreateIncident() {
  errors.value = {}
  success_msg.value = ''
  fetchError.value = ''

  if (title.value.trim() === '') {
    errors.value.title = "Title can't be empty"
  }
  if (service.value.trim() === '') {
    errors.value.service = "Service can't be empty"
  }
  if (Object.keys(errors.value).length > 0) {
    return
  }

  const incReq: CreateIncidentRequest = {
    title: title.value,
    service: service.value,
    severity: severity.value,
  }
  try {
    const res = await createIncident(incReq)
    success_msg.value = "A new Incident " + res.id + " is successfully created"
  } catch (e) {
    fetchError.value = (e as Error).message
  }
}
</script>

<template>
  <main>
    <div class="page">
      <RouterLink :to="{ name: 'incidents' }" class="back mono">← Back to incident</RouterLink>
      <p class="eyebrow">Declare</p>
      <h1 class="page-title create-title">Open an Incident</h1>

      <div class="detail-grid">
        <div class="detail-main">
          <form @submit.prevent="handleCreateIncident">
            <p class="error">{{ fetchError }}</p>

            <div class="field">
              <label class="field-label">Title</label>
              <input class="input" type="text" v-model="title" placeholder="Short, specific summary of what's broken">
              <p class="error">{{ errors?.title }}</p>
            </div>

            <div class="field">
              <label class="field-label">Service</label>
              <input class="input" type="text" v-model="service" placeholder="e.g Payment">
              <p class="error">{{ errors?.service }}</p>
            </div>

            <div class="field">
              <label class="field-label">Severity</label>
              <div class="sev-picker">
                <label
                  v-for="opt in severities"
                  :key="opt.value"
                  class="sev-option"
                  :class="{ selected: opt.value === severity }"
                >
                  <input type="radio" name="severity" :value="opt.value" v-model="severity">
                  <span class="badge" :class="'sev-' + opt.value">{{ opt.value }}</span>
                  <span class="sev-desc">{{ opt.desc }}</span>
                </label>
              </div>
            </div>

            <div class="create-actions">
              <RouterLink :to="{ name: 'incidents' }" class="btn">Cancel</RouterLink>
              <button class="btn btn-primary" type="submit">Open Incident</button>
            </div>
            <p class="success-message">{{ success_msg }}</p>
          </form>
        </div>

        <aside class="detail-side">
          <CurrentOnCallPanel :on-calls="currentOnCalls" :error-msg="onCallError" />
        </aside>
      </div>
    </div>
  </main>
</template>

<style scoped>
.back {
  color: var(--color-text-dim);
  display: inline-block;
  font-size: 13px;
  margin-bottom: 20px;
}
.back:hover { color: var(--color-text-bright); }

.create-title { margin-bottom: 24px; }

/* Two-column layout: main form + sidebar — same as IncidentDetailView */
.detail-grid {
  display: flex;
  gap: 20px;
}
.detail-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 20px;
}

.sev-picker { display: flex; flex-direction: column; gap: 8px; }
.sev-option {
  display: flex;
  align-items: center;
  background-color: var(--color-input);
  border: 1px solid var(--color-border);
  border-radius: 5px;
  cursor: pointer;
  gap: 12px;
  padding: 12px;
}
.sev-option input { accent-color: var(--color-accent); }
.sev-desc { color: var(--color-text-dim); font-size: 13px; }
.sev-option.selected { border-color: var(--color-accent); }

.create-actions {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 10px;
}
.success-message { color: var(--color-success, #3fb950); font-size: 13px; }

@media (max-width: 768px) {
  .detail-grid { flex-direction: column; }
  .detail-side { width: 100%; }
}
</style>