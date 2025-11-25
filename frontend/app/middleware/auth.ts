import { useAuthStore } from '../store/auth'

export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) {
    return
  }

  const auth = useAuthStore()
  auth.loadFromLocalStorage()

  if (!auth.user?.token && to.path === '/dashboard') {
    window.location.assign('/login')
  } else if (auth.user?.token && to.path === '/login') {
    window.location.assign('/dashboard')
  }
})
