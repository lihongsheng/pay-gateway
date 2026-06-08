<!-- SchemaNode.vue：递归渲染 Schema 节点 -->
<template>
  <div v-if="!schema" />

  <!-- ===== Object 节点 ===== -->
  <el-card v-else-if="schema.type === 'object'" class="schema-object" shadow="never">
    <template #header v-if="path.length > 0">{{ title }}</template>
    <div class="schema-form">
      <schema-node
          v-for="item in entries"
          :key="keyOf(childPath(item.key))"
          :schema="item.schema"
          :path="childPath(item.key)"
          :root="root"
          :errors="childErr(item.key)"
          :required="isRequiredKey(item.key)"
          :on-update="onUpdate"
      />
    </div>
  </el-card>

  <!-- ===== Array of objects ===== -->
  <el-form-item v-else-if="schema.type === 'array' && schema.items && schema.items.type === 'object'" :label="title" :required="required">
    <div style="width:100%">
      <div v-for="(item, idx) in arrayValue" :key="idx" class="schema-array-item">
        <div class="schema-array-item-header">
          <span>{{ schema.items.title || '项' }} #{{ idx + 1 }}</span>
          <el-button type="danger" size="small" link @click="removeItem(idx)">删除</el-button>
        </div>
        <schema-node
            :schema="schema.items"
            :path="arrItemPath(idx)"
            :root="root"
            :errors="childErr(idx)"
            :on-update="onUpdate"
        />
      </div>
      <el-button type="primary" plain size="small" @click="addItem">+ 新增{{ schema.items.title || '项' }}</el-button>
    </div>
  </el-form-item>

  <!-- ===== Array of primitives ===== -->
  <el-form-item v-else-if="schema.type === 'array'" :label="title" :error="currentError" :required="required">
    <div style="width:100%">
      <div v-for="(item, idx) in arrayValue" :key="idx" style="display:flex;gap:8px;margin-bottom:6px;">
        <el-input
            class="full-width"
            :model-value="arrayValue[idx]"
            @update:model-value="val => updatePrimitiveArray(idx, val)"
        />
        <el-button type="danger" size="small" link @click="removeItem(idx)">删除</el-button>
      </div>
      <el-button type="primary" plain size="small" @click="addItem">+ 新增</el-button>
    </div>
  </el-form-item>

  <!-- ===== 叶子字段 ===== -->
  <el-form-item v-else :label="title" :error="currentError" :required="required" :class="{ 'has-example': hasExample }">
    <el-input
        v-if="inputTypeV === 'text'"
        class="full-width"
        :model-value="currentValue"
        :placeholder="placeholder"
        clearable
        @update:model-value="val => doUpdate(val)"
    />
    <el-input
        v-else-if="inputTypeV === 'password'"
        class="full-width"
        type="password"
        show-password
        :model-value="currentValue"
        :placeholder="placeholder"
        clearable
        @update:model-value="val => doUpdate(val)"
    />
    <el-input
        v-else-if="inputTypeV === 'textarea'"
        class="full-width"
        type="textarea"
        :rows="textareaRows"
        :model-value="currentValue"
        :placeholder="placeholder"
        @update:model-value="val => doUpdate(val)"
    />
    <el-input-number
        v-else-if="inputTypeV === 'number'"
        :model-value="currentValue"
        @update:model-value="val => doUpdate(val)"
    />
    <el-select
        v-else-if="inputTypeV === 'select'"
        class="full-width"
        :model-value="currentValue"
        clearable
        :placeholder="selectPlaceholder"
        @update:model-value="val => doUpdate(val)"
    >
      <el-option
          v-for="c in choiceItems"
          :key="String(c.value)"
          :label="c.label"
          :value="c.value"
      />
    </el-select>
    <el-radio-group
        v-else-if="inputTypeV === 'radio'"
        :model-value="currentValue"
        @update:model-value="val => doUpdate(val)"
    >
      <el-radio v-for="c in choiceItems" :key="String(c.value)" :label="c.value">
        {{ c.label }}
      </el-radio>
    </el-radio-group>
    <el-checkbox-group
        v-else-if="inputTypeV === 'checkbox'"
        :model-value="currentValue || []"
        @update:model-value="val => doUpdate(val)"
    >
      <el-checkbox v-for="c in choiceItems" :key="String(c.value)" :label="c.value">
        {{ c.label }}
      </el-checkbox>
    </el-checkbox-group>
    <el-switch
        v-else-if="inputTypeV === 'switch'"
        :model-value="Boolean(currentValue)"
        @update:model-value="val => doUpdate(val)"
    />
    <el-input
        v-else
        class="full-width"
        :model-value="currentValue"
        clearable
        @update:model-value="val => doUpdate(val)"
    />
    <div v-if="schema.description" class="hint">{{ schema.description }}</div>
  </el-form-item>
</template>

<script>
import { computed } from 'vue'

function getByPath(root, path) {
  return path.reduce((obj, key) => (obj == null ? obj : obj[key]), root)
}

