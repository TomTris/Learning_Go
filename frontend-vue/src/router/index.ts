import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: { name: 'incidents' } },
        { path: 'incidents', name: 'incidents', component: () => import('@/views/IncidentsView.vue') },
        { path: 'incidents/new', name: 'incidents-new', component: () => import('@/views/IncidentCreateView.vue') },
        { path: 'incidents/:id', name: 'incident-detail', component: () => import('@/views/IncidentDetailView.vue') },
      ],
    },
    {
      path: '/',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        { path: 'log-in', name: 'log-in', component: () => import('@/views/LoginView.vue') },
        { path: 'registration', name: 'registration', component: () => import('@/views/RegistrationView.vue') },],
    },
    { path: '/:pathMatch(.*)*', redirect: { name: 'incidents' } },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuth()

  if (auth.user === null) {
    try {
      await auth.load()
    } catch {
      // not authenticated, ignores
    }
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'log-in' }
  }
  if (to.name === 'log-in' && auth.isAuthenticated) {
    return { name: 'incidents' }
  }
})
export default router
