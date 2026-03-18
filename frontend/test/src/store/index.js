import { ref, inject, reactive } from 'vue'

const token = ref(localStorage.getItem('token'))
const isLoggedIn = ref(!!localStorage.getItem('token'))

const store = reactive({
  token,
  isLoggedIn,
  
  setToken(newToken) {
    token.value = newToken
    isLoggedIn.value = !!newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  },
  
  clearToken() {
    this.setToken(null)
  }
})

export const userStore = store
export const userStoreKey = 'userStore'

export function useUserStore() {
  const injectedStore = inject(userStoreKey)
  if (!injectedStore) {
    throw new Error('User store not provided')
  }
  return injectedStore
}