<!-- payment-edit.vue -->
<template>
  <div class="payment-edit-container">
    <el-form ref="formRef" :model="formData" label-width="120px" :rules="formRules">

      <el-form-item label="通道名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入通道名称" />
      </el-form-item>

      <el-form-item label="支付渠道" v-if="formData.channel_name">
        <el-input v-model="formData.channel_name" disabled />
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input v-model="formData.remark" type="textarea" placeholder="请输入备注" />
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="formData.status">
          <el-radio :label="1">启用</el-radio>
          <el-radio :label="2">停用</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="支付限制" prop="max_limit">
        <el-input-number v-model="formData.max_limit" placeholder="单日最大支付笔数" />
      </el-form-item>

      <!-- 动态渲染支付配置表单（基于 channel_option 的 schema） -->
      <el-form-item label="支付配置" prop="channel_config" class="full-width-config">
        <div class="config-form-wrapper" v-if="configSchema" style="min-width: 100%;">
          <el-form label-width="150px" class="config-sub-form">
            <schema-node
                :schema="configSchema"
                :path="[]"
                :root="channelConfigForm"
                :errors="dynamicErrors"
                :on-update="updateConfigField"
            />
          </el-form>
        </div>
      </el-form-item>

      <!-- 支付方式选择 -->
      <el-form-item label="支付方式" prop="payment_method" class="full-width-config">
        <div class="payment-method-wrapper" v-if="paymentMethodConfig.length" style="min-width: 100%;">
          <div v-for="method in paymentMethodConfig" :key="method.method" class="method-card">
            <h4>{{ method.label }}</h4>
            <el-checkbox-group v-model="selectedPaymentMethods" @change="handlePaymentMethodChange">
              <el-checkbox
                  v-for="product in method.product"
                  :key="product.product"
                  :label="`${method.method}_${product.product}`"
                  :checked="product.used"
              >{{ product.label }}</el-checkbox>
            </el-checkbox-group>
          </div>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">保存</el-button>
        <el-button @click="handleCancel">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { savePaymentAccount } from '@/plugin/payment/api/payment_account.js'
import SchemaNode from './SchemaNode.vue'

const emit = defineEmits(['saveSuccess', 'cancel'])

// ============================================================
// 工具函数（schema 操作）
// ============================================================

function getByPath(root, path) {
  return path.reduce((obj, key) => (obj == null ? obj : obj[key]), root)
}

function setByPath(root, path, value) {
  let obj = root
  for (let i = 0; i < path.length - 1; i++) {
    const key = path[i]
    if (obj[key] === undefined || obj[key] === null || typeof obj[key] !== 'object') {
      const nextKey = path[i + 1]
      obj[key] = typeof nextKey === 'number' ? [] : {}
    }
    obj = obj[key]
  }
  obj[path[path.length - 1]] = value
}

function orderedEntries(schema) {
  if (!schema || !schema.properties) return []
  const properties = schema.properties
  const order = schema['x-order'] || Object.keys(properties)
  const names = [...order, ...Object.keys(properties).filter(k => !order.includes(k))]
  return names.filter(k => properties[k]).map(k => ({ key: k, schema: properties[k] }))
}

function defaultValueFor(schema) {
  if (!schema) return ''
  if (schema.default !== undefined) return JSON.parse(JSON.stringify(schema.default))
  if (schema.type === 'object') return {}
  if (schema.type === 'array') return []
  if (schema.type === 'boolean') return false
  if (schema.type === 'integer' || schema.type === 'number') return undefined
  return ''
}

function initDefaults(schema, root, path = []) {
  if (!schema) return
  if (schema.type === 'object') {
    if (path.length > 0) {
      const cur = getByPath(root, path)
      if (cur === undefined || cur === null || typeof cur !== 'object' || Array.isArray(cur))
        setByPath(root, path, {})
    }
    for (const item of orderedEntries(schema)) initDefaults(item.schema, root, [...path, item.key])
    return
  }
  if (schema.type === 'array') {
    const cur = getByPath(root, path)
    if (!Array.isArray(cur)) setByPath(root, path, Array.isArray(schema.default) ? [...schema.default] : [])
    return
  }
  setByPath(root, path, defaultValueFor(schema))
}

