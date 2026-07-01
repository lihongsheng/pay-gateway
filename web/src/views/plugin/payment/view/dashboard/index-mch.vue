<!-- views/dashboard/system/index-mch.vue -->
<template>
  <div class="system-dashboard">
    <!-- 全局时间选择器 - 放在最上方 -->
    <el-card shadow="never" class="date-range-card">
      <div class="global-date-range">
        <span class="date-label">统计时间：</span>
        <el-date-picker
          v-model="globalDateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          size="default"
          :shortcuts="shortcuts"
          @change="handleGlobalDateChange"
          class="date-picker"
        />
        <el-button type="primary" @click="refreshAllData" :loading="refreshing">
          <el-icon><Refresh /></el-icon>刷新数据
        </el-button>
      </div>
    </el-card>

    <!-- 头部统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6" v-for="stat in statistics" :key="stat.label">
        <el-card shadow="never" class="stat-card">
          <div class="stat-content">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-compare">
              <span :class="stat.trend > 0 ? 'up' : 'down'">
                <el-icon>
                  <ArrowUp v-if="stat.trend > 0" />
                  <ArrowDown v-else />
                </el-icon>
                {{ Math.abs(stat.trend) }}%
              </span>
              <span class="compare-text">较昨日</span>
            </div>
          </div>
          <div class="stat-icon" :style="{ backgroundColor: stat.color + '20', color: stat.color }">
            <el-icon :size="24"><component :is="stat.icon" /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 - 第一行 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">近7日成功支付金额趋势</span>
              <span class="date-display">{{ formatDateRange }}</span>
            </div>
          </template>
          <div class="chart-container">
            <v-chart class="chart" :option="amountChartOption" autoresize />
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">近7日成功支付笔数趋势</span>
              <span class="date-display">{{ formatDateRange }}</span>
            </div>
          </template>
          <div class="chart-container">
            <v-chart class="chart" :option="countChartOption" autoresize />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 - 第二行 (应用交易数据) -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <el-card shadow="never" class="chart-card">
          <template #header>
            <div class="card-header">
              <span class="title">应用交易数据</span>
              <div class="header-controls">
                <!-- 应用下拉框（远程搜索） -->
                <el-select
                  v-model="selectedApp"
                  placeholder="选择应用"
                  filterable
                  remote
                  :remote-method="searchAppRemote"
                  :loading="appSearchLoading"
                  clearable
                  size="small"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="item in filteredAppOptions"
                    :key="item.app_no"
                    :label="item.app_name"
                    :value="item.app_no"
                  />
                </el-select>
                <el-date-picker
                  v-model="appDateRange"
                  type="daterange"
                  range-separator="至"
                  start-placeholder="开始日期"
                  end-placeholder="结束日期"
                  size="small"
                  style="width: 240px; margin-right: 10px"
                />
                <el-button type="primary" size="small" @click="searchAppData" :loading="appLoading">
                  <el-icon><Search /></el-icon>查询
                </el-button>
                <el-radio-group v-model="appChartType" size="small" style="margin-left: 10px">
                  <el-radio-button label="amount">支付金额</el-radio-button>
                  <el-radio-button label="count">支付笔数</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </template>
          <div class="chart-container">
            <v-chart class="chart" :option="appChartOption" autoresize />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 数据表格 - 带高级搜索 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span class="title">交易明细</span>
          <el-button type="primary" link @click="showSearch = !showSearch">
            {{ showSearch ? '收起搜索' : '高级搜索' }}
            <el-icon><ArrowUp v-if="showSearch" /><ArrowDown v-else /></el-icon>
          </el-button>
        </div>
      </template>

      <!-- 高级搜索栏 -->
      <div v-if="showSearch" class="search-bar">
        <el-form :inline="true" :model="tableSearch" class="search-form">
          <!-- 应用下拉框（远程搜索） -->
          <el-form-item label="应用">
            <el-select
              v-model="tableSearch.appNo"
              placeholder="请选择应用"
              filterable
              remote
              :remote-method="searchAppRemote"
              :loading="appSearchLoading"
              clearable
              style="width: 200px"
              @change="handleAppChange"
            >
              <el-option
                v-for="item in filteredAppOptions"
                :key="item.app_no"
                :label="item.app_name"
                :value="item.app_no"
              />
            </el-select>
          </el-form-item>

          <!-- 支付账户下拉框（静态加载，label 使用 name） -->
          <el-form-item label="支付账户">
            <el-select
              v-model="tableSearch.accountNo"
              placeholder="请选择支付账户"
              filterable
              clearable
              style="width: 200px"
              :disabled="!tableSearch.appNo"
            >
              <el-option
                v-for="item in accountOptions"
                :key="item.account_no"
                :label="item.name"
                :value="item.account_no"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="日期范围">
            <el-date-picker
              v-model="tableSearch.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 240px"
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="searchTableData" :loading="tableLoading">
              <el-icon><Search /></el-icon>查询
            </el-button>
            <el-button @click="resetTableSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table :data="tableData" stripe style="width: 100%" v-loading="tableLoading">
        <el-table-column prop="statistic_date" label="日期" width="120" />
        <el-table-column prop="app_name" label="应用名称" min-width="150">
          <template #default="{ row }">
            <el-tag size="small" type="info" v-if="row.app_name">{{ row.app_name }}</el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="account_name" label="支付账户" min-width="150">
          <template #default="{ row }">
            <span v-if="row.account_name">{{ row.account_name }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_order" label="总订单数" width="120" align="right">
          <template #default="{ row }">
            {{ formatNumber(row.total_order) }}
          </template>
        </el-table-column>
        <el-table-column prop="success_order" label="成功笔数" width="120" align="right">
          <template #default="{ row }">
            <span class="success-text">{{ formatNumber(row.success_order) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="success_rate" label="成功率" width="100" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="parseFloat(row.success_rate)"
              :color="getRateColor(row.success_rate)"
              :stroke-width="8"
              :show-text="false"
              style="width: 60px; display: inline-block"
            />
            <span class="rate-text">{{ row.success_rate }}%</span>
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="总金额(元)" width="130" align="right">
          <template #default="{ row }">
            ¥ {{ formatAmount(row.total_amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="success_amount" label="成功金额(元)" width="130" align="right">
          <template #default="{ row }">
            <span class="success-text">¥ {{ formatAmount(row.success_amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="refund_amount" label="退款金额(元)" width="130" align="right">
          <template #default="{ row }">
            <span class="refund-text">¥ {{ formatAmount(row.refund_amount) }}</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowUp, ArrowDown, Refresh, Search } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent
} from 'echarts/components'
import { dashboardMchAppIndex, dashboardMchIndex, dashboardTotalRequest, dashboardSearch, appList, accountList } from '@/api/payment'

// 注册 ECharts 组件
use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent
])

// 全局日期范围
const globalDateRange = ref(getLast7Days())
const refreshing = ref(false)
const rawData = ref([]) // 存储原始数据
const showSearch = ref(false) // 是否显示高级搜索

// 统计数据
const statistics = ref([
  {
    label: '总请求单数',
    value: '0',
    icon: 'Document',
    color: '#E6A23C',
    trend: 0
  },
  {
    label: '渠道支付笔数',
    value: '0',
    icon: 'Document',
    color: '#E6A23C',
    trend: 0
  },
  {
    label: '支付成功金额',
    value: '¥ 0',
    icon: 'SuccessFilled',
    color: '#67C23A',
    trend: 0
  },
  {
    label: '渠道支付成功率',
    value: '0%',
    icon: 'DataLine',
    color: '#F56C6C',
    trend: 0
  }
])

// 应用选项 - 用于应用交易数据图表和交易明细
const appOptions = ref([])           // 全部应用（初始加载）
const filteredAppOptions = ref([])   // 实际显示的下拉选项（根据远程搜索动态变化）
const appSearchLoading = ref(false)  // 远程搜索加载状态

// 应用搜索条件 - 用于应用交易数据图表
const selectedApp = ref('')
const appDateRange = ref(getLast7Days())
const appLoading = ref(false)
const appRawData = ref([]) // 存储应用原始数据
const appChartType = ref('amount') // amount 或 count

// 表格搜索条件 - 用于交易明细
const tableSearch = ref({
  appNo: '',
  accountNo: '',
  dateRange: getLast7Days()
})

// 账户选项 - 用于交易明细的账户下拉框
const accountOptions = ref([])
const accountLoading = ref(false)

// 表格数据
const tableData = ref([])
const tableLoading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 格式化日期范围显示
const formatDateRange = computed(() => {
  if (!globalDateRange.value || !globalDateRange.value[0] || !globalDateRange.value[1]) return ''
  const start = formatDate(globalDateRange.value[0])
  const end = formatDate(globalDateRange.value[1])
  return `${start} 至 ${end}`
})

// 快捷日期选项
const shortcuts = [
  { text: '最近7天', value: () => getLast7Days() },
  { text: '最近30天', value: () => getLast30Days() },
  { text: '本月', value: () => [new Date(new Date().getFullYear(), new Date().getMonth(), 1), new Date()] },
  { text: '上月', value: () => [new Date(new Date().getFullYear(), new Date().getMonth() - 1, 1), new Date(new Date().getFullYear(), new Date().getMonth(), 0)] }
]

// 图表选项（保持不变）
const amountChartOption = ref({
  tooltip: {
    trigger: 'axis',
    formatter: function(params) {
      if (!params || params.length === 0) return ''
      const value = params[0].value || 0
      return `${params[0].name}<br/>成功金额: ¥ ${formatAmount(value)}`
    }
  },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: [],
    axisLabel: { rotate: 0 }
  },
  yAxis: {
    type: 'value',
    name: '金额(元)',
    axisLabel: {
      formatter: (value) => '¥ ' + (value / 100).toFixed(2)
    },
    min: 0
  },
  series: [
    {
      name: '成功金额',
      type: 'line',
      data: [],
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      lineStyle: { width: 3, color: '#409EFF' },
      areaStyle: { color: 'rgba(64, 158, 255, 0.1)' },
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      }
    }
  ],
  dataZoom: [{ type: 'inside' }, { type: 'slider' }]
})

const countChartOption = ref({
  tooltip: {
    trigger: 'axis',
    formatter: function(params) {
      if (!params || params.length === 0) return ''
      const value = params[0].value || 0
      return `${params[0].name}<br/>成功笔数: ${value} 笔`
    }
  },
  grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
  xAxis: {
    type: 'category',
    data: []
  },
  yAxis: {
    type: 'value',
    name: '笔数',
    min: 0
  },
  series: [
    {
      name: '成功笔数',
      type: 'line',
      data: [],
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      lineStyle: { width: 3, color: '#67C23A' },
      areaStyle: { color: 'rgba(103, 194, 58, 0.1)' },
      markPoint: {
        data: [
          { type: 'max', name: '最大值' },
          { type: 'min', name: '最小值' }
        ]
      }
    }
  ],
  dataZoom: [{ type: 'inside' }, { type: 'slider' }]
})

const appChartOption = ref({
  tooltip: {
    trigger: 'axis',
    formatter: function(params) {
      if (!params || params.length === 0) return ''
      let result = params[0].name + '<br/>'
      params.forEach(p => {
        const value = appChartType.value === 'amount'
          ? `¥ ${formatAmount(p.value)}`
          : `${p.value} 笔`
        result += `${p.marker} ${p.seriesName}: ${value}<br/>`
      })
      return result
    }
  },
  legend: {
    data: [],
    bottom: 0,
    type: 'scroll',
    pageIconColor: '#409EFF',
    pageTextStyle: { color: '#606266' }
  },
  grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true },
  xAxis: {
    type: 'category',
    data: [],
    axisLabel: { rotate: 30 }
  },
  yAxis: {
    type: 'value',
    name: computed(() => appChartType.value === 'amount' ? '金额(元)' : '笔数'),
    axisLabel: {
      formatter: (value) => appChartType.value === 'amount'
        ? '¥ ' + (value / 100).toFixed(2)
        : value
    },
    min: 0
  },
  series: [],
  dataZoom: [{ type: 'inside' }, { type: 'slider' }]
})

// 工具函数
function getLast7Days() {
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 6)
  // 设置时间为当天的开始和结束
  start.setHours(0, 0, 0, 0)
  end.setHours(23, 59, 59, 999)
  return [start, end]
}

function getLast30Days() {
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 29)
  start.setHours(0, 0, 0, 0)
  end.setHours(23, 59, 59, 999)
  return [start, end]
}

