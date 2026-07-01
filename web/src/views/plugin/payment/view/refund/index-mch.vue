<template>
  <div class="page-wrap">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchInfo" class="search-form">
        <el-form-item label="应用编号">
          <el-select v-model="searchInfo.app_no" filterable placeholder="请选择应用编号" clearable remote :remote-method="remoteSearchAppVenuer" :loading="venuerAppLoading" style="width:220px" @change="handleAppNoChange">
            <el-option v-for="item in venuerAppList" :key="item.app_no" :label="item.app_no + ' (' + item.app_name + ')'" :value="item.app_no">
              <span>{{ item.app_no }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.app_name }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="订单号">
          <el-input v-model="searchInfo.order_no" placeholder="请输入订单号" clearable style="width:180px" @keyup.enter="onSubmit" @clear="onSubmit" />
        </el-form-item>
        <el-form-item label="交易号">
          <el-input v-model="searchInfo.trade_no" placeholder="请输入交易号" clearable style="width:180px" @keyup.enter="onSubmit" @clear="onSubmit" />
        </el-form-item>
        <el-form-item label="退款单号">
          <el-input v-model="searchInfo.refund_no" placeholder="请输入退款单号" clearable style="width:180px" @keyup.enter="onSubmit" @clear="onSubmit" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable style="width:120px">
            <el-option v-for="(label, value) in statusOptions" :key="value" :label="label" :value="Number(value)" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker v-model="searchInfo.start_time" type="datetime" placeholder="选择开始时间" value-format="YYYY-MM-DDTHH:mm:ssZ" format="YYYY-MM-DD HH:mm:ss" style="width:200px" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker v-model="searchInfo.end_time" type="datetime" placeholder="选择结束时间" value-format="YYYY-MM-DDTHH:mm:ssZ" format="YYYY-MM-DD HH:mm:ss" style="width:200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit"><el-icon><Search /></el-icon>搜索</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据栏 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" row-key="trade_no" v-loading="loading" stripe border>
        <el-table-column prop="refund_no" label="退款单号" min-width="180" show-overflow-tooltip />
        <el-table-column prop="mch_no" label="商户编号" min-width="120" show-overflow-tooltip />
        <el-table-column prop="mch_name" label="商户名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="app_no" label="应用编号" min-width="120" show-overflow-tooltip />
        <el-table-column prop="app_name" label="应用名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="refund_trade_no" label="退款交易号" min-width="180" show-overflow-tooltip />
        <el-table-column prop="order_no" label="订单号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="account_no" label="支付账户" min-width="120" show-overflow-tooltip />
        <el-table-column prop="account_name" label="账户名称" min-width="120" show-overflow-tooltip />
        <el-table-column label="支付金额" width="100" align="right">
          <template #default="{ row }">
            <span style="color:#67C23A;font-weight:500">¥{{ (row.payment_amount / 100).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="退款金额" width="100" align="right">
          <template #default="{ row }">
            <span style="color:#E6A23C;font-weight:500">¥{{ (row.refund_amount / 100).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small" effect="dark">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="退款原因" min-width="150" show-overflow-tooltip />
        <el-table-column label="创建时间" width="170" show-overflow-tooltip>
          <template #default="{ row }">{{ formatDisplayDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="完成时间" width="170" show-overflow-tooltip>
          <template #default="{ row }">{{ row.success_time ? formatDisplayDate(row.success_time) : '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button v-permission="'refund:view'" size="small" type="primary" link @click="viewRefundDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination background layout="total,sizes,prev,pager,next,jumper"
                       :total="total" v-model:current-page="page" v-model:page-size="pageSize"
                       :page-sizes="[10, 20, 50, 100]" @current-change="handleCurrentChange" @size-change="handleSizeChange" />
      </div>
    </el-card>

    <!-- 退款详情侧边栏 -->
    <el-drawer
        v-model="detailDrawerVisible"
        title="退款详情"
        direction="rtl"
        size="40%"
        destroy-on-close
    >
      <div v-if="currentDetail" class="p-4 detail-drawer">
        <!-- 基础信息 -->
        <div class="mb-6">
          <h3 class="text-lg font-bold mb-3 text-gray-700">基本信息</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="退款单号">{{ currentDetail.refund_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="退款交易号">{{ currentDetail.refund_trade_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="订单交易号">{{ currentDetail.trade_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="订单号">{{ currentDetail.order_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="商户编号">{{ currentDetail.mch_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="商户名称">{{ currentDetail.mch_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="应用编号">{{ currentDetail.app_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="应用名称">{{ currentDetail.app_name || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 金额信息 -->
        <div class="mb-6">
          <h3 class="text-lg font-bold mb-3 text-gray-700">金额信息</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="支付金额">
              <span class="font-bold text-green-600">¥{{ (currentDetail.payment_amount / 100).toFixed(2) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="退款金额">
              <span class="font-bold text-orange-600">¥{{ (currentDetail.refund_amount / 100).toFixed(2) }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 状态信息 -->
        <div class="mb-6">
          <h3 class="text-lg font-bold mb-3 text-gray-700">状态信息</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="退款状态">
              <el-tag :type="getStatusTagType(currentDetail.status)" size="small">
                {{ getStatusText(currentDetail.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="通知状态">
              <el-tag :type="currentDetail.notify_status === 2 ? 'success' : 'danger'" size="small">
                {{ getNotifyText(currentDetail.notify_status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="支付账户">{{ currentDetail.account_no || '-' }}</el-descriptions-item>
            <el-descriptions-item label="账户名称">{{ currentDetail.account_name || '-' }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 时间信息 -->
        <div class="mb-6">
          <h3 class="text-lg font-bold mb-3 text-gray-700">时间信息</h3>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="创建时间">{{ formatDisplayDate(currentDetail.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="完成时间">
              {{ currentDetail.success_time ? formatDisplayDate(currentDetail.success_time) : '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 其他信息 -->
        <div class="mb-6">
          <h3 class="text-lg font-bold mb-3 text-gray-700">其他信息</h3>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item label="退款发起来源">
              <div class="p-2 bg-gray-50 rounded text-gray-700">
                {{ currentDetail.refund_from || '无' }}
              </div>
            </el-descriptions-item>
            <el-descriptions-item label="退款原因">
              <div class="p-2 bg-gray-50 rounded text-gray-700">
                {{ currentDetail.reason || '无' }}
              </div>
            </el-descriptions-item>
            <el-descriptions-item label="失败提示">
              <div class="p-2 bg-gray-50 rounded text-gray-700">
                {{ currentDetail.third_msg || '无' }}
              </div>
            </el-descriptions-item>
            <el-descriptions-item label="渠道交易号">
              {{ currentDetail.out_mch_trade_no || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="回调地址">
              <div class="p-2 bg-gray-50 rounded text-gray-700 font-mono text-sm">
                {{ currentDetail.notify_url || '无' }}
              </div>
            </el-descriptions-item>
            <el-descriptions-item label="扩展参数">
              <div class="p-2 bg-gray-50 rounded">
                <pre v-if="currentDetail.ext_params" class="text-gray-700 text-sm whitespace-pre-wrap break-all">{{ currentDetail.ext_params }}</pre>
                <span v-else class="text-gray-500">无</span>
              </div>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      <template #footer>
        <div style="flex: auto" class="px-4 py-3 border-t border-gray-200">
          <el-button type="primary" @click="detailDrawerVisible = false">关闭</el-button>
        </div>
      </template>
    </el-drawer>

    <!-- 退款对话框 -->
    <!--    <el-dialog-->
    <!--        v-model="refundDialogVisible"-->
    <!--        title="发起退款"-->
    <!--        width="500px"-->
    <!--        destroy-on-close-->
    <!--    >-->
    <!--      <el-form-->
    <!--          :model="refundForm"-->
    <!--          :rules="refundRules"-->
    <!--          ref="refundFormRef"-->
    <!--          label-width="100px"-->
    <!--      >-->
    <!--        <el-form-item label="商户">-->
    <!--          <el-select-->
    <!--              v-model="refundForm.mch_no"-->
    <!--              filterable-->
    <!--              placeholder="请选择商户"-->
    <!--              clearable-->
    <!--              remote-->
    <!--              style="width: 100%"-->
    <!--              :remote-method="remoteSearchMchVenuer"-->
    <!--              :loading="venuerMchLoading"-->
    <!--              @change="handleRefundMchChange"-->
    <!--          >-->
    <!--            <el-option-->
    <!--                v-for="item in venuerMchList"-->
    <!--                :key="item.mch_no"-->
    <!--                :label="item.mch_no + ' (' + item.mch_name + ')'"-->
    <!--                :value="item.mch_no"-->
    <!--            >-->
    <!--              <span>{{ item.mch_no }}</span>-->
    <!--              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.mch_name }}</span>-->
    <!--            </el-option>-->
    <!--          </el-select>-->
    <!--        </el-form-item>-->
    <!--        <el-form-item label="应用编号">-->
    <!--          <el-select-->
    <!--              v-model="refundForm.app_no"-->
    <!--              filterable-->
    <!--              placeholder="请选择应用编号"-->
    <!--              clearable-->
    <!--              remote-->
    <!--              style="width: 100%"-->
    <!--              :remote-method="remoteSearchRefundApp"-->
    <!--              :loading="refundAppLoading"-->
    <!--          >-->
    <!--            <el-option-->
    <!--                v-for="item in refundAppList"-->
    <!--                :key="item.app_no"-->
    <!--                :label="item.app_no + ' (' + item.name + ')'"-->
    <!--                :value="item.app_no"-->
    <!--            >-->
    <!--              <span>{{ item.app_no }}</span>-->
    <!--              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.name }}</span>-->
    <!--            </el-option>-->
    <!--          </el-select>-->
    <!--        </el-form-item>-->
    <!--        <el-form-item label="订单号" prop="order_no">-->
    <!--          <el-input-->
    <!--              v-model="refundForm.order_no"-->
    <!--              placeholder="请输入订单号"-->
    <!--              clearable-->
    <!--              @blur="handleOrderNoBlur"-->
    <!--          />-->
    <!--        </el-form-item>-->
    <!--        <el-form-item label="退款金额" prop="amount">-->
    <!--          <el-input-number-->
    <!--              v-model="refundForm.amount"-->
    <!--              :min="0.01"-->
    <!--              :max="maxRefundAmount"-->
    <!--              :step="0.01"-->
    <!--              :precision="2"-->
    <!--              controls-position="right"-->
    <!--              class="w-full"-->
    <!--          >-->
    <!--            <template #prefix>¥</template>-->
    <!--          </el-input-number>-->
    <!--          <div class="text-xs text-gray-500 mt-1" v-if="maxRefundAmount > 0">-->
    <!--            最大可退金额: ¥{{ maxRefundAmount.toFixed(2) }}-->
    <!--          </div>-->
    <!--        </el-form-item>-->
    <!--        <el-form-item label="退款原因" prop="reason">-->
    <!--          <el-input-->
    <!--              v-model="refundForm.reason"-->
    <!--              type="textarea"-->
    <!--              :rows="3"-->
    <!--              placeholder="请输入退款原因"-->
    <!--              maxlength="200"-->
    <!--              show-word-limit-->
    <!--          />-->
    <!--        </el-form-item>-->
    <!--      </el-form>-->
    <!--      <template #footer>-->
    <!--        <div class="dialog-footer">-->
    <!--          <el-button @click="refundDialogVisible = false">取消</el-button>-->
    <!--          <el-button type="warning" @click="handleRefundSubmit" :loading="refundLoading">-->
    <!--            确认退款-->
    <!--          </el-button>-->
    <!--        </div>-->
    <!--      </template>-->
    <!--    </el-dialog>-->
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { refundList, tradeRefund, refundDetail, refundAmount, appList } from '@/api/payment'

// 响应式数据
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 搜索条件
const searchInfo = reactive({
  mch_no: '',
  app_no: '',
  trade_no: '',
  order_no: '',
  refund_no: '',
  status: '',
  start_time: '',
  end_time: ''
})

// 状态选项映射（来自 refund.pb.go）
const statusOptions = {
  0: '全部',
  1: '已创建',
  2: '处理中',
  3: '成功',
  4: '失败',
  5: '异常',
  6: '已关闭',
  7: '已取消'
}

// 商户和应用搜索
const venuerMchLoading = ref(false)
const venuerMchList = ref([])
const venuerAppLoading = ref(false)
const venuerAppList = ref([])

// 退款对话框相关
const refundDialogVisible = ref(false)
const refundFormRef = ref(null)
const refundLoading = ref(false)
const refundAppLoading = ref(false)
const refundAppList = ref([])
const refundForm = reactive({
  mch_no: '',
  app_no: '',
  order_no: '',
  amount: 0,
  reason: ''
})
const maxRefundAmount = ref(0)

// 详情侧边栏相关
const detailDrawerVisible = ref(false)
const currentDetail = ref(null)

// 退款表单验证规则
const refundRules = {
  app_no: [
    { required: true, message: '请选择应用编号', trigger: 'blur' }
  ],
  order_no: [
    { required: true, message: '请输入订单号', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入退款金额', trigger: 'blur' },
    { validator: (rule, value, callback) => {
        if (value <= 0) {
          callback(new Error('退款金额必须大于0'))
        } else if (value > maxRefundAmount.value) {
          callback(new Error(`退款金额不能超过可退金额 ¥${maxRefundAmount.value.toFixed(2)}`))
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

// 组件挂载时初始化数据
const initData = async () => {
  try {
      // 设置默认时间范围（最近7天）
      const endTime = new Date()
      const startTime = new Date()
      startTime.setDate(startTime.getDate() - 7)
      searchInfo.start_time = formatISO8601(startTime)
      searchInfo.end_time = formatISO8601(endTime)
      // 然后根据默认商户获取退款列表
      await getRefundListFunc()
  } catch (error) {
    console.error('初始化数据失败:', error)
    ElMessage.error('初始化数据失败')
  }
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

// 格式化日期为选择器显示的格式
const formatDateForPicker = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 格式化ISO 8601日期为显示格式
const formatDisplayDate = (dateString) => {
  if (!dateString) return ''
  if (dateString === "0001-01-01T00:00:00Z") {
    // 如果包含T，则将其替换为空格
    return '-'
  }
  try {
    // 如果已经是ISO格式，直接解析
    const date = new Date(dateString)
    if (isNaN(date.getTime())) {
      return dateString
    }
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  } catch (error) {
    console.error('日期格式化错误:', error)
    return dateString
  }
}

// 远程搜索应用
const remoteSearchAppVenuer = async (query) => {
  venuerAppLoading.value = true
  try {
    const res = await appList({
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

// 退款对话框中的远程搜索应用
const remoteSearchRefundApp = async (query) => {
  refundAppLoading.value = true
  try {
    const res = await appList({
      mch_no: refundForm.mch_no,
      name: query,
      page: 1,
      page_size: 10
    })
    if (res.code === 0) {
      refundAppList.value = res.data.list || []
    } else {
      refundAppList.value = []
    }
  } catch (error) {
    console.error('搜索应用失败:', error)
    ElMessage.error('搜索应用失败')
  } finally {
    refundAppLoading.value = false
  }
}

// 商户选择变化事件
const handleMchNoChange = (value) => {
  if (value) {
    page.value = 1
    // 清空应用列表，重新加载
    venuerAppList.value = []
    searchInfo.app_no = ''
    getRefundListFunc()
  }
}

// 退款对话框中的商户变化事件
const handleRefundMchChange = (value) => {
  if (value) {
    // 清空应用列表和订单信息
    refundAppList.value = []
    refundForm.app_no = ''
    refundForm.order_no = ''
    maxRefundAmount.value = 0
  }
}

// 应用选择变化事件
const handleAppNoChange = (value) => {
  if (value) {
    page.value = 1
    getRefundListFunc()
  }
}

// 订单号失焦事件
const handleOrderNoBlur = async () => {
  if (refundForm.mch_no && refundForm.app_no && refundForm.order_no) {
    try {
      const response = await refundAmount({
        mch_no: refundForm.mch_no,
        app_no: refundForm.app_no,
        order_no: refundForm.order_no
      })

      if (response.code === 200 || response.code === 0) {
        const amount = response.data || 0
        maxRefundAmount.value = amount / 100
        if (maxRefundAmount.value > 0) {
          refundForm.amount = maxRefundAmount.value
        }
      }
    } catch (error) {
      console.error('查询可退金额失败:', error)
      maxRefundAmount.value = 0
    }
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
// 获取退款列表
const getRefundListFunc = async () => {
  loading.value = true
  try {
    const params = {
      ...searchInfo,
      page: page.value,
      page_size: pageSize.value
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

    console.log('请求参数:', params) // 调试用

    const res = await refundList(params)
    if (res.code === 200 || res.code === 0) {
      const list = res.data?.list || []  // 如果 list 是 null 或 undefined，使用空数组
      tableData.value = Array.isArray(list) ? list : []  // 双重保险
      total.value = res.data?.total || 0
    } else {
      tableData.value = []
      total.value = 0
      ElMessage.error(res.msg || res.message || '获取退款列表失败')
    }
  } catch (error) {
    console.error('获取退款列表失败:', error)
    ElMessage.error('获取退款列表失败')
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 分页处理
const handleCurrentChange = (val) => {
  page.value = val
  getRefundListFunc()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getRefundListFunc()
}

// 搜索和重置
const onSubmit = () => {
  page.value = 1
  getRefundListFunc()
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
      searchInfo[key] = ''
    }
  })
  page.value = 1
  getRefundListFunc()
}

// 显示退款对话框
const showRefundDialog = () => {
  refundForm.mch_no = searchInfo.mch_no || ''
  refundForm.app_no = searchInfo.app_no || ''
  refundForm.order_no = ''
  refundForm.amount = 0
  refundForm.reason = ''
  maxRefundAmount.value = 0
  refundAppList.value = []

  if (refundForm.mch_no) {
    // 加载当前商户的应用列表
    remoteSearchRefundApp('')
  }

  refundDialogVisible.value = true
}

const getNotifyText = (status) => {
  if (status === 0) {
    return '无需回调'
  } else if (status === 1) {
    return '待通知'
  } else if (status === 2) {
    return '通知成功'
  } else if (status === 3) {
    return '通知失败'
  }
  return '未知状态'
}

// 提交退款申请
const handleRefundSubmit = async () => {
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

    const res = await tradeRefund(params)
    if (res.code === 0) {
      ElMessage.success('退款申请提交成功')
      refundDialogVisible.value = false
      // 刷新列表
      getRefundListFunc()
    } else {
      ElMessage.error(res.msg || res.message || '退款申请失败')
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

// 查看退款详情
const viewRefundDetail = async (row) => {
  try {
    const response = await refundDetail({
      mch_no: row.mch_no,
      app_no: row.app_no,
      trade_no: row.refund_trade_no
    })

    if (response.code === 0) {
      currentDetail.value = response.data
      detailDrawerVisible.value = true
    } else {
      ElMessage.error(response.msg || '获取详情失败')
    }
  } catch (error) {
    console.error('获取退款详情失败:', error)
    ElMessage.error('获取详情失败，请稍后重试')
  }
}

// 获取状态文本
const getStatusText = (status) => {
  return statusOptions[status] || '未知状态'
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  const typeMap = {
    0: 'info',      // 未知
    1: 'info',      // 已创建
    2: 'warning',   // 处理中
    3: 'success',   // 成功
    4: 'danger',    // 失败
    5: 'danger',    // 异常
    6: 'danger',    // 已关闭
    7: 'warning'    // 已取消
  }
  return typeMap[status] || 'info'
}
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.search-card { }
.search-form { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.search-form .el-form-item { margin-bottom: 0; }
.table-card { }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 16px; }
.detail-drawer { max-height: calc(100vh - 120px); overflow-y: auto; }
</style>