function clearObject(obj) {
  if (obj && typeof obj === 'object') Object.keys(obj).forEach(k => delete obj[k])
}

// ============================================================
// Schema 兼容函数（legacy Options → 自动生成 schema）
// ============================================================

function legacyType(t) {
  if (t === 'Int') return 'integer'
  if (t === 'Bool') return 'boolean'
  if (t === 'Array') return 'array'
  if (t === 'Object') return 'object'
  return 'string'
}

function castDefault(type, value) {
  if (value === undefined || value === null || value === '') return undefined
  if (type === 'Int') return Number(value)
  if (type === 'Bool') return value === true || value === 'true'
  if (type === 'Array') return Array.isArray(value) ? value : []
  return value
}

function castValue(type, value) {
  if (type === 'Int') return Number(value)
  if (type === 'Bool') return value === true || value === 'true'
  return value
}

function legacyOptionToNode(option) {
  if (Array.isArray(option.children) && option.children.length > 0) {
    const obj = buildLegacyObject(option.children)
    obj.title = option.label
    obj.description = option.description || undefined
    return obj
  }
  if (option.type === 'Array' && Array.isArray(option.item_options) && option.item_options.length > 0) {
    const itemSchema = buildLegacyObject(option.item_options)
    itemSchema.title = option.item_label || `${option.label}项`
    return {
      type: 'array',
      title: option.label,
      description: option.description || undefined,
      items: itemSchema,
    }
  }
  const node = {
    type: legacyType(option.type),
    title: option.label,
    description: option.description || undefined,
    pattern: option.validate_reg || undefined,
    default: castDefault(option.type, option.default),
    'x-ui': { input_type: option.input_type, validate_type: option.validate_type },
  }
  if (Array.isArray(option.values) && option.values.length > 0) {
    node.oneOf = option.values.map(v => ({ const: castValue(v.type, v.value), title: v.label }))
  }
  if (node.type === 'array') node.items = { type: 'string' }
  return node
}

function buildLegacyObject(options) {
  const properties = {}
  const required = []
  const order = []
  for (const o of options) {
    if (!o.name) continue
    order.push(o.name)
    if (o.require) required.push(o.name)
    properties[o.name] = legacyOptionToNode(o)
  }
  return { type: 'object', properties, required, 'x-order': order }
}

function legacyOptionsToSchema(options) {
  const schema = { type: 'object', title: '支付配置', properties: {}, required: [], 'x-order': [] }
  for (const option of options || []) {
    if (!option.name) continue
    schema['x-order'].push(option.name)
    if (option.require) schema.required.push(option.name)
    schema.properties[option.name] = legacyOptionToNode(option)
  }
  return schema
}

// ============================================================
// 表单状态
// ============================================================

const formRef = ref(null)
const submitLoading = ref(false)

const formData = reactive({
  id: 0,
  name: '',
  remark: '',
  app_no: '',
  channel: '',
  channel_name: '',
  status: 1,
  channel_config: '',
  payment_method: [],
  max_limit: 0,
})

const channelConfigForm = reactive({})
const dynamicErrors = reactive({})
const configSchema = ref(null)

const paymentMethodConfig = ref([])
const selectedPaymentMethods = ref([])
const channelOption = ref({})

