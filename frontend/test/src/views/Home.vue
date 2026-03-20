<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

// 固定的训练动作选项
const actionOptions = [
  { label: '卧推', value: '卧推' },
  { label: '引体向上', value: '引体向上' },
  { label: '深蹲', value: '深蹲' },
  { label: '硬拉', value: '硬拉' },
  { label: '俯卧撑', value: '俯卧撑' },
  { label: '哑铃弯举', value: '哑铃弯举' },
  { label: '双杠臂屈伸', value: '双杠臂屈伸' },
  { label: '腹肌训练', value: '腹肌训练' }
]

const router = useRouter()
const articles = ref([])
const loading = ref(false)
const error = ref('')
const showCreateForm = ref(false)
const newArticle = reactive({
  title: '',
  content: '',
  preview: '',
  fitness_actions: []
})

const resetForm = () => {
  newArticle.title = ''
  newArticle.content = ''
  newArticle.preview = ''
  newArticle.fitness_actions = []
}

const addAction = () => {
  newArticle.fitness_actions.push({
    action_name: '',
    remark: '',
    action_groups: [
      { group_index: 1, weight: null, rep_num: null, remark: '' }
    ]
  })
}

const removeAction = (index) => {
  newArticle.fitness_actions.splice(index, 1)
}

const addGroup = (actionIndex) => {
  const groups = newArticle.fitness_actions[actionIndex].action_groups
  groups.push({
    group_index: groups.length + 1,
    weight: null,
    rep_num: null,
    remark: ''
  })
}

const removeGroup = (actionIndex, groupIndex) => {
  const groups = newArticle.fitness_actions[actionIndex].action_groups
  groups.splice(groupIndex, 1)
  groups.forEach((group, idx) => {
    group.group_index = idx + 1
  })
}

const fetchArticles = async () => {
  loading.value = true
  error.value = ''

  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get('http://localhost:3000/api/my/articles', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    articles.value = response.data
  } catch (err) {
    error.value = err.response?.data?.message || '获取记录失败'
  } finally {
    loading.value = false
  }
}

const createArticle = async () => {
  if (!newArticle.title) {
    error.value = '请填写标题'
    return
  }

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  loading.value = true
  error.value = ''

  try {
    const content = newArticle.content || '暂无训练心得'
    const preview = newArticle.preview || content.substring(0, 50)
    
    const fitnessActions = newArticle.fitness_actions
      .filter(a => a.action_name)
      .map(action => ({
        action_name: action.action_name,
        remark: action.remark,
        action_groups: action.action_groups.map(group => ({
          group_index: group.group_index,
          weight: group.weight || 0,
          rep_num: group.rep_num || 0,
          remark: group.remark
        }))
      }))
    
    await axios.post('http://localhost:3000/api/articles', {
      title: newArticle.title,
      content: content,
      preview: preview,
      fitness_actions: fitnessActions
    }, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    resetForm()
    showCreateForm.value = false
    await fetchArticles()
  } catch (err) {
    error.value = err.response?.data?.message || err.response?.data?.error || '发布失败'
  } finally {
    loading.value = false
  }
}

const navigateToArticle = (article) => {
  const id = article.id || article.ID
  router.push(`/article/${id}`)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchArticles()
})
</script>

