<!-- src/views/payment/trade/index.vue -->
<template>
  <div class="gva-container">
    <!-- 页面标题 -->
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-bold">支付订单管理</h2>
    </div>

    <!-- 搜索区域 -->
    <div class="search-box bg-white p-4 rounded-md shadow mb-4">
      <el-form :inline="true" :model="searchInfo" class="search-form">
        <el-row :gutter="10">
          <el-col :span="6">
            <el-form-item label="商户">
              <el-select
                  v-model="searchInfo.mch_no"
                  filterable
                  placeholder="请选择商户"
                  clearable
                  remote
                  :remote-method="remoteSearchMchVenuer"
                  :loading="venuerMchLoading"
                  @change="handleMchNoChange"
              >
                <el-option
                    v-for="item in venuerMchList"
                    :key="item.mch_no"
                    :label="item.mch_no + ' (' + item.mch_name + ')'"
                    :value="item.mch_no"
                >
                  <span>{{ item.mch_no }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ item.mch_name }}</span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="应用编号">
              <el-select
                  v-model="searchInfo.app_no"
                  filterable
                  placeholder="请选择应用编号"
                  clearable
                  remote
                  :remote-method="remoteSearchAppVenuer"
                  :loading="venuerAppLoading"
                  @change="handleAppNoChange"
              >
                <el-option
                    v-for="item in venuerAppList"
                    :key="item.app_no"
                    :label="item.app_no + ' (' + item.app_name + ')'"
                    :value="item.app_no"
                >
                  <span>{{ item.app_no }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ item.app_name }}</span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="订单号">
              <el-input
                  v-model="searchInfo.order_no"
                  placeholder="请输入订单号"
                  clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="支付单号">
              <el-input
                  v-model="searchInfo.trade_no"
                  placeholder="请输入支付单号"
                  clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="6">
            <el-form-item label="状态">
              <el-select
                  v-model="searchInfo.status"
                  placeholder="请选择状态"
                  clearable
                  class="w-full"
              >
                <el-option
                    v-for="(label, value) in statusOptions"
                    :key="value"
                    :label="label"
                    :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="开始时间">
              <el-date-picker
                  v-model="searchInfo.start_time"
                  type="datetime"
                  placeholder="选择开始时间"
                  value-format="YYYY-MM-DDTHH:mm:ssZ"
                  format="YYYY-MM-DD HH:mm:ss"
                  class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="结束时间">
              <el-date-picker
                  v-model="searchInfo.end_time"
                  type="datetime"
                  placeholder="选择结束时间"
                  value-format="YYYY-MM-DDTHH:mm:ssZ"
                  format="YYYY-MM-DD HH:mm:ss"
                  class="w-full"
              />
            </el-form-item>
          </el-col>
          <el-col :span="6" class="text-right">
            <el-button type="primary" @click="onSubmit">查询</el-button>
            <el-button @click="onReset">重置</el-button>
          </el-col>
        </el-row>
      </el-form>
    </div>

    <!-- 表格区域 -->
    <div class="table-container bg-white p-4 rounded-md shadow">
      <el-table
          :data="tableData"
          style="width: 100%"
          row-key="trade_no"
          v-loading="loading"
          border
      >
        <el-table-column prop="trade_no" label="交易号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="order_no" label="订单号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="mch_no" label="商户编号" min-width="120" show-overflow-tooltip />
        <el-table-column prop="mch_name" label="商户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="app_no" label="应用编号" min-width="120" show-overflow-tooltip />
        <el-table-column prop="app_name" label="应用名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="account_no" label="支付账户" min-width="120" show-overflow-tooltip />
        <el-table-column prop="account_name" label="账户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="payment_method" label="支付方式" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getPaymentMethodText(row.payment_method) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="payment_product" label="支付产品" width="100">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ getPaymentProductText(row.payment_product) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="100">
          <template #default="{ row }">
            <span class="font-bold text-green-600">¥{{ (row.amount / 100).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="商品标题" min-width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="order_time" label="下单时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.order_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <div class="flex flex-wrap gap-1">
              <el-button
                  size="small"
                  type="primary"
                  link
                  @click="viewOrderDetail(row)"
              >
                查看详情
              </el-button>
              <el-button
                  v-if="row.status === 3 || row.status === 7"
                  size="small"
                  type="warning"
                  link
                  @click="openRefundDialog(row)"
              >
                退款
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination mt-4">
        <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 20, 30, 50]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 订单详情组件 -->
    <TradeDetail
        v-model:visible="detailDrawerVisible"
        :orderData="selectedOrder"
        @refresh="getTradeListFunc"
    />

    <!-- 退款弹窗 -->
    <el-dialog
        v-model="refundDialogVisible"
        title="申请退款"
        width="500px"
        destroy-on-close
    >
      <el-form
          :model="refundForm"
          :rules="refundRules"
          ref="refundFormRef"
          label-width="100px"
      >
        <el-form-item label="订单号">
          <el-input v-model="refundForm.order_no" disabled />
        </el-form-item>
        <el-form-item label="商户编号">
          <el-input v-model="refundForm.mch_no" disabled />
        </el-form-item>
        <el-form-item label="应用编号">
          <el-input v-model="refundForm.app_no" disabled />
        </el-form-item>
        <el-form-item label="退款金额" prop="amount">
          <el-input-number
              v-model="refundForm.amount"
              :min="0.00"
              :max="currentOrderAmount"
              :step="0.01"
              :precision="2"
              controls-position="right"
              class="w-full"
          >
            <template #prefix>¥</template>
          </el-input-number>
          <div class="text-xs text-gray-500 mt-1">
            可退金额: ¥{{ (currentOrderAmount).toFixed(2) }}
          </div>
        </el-form-item>
        <el-form-item label="退款原因" prop="reason">
          <el-input
              v-model="refundForm.reason"
              type="textarea"
              :rows="3"
              placeholder="请输入退款原因"
              maxlength="20"
              show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="refundDialogVisible = false">取消</el-button>
          <el-button type="warning" @click="handleRefund" :loading="refundLoading">
            确认退款
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {getAvailableRefundAmount, getTradeList, refund as refundApi} from '@/plugin/payment/api/trade.js'
import { getMerchantList } from "@/plugin/payment/api/merchant.js"
import { getApplicationList } from '@/plugin/payment/api/application.js'
import TradeDetail from '@/plugin/payment/components/trade.vue'

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 搜索条件
const searchInfo = reactive({
  mch_no: '',
  order_no: '',
  app_no: '',
  trade_no: '',
  status: '',
  start_time: '',
  end_time: ''
})

// 状态选项映射（来自 payment.pb.go）
const statusOptions = {
  0: '全部',
  1: '初始化',
  2: '预支付',
  3: '成功',
  4: '失败',
  5: '取消',
  6: '关闭',
  7: '退款',
  8: '等待确认',
  9: '等待支付',
  10: '支付超时',
  11: '临时支付失败',
}

// 支付方式映射
const paymentMethodOptions = {
  'wechat': '微信支付',
  'alipay': '支付宝'
}

// 支付产品映射
const paymentProductOptions = {
  'JSAPI': 'JSAPI支付',
  'APP': 'APP支付',
  'MINI': '小程序支付',
  'H5': 'H5支付',
  'NATIVE': '扫码支付'
}

// 详情侧边栏
const detailDrawerVisible = ref(false)
const selectedOrder = ref(null)

// 退款相关
const refundDialogVisible = ref(false)
const refundFormRef = ref(null)
const refundLoading = ref(false)
const refundForm = reactive({
  order_no: '',
  mch_no: '',
  app_no: '',
  amount: 0,
  reason: ''
})
const currentOrderAmount = ref(0)

// 退款表单验证规则
const refundRules = {
  amount: [
    { required: true, message: '请输入退款金额', trigger: 'blur' },
    { validator: (rule, value, callback) => {
        if (value <= 0) {
          callback(new Error('退款金额必须大于0'))
        } else if (value > currentOrderAmount.value) {
          callback(new Error(`退款金额不能超过可退金额 ¥${currentOrderAmount.value.toFixed(2)}`))
        } else {
          callback()
        }
      }, trigger: 'blur' }
  ],
  reason: [
    { required: true, message: '请输入退款原因', trigger: 'blur' },
    { min: 5, message: '退款原因至少5个字符', trigger: 'blur' }
  ]
}

const venuerMchLoading = ref(false)
const venuerMchList = ref([])
const venuerAppLoading = ref(false)
const venuerAppList = ref([])

// 组件挂载时初始化数据
const initData = async () => {
  try {
    // 先获取商户列表
    const merchantRes = await getMerchantList({ page: 1, page_size: 10 })
    if (merchantRes.code === 0 && merchantRes.data.list && merchantRes.data.list.length > 0) {
      venuerMchList.value = merchantRes.data.list
      // 默认选择第一个商户
      const firstMerchant = venuerMchList.value[0]
      searchInfo.mch_no = firstMerchant.mch_no
      // 设置默认时间范围（最近7天）
      const endTime = new Date()
      const startTime = new Date()
      startTime.setDate(startTime.getDate() - 7)
      searchInfo.start_time = formatISO8601(startTime)
      searchInfo.end_time = formatISO8601(endTime)
      // 然后根据默认商户获取支付列表
      await getTradeListFunc()
    } else {
      venuerMchList.value = []
      ElMessage.warning('未获取到商户数据')
    }
  } catch (error) {
    console.error('初始化数据失败:', error)
    ElMessage.error('初始化数据失败')
  }
}

// 格式化日期为选择器格式
const formatDateForPicker = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 远程搜索商户
const remoteSearchMchVenuer = async (query) => {
  venuerMchLoading.value = true
  try {
    const res = await getMerchantList({ mch_name: query, page: 1, page_size: 10 })
    if (res.code === 0) {
      venuerMchList.value = res.data.list || []
    } else {
      venuerMchList.value = []
    }
  } catch (error) {
    console.error('搜索商户失败:', error)
    ElMessage.error('搜索商户失败')
  } finally {
    venuerMchLoading.value = false
  }
}

// 远程搜索应用
const remoteSearchAppVenuer = async (query) => {
  venuerAppLoading.value = true
  try {
    const res = await getApplicationList({
      mch_no: searchInfo.mch_no,
      name: query,
      page: 1,
      page_size: 10
    })
    if (res.code === 0) {
      venuerAppList.value = res.data.list || []
    } else {
      venuerAppList.value = []
    }
  } catch (error) {
    console.error('搜索应用失败:', error)
    ElMessage.error('搜索应用失败')
  } finally {
    venuerAppLoading.value = false
  }
}

// 商户选择变化事件
const handleMchNoChange = (value) => {
  if (value) {
    page.value = 1
    // 清空应用列表，重新加载
    venuerAppList.value = []
    searchInfo.app_no = ''
    getTradeListFunc()
  }
}

// 应用选择变化事件
const handleAppNoChange = (value) => {
  if (value) {
    page.value = 1
    getTradeListFunc()
  }
}

// 组件挂载时获取数据
onMounted(() => {
  initData()
})


// 确保日期字符串为 ISO 8601 格式
const ensureISOFormat = (dateStr) => {
  if (!dateStr) return ''

  // 如果已经是 ISO 格式（包含T和Z），直接返回
  if (dateStr.includes('T') && dateStr.includes('Z')) {
    return dateStr
  }

  // 如果是 YYYY-MM-DD HH:mm:ss 格式，转换为 ISO 格式
  const date = new Date(dateStr)
  return formatISO8601(date)
}


// 获取支付订单列表
const getTradeListFunc = async () => {
  loading.value = true
  try {
    const params = {
      ...searchInfo,
      page: page.value,
      page_size: pageSize.value
    }

    // 如果商户编号为空，则不请求数据
    if (!params.mch_no) {
      tableData.value = []
      total.value = 0
      return
    }

    // 2. 时间范围校验：如果同时传了开始和结束时间
    if (params.start_time && params.end_time) {
      const startDate = new Date(params.start_time)
      const endDate = new Date(params.end_time)

      // 计算时间差（毫秒）
      const diffMs = endDate - startDate
      // 转换为天数
      const diffDays = diffMs / (1000 * 60 * 60 * 24)

      // 如果超过60天（约两个月）
      if (diffDays > 60) {
        ElMessage.warning('查询时间范围不能超过两个月，已自动调整为最近60天')
        // 自动调整结束时间为开始时间+60天
        const newEndDate = new Date(startDate)
        newEndDate.setDate(startDate.getDate() + 60)
        params.end_time = formatISO8601(newEndDate)
      }
    }

    // 3. 格式化时间为 ISO 8601 格式
    if (params.start_time) {
      // 如果已经是字符串，确保格式正确
      params.start_time = ensureISOFormat(params.start_time)
    }
    if (params.end_time) {
      params.end_time = ensureISOFormat(params.end_time)
    }

    const res = await getTradeList(params)
    if (res.code === 0) {
      const list = res.data?.list || []  // 如果 list 是 null 或 undefined，使用空数组
      tableData.value = Array.isArray(list) ? list : []  // 双重保险
      total.value = res.data?.total || 0
    } else {
      tableData.value = []
      total.value = 0
      ElMessage.error(res.message || '获取支付订单列表失败')
    }
  } catch (error) {
    console.error('获取支付订单列表失败:', error)
    ElMessage.error('获取支付订单列表失败')
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 分页处理
const handleCurrentChange = (val) => {
  page.value = val
  getTradeListFunc()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getTradeListFunc()
}

// 搜索和重置
const onSubmit = () => {
  page.value = 1
  getTradeListFunc()
}

const onReset = () => {
  Object.keys(searchInfo).forEach(key => {
    if (key === 'mch_no') {
      // 重置时保持第一个商户的选择
      if (venuerMchList.value.length > 0) {
        searchInfo.mch_no = venuerMchList.value[0].mch_no
      } else {
        searchInfo.mch_no = ''
      }
    } else if (key === 'start_time' || key === 'end_time') {
      // 重置时间范围到最近7天
      const endTime = new Date()
      const startTime = new Date()
      startTime.setDate(startTime.getDate() - 7)
      if (key === 'start_time') {
        searchInfo.start_time = formatISO8601(startTime)
      } else {
        searchInfo.end_time = formatISO8601(endTime)
      }
    } else {
      searchInfo[key] = null
    }
  })
  page.value = 1
  getTradeListFunc()
}

// 查看订单详情
const viewOrderDetail = (row) => {
  selectedOrder.value = row
  detailDrawerVisible.value = true
}

// 处理从详情组件传来的退款请求
// const handleOpenRefund = (order) => {
//   refundForm.order_no = order.order_no
//   refundForm.mch_no = order.mch_no
//   refundForm.app_no = order.app_no
//   refundForm.amount = (order.amount / 100).toFixed(2)
//   refundForm.reason = ''
//   currentOrderAmount.value = order.amount / 100
//   refundDialogVisible.value = true
// }

// 打开退款对话框
const openRefundDialog = async (row) => {
  refundForm.order_no = row.order_no
  refundForm.mch_no = row.mch_no
  refundForm.app_no = row.app_no
  try {
    const res = await getAvailableRefundAmount({order_no: row.order_no, mch_no: row.mch_no, app_no: row.app_no})
    if (res.code === 0) {
      const refundAmount = res.data || 0
      refundForm.amount = (refundAmount / 100).toFixed(2)
      refundForm.reason = ''
      currentOrderAmount.value = refundAmount / 100
      refundDialogVisible.value = true
    } else {
      ElMessage.error(res.message || '获取可退款金额失败')
    }
  } catch (error) {
    console.error('获取可退款金额失败:', error)
    ElMessage.error('获取可退款金额失败')
  }
}

// 处理退款
const handleRefund = async () => {
  if (!refundFormRef.value) return
  try {
    await refundFormRef.value.validate()
    await ElMessageBox.confirm(
        `确认退款 ¥${refundForm.amount.toFixed(2)} 吗？`,
        '退款确认',
        {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )
    refundLoading.value = true
    const params = {
      order_no: refundForm.order_no,
      mch_no: refundForm.mch_no,
      app_no: refundForm.app_no,
      amount: Math.round(refundForm.amount * 100), // 转为分
      reason: refundForm.reason
    }

    const res = await refundApi(params)
    if (res.code === 0) {
      ElMessage.success('退款申请提交成功')
      refundDialogVisible.value = false
      // 刷新列表
      getTradeListFunc()
    } else {
      ElMessage.error(res.message || '退款申请失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('退款失败:', error)
      ElMessage.error('退款失败')
    }
  } finally {
    refundLoading.value = false
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取状态文本
const getStatusText = (status) => {
  return statusOptions[status] || '未知状态'
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const typeMap = {
    0: 'info',
    1: 'info',
    2: 'warning',
    3: 'success',
    4: 'danger',
    5: 'danger',
    6: 'danger',
    7: 'warning',
    8: 'warning',
    9: 'warning',
    10: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取支付方式文本
const getPaymentMethodText = (method) => {
  return method || '-'
}

// 获取支付产品文本
const getPaymentProductText = (product) => {
  return product || '-'
}


// 格式化日期为ISO 8601格式（后端要求的格式）
const formatISO8601 = (date) => {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`
}

</script>

<style scoped>
.gva-container {
  padding: 20px;
}

.search-box {
  margin-bottom: 20px;
}

.table-container {
  margin-bottom: 20px;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-table .el-table__row .cell) {
  word-break: break-all;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
