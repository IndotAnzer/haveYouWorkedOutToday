<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { useUserStore } from '../store'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const article = ref(null)
const loading = ref(false)
const error = ref('')
const likes = ref(0)
const isLiked = ref(false)
const comments = ref([])
const commentContent = ref('')
const submittingComment = ref(false)
const commentError = ref('')
const replyContent = ref('')
const submittingReply = ref(false)
const replyError = ref('')
const replyingTo = ref(null) // { commentId, replyId, author }

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

    const response = await axios.get(`http://localhost:3001/api/articles/${id}`, {
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
    const response = await axios.get(`http://localhost:3001/api/articles/${id}/like`, {
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
    if (!token) {
      router.push('/login')
      return
    }
    const response = await axios.post(`http://localhost:3001/api/articles/${id}/like`, {}, {
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

const fetchComments = async () => {
  const id = route.params.id
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get(`http://localhost:3001/api/articles/${id}/comments`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    comments.value = (response.data.comments || []).map(comment => ({
      ...comment,
      replies: comment.replies || comment.Replies || []
    }))
  } catch (err) {
    console.error('获取评论失败:', err)
  }
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    commentError.value = '评论内容不能为空'
    return
  }

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const id = route.params.id
  submittingComment.value = true
  commentError.value = ''

  try {
    const response = await axios.post(`http://localhost:3001/api/articles/${id}/comments`, {
      content: commentContent.value
    }, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    // 添加新评论到列表
    comments.value.push(response.data.comment)
    commentContent.value = ''
  } catch (err) {
    commentError.value = err.response?.data?.error || '发布评论失败'
  } finally {
    submittingComment.value = false
  }
}

const startReply = (commentId, replyId = null, author = '') => {
  replyingTo.value = { commentId, replyId, author }
  replyContent.value = ''
  replyError.value = ''
  // 滚动到回复输入框
  setTimeout(() => {
    const replyInput = document.querySelector('.reply-input')
    if (replyInput) {
      replyInput.focus()
    }
  }, 100)
}

const cancelReply = () => {
  replyingTo.value = null
  replyContent.value = ''
  replyError.value = ''
}

const submitReply = async () => {
  if (!replyContent.value.trim()) {
    replyError.value = '回复内容不能为空'
    return
  }

  if (!replyingTo.value) return

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const id = route.params.id
  const { commentId, replyId } = replyingTo.value
  submittingReply.value = true
  replyError.value = ''

  try {
    const response = await axios.post(`http://localhost:3001/api/articles/${id}/comments/${commentId}/replies`, {
      content: replyContent.value,
      parent_reply_id: replyId
    }, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    // 找到对应的评论并添加回复
    const comment = comments.value.find(c => c.id === commentId || c.ID === commentId)
    if (comment) {
      if (!comment.replies) comment.replies = []
      comment.replies.push(response.data.reply)
    }

    cancelReply()
  } catch (err) {
    replyError.value = err.response?.data?.error || '发布回复失败'
  } finally {
    submittingReply.value = false
  }
}

const deleteComment = async (commentId) => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const id = route.params.id
  try {
    await axios.delete(`http://localhost:3001/api/articles/${id}/comments/${commentId}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    comments.value = comments.value.filter(c => c.id !== commentId && c.ID !== commentId)
  } catch (err) {
    console.error('删除评论失败:', err)
    alert('删除评论失败')
  }
}

const deleteReply = async (commentId, replyId) => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  const id = route.params.id
  try {
    await axios.delete(`http://localhost:3001/api/articles/${id}/comments/${commentId}/replies/${replyId}`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    const comment = comments.value.find(c => c.id === commentId || c.ID === commentId)
    if (comment && comment.replies) {
      comment.replies = comment.replies.filter(r => r.id !== replyId && r.ID !== replyId)
    }
  } catch (err) {
    console.error('删除回复失败:', err)
    alert('删除回复失败')
  }
}

onMounted(() => {
  fetchArticle()
  fetchComments()
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
          <div v-for="action in article.fitness_actions" :key="action.id || action.ID" class="action-card">
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
              <div v-for="group in action.action_groups" :key="group.id || group.ID" class="group-row">
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

      <!-- 评论区 -->
      <div class="comments-section">
        <div class="section-header">
          <h2>💬 评论</h2>
          <span class="comment-count">{{ comments.length }} 条评论</span>
        </div>

        <!-- 评论输入框 -->
        <div class="comment-input-container">
          <textarea
            v-model="commentContent"
            placeholder="写下你的评论..."
            class="comment-input"
            rows="3"
          ></textarea>
          <div v-if="commentError" class="comment-error">{{ commentError }}</div>
          <button
            @click="submitComment"
            class="submit-comment-btn"
            :disabled="submittingComment || !commentContent.trim()"
          >
            {{ submittingComment ? '发布中...' : '发布评论' }}
          </button>
        </div>

        <!-- 评论列表 -->
        <div v-if="comments.length > 0" class="comments-list">
          <div v-for="comment in comments" :key="comment?.id || comment?.ID || Math.random()" class="comment-item">
            <div class="comment-header">
              <span class="comment-author">{{ (comment.user?.username) || '匿名用户' }}</span>
              <span class="comment-date">{{ formatDate(comment.created_at || comment.CreatedAt) }}</span>
            </div>
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-actions">
              <button 
                @click="startReply(comment.id || comment.ID)" 
                class="reply-button"
              >
                回复
              </button>
              <button 
                v-if="userStore.userId && String(comment.user_id) === String(userStore.userId)"
                @click="deleteComment(comment.id || comment.ID)" 
                class="delete-button"
              >
                删除
              </button>
            </div>
            <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
              <div v-for="reply in comment.replies" :key="reply.id || reply.ID" class="reply-item">
                <div class="reply-header">
                  <span class="reply-author">{{ (reply.user?.username) || '匿名用户' }}</span>
                  <span v-if="reply.parent_reply_id" class="reply-to">
                    回复 {{ (reply.parent_reply?.user?.username) || '用户' }}
                  </span>
                  <span class="reply-date">{{ formatDate(reply.created_at || reply.CreatedAt) }}</span>
                </div>
                <div class="reply-content">{{ reply.content }}</div>
                <div class="reply-actions">
                  <button 
                    @click="startReply(comment.id || comment.ID, reply.id || reply.ID, (reply.user?.username) || '用户')" 
                    class="reply-button"
                  >
                    回复
                  </button>
                  <button 
                    v-if="userStore.userId && String(reply.user_id) === String(userStore.userId)"
                    @click="deleteReply(comment.id || comment.ID, reply.id || reply.ID)" 
                    class="delete-button"
                  >
                    删除
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="empty-comments">
          <span class="empty-icon">💭</span>
          <p>还没有评论，快来发表你的看法吧！</p>
        </div>

        <!-- 回复输入框 -->
        <div v-if="replyingTo" class="reply-input-container">
          <div class="reply-input-header">
            <h4>{{ replyingTo.author ? `回复 ${replyingTo.author}` : '回复评论' }}</h4>
            <button @click="cancelReply" class="cancel-reply-btn">取消</button>
          </div>
          <textarea
            v-model="replyContent"
            placeholder="写下你的回复..."
            class="reply-input"
            rows="2"
          ></textarea>
          <div v-if="replyError" class="reply-error">{{ replyError }}</div>
          <button
            @click="submitReply"
            class="submit-reply-btn"
            :disabled="submittingReply || !replyContent.trim()"
          >
            {{ submittingReply ? '发布中...' : '发布回复' }}
          </button>
        </div>
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

.comments-section {
  padding: 28px 32px;
  border-top: 1px solid #f0f0f0;
}

.comment-count {
  background: linear-gradient(135deg, #ff6b6b15 0%, #ee5a5a15 100%);
  color: #ff6b6b;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.comment-input-container {
  margin: 24px 0;
  background: #fafafa;
  border: 2px solid #e8e8e8;
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
}

.comment-input-container:hover {
  border-color: #667eea;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.1);
}

.comment-input {
  width: 100%;
  border: none;
  background: transparent;
  resize: none;
  font-size: 15px;
  line-height: 1.6;
  color: #333;
  outline: none;
  margin-bottom: 12px;
}

.comment-input::placeholder {
  color: #999;
}

.comment-error {
  color: #ff6b6b;
  font-size: 13px;
  margin-bottom: 12px;
}

.submit-comment-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.submit-comment-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.submit-comment-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  background: #fafafa;
  border-radius: 16px;
  padding: 20px;
  border: 1px solid #e8e8e8;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.comment-author {
  font-weight: 600;
  color: #667eea;
  font-size: 14px;
}

.comment-date {
  font-size: 12px;
  color: #999;
}

.comment-content {
  font-size: 15px;
  line-height: 1.6;
  color: #444;
  margin-bottom: 12px;
  white-space: pre-wrap;
}

.replies-list {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.reply-item {
  font-size: 14px;
  line-height: 1.5;
  color: #666;
}

.reply-author {
  font-weight: 600;
  color: #764ba2;
  margin-right: 8px;
}

.empty-comments {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 16px;
}

@media (max-width: 768px) {
  .comments-section {
    padding: 20px;
  }

  .comment-input-container {
    padding: 16px;
  }

  .comment-item {
    padding: 16px;
  }
}

.comment-actions,
.reply-actions {
  margin-top: 8px;
  display: flex;
  align-items: center;
}

.reply-button {
  background: transparent;
  border: none;
  color: #667eea;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.reply-button:hover {
  background: rgba(102, 126, 234, 0.1);
}

.delete-button {
  background: transparent;
  border: none;
  color: #ff6b6b;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
  margin-left: 8px;
}

.delete-button:hover {
  background: rgba(255, 107, 107, 0.1);
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.reply-to {
  font-size: 12px;
  color: #999;
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 10px;
}

.reply-date {
  font-size: 11px;
  color: #999;
  margin-left: auto;
}

.reply-input-container {
  margin: 20px 0;
  background: #f8f9fa;
  border: 2px solid #e8e8e8;
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
}

.reply-input-container:hover {
  border-color: #667eea;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.1);
}

.reply-input-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.reply-input-header h4 {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.cancel-reply-btn {
  background: transparent;
  border: none;
  color: #999;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.cancel-reply-btn:hover {
  background: #f5f5f5;
  color: #666;
}

.reply-input {
  width: 100%;
  border: none;
  background: transparent;
  resize: none;
  font-size: 14px;
  line-height: 1.5;
  color: #333;
  outline: none;
  margin-bottom: 12px;
}

.reply-input::placeholder {
  color: #999;
}

.reply-error {
  color: #ff6b6b;
  font-size: 12px;
  margin-bottom: 12px;
}

.submit-reply-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.submit-reply-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.submit-reply-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
