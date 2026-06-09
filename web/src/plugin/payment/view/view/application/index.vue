<template>
  <div class="gva-container">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-lg font-bold">应用管理</h2>
    </div>

    <div class="search-box bg-white p-4 rounded-md shadow mb-4">
      <el-form :inline="true" :model="searchInfo" class="search-form">
        <el-row :gutter="10">
          <el-col :span="6">
            <el-form-item label="应用编号">
              <el-input
                v-model="searchInfo.app_no"
                placeholder="请输入应用编号"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="应用名称">
              <el-input
                v-model="searchInfo.name"
                placeholder="请输入应用名称"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="商户编号">
              <el-input
                v-model="searchInfo.mch_no"
                placeholder="请输入商户编号"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="状态">
              <el-select
                v-model="searchInfo.status"
                placeholder="请选择状态"
                clearable
                class="w-full"
              >
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="10">
          <el-col :span="24" class="text-right">
            <el-button type="primary" @click="onSubmit">查询</el-button>
            <el-button @click="onReset">重置</el-button>
            <el-button type="primary" @click="addApplication">新增应用</el-button>
          </el-col>
        </el-row>
      </el-form>
    </div>

    <div class="table-container bg-white p-4 rounded-md shadow">
      <el-table
        :data="tableData"
        style="width: 100%"
        row-key="id"
        v-loading="loading"
        border
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="app_no" label="应用编号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="name" label="应用名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="mch_no" label="商户编号" min-width="150" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              link
              @click="getApplication(row.id)"
            >
              查看
            </el-button>
            <el-button
              size="small"
              type="primary"
              link
              @click="editApplication(row)"
            >
              编辑
            </el-button>
            <el-button
              size="small"
              :type="row.status === 1 ? 'danger' : 'success'"
              link
              @click="changeStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-popover v-model:visible="row.visible" placement="top" width="160">
              <p>确定要删除此应用吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button size="small" @click="row.visible = false">取消</el-button>
                <el-button size="small" type="primary" @click="deleteApplicationFunc(row)">确定</el-button>
              </div>
              <template #reference>
                <el-button size="small" type="danger" link>删除</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>

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

    <!-- 新增/编辑弹窗 -->
    <!-- 在模板中替换 Dialog 为 Drawer -->
    <el-drawer
        v-model="drawerVisible"
        :title="dialogTitle"
        size="700px"
        direction="rtl"
        destroy-on-close
        :with-header="true"
        :z-index="1000"
    >
      <ApplicationForm
          v-model="formData"
          ref="formRef"
      />
      <template #footer>
        <div class="flex justify-end gap-2 mt-4">
          <el-button @click="closeDrawer">取消</el-button>
          <el-button type="primary" @click="enterDialog">确认</el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
/**
 * 应用管理页面
 * @component ApplicationView
 * @description 应用管理功能页面
 */

import { ref, reactive, onMounted } from 'vue'
import {
  getApplicationList,
  createApplication,
  updateApplication,
  deleteApplication,
  updateApplicationStatus
} from '@/plugin/payment/api/application.js'
import { ElMessage, ElMessageBox } from 'element-plus'
import ApplicationForm from "@/plugin/payment/form/application-form.vue";


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

// 表单验证规则
const rules = reactive({
  app_name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' }
  ],
  mch_no: [
    { required: true, message: '请输入商户编号', trigger: 'blur' }
  ],
  notifyUrl: [
    { required: true, message: '请输入通知URL', trigger: 'blur' }
  ]
})

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
    const res = await getApplicationList(params)
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
  formData.value = {
    id: 0,
    app_name: '',
    mch_no: '',
    status: 1,
    is_customer_domain: 0,
    desc: '',
    payment_title: '',
    pay_icon: '',
    customer_domain: '',
    proxy_host: '',
    proxy_port: 3221,
    proxy_user: '',
    proxy_pwd: ''
  }
  drawerVisible.value = true  // 打开侧边栏
}

// 编辑应用
const editApplication = (row) => {
  dialogTitle.value = '编辑应用'
  formData.value = { ...row }
  dialogVisible.value = true
}

// 修改查看应用函数
const getApplication = async (id) => {
  try {
    const res = await getApplicationList({ id })
    if (res.code === 0 && res.data.list && res.data.list.length > 0) {
      const detail = res.data.list[0]
      editApplication(detail)
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

    const res = await updateApplicationStatus({
      id: row.id,
      status: row.status === 1 ? 0 : 1
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
  dialogVisible.value = false
  formData.value = {}
}
// 修改 enterDialog 函数
const enterDialog = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      let res
      if (formData.value.id) {
        // 更新应用
        res = await updateApplication(formData.value)
      } else {
        // 创建应用
        res = await createApplication(formData.value)
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
  })
}

// 初始化
onMounted(() => {
  getApplicationListFunc()
})
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
