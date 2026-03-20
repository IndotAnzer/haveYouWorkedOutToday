<script setup>
import { ref, onMounted, computed, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import * as echarts from 'echarts'

const router = useRouter()
const chartRef = ref(null)
const chartInstance = ref(null)
const loadChartRef = ref(null)
const loadChartInstance = ref(null)
const loading = ref(false)
const error = ref('')
const statsData = ref([])
const heatmapData = ref([])
const dateRange = ref({
  start: '',
  end: ''
})
const selectedAction = ref('')

// 计算默认日期范围（过去7天）
const defaultDateRange = computed(() => {
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 6)
  
  const formatDate = (date) => {
    return date.toISOString().split('T')[0]
  }
  
  return {
    start: formatDate(start),
    end: formatDate(end)
  }
})

// 获取所有动作名称
const actionNames = computed(() => {
  if (statsData.value.length === 0) {
    return []
  }
  return [...new Set(statsData.value.flatMap(item => item.action.map(a => a.action_name)))]
})

// 初始化图表
const initChart = () => {
  console.log('initChart 被调用')
  console.log('chartRef.value:', chartRef.value)
  
  if (chartRef.value) {
    console.log('chartRef.value 存在，初始化图表')
    
    // 设置容器的宽度和高度
    chartRef.value.style.width = '100%'
    chartRef.value.style.height = '500px'
    
    // 初始化图表
    chartInstance.value = echarts.init(chartRef.value)
    console.log('chartInstance.value 初始化完成:', chartInstance.value)
    
    // 更新图表数据
    updateChart()
  } else {
    console.log('chartRef.value 不存在，无法初始化图表')
  }
}

// 更新图表数据
const updateChart = () => {
  console.log('updateChart 被调用')
  console.log('chartInstance.value:', chartInstance.value)
  console.log('chartRef.value:', chartRef.value)
  console.log('statsData.value.length:', statsData.value.length)
  
  if (!chartInstance.value || statsData.value.length === 0) {
    console.log('图表实例不存在或数据为空，跳过更新')
    return
  }
  
  // 提取日期列表（横轴）
  const dates = statsData.value.map(item => item.date)
  console.log('dates:', dates)
  
  // 提取所有动作名称（用于系列）
  const actionNames = [...new Set(statsData.value.flatMap(item => item.action.map(a => a.action_name)))]
  console.log('actionNames:', actionNames)
  
  // 构建系列数据（每个动作对应一个系列）
  const series = actionNames.map((actionName, index) => {
    console.log(`构建系列 ${index}: ${actionName}`)
    
    const data = dates.map((date, dateIndex) => {
      const dayData = statsData.value.find(item => item.date === date)
      console.log(`  日期 ${dateIndex}: ${date}, dayData:`, dayData)
      
      if (!dayData) {
        console.log(`    没有找到该日期的数据，返回0`)
        return 0
      }
      
      const actionData = dayData.action.find(a => a.action_name === actionName)
      console.log(`    actionData:`, actionData)
      
      const result = actionData ? actionData.num : 0
      console.log(`    结果: ${result}`)
      return result
    })
    
    console.log(`  系列 ${index} 的数据:`, data)
    
    return {
      name: actionName,
      type: 'bar', // 柱状图
      data: data,
      xAxisIndex: 0, // 指定使用第一个xAxis
      yAxisIndex: 0,  // 指定使用第一个yAxis
      label: {
        show: true,
        position: 'top',
        formatter: function(params) {
          return actionName
        },
        color: '#333',
        fontSize: 12,
        fontWeight: 'bold'
      }
    }
  })
  console.log('series:', series)
  
  // 配置 ECharts 实例
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    legend: {
      data: actionNames, // 动作名称作为图例
      top: 30
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates, // 日期作为横轴
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: {
      type: 'value',
      name: '组数', // 纵轴为组数
      minInterval: 1
    },
    series: series
  }
  
  console.log('option:', option)
  
  // 使用setTimeout确保DOM完全渲染后再设置option
  setTimeout(() => {
    chartInstance.value.setOption(option)
  }, 100)
}

// 处理窗口大小变化
const handleResize = () => {
  if (chartInstance.value) {
    chartInstance.value.resize()
  }
  if (loadChartInstance.value) {
    loadChartInstance.value.resize()
  }
}

