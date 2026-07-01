<template>
  <div class="page-wrap">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchInfo" class="search-form">
        <el-form-item label="应用编号">
          <el-input v-model="searchInfo.app_no" placeholder="请输入应用编号" clearable style="width:180px" @keyup.enter="onSubmit" @clear="onSubmit">
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item label="应用名称">
          <el-input v-model="searchInfo.app_name" placeholder="请输入应用名称" clearable style="width:180px" @keyup.enter="onSubmit" @clear="onSubmit" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable style="width:120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit"><el-icon><Search /></el-icon>搜索</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据栏 -->
    <el-card shadow="never" class="table-card">
      <div class="table-toolbar">
        <el-button v-permission="'application:add'" type="primary" @click="addApplication">
          <el-icon><Plus /></el-icon>新增应用
        </el-button>
      </div>
      <el-table :data="tableData" row-key="id" v-loading="loading" stripe border>
        <el-table-column prop="app_no" label="应用编号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="app_name" label="应用名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="mch_no" label="商户编号" min-width="150" show-overflow-tooltip />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small" effect="dark">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170" show-overflow-tooltip>
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="loadApplicationQrcode(row.app_no, row.app_name)">二维码</el-button>
            <el-button size="small" type="primary" link @click="getApplication(row.app_no)">查看</el-button>
            <el-button v-permission="'application:edit'" size="small" type="primary" link @click="editApplication(row)">编辑</el-button>
            <el-button v-permission="'application:status'" size="small" :type="row.status === 1 ? 'danger' : 'success'" link @click="changeStatus(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="primary" link @click="showPaymentAccount(row.app_no, row.app_name, row.multi_channel)">支付配置</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination background layout="total,sizes,prev,pager,next,jumper"
                       :total="total" v-model:current-page="page" v-model:page-size="pageSize"
                       :page-sizes="[10, 20, 50, 100]" @current-change="handleCurrentChange" @size-change="handleSizeChange" />
      </div>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <!-- 在模板中替换 Dialog 为 Drawer -->
    <el-drawer
        v-model="drawerVisible"
        :title="dialogTitle"
        size="900px"
        direction="rtl"
        destroy-on-close
        :with-header="true"
        :z-index="1000"
    >
      <ApplicationForm
          v-model="formData"
          :viewMode="viewMode"
          ref="formRef"
      />
      <template #footer v-if="isViewMode">
        <div class="flex justify-end gap-2 mt-4">
          <el-button @click="closeDrawer" >取消</el-button>
          <el-button type="primary" @click="enterDialog">确认</el-button>
        </div>
      </template>
    </el-drawer>
    <el-drawer
        v-model="paymentAccountVisible"
        size="900px"
        direction="rtl"
        destroy-on-close
        :with-header="true"
        :z-index="1000"
        :title="paymentAccountTitle"
    >
      <PaymentAccount
          ref="PaymentAccountRef"
          :appNo="currentAppNo"
          :multiChannel="multiChannel"
      />
    </el-drawer>

    <!-- 二维码弹窗 -->
    <el-dialog
        v-model="qrcodeDialogVisible"
        :title="'支付账户二维码-' + qrAppNo + '(' + qrAppName + ')'"
        width="600px"
        destroy-on-close
    >
      <div class="qrcode-dialog">
        <!-- 链接展示区域 -->
        <div class="mb-6">
          <div class="text-sm font-medium text-gray-700 mb-2">应用链接：</div>
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
          <div class="text-sm font-medium text-gray-700 mb-2">二维码：</div>
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
 * 应用管理页面
 * @component ApplicationView
 * @description 应用管理功能页面
 */

import { ref, reactive, onMounted, nextTick } from 'vue'
import { appList, appCreate, appUpdate, appStatus, appDetail, appQrcode } from '@/api/payment'
import { ElMessage, ElMessageBox } from 'element-plus'
import ApplicationForm from "../../form/application-form-mch.vue"
import PaymentAccount  from "../../components/payment-account.vue"
import { useRoute } from 'vue-router'
import { DocumentCopy, Search, Plus } from '@element-plus/icons-vue'
import QRCode from 'qrcode'

const route = useRoute()  // 添加路由实例

// 渠道配置
const paymentAccountVisible = ref(false)
const multiChannel = ref(false)
const currentAppNo = ref('');
const paymentAccountTitle = ref('')

function showPaymentAccount(appNo, title, multiChannelValue) {
  paymentAccountVisible.value = true
  // 支付参数配置
  currentAppNo.value = appNo
  multiChannel.value = multiChannelValue === 1
  paymentAccountTitle.value = title + "-支付配置"
}


// 响应式数据
const loading = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 搜索条件
const searchInfo = reactive({
  app_no: '',
  app_name: '',
  mch_no: '',
  status: null
})

// 将 dialog 相关变量替换为 drawer 相关变量
const drawerVisible = ref(false)  // 侧边栏显示状态
const dialogTitle = ref('')
const formData = ref({})
const formRef = ref(null)


// 二维码弹窗相关
const qrcodeDialogVisible = ref(false)
const currentQrcodeLink = ref('')
const qrcodeCanvas = ref(null)
const qrcodeSize = ref(200)
const qrAppNo = ref('')
const qrAppName = ref('')

// 在 onMounted 中初始化
onMounted(() => {
  // 检查 URL 参数中是否有 mch_no
  if (route.query.mch_no) {
    searchInfo.mch_no = route.query.mch_no
  }
  getApplicationListFunc()
})

