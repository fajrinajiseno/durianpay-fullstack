export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) {
    return
  }

  const auth = useAuthStore()

  if (!auth.getUser()?.token && to.path === '/dashboard') {
    window.location.assign('/login')
  } else if (auth.getUser()?.token && to.path === '/login') {
    window.location.assign('/dashboard')
  }
})