function setByPath(root, path, value) {
  let obj = root
  for (let i = 0; i < path.length - 1; i++) {
    const key = path[i]
    if (obj[key] === undefined || obj[key] === null || typeof obj[key] !== 'object') {
      obj[key] = typeof path[i + 1] === 'number' ? [] : {}
    }
    obj = obj[key]
  }
  obj[path[path.length - 1]] = value
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

// formatExample: 把示例值格式化成 placeholder 文案
//  - 字符串过长（>60）取首行 + 截断
//  - 对象/数组用紧凑 JSON
function formatExample(v) {
  if (typeof v === 'string') {
    const firstLine = v.split('\n')[0]
    if (firstLine.length > 60) return firstLine.slice(0, 57) + '...'
    if (v.includes('\n')) return firstLine + '...'
    return firstLine
  }
  if (typeof v === 'number' || typeof v === 'boolean') return String(v)
  try {
    const s = JSON.stringify(v)
    return s.length > 60 ? s.slice(0, 57) + '...' : s
  } catch {
    return String(v)
  }
}

export default {
  name: 'SchemaNode',
  props: {
    schema: Object,
    path: { type: Array, default: () => [] },
    root: Object,
    errors: { type: Object, default: () => ({}) },
    required: { type: Boolean, default: false },
    onUpdate: { type: Function, default: null },
  },
  computed: {
    entries() {
      if (!this.schema || !this.schema.properties) return []
      const order = this.schema['x-order'] || Object.keys(this.schema.properties)
      const names = [...order, ...Object.keys(this.schema.properties).filter(k => !order.includes(k))]
      return names.filter(k => this.schema.properties[k]).map(k => ({ key: k, schema: this.schema.properties[k] }))
    },
    title() {
      return this.schema?.title || (this.path.length ? this.path[this.path.length - 1] : '配置')
    },
    currentValue() {
      return getByPath(this.root, this.path)
    },
    currentError() {
      let e = this.errors
      for (const k of this.path) {
        if (e == null) return ''
        e = e[k]
      }
      return typeof e === 'string' ? e : ''
    },
    inputTypeV() {
      const s = this.schema
      if (!s) return 'text'
      if (s.oneOf?.length || s.enum?.length) return s['x-ui']?.input_type || 'select'
      if (s.type === 'boolean') return s['x-ui']?.input_type || 'switch'
      if (s.type === 'integer' || s.type === 'number') return 'number'
      return s['x-ui']?.input_type || 'text'
    },
    placeholder() {
      // 优先用 x-ui.placeholder；其次用 schema.examples[0] 拼成「示例：xxx」；再回退到「请输入 + title」
      const ui = this.schema['x-ui']?.placeholder
      if (ui) return ui
      const example = Array.isArray(this.schema.examples) ? this.schema.examples[0] : undefined
      if (example !== undefined && example !== null && example !== '') {
        return `示例：${formatExample(example)}`
      }
      return '请输入' + (this.schema.title || '')
    },
    selectPlaceholder() {
      const ui = this.schema['x-ui']?.placeholder
      if (ui) return ui
      const example = Array.isArray(this.schema.examples) ? this.schema.examples[0] : undefined
      if (example !== undefined && example !== null && example !== '') {
        return `示例：${formatExample(example)}`
      }
      return '请选择' + (this.schema.title || '')
    },
    textareaRows() {
      return this.schema['x-ui']?.rows || 6
    },
    choiceItems() {
      if (this.schema?.oneOf?.length) return this.schema.oneOf.map(x => ({ value: x.const, label: x.title || x.const }))
      if (this.schema?.enum?.length) return this.schema.enum.map(x => ({ value: x, label: String(x) }))
      return []
    },
    hasExample() {
      return Array.isArray(this.schema?.examples) && this.schema.examples.length > 0
    },
    arrayValue() {
      const v = getByPath(this.root, this.path)
      return Array.isArray(v) ? v : []
    },
  },
  methods: {
    keyOf(p) { return p.join('.') },
    childPath(key) { return [...this.path, key] },
    arrItemPath(idx) { return [...this.path, idx] },
    isRequiredKey(key) {
      return Array.isArray(this.schema?.required) && this.schema.required.includes(key)
    },
    childErr(key) {
      let e = this.errors
      for (const k of this.path) {
        if (e == null) return {}
        e = e[k]
      }
      return e?.[key] || {}
    },
    doUpdate(val) {
      if (this.onUpdate) this.onUpdate(this.path, val)
    },
    addItem() {
      const arr = [...this.arrayValue]
      arr.push(defaultValueFor(this.schema?.items))
      if (this.schema?.items?.type === 'object') {
        arr[arr.length - 1] = {}
      }
      this.doUpdate(arr)
    },
    removeItem(idx) {
      const arr = [...this.arrayValue]
      arr.splice(idx, 1)
      this.doUpdate(arr)
    },
    updatePrimitiveArray(idx, val) {
      const arr = [...this.arrayValue]
      arr[idx] = val
      this.doUpdate(arr)
    },
  },
}
</script>

<style scoped>
.schema-object { margin-bottom: 16px; }
.schema-object .el-card__header { padding: 10px 14px; font-weight: 600; }
.schema-form { padding-top: 8px; }
.schema-array-item { border: 1px dashed #dcdfe6; padding: 8px 12px; border-radius: 4px; margin-bottom: 8px; }
.schema-array-item-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.hint { color: #909399; font-size: 12px; margin-top: 4px; }
.full-width { width: 100%; }
</style>