<template>
  <div class="home-container">
    <div class="page-header">
      <div class="header-content">
        <h2>🏋️ 训练动态</h2>
        <p>记录每一次突破，见证每一份成长</p>
      </div>
      <div class="header-buttons">
        <router-link to="/stats" class="stats-btn">
          📊 统计分析
        </router-link>
        <button @click="showCreateForm = !showCreateForm; resetForm()" class="create-btn" :class="{ active: showCreateForm }">
          <span v-if="!showCreateForm">✨ 记录一下</span>
          <span v-else>✕ 取消</span>
        </button>
      </div>
    </div>

    <transition name="slide-fade">
      <div v-if="showCreateForm" class="create-form">
        <div class="form-header">
          <h3>📝 记录训练</h3>
          <p>分享你的训练成果</p>
        </div>
        
        <div class="form-body">
          <div class="form-group">
            <label>标题 <span class="required">*</span></label>
            <input 
              type="text" 
              v-model="newArticle.title" 
              placeholder="例如：今日胸肌训练 💪" 
              :disabled="loading"
            >
          </div>
          
          <div class="form-group">
            <label>训练心得</label>
            <textarea 
              v-model="newArticle.content" 
              placeholder="记录你的训练感受、进步或心得..." 
              rows="3" 
              :disabled="loading"
            ></textarea>
          </div>

          <div class="actions-section">
            <div class="actions-header">
              <label>🎯 训练动作</label>
              <button type="button" @click="addAction" class="add-action-btn">
                <span>+</span> 添加动作
              </button>
            </div>

            <div v-if="newArticle.fitness_actions.length === 0" class="no-actions">
              <span class="empty-icon">🏋️</span>
              <p>暂无动作，点击上方按钮添加</p>
            </div>

            <transition-group name="list" tag="div">
              <div v-for="(action, aIndex) in newArticle.fitness_actions" :key="aIndex" class="action-item">
                <div class="action-header">
                  <span class="action-badge">动作 {{ aIndex + 1 }}</span>
                  <button type="button" @click="removeAction(aIndex)" class="remove-action-btn">删除</button>
                </div>
                
                <div class="action-form-row">
                  <div class="form-group action-name">
                    <label>动作名称 <span class="required">*</span></label>
                    <select 
                      v-model="action.action_name" 
                      :disabled="loading"
                      class="action-select"
                    >
                      <option value="">请选择动作</option>
                      <option v-for="option in actionOptions" :key="option.value" :value="option.value">
                        {{ option.label }}
                      </option>
                    </select>
                  </div>
                  <div class="form-group action-remark">
                    <label>备注</label>
                    <input 
                      type="text" 
                      v-model="action.remark" 
                      placeholder="可选"
                      :disabled="loading"
                    >
                  </div>
                </div>

                <div class="groups-section">
                  <div class="groups-header">
                    <label>📊 组数详情</label>
                    <button type="button" @click="addGroup(aIndex)" class="add-group-btn">+ 添加组</button>
                  </div>

                  <div class="groups-table-wrapper">
                    <table class="groups-table">
                      <thead>
                        <tr>
                          <th>组数</th>
                          <th>重量(kg)</th>
                          <th>次数</th>
                          <th>备注</th>
                          <th></th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="(group, gIndex) in action.action_groups" :key="gIndex">
                          <td><span class="group-badge">第{{ group.group_index }}组</span></td>
                          <td>
                            <input 
                              type="number" 
                              v-model.number="group.weight" 
                              placeholder="重量"
                              min="0"
                              :disabled="loading"
                            >
                          </td>
                          <td>
                            <input 
                              type="number" 
                              v-model.number="group.rep_num" 
                              placeholder="次数"
                              min="0"
                              :disabled="loading"
                            >
                          </td>
                          <td>
                            <input 
                              type="text" 
                              v-model="group.remark" 
                              placeholder="备注"
                              :disabled="loading"
                            >
                          </td>
                          <td>
                            <button 
                              type="button" 
                              @click="removeGroup(aIndex, gIndex)" 
                              class="remove-group-btn"
                              :disabled="action.action_groups.length <= 1"
                            >
                              ✕
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </transition-group>
          </div>

          <div v-if="error" class="error-message">
            ⚠️ {{ error }}
          </div>
          
          <button @click="createArticle" class="submit-btn" :disabled="loading">
            <span v-if="!loading">🚀 发布</span>
            <span v-else>发布中...</span>
          </button>
        </div>
      </div>
    </transition>

    <div v-if="loading && !showCreateForm" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="error && !showCreateForm" class="error-state">
      <span class="error-icon">😢</span>
      <p>{{ error }}</p>
    </div>

    <div v-else-if="articles.length === 0" class="empty-state">
      <span class="empty-icon">🏋️</span>
      <h3>还没有训练记录</h3>
      <p>点击"记录一下"开始你的健身之旅</p>
    </div>

    <div v-else class="article-grid">
      <transition-group name="list">
        <div 
          v-for="article in articles" 
          :key="article.id || article.ID" 
          class="article-card"
          @click="navigateToArticle(article)"
        >
          <div class="card-header">
            <h3>{{ article.title }}</h3>
            <span class="card-date">{{ formatDate(article.created_at || article.CreatedAt) }}</span>
          </div>
          <p class="card-content">{{ article.content?.substring(0, 80) }}...</p>
          <div class="card-footer">
            <span class="card-action" v-if="article.fitness_actions?.length">
              🎯 {{ article.fitness_actions.length }} 个动作
            </span>
            <span class="card-likes">❤️ {{ article.likes || 0 }}</span>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  max-width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 20px 24px;
  border-radius: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.header-content h2 {
  font-size: 24px;
  color: #333;
  margin: 0 0 4px 0;
  font-weight: 700;
}

.header-content p {
  font-size: 14px;
  color: #888;
  margin: 0;
}

.header-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
}

.stats-btn {
  background: linear-gradient(135deg, #52c41a 0%, #38b000 100%);
  color: white;
  border: none;
  padding: 14px 28px;
  border-radius: 16px;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(82, 196, 26, 0.4);
  text-decoration: none;
  display: inline-block;
}

.stats-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(82, 196, 26, 0.5);
}

.create-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 28px;
  border-radius: 16px;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.create-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.5);
}

.create-btn.active {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a5a 100%);
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.4);
}

.create-form {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  margin-bottom: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.form-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 24px;
  color: white;
}

.form-header h3 {
  margin: 0 0 4px 0;
  font-size: 20px;
}

.form-header p {
  margin: 0;
  opacity: 0.9;
  font-size: 14px;
}

