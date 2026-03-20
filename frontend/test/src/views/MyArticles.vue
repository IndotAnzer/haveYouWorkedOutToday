<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const articles = ref([])
const loading = ref(false)
const error = ref('')
const deletingId = ref(null)
const showDeleteConfirm = ref(false)
const articleToDelete = ref(null)

const fetchMyArticles = async () => {
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get('http://localhost:3001/api/my/articles', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    articles.value = response.data
  } catch (err) {
    error.value = err.response?.data?.error || '获取文章失败'
  } finally {
    loading.value = false
  }
}

const navigateToArticle = (article) => {
  const id = article.ID || article.id
  router.push(`/article/${id}`)
}

const confirmDelete = (article, event) => {
  event.stopPropagation()
  articleToDelete.value = article
  showDeleteConfirm.value = true
}

const cancelDelete = () => {
  showDeleteConfirm.value = false
  articleToDelete.value = null
}

const deleteArticle = async () => {
  if (!articleToDelete.value) return
  
  try {
    const token = localStorage.getItem('token')
    const id = articleToDelete.value.ID || articleToDelete.value.id
    deletingId.value = id
    
    await axios.delete(`http://localhost:3001/api/articles/${id}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    await fetchMyArticles()
    showDeleteConfirm.value = false
    articleToDelete.value = null
  } catch (err) {
    error.value = err.response?.data?.error || '删除失败'
  } finally {
    deletingId.value = null
  }
}

const goHome = () => {
  router.push('/')
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getActionCount = (article) => {
  return article.fitness_actions?.length || 0
}

onMounted(() => {
  fetchMyArticles()
})
</script>

<template>
  <div class="my-articles-page">
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
    </div>

    <div class="my-articles-container">
      <div class="page-header">
        <div class="header-content">
          <div class="header-icon">📝</div>
          <div class="header-text">
            <h2>我的记录</h2>
            <p>查看和管理你的训练记录</p>
          </div>
        </div>
        <button @click="goHome" class="home-button">
          <span class="btn-icon">🏠</span>
          <span>返回首页</span>
        </button>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <span class="error-icon">😢</span>
        <h3>出错了</h3>
        <p>{{ error }}</p>
        <button @click="fetchMyArticles" class="retry-btn">
          <span>🔄</span> 重试
        </button>
      </div>

      <div v-else-if="articles.length === 0" class="empty-state">
        <div class="empty-icon-wrapper">
          <span class="empty-icon">🏋️</span>
        </div>
        <h3>还没有训练记录</h3>
        <p>开始记录你的第一次训练吧！</p>
        <button @click="goHome" class="create-button">
          <span class="btn-icon">✨</span>
          <span>去记录</span>
        </button>
      </div>

      <div v-else class="article-grid">
        <transition-group name="list" tag="div" class="article-grid-inner">
          <div 
            v-for="article in articles" 
            :key="article.ID" 
            class="article-card"
            @click="navigateToArticle(article)"
          >
            <div class="card-header">
              <h3 class="card-title">{{ article.title }}</h3>
              <button 
                @click="confirmDelete(article, $event)" 
                class="delete-button"
                :disabled="deletingId === article.ID"
              >
                <span v-if="deletingId === article.ID" class="deleting-spinner"></span>
                <span v-else>🗑️</span>
              </button>
            </div>

            <div class="card-body">
              <p v-if="article.content" class="card-content">
                {{ article.content.substring(0, 80) }}{{ article.content.length > 80 ? '...' : '' }}
              </p>
              <p v-else class="card-content empty-text">
                暂无训练心得
              </p>
            </div>

            <div class="card-footer">
              <div class="card-stats">
                <span class="stat-item">
                  <span class="stat-icon">🎯</span>
                  {{ getActionCount(article) }} 个动作
                </span>
              </div>
              <div class="card-date">
                <span class="date-icon">📅</span>
                {{ formatDate(article.CreatedAt || article.created_at) }}
              </div>
            </div>

            <div class="card-hover-effect"></div>
          </div>
        </transition-group>
      </div>
    </div>

    <transition name="modal">
      <div v-if="showDeleteConfirm" class="modal-overlay" @click="cancelDelete">
        <div class="modal-container" @click.stop>
          <div class="modal-icon">⚠️</div>
          <h3 class="modal-title">确认删除</h3>
          <p class="modal-message">
            确定要删除记录「{{ articleToDelete?.title }}」吗？<br>
            此操作无法撤销。
          </p>
          <div class="modal-actions">
            <button @click="cancelDelete" class="modal-btn cancel-btn">
              取消
            </button>
            <button @click="deleteArticle" class="modal-btn confirm-btn">
              确认删除
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.my-articles-page {
  min-height: calc(100vh - 180px);
  padding: 40px 20px;
  position: relative;
  overflow: hidden;
}

.background-decoration {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  overflow: hidden;
  z-index: 0;
}

.circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.08;
  animation: float 25s infinite ease-in-out;
}

.circle-1 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  top: -250px;
  right: -150px;
  animation-delay: 0s;
}

.circle-2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -200px;
  left: -150px;
  animation-delay: 7s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(30px, -30px) scale(1.05);
  }
  50% {
    transform: translate(-30px, 30px) scale(0.95);
  }
  75% {
    transform: translate(30px, 30px) scale(1.02);
  }
}

.my-articles-container {
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
  padding: 24px 32px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
  animation: slideDown 0.5s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.header-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.header-icon {
  font-size: 48px;
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.header-text h2 {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 4px;
}

.header-text p {
  color: #6b7280;
  font-size: 14px;
}

.home-button {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.home-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.btn-icon {
  font-size: 18px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loading-state p {
  color: #6b7280;
  font-size: 16px;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
}

.error-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.error-state h3 {
  font-size: 24px;
  color: #1a1a2e;
  margin-bottom: 10px;
}

.error-state p {
  color: #6b7280;
  margin-bottom: 24px;
}

.retry-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  transform: translateY(-2px);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
}

.empty-icon-wrapper {
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(102, 126, 234, 0.4);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 0 0 20px rgba(102, 126, 234, 0);
  }
}

.empty-icon {
  font-size: 50px;
}

.empty-state h3 {
  font-size: 24px;
  color: #1a1a2e;
  margin-bottom: 10px;
}

.empty-state p {
  color: #6b7280;
  margin-bottom: 24px;
}

.create-button {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 28px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.create-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.article-grid-inner {
  display: contents;
}

.article-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.8);
}

.article-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

.card-hover-effect {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2, #f093fb, #f5576c);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.article-card:hover .card-hover-effect {
  opacity: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.card-title {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.4;
  flex: 1;
  margin-right: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.delete-button {
  background: rgba(239, 68, 68, 0.1);
  border: none;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.delete-button:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.2);
  transform: scale(1.1);
}

.delete-button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.deleting-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #ef4444;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.card-body {
  margin-bottom: 16px;
}

.card-content {
  color: #6b7280;
  font-size: 14px;
  line-height: 1.6;
}

.empty-text {
  font-style: italic;
  opacity: 0.6;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #f3f4f6;
}

.card-stats {
  display: flex;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
  background: #f3f4f6;
  padding: 6px 12px;
  border-radius: 8px;
}

.stat-icon {
  font-size: 14px;
}

.card-date {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #9ca3af;
}

.date-icon {
  font-size: 14px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.4s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.list-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

@media (max-width: 768px) {
  .my-articles-page {
    padding: 20px;
  }

  .page-header {
    flex-direction: column;
    gap: 20px;
    padding: 20px;
    text-align: center;
  }

  .header-content {
    flex-direction: column;
    gap: 12px;
  }

  .header-text h2 {
    font-size: 24px;
  }

  .article-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .article-card {
    padding: 20px;
  }
}

@media (max-width: 480px) {
  .header-icon {
    font-size: 40px;
  }

  .card-title {
    font-size: 16px;
  }

  .card-footer {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-container {
  background: white;
  border-radius: 20px;
  padding: 32px;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  text-align: center;
}

.modal-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.modal-title {
  font-size: 22px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 12px 0;
}

.modal-message {
  color: #6b7280;
  font-size: 15px;
  line-height: 1.6;
  margin: 0 0 24px 0;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.modal-btn {
  flex: 1;
  padding: 12px 24px;
  border-radius: 12px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.cancel-btn:hover {
  background: #e8e8e8;
}

.confirm-btn {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a5a 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.4);
}

.confirm-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 107, 107, 0.5);
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9);
}
</style>