function formatDate(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatDateTime(date) {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`
}

const formatNumber = (num) => {
  if (!num && num !== 0) return '0'
  return num.toLocaleString()
}

const formatAmount = (amount) => {
  if (!amount && amount !== 0) return '0.00'
  return (amount / 100).toFixed(2)
}

const getRateColor = (rate) => {
  const numRate = parseFloat(rate)
  if (numRate >= 99) return '#67C23A'
  if (numRate >= 95) return '#409EFF'
  if (numRate >= 90) return '#E6A23C'
  return '#F56C6C'
}

// 应用远程搜索方法
const searchAppRemote = async (query) => {
  if (query !== '') {
    appSearchLoading.value = true
    try {
      const res = await appList({
        app_name: query,
        page: 1,
        page_size: 20
      })
      if (res.code === 0) {
        filteredAppOptions.value = res.data.list || []
      } else {
        filteredAppOptions.value = []
      }
    } catch (error) {
      console.error('搜索应用失败:', error)
      ElMessage.error('搜索应用失败')
      filteredAppOptions.value = []
    } finally {
      appSearchLoading.value = false
    }
  } else {
    // 输入为空时，显示全部应用
    filteredAppOptions.value = appOptions.value
  }
}

// 获取全部应用列表（初始化时调用）
const fetchAppOptions = async () => {
  try {
    const res = await appList({ page: 1, page_size: 100 })
    if (res.code === 0) {
      appOptions.value = res.data.list || []
      filteredAppOptions.value = appOptions.value

      // 默认选中第一个应用用于应用交易数据图表
      if (appOptions.value.length > 0 && !selectedApp.value) {
        selectedApp.value = appOptions.value[0].app_no
        // 自动加载应用数据
        searchAppData()
      }
    }
  } catch (error) {
    console.error('获取应用列表失败:', error)
    ElMessage.error('获取应用列表失败')
  }
}

// 处理应用变化 - 加载账户列表（用于交易明细）
const handleAppChange = async (appNo) => {
  tableSearch.value.accountNo = ''

  if (!appNo) {
    accountOptions.value = []
    return
  }

  accountLoading.value = true
  try {
    const res = await accountList({
      app_no: appNo,
      page: 1,
      page_size: 100
    })
    if (res.code === 0) {
      accountOptions.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取账户列表失败:', error)
    ElMessage.error('获取账户列表失败')
  } finally {
    accountLoading.value = false
  }
}

// 处理全局日期变化
const handleGlobalDateChange = () => {
  refreshAllData()
}

// 刷新所有数据
const refreshAllData = async () => {
  refreshing.value = true
  try {
    await fetchSystemData()
    // 加载应用数据（如果有选中的应用）
    if (selectedApp.value) {
      await searchAppData()
    }
    await searchTableData() // 刷新表格数据
  //  ElMessage.success('数据刷新成功')
  } catch (error) {
    console.error('数据刷新失败:', error)
    ElMessage.error('数据刷新失败')
  } finally {
    refreshing.value = false
  }
}

// 查询应用数据 - 用于应用交易数据图表
const searchAppData = async () => {
  if (!selectedApp.value) {
    appRawData.value = []
    updateAppChart()
    return
  }

  appLoading.value = true
  try {
    const params = {
      app_no: selectedApp.value,
      start_time: formatDateTime(appDateRange.value[0]),
      end_time: formatDateTime(appDateRange.value[1])
    }

    const res = await dashboardMchAppIndex(params)
    if (res.code === 0) {
      appRawData.value = Array.isArray(res.data) ? res.data : (res.data.list || [])
      updateAppChart()
    }
  } catch (error) {
    console.error('获取应用数据失败:', error)
    ElMessage.error('获取应用数据失败')
  } finally {
    appLoading.value = false
  }
}

// 查询表格数据 - 用于交易明细
const searchTableData = async () => {
  tableLoading.value = true
  try {
    const params = {
      app_no: tableSearch.value.appNo || undefined,
      account_no: tableSearch.value.accountNo || undefined,
      start_time: tableSearch.value.dateRange ? formatDateTime(tableSearch.value.dateRange[0]) : formatDateTime(globalDateRange.value[0]),
      end_time: tableSearch.value.dateRange ? formatDateTime(tableSearch.value.dateRange[1]) : formatDateTime(globalDateRange.value[1]),
      page: currentPage.value,
      page_size: pageSize.value
    }

    const res = await dashboardSearch(params)
    if (res.code === 0) {
      const list = Array.isArray(res.data) ? res.data : (res.data.list || [])

      tableData.value = list.map(item => ({
        ...item,
        success_rate: item.total_order
          ? ((item.success_order / item.total_order) * 100).toFixed(2)
          : 0
      }))
      total.value = list.length
    }
  } catch (error) {
    console.error('获取表格数据失败:', error)
    ElMessage.error('获取表格数据失败')
  } finally {
    tableLoading.value = false
  }
}

// 重置表格搜索
const resetTableSearch = () => {
  tableSearch.value = {
    appNo: '',
    accountNo: '',
    dateRange: getLast7Days()
  }
  accountOptions.value = []
  currentPage.value = 1
  searchTableData()
}

// 监听应用图表类型变化（不重新请求接口）
watch(appChartType, () => {
  if (appRawData.value.length > 0) {
    updateAppChart()
  }
})

// 获取系统数据
const fetchSystemData = async () => {
  try {
    const params = {
      start_time: formatDateTime(globalDateRange.value[0]),
      end_time: formatDateTime(globalDateRange.value[1])
    }
    const res = await dashboardMchIndex(params)
    const totalRequestRes = await dashboardTotalRequest(params)
    if (res.code === 0) {
      rawData.value = Array.isArray(res.data) ? res.data : (res.data.list || [])
      updateSystemCharts(rawData.value)
      updateStatistics(rawData.value, totalRequestRes?.data || 0)
    }
  } catch (error) {
    console.error('获取系统数据失败:', error)
    throw error
  }
}

// 更新系统图表
const updateSystemCharts = (data) => {
  // 与原来相同
  if (!data || data.length === 0) {
    amountChartOption.value.xAxis.data = []
    amountChartOption.value.series[0].data = []
    countChartOption.value.xAxis.data = []
    countChartOption.value.series[0].data = []
    return
  }

  const dateMap = new Map()
  data.forEach(item => {
    const date = item.statistic_date
    if (!dateMap.has(date)) {
      dateMap.set(date, {
        success_amount: 0,
        success_order: 0
      })
    }
    const dayData = dateMap.get(date)
    dayData.success_amount += item.success_amount || 0
    dayData.success_order += item.success_order || 0
  })

  const sortedDates = Array.from(dateMap.keys()).sort()
  const amountData = sortedDates.map(date => dateMap.get(date).success_amount)
  const countData = sortedDates.map(date => dateMap.get(date).success_order)

  amountChartOption.value.xAxis.data = sortedDates
  amountChartOption.value.series[0].data = amountData
  countChartOption.value.xAxis.data = sortedDates
  countChartOption.value.series[0].data = countData
}

// 更新应用图表（按账户分组）
const updateAppChart = () => {
  const data = appRawData.value
  if (!data || data.length === 0) {
    appChartOption.value.xAxis.data = []
    appChartOption.value.series = []
    appChartOption.value.legend.data = []
    return
  }

  // 按账户分组处理数据
  const accountMap = new Map()
  const dateSet = new Set()

  data.forEach(item => {
    const date = item.statistic_date
    const accountNo = item.account_no || 'default'
    const accountName = item.account_name || accountNo

    dateSet.add(date)

    if (!accountMap.has(accountNo)) {
      accountMap.set(accountNo, {
        name: accountName,
        data: new Map()
      })
    }

    const accountData = accountMap.get(accountNo)
    accountData.data.set(date, {
      amount: item.success_amount || 0,
      count: item.success_order || 0
    })
  })

  const sortedDates = Array.from(dateSet).sort()

  // 生成系列数据 - 使用柔和多样的颜色
  const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#9B59B6', '#3498DB', '#E67E22', '#1ABC9C', '#E74C3C']

  const series = []
  let colorIndex = 0
  accountMap.forEach((account) => {
    series.push({
      name: account.name,
      type: 'line',
      data: sortedDates.map(date => {
        const dayData = account.data.get(date)
        return appChartType.value === 'amount'
          ? (dayData?.amount || 0)
          : (dayData?.count || 0)
      }),
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      lineStyle: { width: 2, color: colors[colorIndex % colors.length] }
    })
    colorIndex++
  })

  appChartOption.value.xAxis.data = sortedDates
  appChartOption.value.series = series
  appChartOption.value.legend.data = series.map(s => s.name)

  // 更新Y轴名称
  appChartOption.value.yAxis.name = appChartType.value === 'amount' ? '金额(元)' : '笔数'
}

const updateStatistics = (data) => {
  if (!data || data.length === 0) {
    statistics.value[0].value =  '0'
    statistics.value[1].value =  '0'
    statistics.value[2].value =  '¥ 0'
    statistics.value[3].value =  '0%'
    return
  }

  const totalAmount = data.reduce((sum, item) => sum + (item.total_amount || 0), 0)
  const successAmount = data.reduce((sum, item) => sum + (item.success_amount || 0), 0)
  const totalOrders = data.reduce((sum, item) => sum + (item.total_order || 0), 0)
  const successRequests = data.reduce((sum, item) => sum + (item.success_order || 0), 0)

  statistics.value[0].value = formatNumber(totalRequests)
  statistics.value[1].value = formatNumber(totalOrders)
  statistics.value[2].value = '¥ ' + formatAmount(successAmount)
  statistics.value[3].value = totalOrders
      ? ((successRequests / totalOrders) * 100).toFixed(2) + '%'
      : '0%'
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  searchTableData()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  searchTableData()
}

// 监听全局日期变化
watch(globalDateRange, () => {
  refreshAllData()
}, { deep: true })

// 初始化
onMounted(async () => {
  try {
    // 获取应用列表 - 只在初始化时加载一次
    await fetchAppOptions()

    // 加载数据
    await refreshAllData()
  } catch (error) {
    console.error('初始化数据失败:', error)
    ElMessage.error('初始化数据失败')
  }
})
</script>

<style scoped lang="scss">
.system-dashboard {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .date-range-card {
    .global-date-range {
      display: flex;
      align-items: center;
      gap: 16px;

      .date-label {
        font-size: 14px;
        color: #606266;
        font-weight: 500;
      }

      .date-picker {
        width: 380px;
      }
    }
  }

  .stat-cards {
  }

  .stat-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 8px 16px rgba(0,0,0,0.1);
    }

    .stat-content {
      flex: 1;

      .stat-label {
        font-size: 14px;
        color: #909399;
        margin-bottom: 8px;
      }

      .stat-value {
        font-size: 28px;
        font-weight: bold;
        margin-bottom: 8px;
      }

      .stat-compare {
        font-size: 12px;
        display: flex;
        align-items: center;
        gap: 4px;

        .up {
          color: #67C23A;
          display: flex;
          align-items: center;
        }

        .down {
          color: #F56C6C;
          display: flex;
          align-items: center;
        }

        .compare-text {
          color: #909399;
        }
      }
    }

    .stat-icon {
      width: 60px;
      height: 60px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .chart-row {
  }

  .chart-card {
    :deep(.el-card__header) {
      padding: 12px 20px;
      border-bottom: 1px solid #ebeef5;
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 10px;

      .title {
        font-size: 16px;
        font-weight: 500;
      }

      .date-display {
        font-size: 13px;
        color: #909399;
      }

      .header-controls {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 10px;
      }
    }

    .chart-container {
      height: 400px;

      .chart {
        width: 100%;
        height: 100%;
      }
    }
  }

  .table-card {
    :deep(.el-card__header) {
      padding: 12px 20px;
      border-bottom: 1px solid #ebeef5;
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        font-size: 16px;
        font-weight: 500;
      }
    }

    .search-bar {
      padding: 20px;
      background-color: #f8f9fa;
      border-bottom: 1px solid #ebeef5;

      .search-form {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
        align-items: flex-end;
      }
    }

    .success-text {
      color: #67C23A;
      font-weight: 500;
    }

    .refund-text {
      color: #F56C6C;
    }

    .text-muted {
      color: #909399;
    }

    .rate-text {
      margin-left: 4px;
      font-size: 12px;
    }

    .pagination-wrap {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }
  }
}
</style>
