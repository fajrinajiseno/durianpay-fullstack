import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface User {
  email: string
  token: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)

  const USER_KEY = 'auth.user'

  // load from localStorage (client only)
  function loadFromLocalStorage() {
    if (import.meta.server) return
    try {
      const rawUser = localStorage.getItem(USER_KEY)
      if (rawUser) user.value = JSON.parse(rawUser)
    } catch (e) {
      console.warn('auth: loadFromLocalStorage failed', e)
    }
  }

  function saveToLocalStorage() {
    if (import.meta.server) return
    try {
      if (user.value) localStorage.setItem(USER_KEY, JSON.stringify(user.value))
      else localStorage.removeItem(USER_KEY)
    } catch (e) {
      console.warn('auth: saveToLocalStorage failed', e)
    }
  }

  function setUser(u: User | null) {
    user.value = u
    saveToLocalStorage()
  }

  function logout() {
    setUser(null)
  }

  return {
    user,
    logout,
    loadFromLocalStorage,
    setUser
  }
})
