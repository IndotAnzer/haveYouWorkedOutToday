<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const loading = ref(false)
const error = ref('')
const likes = ref(0)
const isLiked = ref(false)

const fetchArticle = async () => {
  const id = route.params.id
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get(`http://localhost:3000/api/articles/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    article.value = response.data
    await fetchLikes()
  } catch (err) {
    error.value = err.response?.data?.message || '获取记录失败'
  } finally {
    loading.value = false
  }
}

const fetchLikes = async () => {
  const id = route.params.id
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get(`http://localhost:3000/api/articles/${id}/like`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    likes.value = response.data.count
    isLiked.value = response.data.is_liked
  } catch (err) {
    console.error('获取点赞信息失败:', err)
  }
}

const likeArticle = async () => {
  const id = route.params.id
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(`http://localhost:3000/api/articles/${id}/like`, {}, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    if (response.data.action === 'liked') {
      isLiked.value = true
      likes.value = parseInt(likes.value) + 1
    } else if (response.data.action === 'unliked') {
      isLiked.value = false
      likes.value = Math.max(0, parseInt(likes.value) - 1)
    }
  } catch (err) {
    error.value = err.response?.data?.message || '点赞失败'
  }
}

const goBack = () => {
  router.push('/')
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { 
    year: 'numeric',
    month: 'long', 
    day: 'numeric',
    weekday: 'long'
  })
}

onMounted(() => {
  fetchArticle()
})
</script>

<template>
  <div class="article-container">
    <button @click="goBack" class="back-button">
      <span>←</span> 返回首页
    </button>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <span class="error-icon">😢</span>
      <h3>出错了</h3>
      <p>{{ error }}</p>
      <button @click="goBack" class="back-btn">返回首页</button>
    </div>

    <div v-else-if="article" class="article-content">
      <div class="article-header">
        <h1>{{ article.title }}</h1>
        <div class="article-meta">
          <span class="meta-date">📅 {{ formatDate(article.created_at || article.CreatedAt) }}</span>
        </div>
      </div>

      <div v-if="article.content" class="article-body">
        <p>{{ article.content }}</p>
      </div>

      <div v-if="article.fitness_actions && article.fitness_actions.length > 0" class="fitness-actions">
        <div class="section-header">
          <h2>🎯 训练动作</h2>
          <span class="action-count">{{ article.fitness_actions.length }} 个动作</span>
        </div>

        <div class="actions-grid">
          <div v-for="action in article.fitness_actions" :key="action.ID" class="action-card">
            <div class="action-header">
              <h3>{{ action.action_name }}</h3>
              <span v-if="action.remark" class="action-remark">{{ action.remark }}</span>
            </div>

            <div v-if="action.action_groups && action.action_groups.length > 0" class="groups-container">
              <div class="groups-header">
                <span>组数</span>
                <span>重量</span>
                <span>次数</span>
                <span>备注</span>
              </div>
              <div v-for="group in action.action_groups" :key="group.ID" class="group-row">
                <span class="group-index">第 {{ group.group_index }} 组</span>
                <span class="group-weight">{{ group.weight }} kg</span>
                <span class="group-reps">{{ group.rep_num }} 次</span>
                <span class="group-remark">{{ group.remark || '-' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="article-footer">
        <button 
          @click="likeArticle" 
          class="like-button"
          :class="{ 'liked': isLiked }"
        >
          <span class="like-icon">{{ isLiked ? '❤️' : '🤍' }}</span>
          <span class="like-count">{{ likes }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.article-container {
  max-width: 100%;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  color: #555;
  border: none;
  padding: 12px 20px;
  border-radius: 14px;
  font-weight: 500;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
}

.back-button:hover {
  background: white;
  transform: translateX(-4px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.loading-state {
  text-align: center;
  padding: 80px 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #e8e8e8;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state {
  text-align: center;
  padding: 60px 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
}

.error-state .error-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.error-state h3 {
  font-size: 22px;
  color: #333;
  margin: 0 0 8px 0;
}

.error-state p {
  color: #888;
  margin: 0 0 24px 0;
}

.back-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 28px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.back-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.article-content {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.article-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 32px;
  color: white;
}

.article-header h1 {
  font-size: 28px;
  margin: 0 0 12px 0;
  font-weight: 700;
}

.article-meta {
  display: flex;
  gap: 16px;
  opacity: 0.9;
  font-size: 14px;
}

.article-body {
  padding: 28px 32px;
  font-size: 16px;
  line-height: 1.8;
  color: #444;
  border-bottom: 1px solid #f0f0f0;
}

.article-body p {
  margin: 0;
  white-space: pre-wrap;
}

.fitness-actions {
  padding: 28px 32px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 20px;
  color: #333;
  margin: 0;
  font-weight: 600;
}

.action-count {
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
  color: #667eea;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.actions-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.action-card {
  background: #fafafa;
  border: 2px solid #e8e8e8;
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
}

.action-card:hover {
  border-color: #667eea;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.1);
}

.action-header {
  margin-bottom: 16px;
}

.action-header h3 {
  font-size: 18px;
  color: #667eea;
  margin: 0 0 6px 0;
  font-weight: 600;
}

.action-remark {
  font-size: 13px;
  color: #888;
  background: white;
  padding: 4px 10px;
  border-radius: 8px;
}

.groups-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e8e8e8;
}

.groups-header {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1.5fr;
  padding: 12px 16px;
  background: #f5f5f5;
  font-size: 13px;
  font-weight: 600;
  color: #666;
}

.group-row {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr 1.5fr;
  padding: 14px 16px;
  border-top: 1px solid #e8e8e8;
  font-size: 14px;
  transition: background 0.2s ease;
}

.group-row:hover {
  background: #fafafa;
}

.group-index {
  color: #667eea;
  font-weight: 600;
}

.group-weight {
  color: #52c41a;
  font-weight: 600;
}

.group-reps {
  color: #ff6b6b;
  font-weight: 600;
}

.group-remark {
  color: #888;
}

.article-footer {
  padding: 24px 32px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: flex-end;
}

.like-button {
  display: flex;
  align-items: center;
  gap: 10px;
  background: white;
  border: 2px solid #d9d9d9;
  color: #8c8c8c;
  padding: 14px 28px;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 15px;
  font-weight: 600;
}

.like-button:hover {
  background: #fafafa;
  border-color: #bfbfbf;
  transform: scale(1.02);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.like-button.liked {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a5a 100%);
  color: white;
  border-color: transparent;
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.4);
}

.like-icon {
  font-size: 20px;
}

.like-count {
  background: rgba(0, 0, 0, 0.1);
  padding: 4px 10px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
}

.like-button.liked .like-count {
  background: rgba(255, 255, 255, 0.2);
}

@media (max-width: 768px) {
  .article-header {
    padding: 24px;
  }

  .article-header h1 {
    font-size: 22px;
  }

  .article-body,
  .fitness-actions,
  .article-footer {
    padding: 20px;
  }

  .groups-header,
  .group-row {
    grid-template-columns: 1fr 1fr 1fr;
    gap: 8px;
  }

  .groups-header span:last-child,
  .group-row .group-remark {
    display: none;
  }

  .like-button {
    padding: 12px 20px;
  }
}

@media (max-width: 480px) {
  .article-header h1 {
    font-size: 18px;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