const formRules = reactive({
  name: [{ required: true, message: '请输入通道名称', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
})

// ============================================================
// 校验
// ============================================================

function validateSchemaNode(schema, root, path = [], parentRequired = []) {
  const name = path[path.length - 1]
  const value = path.length === 0 ? root : getByPath(root, path)
  const required = parentRequired.includes(name)

  if (schema.type === 'object') {
    const errs = {}
    for (const item of orderedEntries(schema)) {
      const ce = validateSchemaNode(item.schema, root, [...path, item.key], schema.required || [])
      if (ce && (typeof ce === 'string' || Object.keys(ce).length > 0)) errs[item.key] = ce
    }
    return errs
  }

  if (schema.type === 'array') {
    const errs = {}
    const arr = Array.isArray(value) ? value : []
    if (required && arr.length === 0) return `${schema.title || name}至少要有一项`
    if (schema.items) {
      arr.forEach((_, idx) => {
        const ce = validateSchemaNode(schema.items, root, [...path, idx], [])
        if (ce && (typeof ce === 'string' || Object.keys(ce).length > 0)) errs[idx] = ce
      })
    }
    return errs
  }

  const label = schema.title || name
  if (required && (value === undefined || value === null || value === '')) return `${label}是必填项`
  if (value === undefined || value === null || value === '') return {}
  if (schema.type === 'integer' && !Number.isInteger(Number(value))) return `${label}必须是整数`
  if (schema.type === 'number' && Number.isNaN(Number(value))) return `${label}必须是数字`
  if (schema.pattern && !new RegExp(schema.pattern).test(String(value))) return `${label}格式不正确`

  const vt = schema['x-ui']?.validate_type || schema.format
  if (vt === 'Email' && !/^\S+@\S+\.\S+$/.test(String(value))) return `${label}格式不正确`
  if (vt === 'Phone' && !/^1[3-9]\d{9}$/.test(String(value))) return `${label}必须是11位手机号码`
  if (vt === 'Domain' && !/^(?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/.test(String(value))) return `${label}格式不正确`
  if ((vt === 'Url' || schema.format === 'url') && !/^https?:\/\/.+/.test(String(value))) return `${label}格式不正确`
  if (vt === 'RsaPrivate' && (!String(value).includes('-----BEGIN') || !String(value).includes('PRIVATE KEY-----'))) return `${label}格式不正确，必须是有效的RSA私钥`
  if (vt === 'RsaPublic' && (!String(value).includes('-----BEGIN') || !String(value).includes('PUBLIC KEY-----'))) return `${label}格式不正确，必须是有效的RSA公钥`
  if (vt === 'RsaCert' && (!String(value).includes('-----BEGIN CERTIFICATE-----') || !String(value).includes('-----END CERTIFICATE-----'))) return `${label}格式不正确，必须是有效的证书`
  return {}
}

function flattenErrors(errors, prefix = '') {
  return Object.entries(errors).flatMap(([k, v]) => {
    const p = prefix ? `${prefix}.${k}` : k
    if (typeof v === 'string') return [[p, v]]
    if (v && typeof v === 'object') return flattenErrors(v, p)
    return []
  })
}

// ============================================================
// 业务函数
// ============================================================

function updateConfigField(path, value) {
  if (path.length === 0) {
    clearObject(channelConfigForm)
    if (value && typeof value === 'object') Object.assign(channelConfigForm, value)
    return
  }
  if (path.length === 1) {
    channelConfigForm[path[0]] = value
    return
  }
  setByPath(channelConfigForm, path, value)
  // 触发顶层依赖刷新
  channelConfigForm[path[0]] = channelConfigForm[path[0]]
}

function handlePaymentMethodChange() {
  formData.payment_method = selectedPaymentMethods.value
      .filter(Boolean)
      .map(item => {
        const [method, product] = item.split('_')
        return { method, product }
      })
}

function deepMerge(target, source) {
  for (const k of Object.keys(source)) {
    if (source[k] !== null && typeof source[k] === 'object' && !Array.isArray(source[k])) {
      if (!target[k] || typeof target[k] !== 'object') target[k] = {}
      deepMerge(target[k], source[k])
    } else {
      target[k] = source[k]
    }
  }
}

function resetFormData() {
  Object.assign(formData, {
    id: 0,
    name: '',
    remark: '',
    app_no: '',
    channel: '',
    channel_name: '',
    status: 1,
    channel_config: '',
    payment_method: [],
    max_limit: 0,
  })
  clearObject(channelConfigForm)
  clearObject(dynamicErrors)
  configSchema.value = null
  paymentMethodConfig.value = []
  selectedPaymentMethods.value = []
  channelOption.value = {}
}

function initForm(data) {
  resetFormData()
  if (!data) return

  formData.id = Number(data.id) || 0
  formData.name = data.name || ''
  formData.remark = data.remark || ''
  formData.app_no = data.app_no || ''
  formData.channel = data.channel || ''
  formData.channel_name = data.channel_name || ''
  formData.status = Number(data.status) || 1
  formData.payment_method = Array.isArray(data.payment_method) ? data.payment_method : []
  formData.max_limit = Number(data.max_limit) || 0

  channelOption.value = data.channel_option || {}

  // 基于 channel_option.schema 渲染
  const nextSchema = channelOption.value?.schema || null
  configSchema.value = nextSchema

  clearObject(channelConfigForm)
  if (nextSchema) initDefaults(nextSchema, channelConfigForm)

  if (data.channel_config) {
    try {
      const cfg = typeof data.channel_config === 'string' ? JSON.parse(data.channel_config) : data.channel_config
      if (cfg && typeof cfg === 'object') deepMerge(channelConfigForm, cfg)
    } catch (e) {
      console.error('解析 channel_config 失败:', e)
    }
  }

  paymentMethodConfig.value = Array.isArray(data.payment_method_config) ? data.payment_method_config : []
  selectedPaymentMethods.value = []
  if (Array.isArray(formData.payment_method)) {
    formData.payment_method.forEach(item => {
      if (item.method && item.product) selectedPaymentMethods.value.push(`${item.method}_${item.product}`)
    })
  }
  handlePaymentMethodChange()
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()

    clearObject(dynamicErrors)
    const errs = configSchema.value ? validateSchemaNode(configSchema.value, channelConfigForm) : {}
    Object.assign(dynamicErrors, errs)
    const flat = flattenErrors(dynamicErrors)
    if (flat.length > 0) {
      ElMessage.error(flat[0][1])
      return
    }

    submitLoading.value = true
    formData.channel_config = JSON.stringify({ ...channelConfigForm })
    handlePaymentMethodChange()

    const submitData = {
      id: Number(formData.id) || 0,
      name: formData.name || '',
      remark: formData.remark || '',
      app_no: formData.app_no || '',
      channel: formData.channel || '',
      channel_name: formData.channel_name || '',
      status: Number(formData.status) || 1,
      channel_config: formData.channel_config,
      payment_method: Array.isArray(formData.payment_method) ? formData.payment_method : [],
      max_limit: Number(formData.max_limit) || 0,
    }

    const res = await savePaymentAccount(submitData)
    if (res.code === 0) {
      ElMessage.success('保存成功')
      emit('saveSuccess')
      handleCancel()
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } catch (error) {
    console.error('提交失败:', error)
    if (error instanceof Error && error.name !== 'ValidationError') {
      ElMessage.error('保存失败')
    }
  } finally {
    submitLoading.value = false
  }
}

function handleCancel() {
  emit('cancel')
}

defineExpose({ initForm })
</script>

<style scoped>
.payment-edit-container {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
}
.config-form-wrapper {
  margin-top: 10px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}
.config-sub-form {
  margin-top: 10px;
}
.method-card {
  margin-bottom: 15px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}
.method-card h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #303133;
}
.el-checkbox {
  margin-right: 15px;
  margin-bottom: 8px;
}
:deep(.el-form-item.is-required .el-form-item__label:before) {
  content: '*';
  color: #f56c6c;
  margin-right: 4px;
}
:deep(.el-form-item.is-error .el-input__inner),
:deep(.el-form-item.is-error .el-textarea__inner) {
  border-color: #f56c6c;
}
:deep(.el-form-item.is-error .el-input__inner:focus),
:deep(.el-form-item.is-error .el-textarea__inner:focus) {
  border-color: #f56c6c;
  box-shadow: 0 0 0 2px rgba(245, 108, 108, 0.2);
}
</style>