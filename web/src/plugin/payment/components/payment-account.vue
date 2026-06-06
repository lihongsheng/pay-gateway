<!-- payment-account.vue -->
<template>
  <!-- 移除 el-drawer 包装，只保留内容 -->
  <div>
    <div class="payment-account-header">
      <el-tabs v-model="viewData.currentStep" type="card" @tab-click="handleTabClick" style="width: 80%;">
        <el-tab-pane :name="0" label="支付通道列表"></el-tab-pane>
        <el-tab-pane :name="1" label="支付通道配置"></el-tab-pane>
      </el-tabs>
    </div>

    <div class="payment-account-content">
      <div v-if="viewData.currentStep === 0">
        <!-- 支付通道列表内容 -->
        <el-table
            :data="paymentChannels"          style="width: 100%"
            border
            v-loading="loading"
        >
          <el-table-column prop="account_no" label="通道编号" min-width="150" />
          <el-table-column prop="name" label="通道名称" min-width="120" />
          <el-table-column prop="channelName" label="支付渠道" min-width="120">
            <template #default="{ row }">
              {{ row.channel_name }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '启用' : '停用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="验证状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.validate_status === 1 ? 'success' : 'danger'">
                {{ row.validate_status === 1 ? '已验证' : '未验证' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button
                  size="small"
                  type="primary"
                  link
                  @click="editChannel(row.id)"
              >
                编辑
              </el-button>
              <el-button
                  size="small"
                  type="primary"
                  link
                  @click="loadApplicationQrcode(row)"
              >
                测试二维码
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-if="viewData.currentStep === 1">
        <!-- 支付通道配置内容 -->
        <div class="channel-config-container">
        <div
              v-for="(config, index) in channelConfigs"
              :key="config.channel"
              class="channel-card"
        >
            <div class="card-header">
              <h3>{{ config.label }}</h3>
<!--              <span class="status">未开通</span>-->
              <el-tag :type="getCardStatusType(config.channel)">
                {{ getCardStatusTitle(config.channel)}}
              </el-tag>
            </div>
            <div class="card-body">
<!--              @click="openConfigForm(config)"-->
              <el-button
                  type="primary"
                  link
                  v-if="isShowCardButton(config.channel)"
                  @click="openConfigForm(config.channel)"
              >
                填写参数 >
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <el-drawer
        v-model="PaymentEditVisible"
        size="900px"
        direction="rtl"
        destroy-on-close
        :with-header="true"
        :z-index="10001"
        :title="PaymentEditTitle"
    >
      <PaymentEdit
          ref="paymentEditRef"
          @cancel="PaymentEditVisible = false"
          @save-success="handleSaveSuccess"
      />
    </el-drawer>


    <!-- 二维码弹窗 -->
    <el-dialog
        v-model="qrcodeDialogVisible"
        :title="'支付账户二维码-' + qrAccountNO + '(' +qrAccountName + ')'"
        width="600px"
        destroy-on-close
    >
      <div class="qrcode-dialog">
        <!-- 链接展示区域 -->
        <div class="mb-6">
          <div class="text-sm font-medium text-gray-700 mb-2">链接：</div>
          <div class="flex items-center">
            <el-input
                v-model="currentQrcodeLink"
                readonly
                class="flex-1 mr-2"
            >
              <template #append>
                <el-button @click="copyQrcodeLink">
                  <el-icon><DocumentCopy /></el-icon>
                </el-button>
              </template>
            </el-input>
          </div>
        </div>

        <!-- 二维码生成区域 -->
        <div class="mb-6">
          <div class="text-sm font-medium text-gray-700 mb-2">订单号：{{ qrOrderNO }}</div>
          <div class="flex flex-col items-center">
            <!-- 二维码容器 -->
            <div class="mb-4 p-4 bg-white border rounded-lg">
              <canvas ref="qrcodeCanvas" :width="qrcodeSize" :height="qrcodeSize"></canvas>
            </div>

            <!-- 二维码尺寸调节 -->
            <div class="w-full max-w-sm">
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm text-gray-600">二维码尺寸</span>
                <span class="text-sm font-medium">{{ qrcodeSize }}px</span>
              </div>
              <el-slider
                  v-model="qrcodeSize"
                  :min="100"
                  :max="500"
                  :step="10"
                  show-input
                  @change="generateQRCode"
              />
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex justify-between pt-4 border-t">
          <div>
            <el-button @click="downloadQRCode('png')" type="primary" plain>
              下载PNG
            </el-button>
            <el-button @click="downloadQRCode('jpg')" type="primary" plain class="ml-2">
              下载JPG
            </el-button>
          </div>
          <el-button @click="qrcodeDialogVisible = false">
            关闭
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
/**
 * 支付配置表单组件
 * @component PaymentAccount
 * @description 支付配置表单组件
 */
import {reactive, ref, watch, onMounted, nextTick} from "vue";
import {getAccountList, getChannelConfig, getChannelPayment, getChannelPaymentQrcode}  from "@/plugin/payment/api/payment_account.js";
import { ElMessage, ElMessageBox } from 'element-plus';
import PaymentEdit   from "@/plugin/payment/form/payment-edit.vue";
import {DocumentCopy} from "@element-plus/icons-vue";  // 添加路由导入
import QRCode from 'qrcode'
// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}
const props = defineProps({
  appNo: {
    type: String,
    default: ''
  },
  multiChannel: {
    type: Boolean,
    default: false
  }
})

// 步骤条切换
const viewData = reactive({
  currentStep: 0, // 当前标签页index
  btnLoading: false,
  appNo: props.appNo, // 应用appId
  open: true, // 抽屉开关
});

function handleTabClick(tab) {
  viewData.currentStep = Number(tab.paneName);
}

function show(appNo) {
  viewData.open = true;
  viewData.appNo = appNo;
}

function onClose() {
  viewData.open = false;
}

const PaymentEditVisible = ref(false);
const PaymentEditTitle = ref('');
const paymentEditRef =  ref(null)

const paymentChannels = ref([]);
const loading = ref(false);
const channelConfigs = ref([]);
const paymentProducts = ref([]);
// 获取支付通道列表
async function fetchPaymentChannels() {
  loading.value = true;
  try {
    const res = await getAccountList({ app_no: viewData.appNo });
    if (res.code === 0) {
      paymentChannels.value = res.data || [];
    } else {
      ElMessage.error('获取支付通道列表失败');
    }
  } catch (error) {
    console.error('获取支付通道列表失败:', error);
    ElMessage.error('获取支付通道列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取支付通道列表
async function fetchChannelConfig() {
  loading.value = true;
  try {
    const res = await getChannelConfig({ app_no: viewData.appNo });
    if (res.code === 0) {
      channelConfigs.value = res.data || [];
    } else {
      ElMessage.error('获取支付通道列表失败');
    }
  } catch (error) {
    console.error('获取支付通道列表失败:', error);
    ElMessage.error('获取支付通道列表失败');
  } finally {
    loading.value = false;
  }
}

// 获取支付通道列表
async function fetchPaymentProducts(channel) {
  loading.value = true;
  try {
    const res = await getChannelPayment({ channel:channel });
    if (res.code === 0) {
      paymentProducts.value = res.data || [];
    } else {
      ElMessage.error('获取产品列表失败');
    }
  } catch (error) {
    console.error('获取支付通道列表失败:', error);
    ElMessage.error('获取产品列表失败');
  } finally {
    loading.value = false;
  }
}

// 添加保存成功处理方法
const handleSaveSuccess = () => {
  // 关闭编辑抽屉
  PaymentEditVisible.value = false;
  // 刷新支付通道列表
  fetchPaymentChannels();
  // 如果需要同时刷新配置列表，也可以一并调用
  //fetchChannelConfig();
}

// 获取卡片状态类型
const getCardStatusType = (channel) => {
  // 检查 paymentChannels 中是否存在该 channel 的记录
  const channelRecord = paymentChannels.value.find(item => item.channel === channel);
  if (channelRecord) {
    return 'success'; // 启用状态
  } else {
    return 'danger'; // 未开通状态
  }
};

// 获取卡片状态标题
const getCardStatusTitle = (channel) => {
  // 检查 paymentChannels 中是否存在该 channel 的记录
  const channelRecord = paymentChannels.value.find(item => item.channel === channel);
  if (channelRecord) {
    return '已有历史配置';
  } else {
    return '未开通';
  }
};

const isShowCardButton = (channel) => {
  // 检查 paymentChannels 中是否存在该 channel 的记录
  const channelRecord = paymentChannels.value.find(item => item.channel === channel);
  if (channelRecord) {
    return  props.multiChannel
  }
  return true
}

// 给payment-edit使用，让其渲染
const channelData = reactive({
      id: 0,
      name: "",
      remark: "",
      app_no: props.appNo,
      channel: "",
      channel_name: "",
      status: 1, // 默认启用
      channel_options: {}, //
      payment_method_config: [], //
      payment_method: [], // 最终提交的支付方式数组
      channel_config: "",
})

// 重置 channelData
const resetChannelData = () => {
  Object.assign(channelData, {
    id: 0,
    name: "",
    remark: "",
    app_no: viewData.appNo,
    channel: "",
    channel_name: "",
    status: 1,
    channel_option: {},
    payment_method_config: [],
    payment_method: [],
    channel_config: "",
  })
}

// 打开配置表单
const openConfigForm = async (channel) => {
  try {
      // 清空之前的表单数据
      resetChannelData()
      await fetchPaymentProducts(channel)
      PaymentEditVisible.value = true
      // 检查是否有历史配置
      const historyConfig = channelConfigs.value.find(item => item.channel === channel)
      if (historyConfig) {
        Object.assign(channelData, {
          channel_option:historyConfig,
          payment_method_config:paymentProducts,
          channel: historyConfig.channel,
          channel_name: historyConfig.label,
          status: 1,
        })
        PaymentEditTitle.value =  historyConfig.label
      }
    // 确保组件已渲染后初始化表单
    await nextTick(() => {
      if (paymentEditRef.value && paymentEditRef.value.initForm) {
        paymentEditRef.value.initForm({...channelData})
      } else {
        console.error('paymentEditRef 或 initForm 不存在')
      }
    })
  } catch (error) {
    console.error('获取渠道配置失败:', error)
    ElMessage.error('获取渠道配置失败')
    PaymentEditTitle.value = ""
    PaymentEditVisible.value = false
  }
}

const editChannel = (id) => {
  try {
    // 检查是否有历史配置
    const historyConfig = paymentChannels.value.find(item => item.id === id)
    if (historyConfig) {

      Object.assign(channelData, {
        id: Number(historyConfig.id),
        name: historyConfig.name,
        channel: historyConfig.channel,
        channel_name: historyConfig.channel_name,
        remark: historyConfig.remark,
        app_no: historyConfig.app_no,
        status: Number(historyConfig.status),
        channel_option: historyConfig.channel_option,
        payment_method_config: historyConfig.payment_method_config,
        channel_config: historyConfig.channel_config,
        max_limit:  Number(historyConfig.max_limit),
      })
    }
    PaymentEditTitle.value = historyConfig.channel_name
    PaymentEditVisible.value = true
    // 确保组件已渲染后初始化表单
    nextTick(() => {
      if (paymentEditRef.value && paymentEditRef.value.initForm) {
        paymentEditRef.value.initForm({...channelData})
      } else {
        console.error('paymentEditRef 或 initForm 不存在')
      }
    })
  } catch (error) {
    console.error('获取渠道配置失败:', error)
    ElMessage.error('获取渠道配置失败')
    PaymentEditTitle.value = ""
    PaymentEditVisible.value = false
  }
  //paymentEditRef.value.open(channelData)
}
// 初始化
onMounted(() => {
  if (viewData.appNo) {
    fetchPaymentChannels();
    fetchChannelConfig();
  }
});
// 如果需要点击步骤标题跳转，可以添加监听
watch(() => viewData.currentStep, (newVal) => {
  // 步骤变化时的处理逻辑
});



// 二维码弹窗相关
const qrcodeDialogVisible = ref(false)
const currentQrcodeLink = ref('')
const qrcodeCanvas = ref(null)
const qrcodeSize = ref(200)
const qrOrderNO = ref('')
const qrAccountNO = ref('')
const qrAccountName = ref('')
// 生成二维码
const generateQRCode = () => {
  if (!currentQrcodeLink.value || !qrcodeCanvas.value) return

  try {
    QRCode.toCanvas(qrcodeCanvas.value, currentQrcodeLink.value, {
      width: qrcodeSize.value,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#ffffff'
      }
    })
  } catch (error) {
    console.error('生成二维码失败:', error)
    ElMessage.error('生成二维码失败')
  }
}

// 复制二维码链接
const copyQrcodeLink = () => {
  if (!currentQrcodeLink.value) return

  navigator.clipboard.writeText(currentQrcodeLink.value)
      .then(() => {
        ElMessage.success('链接已复制到剪贴板')
      })
      .catch(err => {
        console.error('复制失败:', err)
        ElMessage.error('复制失败')
      })
}


// 下载二维码
const downloadQRCode = async (format = 'png') => {
  if (!qrcodeCanvas.value) return

  try {
    const canvas = qrcodeCanvas.value
    const link = document.createElement('a')
    const filename = `app_qrcode_${Date.now()}.${format}`

    link.download = filename
    link.href = canvas.toDataURL(`image/${format}`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    ElMessage.success(`二维码已下载为${filename}`)
  } catch (error) {
    console.error('下载二维码失败:', error)
    ElMessage.error('下载二维码失败')
  }
}

// 修改 loadApplicationQrcode 函数
const loadApplicationQrcode = async (row) => {
  try {
    const res = await getChannelPaymentQrcode(row.app_no, row.account_no)
    if (res.code === 0) {
      qrAccountNO.value = row.account_no
      qrAccountName.value = row.name
      const qrcodeUrl = res.data?.qr_link || ""
      qrOrderNO.value = res.data?.order_no
      // 显示二维码弹窗
      currentQrcodeLink.value = qrcodeUrl
      qrcodeDialogVisible.value = true
      nextTick(() => {
        generateQRCode()
      })
    } else {
      ElMessage.error('获取二维码失败:' + res.msg)
    }
  } catch (error) {
    console.error('获取应用详情失败:', error)
    ElMessage.error('获取应用二维码失败')
  }
}

</script>


<style  scoped>
.title {
  font-size: 16px;
  font-family:
    PingFang SC,
    PingFang SC-Bold;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: 1px;
}
.payment-account-header {
  padding: 10px 0;
}

.payment-account-content {
  padding: 10px 0;
}

.config-form {
  padding: 10px 0;
}

/* 补充卡片样式 */
.channel-config-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding: 10px 0;
}

.channel-card {
  width: 280px;
  border-radius: 8px;
  border: 1px solid #ebeef5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  overflow: hidden;
  background-color: white;
}

.channel-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.card-header {
  padding: 16px 20px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.status {
  font-size: 12px;
  color: #f56c6c;
  padding: 2px 8px;
  background-color: #fef0f0;
  border-radius: 12px;
}

.card-body {
  padding: 16px 20px;
  text-align: center;
}

.card-body .el-button {
  padding: 8px 16px;
  font-size: 14px;
  border-radius: 4px;
  margin-top: 12px;
}
</style>
