<template>
  <div class="page-wrap">
    <!-- 应用选择栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" class="search-form">
        <el-form-item label="应用">
          <el-select
            v-model="selectedAppNo"
            filterable
            placeholder="请选择应用"
            clearable
            remote
            :remote-method="remoteSearchApp"
            :loading="appLoading"
            style="width: 320px"
            @change="handleAppChange"
          >
            <el-option
              v-for="item in appOptions"
              :key="item.app_no"
              :label="item.app_no + ' (' + item.app_name + ')'"
              :value="item.app_no"
            >
              <span>{{ item.app_no }}</span>
              <span style="float: right; color: #8492a6; font-size: 13px">{{ item.app_name }}</span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 支付渠道组件 -->
    <PaymentAccount
      v-if="selectedAppNo"
      :appNo="selectedAppNo"
      :multiChannel="isMultiChannel"
    />
    <el-empty v-else description="请先选择应用" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { appList } from '@/api/payment'
import { ElMessage } from 'element-plus'
import PaymentAccount from '../../components/payment-account.vue'

const selectedAppNo = ref('')
const isMultiChannel = ref(false)
const appLoading = ref(false)
const appOptions = ref([])

const remoteSearchApp = async (query) => {
  appLoading.value = true
  try {
    const res = await appList({ name: query, page: 1, page_size: 20 })
    if (res.code === 0) {
      appOptions.value = res.data?.list || []
    } else {
      appOptions.value = []
    }
  } catch (error) {
    console.error('搜索应用失败:', error)
    ElMessage.error('搜索应用失败')
  } finally {
    appLoading.value = false
  }
}

const handleAppChange = (val) => {
  if (val) {
    const app = appOptions.value.find(item => item.app_no === val)
    isMultiChannel.value = app?.multi_channel === 1
  }
}

onMounted(() => {
  remoteSearchApp('')
})
</script>

<style scoped>
.page-wrap { display: flex; flex-direction: column; gap: 12px; }
.search-card { }
.search-form { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.search-form .el-form-item { margin-bottom: 0; }
</style>
