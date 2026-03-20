import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Home from '../views/Home.vue'
import MyArticles from '../views/MyArticles.vue'
import Article from '../views/Article.vue'
import axios from 'axios'
import { userStore, userStoreKey } from '../store'
import { createRouter, createWebHistory } from 'vue-router'

// 模拟axios
vi.mock('axios')

// 模拟vue-router
vi.mock('vue-router', () => {
  return {
    useRoute: () => ({
      params: {
        id: '1'
      }
    }),
    useRouter: () => ({
      push: vi.fn()
    })
  }
})

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

describe('文章功能测试', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  describe('创建文章', () => {
    it('正常场景：输入有效的标题和内容，创建成功', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          id: 1,
          title: '测试文章',
          content: '测试内容',
          user_id: 1,
          created_at: '2024-01-01T00:00:00Z'
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Home, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      // 先点击"记录一下"按钮显示表单
      await wrapper.find('.create-btn').trigger('click')
      // 等待表单显示
      await wrapper.vm.$nextTick()
      // 模拟输入
      await wrapper.find('input[type="text"]').setValue('测试文章')
      await wrapper.find('textarea').setValue('测试内容')
      // 点击发布按钮
      await wrapper.find('.submit-btn').trigger('click')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/articles', {
        title: '测试文章',
        content: '测试内容',
        preview: '测试内容',
        fitness_actions: []
      }, {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })

    it('异常场景：未登录尝试创建文章', async () => {
      const wrapper = mount(Home, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      // 先点击"记录一下"按钮显示表单
      await wrapper.find('.create-btn').trigger('click')
      // 等待表单显示
      await wrapper.vm.$nextTick()
      // 模拟输入
      await wrapper.find('input[type="text"]').setValue('测试文章')
      await wrapper.find('textarea').setValue('测试内容')
      // 点击发布按钮
      await wrapper.find('.submit-btn').trigger('click')

      expect(axios.post).not.toHaveBeenCalled()
    })

    it('异常场景：标题为空', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const wrapper = mount(Home, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      // 先点击"记录一下"按钮显示表单
      await wrapper.find('.create-btn').trigger('click')
      // 等待表单显示
      await wrapper.vm.$nextTick()
      // 模拟输入
      await wrapper.find('input[type="text"]').setValue('')
      await wrapper.find('textarea').setValue('测试内容')
      // 点击发布按钮
      await wrapper.find('.submit-btn').trigger('click')

      expect(axios.post).not.toHaveBeenCalled()
    })

    it('正常场景：内容为空时使用默认值', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          id: 1,
          title: '测试文章',
          content: '暂无训练心得',
          user_id: 1,
          created_at: '2024-01-01T00:00:00Z'
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Home, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })

      // 先点击"记录一下"按钮显示表单
      await wrapper.find('.create-btn').trigger('click')
      // 等待表单显示
      await wrapper.vm.$nextTick()
      // 模拟输入
      await wrapper.find('input[type="text"]').setValue('测试文章')
      await wrapper.find('textarea').setValue('')
      // 点击发布按钮
      await wrapper.find('.submit-btn').trigger('click')

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/articles', {
        title: '测试文章',
        content: '暂无训练心得',
        preview: '暂无训练心得',
        fitness_actions: []
      }, {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })
  })

  describe('获取所有文章', () => {
    it('正常场景：成功获取所有文章列表', async () => {
      // 设置token
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: [
          {
            id: 1,
            title: '测试文章1',
            content: '测试内容1',
            user_id: 1,
            username: 'testuser',
            created_at: '2024-01-01T00:00:00Z'
          },
          {
            id: 2,
            title: '测试文章2',
            content: '测试内容2',
            user_id: 2,
            username: 'otheruser',
            created_at: '2024-01-02T00:00:00Z'
          }
        ]
      }

      axios.get.mockResolvedValue(mockResponse)

      const wrapper = mount(Home, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })
      await wrapper.vm.$nextTick()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:3001/api/articles', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
      expect(wrapper.vm.articles.length).toBe(2)
    })

    it('边界场景：没有文章时返回空列表', async () => {
      // 设置token
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: []
      }

      axios.get.mockResolvedValue(mockResponse)

      const wrapper = mount(Home, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })
      await wrapper.vm.$nextTick()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:3001/api/articles', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
      expect(wrapper.vm.articles.length).toBe(0)
    })
  })

  describe('获取用户自己的文章', () => {
    it('正常场景：登录后成功获取自己的文章列表', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: [
          {
            id: 1,
            title: '我的文章1',
            content: '我的内容1',
            user_id: 1,
            created_at: '2024-01-01T00:00:00Z'
          }
        ]
      }

      axios.get.mockResolvedValue(mockResponse)

      const wrapper = mount(MyArticles, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })
      await wrapper.vm.$nextTick()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:3001/api/my/articles', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
      expect(wrapper.vm.articles.length).toBe(1)
    })

    it('边界场景：用户没有文章时返回空列表', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: []
      }

      axios.get.mockResolvedValue(mockResponse)

      const wrapper = mount(MyArticles, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })
      await wrapper.vm.$nextTick()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:3001/api/my/articles', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
      expect(wrapper.vm.articles.length).toBe(0)
    })
  })

  describe('删除文章', () => {
    it('正常场景：登录后成功删除自己的文章', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          message: '删除成功'
        }
      }

      axios.delete.mockResolvedValue(mockResponse)

      const wrapper = mount(MyArticles, {
        global: {
          provide: {
            [userStoreKey]: userStore
          }
        }
      })
      // 设置要删除的文章
      wrapper.vm.articleToDelete = {
        id: 1,
        title: '测试文章'
      }
      // 调用删除方法
      await wrapper.vm.deleteArticle()

      expect(axios.delete).toHaveBeenCalledWith('http://localhost:3001/api/articles/1', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })
  })

  describe('获取文章详情', () => {
    it('正常场景：成功获取文章详情', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockArticleResponse = {
        data: {
          id: 1,
          title: '测试文章',
          content: '测试内容',
          user_id: 1,
          username: 'testuser',
          created_at: '2024-01-01T00:00:00Z'
        }
      }

      const mockLikesResponse = {
        data: {
          count: 5,
          is_liked: false
        }
      }

      const mockCommentsResponse = {
        data: {
          comments: []
        }
      }

      axios.get.mockImplementation((url) => {
        if (url.includes('articles/1') && !url.includes('like') && !url.includes('comments')) {
          return Promise.resolve(mockArticleResponse)
        } else if (url.includes('articles/1/like')) {
          return Promise.resolve(mockLikesResponse)
        } else if (url.includes('articles/1/comments')) {
          return Promise.resolve(mockCommentsResponse)
        }
        return Promise.reject(new Error('Unexpected URL'))
      })

      const wrapper = mount(Article, {
        global: {
          provide: {
            [userStoreKey]: userStore
          },
          mocks: {
            $route: {
              params: {
                id: '1'
              }
            },
            $router: {
              push: vi.fn()
            }
          }
        }
      })
      
      // 等待异步操作完成
      await wrapper.vm.$nextTick()
      // 手动调用fetchArticle方法
      await wrapper.vm.fetchArticle()
      await wrapper.vm.$nextTick()
      // 手动调用fetchComments方法
      await wrapper.vm.fetchComments()
      await wrapper.vm.$nextTick()

      expect(axios.get).toHaveBeenCalled()
      expect(wrapper.vm.article).not.toBeNull()
      expect(wrapper.vm.article.title).toBe('测试文章')
      expect(wrapper.vm.likes).toBe(5)
    })
  })
})
