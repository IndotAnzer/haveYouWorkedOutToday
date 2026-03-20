import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Article from '../views/Article.vue'
import axios from 'axios'
import { userStore, userStoreKey } from '../store'

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

describe('点赞和评论功能测试', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  describe('点赞功能', () => {
    it('正常场景：登录后成功点赞文章', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          message: '点赞成功'
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      await wrapper.vm.likeArticle()

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/like', {}, {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })

    it('异常场景：未登录尝试点赞文章', async () => {
      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      await wrapper.vm.likeArticle()

      expect(axios.post).not.toHaveBeenCalled()
    })

    it('正常场景：获取文章点赞数', async () => {
      const mockResponse = {
        data: {
          count: 10
        }
      }

      axios.get.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      await wrapper.vm.fetchLikes()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/like', {
        headers: {
          Authorization: 'Bearer null'
        }
      })
      expect(wrapper.vm.likes).toBe(10)
    })
  })

  describe('评论功能', () => {
    it('正常场景：登录后成功创建评论', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          id: 1,
          content: '测试评论',
          user_id: 1,
          username: 'testuser',
          created_at: '2024-01-01T00:00:00Z',
          replies: []
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      wrapper.vm.commentContent = '测试评论'
      await wrapper.vm.submitComment()

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/comments', {
        content: '测试评论'
      }, {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })

    it('异常场景：未登录尝试创建评论', async () => {
      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      wrapper.vm.commentContent = '测试评论'
      await wrapper.vm.submitComment()

      expect(axios.post).not.toHaveBeenCalled()
    })

    it('异常场景：评论内容为空', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      wrapper.vm.commentContent = ''
      await wrapper.vm.submitComment()

      expect(axios.post).not.toHaveBeenCalled()
    })

    it('正常场景：获取评论和回复', async () => {
      const mockResponse = {
        data: {
          comments: [
            {
              id: 1,
              content: '测试评论',
              user_id: 1,
              user: {
                username: 'testuser'
              },
              created_at: '2024-01-01T00:00:00Z',
              replies: [
                {
                  id: 1,
                  content: '测试回复',
                  user_id: 2,
                  user: {
                    username: 'otheruser'
                  },
                  created_at: '2024-01-01T00:01:00Z'
                }
              ]
            }
          ]
        }
      }

      axios.get.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      await wrapper.vm.fetchComments()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/comments', {
        headers: {
          Authorization: 'Bearer null'
        }
      })
      expect(wrapper.vm.comments.length).toBe(1)
      expect(wrapper.vm.comments[0].replies.length).toBe(1)
    })

    it('正常场景：登录后成功创建回复', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          id: 1,
          content: '测试回复',
          user_id: 1,
          username: 'testuser',
          created_at: '2024-01-01T00:00:00Z'
        }
      }

      axios.post.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      // 设置回复相关状态
      wrapper.vm.replyingTo = { commentId: 1, replyId: null, author: 'testuser' }
      wrapper.vm.replyContent = '测试回复'
      
      await wrapper.vm.submitReply()

      expect(axios.post).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/comments/1/replies', {
        content: '测试回复',
        parent_reply_id: null
      }, {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })

    it('正常场景：登录后成功删除评论', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          message: '删除成功'
        }
      }

      axios.delete.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      await wrapper.vm.deleteComment(1)

      expect(axios.delete).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/comments/1', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })

    it('正常场景：登录后成功删除回复', async () => {
      localStorageMock.setItem('token', 'test-token')
      localStorageMock.setItem('userId', '1')

      const mockResponse = {
        data: {
          message: '删除成功'
        }
      }

      axios.delete.mockResolvedValue(mockResponse)

      const wrapper = mount(Article, {
        props: {
          id: '1'
        },
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

      await wrapper.vm.deleteReply(1, 1)

      expect(axios.delete).toHaveBeenCalledWith('http://localhost:3001/api/articles/1/comments/1/replies/1', {
        headers: {
          Authorization: 'Bearer test-token'
        }
      })
    })
  })
})
