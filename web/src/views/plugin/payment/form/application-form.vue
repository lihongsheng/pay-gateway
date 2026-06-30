<template>
  <el-form
      :model="localFormData"
      :rules="rules"
      ref="formRef"
      label-width="120px"
      @submit.prevent
      class="application-form"
  >
    <el-scrollbar height="calc(100vh - 120px)">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item label="应用名称" prop="app_name">
            <el-input
                v-model="localFormData.app_name"
                placeholder="请输入应用名称"
                :disabled="(isReadOnly)"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20" v-if="isReadOnly">
        <el-col :span="24">
          <el-form-item label="应用编码" prop="app_no">
            <el-input
                v-model="localFormData.app_no"
                placeholder=""
                disabled
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item label="商户编号" prop="mch_no">
            <el-select
                v-model="localFormData.mch_no"
                filterable
                placeholder="请选择商户"
                clearable
                remote
                :disabled="!!localFormData.id"
                style="width: 100%"
                :remote-method="remoteSearchVenuer"
                :loading="venuerLoading"
            >
              <el-option
                  v-for="item in venuerList"
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
      </el-row>

      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item label="应用描述" prop="desc">
            <el-input
                v-model="localFormData.desc"
                type="textarea"
                :rows="3"
                placeholder="请输入应用描述"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20" v-if="isReadOnly">
        <el-col :span="24">
          <el-form-item label="应用秘钥" prop="payment_secret">
            <el-input
                v-model="localFormData.secret"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item label="支付标题" prop="payment_title">
            <el-input
                v-model="localFormData.payment_title"
                placeholder="请输入支付标题"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
      </el-row>

<!--      <el-row :gutter="20">-->
<!--        <el-col :span="24">-->
<!--          <el-form-item label="支付图标" prop="pay_icon">-->
<!--            <el-upload-->
<!--                v-model="localFormData.pay_icon"-->
<!--                action="/api/upload"-->
<!--                list-type="picture"-->
<!--                :limit="1"-->
<!--                :file-list="getFileList()"-->
<!--            >-->
<!--              <el-button type="primary">点击上传</el-button>-->
<!--            </el-upload>-->
<!--          </el-form-item>-->
<!--        </el-col>-->
<!--      </el-row>-->

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="应用状态" prop="status">
            <el-radio-group v-model="localFormData.status"  :disabled="isReadOnly">
              <el-radio :label="1">启用</el-radio>
              <el-radio :label="2">禁用</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="支持多微信等" prop="is_customer_domain" >
            <el-radio-group v-model="localFormData.multi_channel" :disabled="isReadOnly">
              <el-radio :label="0">否</el-radio>
              <el-radio :label="1">是</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="代理模式" prop="is_customer_domain" >
          <el-radio-group v-model="localFormData.is_customer_domain" :disabled="isReadOnly">
            <el-radio :label="0">否</el-radio>
            <el-radio :label="1">是</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-col>
      </el-row>

      <el-row :gutter="20" v-if="localFormData.is_customer_domain === 1">
        <el-col :span="24">
          <el-form-item label="域名地址" prop="customer_domain">
            <el-input
                v-model="localFormData.customer_domain"
                placeholder="请输入合法域名（如：https://demo.com）"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20" v-if="localFormData.is_customer_domain === 1">
        <el-col :span="24">
          <el-form-item label="代理地址" prop="proxy_host">
            <el-input
                v-model="localFormData.proxy_host"
                placeholder="请输入代理IP/域名（如：192.168.1.1 或 proxy.com:8080）"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20" v-if="localFormData.is_customer_domain === 1">
        <el-col :span="12">
          <el-form-item label="代理端口" prop="proxy_port">
            <el-input
                v-model="localFormData.proxy_port"
                type="number"
                placeholder="请输入代理端口号"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="代理用户名" prop="proxy_user">
            <el-input v-model="localFormData.proxy_user" placeholder="请输入代理用户名"  :disabled="isReadOnly"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20" v-if="localFormData.is_customer_domain === 1">
        <el-col :span="24">
          <el-form-item label="代理密码" prop="proxy_pwd">
            <el-input
                v-model="localFormData.proxy_pwd"
                type="password"
                placeholder="请输入代理密码"
                :disabled="isReadOnly"
            />
          </el-form-item>
        </el-col>
      </el-row>
    </el-scrollbar>
  </el-form>
</template>

<script setup>
/**
 * 应用表单组件
 * @component ApplicationForm
 * @description 应用信息表单组件，确保id为int类型提交
 */

import {ref, reactive, watch, onMounted, computed} from 'vue'
import { mchList } from '@/api/system'
import { useRoute } from 'vue-router'  // 添加路由导入
const route = useRoute()  // 添加路由实例

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({})
  },
  viewMode: {
    type: String,
    default: 'view'
  }
})

const emit = defineEmits(['update:modelValue', 'submit', 'validate-success', 'validate-fail'])

// 本地表单数据
const localFormData = reactive({
  id: 0,
  app_name: '',
  secret: '',
  mch_no: '',
  desc: '',
  payment_title: '',
  pay_icon: '',
  status: 1,
  is_customer_domain: 0,
  customer_domain: '',
  proxy_host: '',
  proxy_port: 3221,
  proxy_user: '',
  proxy_pwd: '',
})