// 初始化折线图（总重量）
const initLoadChart = () => {
  console.log('initLoadChart 被调用')
  console.log('loadChartRef.value:', loadChartRef.value)
  
  if (loadChartRef.value) {
    console.log('loadChartRef.value 存在，初始化图表')
    
    // 设置容器的宽度和高度
    loadChartRef.value.style.width = '100%'
    loadChartRef.value.style.height = '400px'
    
    // 初始化图表
    loadChartInstance.value = echarts.init(loadChartRef.value)
    console.log('loadChartInstance.value 初始化完成:', loadChartInstance.value)
    
    // 更新图表数据
    updateLoadChart()
  } else {
    console.log('loadChartRef.value 不存在，无法初始化图表')
  }
}

// 更新折线图数据
const updateLoadChart = () => {
  console.log('updateLoadChart 被调用')
  console.log('loadChartInstance.value:', loadChartInstance.value)
  console.log('loadChartRef.value:', loadChartRef.value)
  console.log('statsData.value.length:', statsData.value.length)
  console.log('selectedAction.value:', selectedAction.value)
  
  if (!loadChartRef.value || !loadChartInstance.value || statsData.value.length === 0 || !selectedAction.value) {
    console.log('图表实例不存在、数据为空或未选择动作，跳过更新')
    return
  }
  
  // 提取日期列表（横轴）
  const dates = statsData.value.map(item => item.date)
  console.log('dates:', dates)
  
  // 构建折线图数据
  const data = dates.map((date) => {
    const dayData = statsData.value.find(item => item.date === date)
    if (!dayData) {
      return 0
    }
    
    const actionData = dayData.action.find(a => a.action_name === selectedAction.value)
    return actionData ? actionData.load : 0
  })
  console.log('折线图数据:', data)
  
  // 配置 ECharts 实例
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: {
      type: 'value',
      name: '总重量 (kg)'
    },
    series: [{
      name: selectedAction.value,
      type: 'line',
      data: data,
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      lineStyle: {
        width: 3,
        color: '#667eea'
      },
      itemStyle: {
        color: '#667eea'
      },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(102, 126, 234, 0.3)' },
          { offset: 1, color: 'rgba(102, 126, 234, 0.1)' }
        ])
      }
    }]
  }
  
  console.log('折线图配置:', option)
  
  // 使用setTimeout确保DOM完全渲染后再设置option
  setTimeout(() => {
    loadChartInstance.value.setOption(option)
  }, 100)
}

// 监听动作选择变化
watch(selectedAction, () => {
  updateLoadChart()
})

const formatDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const generateHeatmapData = () => {
  const data = []
  
  if (!dateRange.value.start || !dateRange.value.end) {
    return data
  }
  
  const startDate = new Date(dateRange.value.start)
  const endDate = new Date(dateRange.value.end)

  console.log('生成热图数据，heatmapData.value:', heatmapData.value)
  console.log('开始日期:', formatDate(startDate), '结束日期:', formatDate(endDate))

  const currentDate = new Date(startDate)
  while (currentDate <= endDate) {
    const dateStr = formatDate(currentDate)
    
    // 处理后端返回的日期格式，提取日期部分
    const trainingData = heatmapData.value.find(item => {
      const backendDate = item.date.split('T')[0]
      return backendDate === dateStr
    })
    console.log(`检查日期 ${dateStr}, 找到数据:`, trainingData)
    
    const hasTraining = trainingData ? trainingData.intensity > 0 : false
    
    data.push({
      date: dateStr,
      hasTraining: hasTraining
    })
    
    currentDate.setDate(currentDate.getDate() + 1)
  }

  console.log('生成的热图数据:', data)
  return data
}

