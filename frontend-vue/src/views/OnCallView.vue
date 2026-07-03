<script setup lang="ts">
import { listOnCalls, createShift, updateShift } from '@/api';
import { useAuth } from '@/stores/auth';
import type { OnCallShiftEntry } from '@/types';
import { computed, onMounted, ref } from 'vue';

const auth = useAuth()
const isAdmin = computed(() => auth.user?.role === 'admin')

const shifts = ref<OnCallShiftEntry[]>([])
const listError = ref('')

// filter state
const filterFrom = ref('')  // datetime-local value
const filterTo = ref('')

// form state — editingId null = create mode, set = update mode
const editingId = ref<string | null>(null)
const fService = ref('')
const fUsername = ref('')
const fStartsAt = ref('')  // datetime-local value
const fEndsAt = ref('')
const formError = ref('')
const successMsg = ref('')

async function loadShifts() {
  listError.value = ''
  try {
    const from = filterFrom.value ? new Date(filterFrom.value).toISOString() : undefined
    const to = filterTo.value ? new Date(filterTo.value).toISOString() : undefined
    shifts.value = await listOnCalls(from, to)
  } catch (e) {
    listError.value = 'Could not load on-call shifts.'
  }
}

onMounted(() => {
  const now = new Date()
  
  loadShifts()
  const startOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0)
  const endOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59)

  filterFrom.value = toLocalInput(startOfDay.toISOString())
  filterTo.value = toLocalInput(endOfDay.toISOString())
})

function applyFilter() {
  if (filterFrom.value && filterTo.value && new Date(filterFrom.value) >= new Date(filterTo.value)) {
    listError.value = 'Filter: From must be before To.'
    return
  }
  loadShifts()
}

function clearFilter() {
  filterFrom.value = ''
  filterTo.value = ''
  loadShifts()
}

function resetForm() {
  editingId.value = null
  fService.value = ''
  fUsername.value = ''
  fStartsAt.value = ''
  fEndsAt.value = ''
  formError.value = ''
}

function startEdit(s: OnCallShiftEntry) {
  editingId.value = s.id
  fService.value = s.service
  fUsername.value = s.username
  fStartsAt.value = toLocalInput(s.starts_at)
  fEndsAt.value = toLocalInput(s.ends_at)
  formError.value = ''
  successMsg.value = ''
}

async function handleSubmit() {
  formError.value = ''
  successMsg.value = ''

  if (!fService.value.trim() || !fUsername.value.trim()) {
    formError.value = 'Service and username are required.'
    return
  }
  if (!fStartsAt.value || !fEndsAt.value) {
    formError.value = 'Start and end are required.'
    return
  }
  const startsAt = new Date(fStartsAt.value)
  const endsAt = new Date(fEndsAt.value)
  if (startsAt >= endsAt) {
    formError.value = 'Start must be before end.'
    return
  }

  const payload = {
    service: fService.value.trim(),
    username: fUsername.value.trim(),
    starts_at: startsAt.toISOString(),
    ends_at: endsAt.toISOString(),
  }

  try {
    if (editingId.value) {
      await updateShift(editingId.value, payload)
      successMsg.value = 'Shift updated.'
    } else {
      await createShift(payload)
      successMsg.value = 'Shift created.'
    }
    resetForm()
    await loadShifts()
  } catch (e) {
    formError.value = (e as Error).message
  }
}

// ISO string -> value for <input type="datetime-local"> (local time, no seconds/zone)
function toLocalInput(iso: string): string {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function fmt(iso: string): string {
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <div class="page">
    <RouterLink :to="{ name: 'incidents' }" class="back mono">← Back to incidents</RouterLink>
    <p class="eyebrow">Schedule</p>
    <h1 class="page-title">On-call</h1>

    <p v-if="!isAdmin" class="error">You need admin access to manage on-call shifts.</p>

    <div class="detail-grid">
      <div class="detail-main">
        <div class="panel">
          <h3 class="panel-title">Shifts</h3>

          <div class="filters">
            <div class="filter">
              <label class="field-label">From</label>
              <input class="input" type="datetime-local" v-model="filterFrom">
            </div>
            <div class="filter">
              <label class="field-label">To</label>
              <input class="input" type="datetime-local" v-model="filterTo">
            </div>
            <div class="filter-actions">
              <button class="btn" @click="clearFilter">Clear</button>
              <button class="btn btn-primary" @click="applyFilter">Filter</button>
            </div>
          </div>

          <p v-if="listError" class="error">{{ listError }}</p>
          <p v-else-if="!shifts.length" class="dim">No shifts scheduled.</p>
          <table v-else class="shift-table">
            <thead>
              <tr>
                <th>Service</th><th>On call</th><th>Start</th><th>End</th>
                <th v-if="isAdmin"></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="s in shifts" :key="s.id">
                <td class="mono">{{ s.service }}</td>
                <td class="mono">{{ s.username }}</td>
                <td class="mono dim">{{ fmt(s.starts_at) }}</td>
                <td class="mono dim">{{ fmt(s.ends_at) }}</td>
                <td v-if="isAdmin">
                  <button class="btn edit-btn" @click="startEdit(s)">Edit</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <aside v-if="isAdmin" class="detail-side">
        <div class="panel">
          <h3 class="panel-title">{{ editingId ? 'Edit shift' : 'New shift' }}</h3>
          <div class="field">
            <label class="field-label">Service</label>
            <input class="input" v-model="fService" placeholder="e.g Payment">
          </div>
          <div class="field">
            <label class="field-label">Username</label>
            <input class="input" v-model="fUsername" placeholder="username on call">
          </div>
          <div class="field">
            <label class="field-label">Starts at</label>
            <input class="input" type="datetime-local" v-model="fStartsAt">
          </div>
          <div class="field">
            <label class="field-label">Ends at</label>
            <input class="input" type="datetime-local" v-model="fEndsAt">
          </div>
          <p class="error">{{ formError }}</p>
          <p class="success-message">{{ successMsg }}</p>
          <div class="form-actions">
            <button v-if="editingId" class="btn" @click="resetForm">Cancel</button>
            <button class="btn btn-primary" @click="handleSubmit">
              {{ editingId ? 'Update' : 'Create' }}
            </button>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.back {
  color: var(--color-text-dim);
  display: inline-block;
  font-size: 13px;
  margin-bottom: 20px;
}

.back:hover {
    color: var(--color-text-bright);
}

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
.detail-side {
    width: 300px;
}

.filters {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  margin-bottom: 18px;
  flex-wrap: wrap;
}
.filter {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.filter-actions {
    display: flex;
    gap: 8px;
}

.shift-table {
    border-collapse: collapse;
    width: 100%;
    font-size: 13px;
}
.shift-table th {
  color: var(--color-text-dim);
  font-size: 11px;
  letter-spacing: 1px;
  text-transform: uppercase;
  text-align: left;
  padding: 6px 8px;
  border-bottom: 1px solid var(--color-border);
}
.shift-table td {
    padding: 8px;
    border-bottom: 1px solid var(--color-border);
}
.shift-table tr:last-child td {
    border-bottom: none;
}

.edit-btn {
    padding: 4px 10px;
    font-size: 12px;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 8px;
}
.success-message {
    color: var(--color-success, #3fb950);
    font-size: 13px;
}

@media (max-width: 768px) {
  .detail-grid {
    flex-direction: column;
}
  .detail-side {
    width: 100%;
}
  .filters {
    flex-direction: column;
    align-items: stretch;
}
}
</style>