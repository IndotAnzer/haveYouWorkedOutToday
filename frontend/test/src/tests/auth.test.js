import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import axios from 'axios'
import { userStore, userStoreKey } from '../store'

// 模拟axios
vi.mock('axios')

// 模拟localStorage
const localStorageMock = (() => {
  let store = {}
  return {
    getItem: (key) => store[key] || null,
    setItem: (key, value) => store[key] = value.toString(),
    clear: () => store = {}
  }
})()

// 检查 window 是否存在
if (typeof window !== 'undefined') {
  Object.defineProperty(window, 'localStorage', { value: localStorageMock })
} else {
  // 在非浏览器环境中，直接设置 global.localStorage
  global.localStorage = localStorageMock
}

// 模拟 vue-router
const mockRouter = {
  push: vi.fn()
}

vi.mock('vue-router', () => ({
  useRouter: () => mockRouter
}))

describe('认证功能测试', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  describe('登录功能', () => {
    it('正常场景：输入正确的用户名和密码，登录成功', async () => {
      const mockResponse = {
        data: {
          token: 'test-token',
          user: {
            id: 1,
            username: 'testuser'
          }
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Login, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      await wrapper.find('#username').setValue('testuser')
      await wrapper.find('#password').setValue('password123')
      await wrapper.find('form').trigger('submit.prevent')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/auth/login', {
        username: 'testuser',
        password: 'password123'
      })
      expect(localStorage.getItem('token')).toBe('test-token')
      expect(localStorage.getItem('userId')).toBe('1')
    })

    it('异常场景：用户名不存在', async () => {
      const mockError = {
        response: {
          data: {
            message: '用户不存在'
          }
        }
      }

      axios.post.mockRejectedValue(mockError)

      const wrapper = mount(Login, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      await wrapper.find('#username').setValue('nonexistent')
      await wrapper.find('#password').setValue('password123')
      await wrapper.find('form').trigger('submit.prevent')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/auth/login', {
        username: 'nonexistent',
        password: 'password123'
      })
    })

    it('异常场景：密码错误', async () => {
      const mockError = {
        response: {
          data: {
            message: '密码错误'
          }
        }
      }

      axios.post.mockRejectedValue(mockError)

      const wrapper = mount(Login, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      await wrapper.find('#username').setValue('testuser')
      await wrapper.find('#password').setValue('wrongpassword')
      await wrapper.find('form').trigger('submit.prevent')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/auth/login', {
        username: 'testuser',
        password: 'wrongpassword'
      })
    })
  })

  describe('注册功能', () => {
    it('正常场景：输入有效的用户名、密码和确认密码，注册成功', async () => {
      const mockResponse = {
        data: {
          message: '注册成功'
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Register, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      await wrapper.find('#username').setValue('newuser')
      await wrapper.find('#password').setValue('password123')
      await wrapper.find('#confirmPassword').setValue('password123')
      await wrapper.find('form').trigger('submit.prevent')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/auth/register', {
        username: 'newuser',
        password: 'password123'
      })
    })

    it('异常场景：用户名已存在', async () => {
      const mockError = {
        response: {
          data: {
            message: '用户名已存在'
          }
        }
      }

      axios.post.mockRejectedValue(mockError)

      const wrapper = mount(Register, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      await wrapper.find('#username').setValue('existinguser')
      await wrapper.find('#password').setValue('password123')
      await wrapper.find('#confirmPassword').setValue('password123')
      await wrapper.find('form').trigger('submit.prevent')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/auth/register', {
        username: 'existinguser',
        password: 'password123'
      })
    })

    it('异常场景：确认密码与密码不一致', async () => {
      const wrapper = mount(Register, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      await wrapper.find('#username').setValue('newuser')
      await wrapper.find('#password').setValue('password123')
      await wrapper.find('#confirmPassword').setValue('differentpassword')
      await wrapper.find('form').trigger('submit.prevent')

      expect(axios.post).not.toHaveBeenCalled()
    })
  })
})