.form-body {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #444;
}

.required {
  color: #ff6b6b;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid #e8e8e8;
  border-radius: 12px;
  font-size: 15px;
  transition: all 0.3s ease;
  box-sizing: border-box;
  background: #fafafa;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.action-select {
  appearance: none;
  background-image: url("data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3E%3Cpath stroke='%23666' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3E%3C/svg%3E");
  background-position: right 12px center;
  background-repeat: no-repeat;
  background-size: 16px;
  padding-right: 40px;
}

.actions-section {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 2px dashed #e8e8e8;
}

.actions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.actions-header label {
  font-size: 16px;
  font-weight: 600;
  color: #444;
}

.add-action-btn {
  background: linear-gradient(135deg, #52c41a 0%, #38b000 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(82, 196, 26, 0.3);
}

.add-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(82, 196, 26, 0.4);
}

.no-actions {
  text-align: center;
  padding: 40px 20px;
  background: #fafafa;
  border: 2px dashed #e0e0e0;
  border-radius: 16px;
}

.no-actions .empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 12px;
}

.no-actions p {
  color: #888;
  margin: 0;
}

.action-item {
  background: #fafafa;
  border: 2px solid #e8e8e8;
  border-radius: 16px;
  padding: 20px;
  margin-bottom: 16px;
  transition: all 0.3s ease;
}

.action-item:hover {
  border-color: #667eea;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.1);
}

.action-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.action-badge {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.remove-action-btn {
  background: transparent;
  color: #ff6b6b;
  border: 2px solid #ff6b6b;
  padding: 6px 14px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.remove-action-btn:hover {
  background: #ff6b6b;
  color: white;
}

.action-form-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.action-name {
  flex: 2;
  margin-bottom: 0;
}

.action-remark {
  flex: 1;
  margin-bottom: 0;
}

.groups-section {
  margin-top: 12px;
}

.groups-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.groups-header label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.add-group-btn {
  background: #f0f0f0;
  color: #555;
  border: none;
  padding: 6px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.add-group-btn:hover {
  background: #e0e0e0;
}

.groups-table-wrapper {
  overflow-x: auto;
  border-radius: 12px;
  border: 2px solid #e8e8e8;
}

.groups-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.groups-table th,
.groups-table td {
  padding: 12px 14px;
  text-align: center;
}

.groups-table th {
  background: #f5f5f5;
  font-weight: 600;
  color: #555;
  border-bottom: 2px solid #e8e8e8;
}

.groups-table td {
  border-top: 1px solid #e8e8e8;
}

.group-badge {
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
  color: #667eea;
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
}

.groups-table input {
  width: 70px;
  padding: 8px 10px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  text-align: center;
  font-size: 14px;
  background: white;
}

.groups-table input:focus {
  outline: none;
  border-color: #667eea;
}

.remove-group-btn {
  background: transparent;
  color: #ccc;
  border: none;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.remove-group-btn:hover:not(:disabled) {
  background: #ff6b6b;
  color: white;
}

.remove-group-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.error-message {
  background: #fff2f0;
  border: 2px solid #ffccc7;
  color: #ff4d4f;
  padding: 14px 18px;
  border-radius: 12px;
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 500;
}

.submit-btn {
  width: 100%;
  background: linear-gradient(135deg, #52c41a 0%, #38b000 100%);
  color: white;
  border: none;
  padding: 16px;
  border-radius: 14px;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(82, 196, 26, 0.3);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(82, 196, 26, 0.4);
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  box-shadow: none;
}

.loading-state {
  text-align: center;
  padding: 60px 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
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
  border-radius: 20px;
}

.error-state .error-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
}

.empty-state .empty-icon {
  font-size: 80px;
  display: block;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 22px;
  color: #333;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #888;
  margin: 0;
}

.article-grid {
  display: grid;
  gap: 20px;
}

.article-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
  border: 2px solid transparent;
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  border-color: #667eea;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.card-header h3 {
  font-size: 18px;
  color: #333;
  margin: 0;
  font-weight: 600;
  flex: 1;
}

.card-date {
  font-size: 12px;
  color: #999;
  background: #f5f5f5;
  padding: 4px 10px;
  border-radius: 8px;
  white-space: nowrap;
  margin-left: 12px;
}

.card-content {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 16px 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.card-action {
  font-size: 13px;
  color: #667eea;
  font-weight: 500;
}

.card-likes {
  font-size: 13px;
  color: #ff6b6b;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-20px);
  opacity: 0;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .header-buttons {
    flex-direction: column;
    width: 100%;
  }

  .stats-btn,
  .create-btn {
    width: 100%;
    text-align: center;
  }

  .action-form-row {
    flex-direction: column;
    gap: 12px;
  }

  .form-body {
    padding: 16px;
  }

  .article-card {
    padding: 18px;
  }

  .card-header h3 {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .groups-table input {
    width: 50px;
    padding: 6px;
    font-size: 13px;
  }
}
</style>