// 获取训练量统计数据
const fetchStats = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }
    
    const params = new URLSearchParams()
    if (dateRange.value.start) params.append('startDate', dateRange.value.start)
    if (dateRange.value.end) params.append('endDate', dateRange.value.end)
    
    const [volumeResponse, frequencyResponse] = await Promise.all([
      axios.get(`http://localhost:3000/api/articles/stats/volume?${params.toString()}`, {
        headers: { Authorization: `Bearer ${token}` }
      }),
      axios.get(`http://localhost:3000/api/articles/stats/frequency?${params.toString()}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
    ])
    
    console.log('后端返回的原始数据:', volumeResponse.data)
    statsData.value = volumeResponse.data
    console.log('statsData.value:', statsData.value)
    console.log('statsData.value.length:', statsData.value.length)
    
    heatmapData.value = frequencyResponse.data.map(item => ({
      date: item.date,
      intensity: parseInt(item.count)
    }))
    console.log('热图数据:', heatmapData.value)
    console.log('热图数据详情:', JSON.stringify(heatmapData.value, null, 2))
    console.log('热图数据数量:', heatmapData.value.length)
    if (heatmapData.value.length > 0) {
      console.log('第一个数据项:', heatmapData.value[0])
      console.log('第一个数据项的date类型:', typeof heatmapData.value[0].date)
      console.log('第一个数据项的date值:', heatmapData.value[0].date)
    }
    
    // 使用nextTick确保DOM已经更新后再初始化图表
    await nextTick()
    await nextTick()
    await nextTick()
    
    // 如果图表实例不存在，初始化图表
    if (!chartInstance.value && statsData.value.length > 0) {
      console.log('初始化图表实例')
      // 再次等待DOM渲染完成
      await nextTick()
      console.log('chartRef.value:', chartRef.value)
      initChart()
    } else {
      console.log('更新图表数据')
      updateChart()
    }
    
    // 初始化折线图
    if (statsData.value.length > 0) {
      await nextTick()
      initLoadChart()
    }
  } catch (err) {
    error.value = err.response?.data?.error || '获取统计数据失败'
  } finally {
    loading.value = false
  }
}

// 重置日期范围
const resetDateRange = () => {
  dateRange.value = defaultDateRange.value
  fetchStats()
}

// 监听日期变化
const handleDateChange = () => {
  fetchStats()
}

onMounted(() => {
  // 设置默认日期范围
  dateRange.value = defaultDateRange.value
  
  // 获取统计数据
  fetchStats()
  
  // 添加窗口大小变化监听
  window.addEventListener('resize', handleResize)
  
  // 清理函数
  return () => {
    window.removeEventListener('resize', handleResize)
    if (chartInstance.value) {
      chartInstance.value.dispose()
    }
    if (loadChartInstance.value) {
      loadChartInstance.value.dispose()
    }
  }
})
</script>

<template>
  <div class="stats-container">
    <div class="page-header">
      <div class="header-content">
        <h2>📊 训练量统计</h2>
        <p>查看你的训练进度和趋势</p>
      </div>
      <button @click="router.push('/')" class="back-btn">
        <span>←</span> 返回首页
      </button>
    </div>
    
    <div class="stats-content">
      <!-- 日期选择器 -->
      <div class="date-range-selector">
        <div class="date-group">
          <label>开始日期</label>
          <input 
            type="date" 
            v-model="dateRange.start" 
            @change="handleDateChange"
            :disabled="loading"
          >
        </div>
        <div class="date-group">
          <label>结束日期</label>
          <input 
            type="date" 
            v-model="dateRange.end" 
            @change="handleDateChange"
            :disabled="loading"
          >
        </div>
        <button @click="resetDateRange" class="reset-btn" :disabled="loading">
          重置
        </button>
      </div>
      
      <!-- 错误信息 -->
      <div v-if="error" class="error-message">
        ⚠️ {{ error }}
      </div>
      
      <!-- 加载状态 -->
      <div v-show="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>
      
      <!-- 图表容器 -->
      <div v-show="!loading" class="chart-container">
        <div v-show="statsData.length === 0" class="empty-state">
          <span class="empty-icon">📊</span>
          <p>该时间段内暂无训练记录</p>
        </div>
        <div v-show="statsData.length > 0" ref="chartRef" class="chart"></div>
      </div>
      
      <!-- 热图容器 -->
      <div v-show="!loading" class="heatmap-section">
        <h3>🔥 训练热图</h3>
        <p>查看选定日期范围内的训练情况</p>
        <div class="heatmap-grid">
          <div 
            v-for="day in generateHeatmapData()" 
            :key="day.date" 
            class="heatmap-cell"
            :class="{ 'has-training': day.hasTraining }"
          >
            <div class="cell-date">{{ day.date.slice(5) }}</div>
            <div v-if="day.hasTraining" class="check-mark">✓</div>
          </div>
        </div>
      </div>
      
      <!-- 总重量折线图 -->
      <div v-show="!loading" class="load-chart-section">
        <h3>📈 重量增长趋势</h3>
        <p>选择动作查看总重量的增长情况</p>
        
        <!-- 动作选择器 -->
        <div class="action-selector">
          <label>选择动作</label>
          <select 
            v-model="selectedAction" 
            :disabled="loading || statsData.length === 0"
          >
            <option value="">-- 请选择动作 --</option>
            <option 
              v-for="actionName in actionNames" 
              :key="actionName" 
              :value="actionName"
            >
              {{ actionName }}
            </option>
          </select>
        </div>
        
        <!-- 折线图容器 -->
        <div v-show="selectedAction" ref="loadChartRef" class="load-chart"></div>
        <div v-show="!selectedAction" class="empty-action">
          <span class="empty-icon">🏋️</span>
          <p>请选择一个动作查看重量增长趋势</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-container {
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

.back-btn {
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
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.08);
}

.back-btn:hover {
  background: white;
  transform: translateX(-4px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

.stats-content {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 24px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.date-range-selector {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.date-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.date-group label {
  font-size: 14px;
  font-weight: 600;
  color: #444;
}

.date-group input {
  padding: 12px 16px;
  border: 2px solid #e8e8e8;
  border-radius: 12px;
  font-size: 15px;
  transition: all 0.3s ease;
  background: #fafafa;
}

.date-group input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.reset-btn {
  align-self: flex-end;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.reset-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.reset-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

.loading-state {
  text-align: center;
  padding: 60px 20px;
  background: #fafafa;
  border-radius: 16px;
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

.chart-container {
  background: #fafafa;
  border-radius: 16px;
  padding: 20px;
  border: 2px solid #e8e8e8;
  min-height: 550px;
}

.chart {
  width: 100%;
  height: 500px;
  min-width: 300px;
  min-height: 300px;
  display: block;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #888;
}

.empty-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

.heatmap-section {
  margin-top: 32px;
  background: #fafafa;
  border-radius: 16px;
  padding: 24px;
  border: 2px solid #e8e8e8;
}

.heatmap-section h3 {
  font-size: 20px;
  color: #333;
  margin: 0 0 8px 0;
  font-weight: 700;
}

.heatmap-section p {
  font-size: 14px;
  color: #888;
  margin: 0 0 20px 0;
}

.heatmap-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 12px;
  margin-top: 20px;
}

.heatmap-cell {
  background: #f5f5f5;
  border-radius: 12px;
  padding: 20px 12px;
  text-align: center;
  border: 2px solid #e8e8e8;
  transition: all 0.3s ease;
  min-height: 80px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 8px;
}

.heatmap-cell.has-training {
  background: linear-gradient(135deg, #9be9a8 0%, #40c463 100%);
  border-color: #30a14e;
  box-shadow: 0 4px 15px rgba(64, 196, 99, 0.3);
}

.cell-date {
  font-size: 14px;
  font-weight: 600;
  color: #666;
}

.heatmap-cell.has-training .cell-date {
  color: #216e39;
}

.check-mark {
  font-size: 32px;
  font-weight: bold;
  color: #216e39;
  animation: checkPop 0.3s ease;
}

@keyframes checkPop {
  0% { transform: scale(0); opacity: 0; }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); opacity: 1; }
}

/* 总重量折线图部分 */
.load-chart-section {
  margin-top: 32px;
  background: #fafafa;
  border-radius: 16px;
  padding: 24px;
  border: 2px solid #e8e8e8;
}

.load-chart-section h3 {
  font-size: 20px;
  color: #333;
  margin: 0 0 8px 0;
  font-weight: 700;
}

.load-chart-section p {
  font-size: 14px;
  color: #888;
  margin: 0 0 20px 0;
}

.action-selector {
  margin-bottom: 20px;
}

.action-selector label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #444;
  margin-bottom: 8px;
}

.action-selector select {
  width: 100%;
  max-width: 300px;
  padding: 12px 16px;
  border: 2px solid #e8e8e8;
  border-radius: 12px;
  font-size: 15px;
  background: #fafafa;
  transition: all 0.3s ease;
}

.action-selector select:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.action-selector select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.load-chart {
  width: 100%;
  height: 400px;
  min-width: 300px;
  min-height: 300px;
  display: block;
  background: white;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #e8e8e8;
}

.empty-action {
  text-align: center;
  padding: 80px 20px;
  color: #888;
  background: white;
  border-radius: 12px;
  border: 1px solid #e8e8e8;
}

.empty-action .empty-icon {
  font-size: 64px;
  display: block;
  margin-bottom: 16px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .date-range-selector {
    flex-direction: column;
  }
  
  .reset-btn {
    align-self: flex-start;
  }
  
  .chart {
    height: 400px;
  }
  
  .stats-content {
    padding: 16px;
  }
}
</style>