// 表单验证规则
const rules = reactive({
  app_name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' }
  ],
  mch_no: [
    { required: true, message: '请输入商户编号', trigger: 'blur' }
  ]
})
// 在 index.vue 的 script 部分添加查看模式标志
const viewMode = ref("")  // 添加查看模式标志
const isViewMode = ref(false)

// 添加 closeDrawer 函数
const closeDrawer = () => {
  drawerVisible.value = false
  // 重置表单
  if (formRef.value) {
    formRef.value.resetFields?.() // 若表单有重置方法则调用
  }
  formData.value = {}
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取应用列表
const getApplicationListFunc = async () => {
  loading.value = true
  try {
    const params = {
      ...searchInfo,
      page: page.value,
      page_size: pageSize.value
    }
    const res = await appList(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取应用列表失败:', error)
    ElMessage.error('获取应用列表失败')
  } finally {
    loading.value = false
  }
}

// 分页处理
const handleCurrentChange = (val) => {
  page.value = val
  getApplicationListFunc()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getApplicationListFunc()
}

// 搜索和重置
const onSubmit = () => {
  page.value = 1
  getApplicationListFunc()
}

const onReset = () => {
  Object.keys(searchInfo).forEach(key => {
    searchInfo[key] = null
  })
  page.value = 1
  getApplicationListFunc()
}

// 修改新增应用函数
const addApplication = () => {
  dialogTitle.value = '新增应用'
  viewMode.value = "add"  // 设置为编辑模式
  isViewMode.value = true
  formData.value = {
    id: 0,
    app_name: '',
    mch_no: route.query.mch_no || '',
    status: 1,
    is_customer_domain: 0,
    desc: '',
    payment_title: '',
    pay_icon: '',
    customer_domain: '',
    proxy_host: '',
    proxy_port: 3221,
    proxy_user: '',
    proxy_pwd: '',
    multi_channel: 0,
  }
  drawerVisible.value = true  // 打开侧边栏
}
/**
 * 生成指定长度的随机字符串 (大小写字母 + 数字)
 * @param {Number} len 长度，默认32位
 * @returns {String} 随机字符串
 */
function randomString(len = 32) {
  // 字符库：大写字母+小写字母+0-9数字，无特殊字符，完美匹配你的应用秘钥secret需求
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charsLen = chars.length;
  let randomStr = '';
  for (let i = 0; i < len; i++) {
    // 随机取字符库的下标
    randomStr += chars.charAt(Math.floor(Math.random() * charsLen));
  }
  return randomStr;
}
// 编辑应用
const editApplication = (row) => {
  dialogTitle.value = '编辑应用'
  viewMode.value = "edit"  // 设置为编辑模式
  isViewMode.value = true
  formData.value = { ...row }
  drawerVisible.value = true
}

// 修改查看应用函数
const getApplication = async (appNo) => {
  try {
    const res = await appDetail({ app_no: appNo })
    if (res.code === 0) {
      const detail = res.data
      viewMode.value = "view"  // 设置为查看模式
      isViewMode.value = false
      dialogTitle.value = '查看应用'
      formData.value = { ...detail }
      drawerVisible.value = true
    }
  } catch (error) {
    console.error('获取应用详情失败:', error)
    ElMessage.error('获取应用详情失败')
  }
}

// 状态切换
const changeStatus = async (row) => {
  try {
    await ElMessageBox.confirm(
        `确定要${row.status === 1 ? '禁用' : '启用'}此应用吗？`,
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )

    const res = await appStatus({
      id: row.id,
      status: row.status === 1 ? 2 : 1
    })

    if (res.code === 0) {
      ElMessage.success(`${row.status === 1 ? '禁用' : '启用'}成功`)
      getApplicationListFunc()
    } else {
      ElMessage.error(`${row.status === 1 ? '禁用' : '启用'}失败`)
    }
  } catch (error) {
    console.log('操作已取消')
  }
}

// 关闭弹窗
const closeDialog = () => {
  drawerVisible.value = false
  formData.value = {}
}
// 修改 enterDialog 函数
const enterDialog = async () => {
  // 验证表单
  const isValid = await formRef.value.validate()
  if (!isValid) return

  try {
    let res
    if (formData.value.id) {
      // 更新应用
      res = await appUpdate(formData.value)
    } else {
      // 创建应用
      res = await appCreate(formData.value)
    }

    if (res.code === 0) {
      ElMessage.success(formData.value.id ? '更新成功' : '创建成功')
      closeDrawer()
      getApplicationListFunc()
    } else {
      ElMessage.error(formData.value.id ? '更新失败' : '创建失败')
    }
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error('操作失败')
  }
}

// 初始化
onMounted(() => {
  getApplicationListFunc()
})


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

// 打开二维码链接
const openQrcodeLink = () => {
  if (!currentQrcodeLink.value) return
  window.open(currentQrcodeLink.value, '_blank')
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
const loadApplicationQrcode = async (appNo, appName) => {
  try {
    const res = await appQrcode({ app_no: appNo })
    if (res.code === 0) {
      qrAppNo.value = appNo
      qrAppName.value = appName
      const qrcodeUrl = res.data
      // 显示二维码弹窗
      currentQrcodeLink.value = qrcodeUrl
      qrcodeDialogVisible.value = true
      nextTick(() => {
        generateQRCode()
      })
    } else {
      ElMessage.error('获取应用二维码失败' + res.msg)
    }
  } catch (error) {
    console.error('获取应用详情失败:', error)
    ElMessage.error('获取应用二维码失败')
  }
}

</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.search-card { }
.search-form { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.search-form .el-form-item { margin-bottom: 0; }
.table-card { }
.table-toolbar { margin-bottom: 12px; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
