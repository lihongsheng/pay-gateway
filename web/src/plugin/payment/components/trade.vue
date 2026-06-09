<!-- src/views/payment/trade/components/TradeDetail.vue -->
<template>
  <el-drawer
      v-model="drawerVisible"
      title="订单详情"
      :size="600"
      direction="rtl"
      destroy-on-close
  >
    <div v-if="orderDetail" class="order-detail-drawer">
      <!-- 订单基本信息卡片 -->
      <div class="order-info-card mb-6">
        <div class="flex items-center mb-4">
          <div class="flex-1">
            <h3 class="text-lg font-bold text-gray-800">{{ orderDetail.subject || '无商品标题' }}</h3>
            <p class="text-sm text-gray-500 mt-1">{{ orderDetail.desc || '无商品描述' }}</p>
          </div>
          <el-tag :type="getStatusTagType(orderDetail.status)" size="large" class="ml-2">
            {{ getStatusText(orderDetail.status) }}
          </el-tag>
        </div>

        <div class="amount-section mb-4">
          <div class="text-sm text-gray-500">支付金额</div>
          <div class="text-2xl font-bold text-green-600">¥{{ (orderDetail.amount / 100).toFixed(2) }}</div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <div class="text-xs text-gray-400">订单号</div>
            <div class="text-sm font-medium">{{ orderDetail.order_no }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-400">交易号</div>
            <div class="text-sm font-medium">{{ orderDetail.trade_no }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-400">第三方流水号</div>
            <div class="text-sm font-medium">{{ orderDetail.out_mch_trade_no || '-' }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-400">支付方式</div>
            <div class="text-sm font-medium">{{ getPaymentMethodText(orderDetail.payment_method) }}</div>
          </div>
        </div>
      </div>

      <!-- 商户信息卡片 -->
      <div class="info-card mb-6">
        <h4 class="text-md font-semibold mb-3 text-gray-700 border-b pb-2">商户信息</h4>
        <div class="space-y-3">
          <div class="flex items-center">
            <span class="info-label">商户编号</span>
            <span class="info-value">{{ orderDetail.mch_no }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">商户名称</span>
            <span class="info-value">{{ orderDetail.mch_name }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">应用编号</span>
            <span class="info-value">{{ orderDetail.app_no }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">应用名称</span>
            <span class="info-value">{{ orderDetail.app_name }}</span>
          </div>
        </div>
      </div>

      <!-- 支付信息卡片 -->
      <div class="info-card mb-6">
        <h4 class="text-md font-semibold mb-3 text-gray-700 border-b pb-2">支付信息</h4>
        <div class="space-y-3">
          <div class="flex items-center">
            <span class="info-label">支付账户</span>
            <span class="info-value">{{ orderDetail.account_no }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">账户名称</span>
            <span class="info-value">{{ orderDetail.account_name }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">支付产品</span>
            <span class="info-value">{{ getPaymentProductText(orderDetail.payment_product) }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">用户OpenID</span>
            <span class="info-value">{{ orderDetail.user_open_id || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- 时间信息卡片 -->
      <div class="info-card mb-6">
        <h4 class="text-md font-semibold mb-3 text-gray-700 border-b pb-2">时间信息</h4>
        <div class="space-y-3">
          <div class="flex items-center">
            <span class="info-label">创建时间</span>
            <span class="info-value">{{ formatDate(orderDetail.created_at) }}</span>
          </div>
          <div class="flex items-center">
            <span class="info-label">下单时间</span>
            <span class="info-value">{{ formatDate(orderDetail.order_time) }}</span>
          </div>
        </div>
      </div>

      <!-- 通知信息卡片 -->
      <div class="info-card mb-6">
        <h4 class="text-md font-semibold mb-3 text-gray-700 border-b pb-2">通知信息</h4>
        <div class="space-y-3">
          <div class="flex items-center">
            <span class="info-label">通知状态</span>
            <span class="info-value">
              <el-tag :type="orderDetail.notify_status === 1 ? 'success' : 'danger'" size="small">
                {{ getNotifyText(orderDetail.notify_status) }}
              </el-tag>
            </span>
          </div>
          <div class="flex items-start">
            <span class="info-label pt-1">通知地址</span>
            <span class="info-value break-all">
              <a v-if="orderDetail.notify_url" :href="orderDetail.notify_url" target="_blank" class="text-blue-500 hover:underline text-sm">
                {{ orderDetail.notify_url }}
              </a>
              <span v-else class="text-gray-400">-</span>
            </span>
          </div>
          <div class="flex items-start">
            <span class="info-label pt-1">跳转链接</span>
            <span class="info-value break-all">
              <a v-if="orderDetail.redirect_url" :href="orderDetail.redirect_url" target="_blank" class="text-blue-500 hover:underline text-sm">
                {{ orderDetail.redirect_url }}
              </a>
              <span v-else class="text-gray-400">-</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 第三方信息卡片 -->
      <div class="info-card mb-6">
        <h4 class="text-md font-semibold mb-3 text-gray-700 border-b pb-2">第三方信息</h4>
        <div class="space-y-3">
          <div class="flex items-center">
            <span class="info-label">第三方代码</span>
            <span class="info-value">{{ orderDetail.third_code || '-' }}</span>
          </div>
          <div class="flex items-start">
            <span class="info-label pt-1">第三方消息</span>
            <span class="info-value break-all">{{ orderDetail.third_msg || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- 其他信息卡片 -->
      <div class="info-card">
        <h4 class="text-md font-semibold mb-3 text-gray-700 border-b pb-2">其他信息</h4>
        <div class="space-y-3">
          <div class="flex items-start">
            <span class="info-label pt-1">扩展参数</span>
            <span class="info-value break-all">
              <pre v-if="orderDetail.ext_params" class="text-xs bg-gray-50 p-2 rounded mt-1 overflow-auto max-h-32">orderDetail.ext_params</pre>
              <span v-else class="text-gray-400">-</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
<!--      <div class="mt-8 pt-6 border-t border-gray-200">-->
<!--        <div class="flex justify-between">-->
<!--          <el-button @click="drawerVisible = false">关闭</el-button>-->
<!--          <div class="space-x-2">-->
<!--            <el-button-->
<!--                v-if="orderDetail.status === 3 && orderDetail.refund_status !== 7"-->
<!--                type="warning"-->
<!--                @click="handleOpenRefund"-->
<!--            >-->
<!--              申请退款-->
<!--            </el-button>-->
<!--&lt;!&ndash;            <el-button&ndash;&gt;-->
<!--&lt;!&ndash;                v-if="orderDetail.status === 3"&ndash;&gt;-->
<!--&lt;!&ndash;                type="success"&ndash;&gt;-->
<!--&lt;!&ndash;                @click="handleResendNotify"&ndash;&gt;-->
<!--&lt;!&ndash;            >&ndash;&gt;-->
<!--&lt;!&ndash;              重发通知&ndash;&gt;-->
<!--&lt;!&ndash;            </el-button>&ndash;&gt;-->
<!--          </div>-->
<!--        </div>-->
<!--      </div>-->
    </div>

    <!-- 加载状态 -->
    <div v-else-if="loading" class="flex justify-center items-center h-64">
      <el-icon class="is-loading" size="24">
        <Loading />
      </el-icon>
      <span class="ml-2 text-gray-500">加载中...</span>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, defineProps, defineEmits, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { getTrade } from '@/plugin/payment/api/trade.js'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  orderData: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:visible', 'open-refund', 'refresh'])

const drawerVisible = ref(false)
const orderDetail = ref(null)
const loading = ref(false)

// 状态选项映射
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

// 监听visible变化
watch(() => props.visible, (val) => {
  drawerVisible.value = val
  if (val && props.orderData) {
    loadOrderDetail()
  } else {
    orderDetail.value = null
  }
})

// 监听drawerVisible变化
watch(() => drawerVisible.value, (val) => {
  emit('update:visible', val)
})

// 加载订单详情
const loadOrderDetail = async () => {
  if (!props.orderData) return

  loading.value = true
  try {
    const res = await getTrade({
      mch_no: props.orderData.mch_no,
      order_no: props.orderData.order_no,
      app_no: props.orderData.app_no
    })

    if (res.code === 0) {
      orderDetail.value = res.data
    } else {
      ElMessage.error(res.message || '获取订单详情失败')
      drawerVisible.value = false
    }
  } catch (error) {
    console.error('获取订单详情失败:', error)
    ElMessage.error('获取订单详情失败')
    drawerVisible.value = false
  } finally {
    loading.value = false
  }
}

// 处理退款点击
// const handleOpenRefund = () => {
//   if (!orderDetail.value) return
//   emit('open-refund', orderDetail.value)
//   drawerVisible.value = false
// }

// 重发通知
// const handleResendNotify = async () => {
//   if (!orderDetail.value) return
//
//   try {
//     await ElMessageBox.confirm(
//         '确认要重发支付通知吗？',
//         '重发通知确认',
//         {
//           confirmButtonText: '确认',
//           cancelButtonText: '取消',
//           type: 'warning'
//         }
//     )
//
//     const res = await resendNotifyApi({
//       mch_no: orderDetail.value.mch_no,
//       order_no: orderDetail.value.order_no
//     })
//
//     if (res.code === 0) {
//       ElMessage.success('通知重发成功')
//       emit('refresh')
//       // 重新加载详情
//       await loadOrderDetail()
//     } else {
//       ElMessage.error(res.message || '通知重发失败')
//     }
//   } catch (error) {
//     if (error !== 'cancel') {
//       console.error('重发通知失败:', error)
//       ElMessage.error('重发通知失败')
//     }
//   }
// }

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 格式化JSON
const formatJson = (jsonString) => {
  if (!jsonString) return ''
  try {
    const obj = JSON.parse(jsonString)
    return JSON.stringify(obj, null, 2)
  } catch {
    return jsonString
  }
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
    10: 'danger',
    11: 'danger',
  }
  return typeMap[status] || 'info'
}

// 获取支付方式文本
const getPaymentMethodText = (method) => {
  return method || '-'
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
// 获取支付产品文本
const getPaymentProductText = (product) => {
  return product || '-'
}
</script>

<style scoped>
.order-detail-drawer {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

.order-info-card {
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  border-radius: 10px;
  padding: 20px;
  border-left: 4px solid #409eff;
}

.amount-section {
  background: white;
  border-radius: 8px;
  padding: 15px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.info-card {
  background: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid #e4e7ed;
}

.info-label {
  display: inline-block;
  width: 100px;
  color: #666;
  font-size: 13px;
  flex-shrink: 0;
}

.info-value {
  flex: 1;
  color: #333;
  font-size: 14px;
  font-weight: 500;
  word-break: break-word;
}

/* 侧边栏标题样式 */
:deep(.el-drawer__header) {
  margin-bottom: 0;
  padding: 20px 20px 15px;
  border-bottom: 1px solid #e4e7ed;
}

:deep(.el-drawer__title) {
  font-weight: 600;
  color: #333;
}

:deep(.el-drawer__body) {
  padding: 0;
}

/* 滚动条样式 */
.order-detail-drawer::-webkit-scrollbar {
  width: 6px;
}

.order-detail-drawer::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.order-detail-drawer::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.order-detail-drawer::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