// 计算属性：是否为只读模式
// 计算属性：是否为只读模式
const isReadOnly = computed(() => props.viewMode === 'view')
// 表单引用
const formRef = ref(null)

// 表单验证规则
const rules = reactive({
  id: [
    { type: 'number', message: 'ID必须为数字类型', trigger: 'blur' }
  ],
  app_name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' },
    { min: 2, message: '应用名称至少2个字符', trigger: 'blur' }
  ],
  mch_no: [
    { required: true, message: '请输入公司编号', trigger: 'blur' },
    { min: 2, message: '公司编号至少2个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择应用状态', trigger: 'change' },
    { type: 'number', message: '状态值必须为数字', trigger: 'change' }
  ],
  customer_domain: [
    {
      required: true,
      message: '自定义域名不能为空',
      trigger: 'blur',
      validator: (rule, value) => {
        if (localFormData.is_customer_domain === 1 && !value) {
          return Promise.reject(new Error('自定义域名不能为空'));
        }
        return Promise.resolve();
      }
    },

    { pattern: /^(https?:\/\/)?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/, message: '请输入合法的域名', trigger: 'blur' }
  ],
  proxy_host: [
    {
      pattern: /^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})|([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})(:\d{1,5})?$/,
      message: '请输入合法的IP或域名（支持端口）',
      trigger: 'blur'
    }
  ],
  // notifyUrl: [
  //   { required: true, message: '请输入通知URL', trigger: 'blur' }
  // ]
})

// 监听父组件传入值，同步到本地数据
watch(
    () => props.modelValue,
    (newVal) => {
      if (newVal) {
        Object.assign(localFormData, {
          ...newVal,
          id: newVal.id ? Number(newVal.id) : 0,
          status: newVal.status !== undefined ? Number(newVal.status) : 1,
          is_customer_domain: newVal.is_customer_domain !== undefined ? Number(newVal.is_customer_domain) : 0,
          proxy_port: newVal.proxy_port ? Number(newVal.proxy_port) : 0
        })
      }
    },
    { immediate: true, deep: true }
)

// 监听本地数据变化，同步到父组件
watch(
    () => localFormData,
    (newVal) => {
      const submitData = {
        ...newVal,
        id: Number(newVal.id) || 0,
        proxy_port: Number(newVal.proxy_port) || 0
      }
      emit('update:modelValue', submitData)
    },
    { deep: true }
)

// 验证表单
const validate = async () => {
  if (!formRef.value) {
    emit('validate-fail', new Error('表单引用不存在'))
    return false
  }
  try {
    localFormData.id = Number(localFormData.id) || 0
    await formRef.value.validate()
    emit('validate-success')
    return true
  } catch (error) {
    emit('validate-fail', error)
    return false
  }
}

// 重置表单方法
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(localFormData, {
    id: 0,
    app_name: '',
    secret: '',
    mch_no: route.query.mch_no || '',
    desc: '',
    payment_title: '',
    pay_icon: '',
    status: 1,
    is_customer_domain: 0,
    customer_domain: '',
    proxy_host: '',
    proxy_port: 0,
    proxy_user: '',
    proxy_pwd: '',
    multi_channel: 0,
  })
}
const venuerLoading = ref(false)
const venuerList = ref([])
// 组件挂载时获取商户列表
onMounted(() => {
  const res =  mchList({mch_no: localFormData.mch_no, page:1, page_size:10})
  if (res.code === 0) {
    venuerList.value = res.data.list || []
  } else {
    venuerList.value = []
  }
})

// 远程搜索场地管理方
const remoteSearchVenuer = async (query) => {
  venuerLoading.value = true
  try {
    if (route.query.mch_no !== '') {
       const res = await mchList({mch_no: route.query.mch_no, page: 1, page_size: 10})
      if (res.code === 0) {
        venuerList.value = res.data.list || []
      } else {
        venuerList.value = []
      }
    } else {
      const res = await mchList({mch_name: query, page: 1, page_size: 10})
      if (res.code === 0) {
        venuerList.value = res.data.list || []
      } else {
        venuerList.value = []
      }
    }

  }finally {
    venuerLoading.value = false
  }
}

// 处理支付图标回显
const getFileList = () => {
  return localFormData.pay_icon ? [{url: localFormData.pay_icon, name: 'pay_icon'}] : [];
}

// 暴露方法给父组件
defineExpose({
  validate,
  resetForm,
  // 暴露获取格式化数据的方法
  getFormData: () => ({
    ...localFormData,
    id: Number(localFormData.id) || 0,
    proxy_port: Number(localFormData.proxy_port) || 0
  }),
})
</script>

<style scoped>
.application-form {
  height: 100%;
  overflow-y: auto;  /* 侧边栏内容过长时滚动 */
  padding-right: 10px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

/* 适配textarea在侧边栏的样式 */
:deep(.el-textarea__inner) {
  resize: vertical;
  min-height: 60px;
}
</style>
