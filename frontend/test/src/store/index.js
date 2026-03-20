import { ref, inject, reactive } from 'vue'

// 检查 localStorage 是否可用
const getLocalStorageItem = (key) => {
  try {
    return localStorage.getItem(key)
  } catch (error) {
    return null
  }
}

const token = ref(getLocalStorageItem('token'))
const isLoggedIn = ref(!!getLocalStorageItem('token'))
const getUserId = () => {
  const id = getLocalStorageItem('userId')
  const parsedId = parseInt(id)
  return isNaN(parsedId) ? null : parsedId
}
const userId = ref(getUserId())
const username = ref(getLocalStorageItem('username'))

const store = reactive({
  token,
  isLoggedIn,
  userId,
  username,
  
  setToken(newToken) {
    token.value = newToken
    isLoggedIn.value = !!newToken
    try {
      if (newToken) {
        localStorage.setItem('token', newToken)
      } else {
        localStorage.removeItem('token')
      }
    } catch (error) {
      console.error('Error accessing localStorage:', error)
    }
  },
  
  setUserInfo(id, name) {
    userId.value = id
    username.value = name
    try {
      if (id) {
        localStorage.setItem('userId', id)
        localStorage.setItem('username', name)
      } else {
        localStorage.removeItem('userId')
        localStorage.removeItem('username')
      }
    } catch (error) {
      console.error('Error accessing localStorage:', error)
    }
  },
  
  clearToken() {
    this.setToken(null)
    this.setUserInfo(null, null)
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