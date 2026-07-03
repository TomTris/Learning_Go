import { defineStore } from 'pinia'
import { whoAmI } from '@/api'
import type { UserContext } from '@/types'

export const useAuth = defineStore('auth', {
    state: () => ({ user: null as UserContext | null }),
    getters: {
        isAuthenticated: (s) => s.user !== null,
    },
    actions: {
        async load() {
            this.user = await whoAmI()
        },
        clear() {
            this.user = null
        },
    },
})