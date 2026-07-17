<template>
  <div class="rule-builder" :class="{ 'is-disabled': disabled }">
    <div v-if="group" class="rule-group">
      <!-- 条件逻辑切换 -->
      <div class="group-header">
        <el-radio-group
          v-model="group.condition"
          :disabled="disabled"
          size="small"
          @change="emitUpdate"
        >
          <el-radio-button
            v-for="logic in conditionOptions"
            :key="logic.id"
            :value="logic.id"
          >{{ logic.label }}</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 规则列表 -->
      <div class="group-rules">
        <template v-for="(rule, index) in group.rules" :key="index">
          <!-- 子条件组 -->
          <div v-if="rule.condition" class="rule-sub-group">
            <RuleBuilder
              :model-value="rule"
              :external-schema="externalSchema"
              :disabled="disabled"
              :depth="depth + 1"
              @update:model-value="(val) => updateRule(index, val)"
              @remove="removeRule(index)"
            />
          </div>

          <!-- 叶子条件 -->
          <div v-else class="rule-item">
            <el-select
              v-model="rule.id"
              :disabled="disabled"
              placeholder="选择字段"
              class="rule-field"
              @change="onFieldChange(rule); emitUpdate()"
            >
              <el-option
                v-for="field in fields"
                :key="field.id"
                :label="field.label"
                :value="field.id"
              />
            </el-select>

            <el-select
              v-model="rule.operator"
              :disabled="disabled"
              placeholder="操作符"
              class="rule-operator"
              @change="emitUpdate()"
            >
              <el-option
                v-for="op in getOperatorsForField(rule.id)"
                :key="op.id"
                :label="op.label"
                :value="op.id"
              />
            </el-select>

            <!-- 值输入 -->
            <template v-if="needsValue(rule.operator)">
              <el-select
                v-if="getFieldInput(rule.id) === 'select'"
                v-model="rule.value"
                :disabled="disabled"
                placeholder="选择值"
                class="rule-value"
                multiple
                filterable
                @change="emitUpdate()"
              >
                <el-option
                  v-for="opt in getFieldValues(rule.id)"
                  :key="opt.value"
                  :label="opt.label"
                  :value="opt.value"
                />
              </el-select>
              <el-input
                v-else
                v-model="rule.value"
                :disabled="disabled"
                placeholder="输入值"
                class="rule-value"
                @change="emitUpdate()"
              />
            </template>

            <el-button
              v-if="!disabled"
              type="danger"
              :icon="Delete"
              circle
              size="small"
              @click="removeRule(index)"
            />
          </div>
        </template>
      </div>

      <!-- 添加按钮 -->
      <div v-if="!disabled" class="group-actions">
        <el-button size="small" @click="addCondition">
          <el-icon><Plus /></el-icon> 添加条件
        </el-button>
        <el-button size="small" @click="addGroup">
          <el-icon><FolderAdd /></el-icon> 添加条件组
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, watch, onMounted } from 'vue'
  import { Plus, Delete, FolderAdd } from '@element-plus/icons-vue'
  import { getRoutingSchema } from '@/api/routing'

  const props = defineProps({
    modelValue: { type: Object, default: () => ({ condition: 'AND', rules: [] }) },
    externalSchema: { type: Object, default: null },
    disabled: { type: Boolean, default: false },
    depth: { type: Number, default: 0 }
  })

  const emit = defineEmits(['update:modelValue', 'remove'])

  const group = ref({ ...props.modelValue })
  const fields = ref([])
  const operators = ref([])
  const conditionOptions = ref([])

  // 监听外部 schema 或 modelValue 变化
  watch(() => props.externalSchema, (schema) => {
    if (schema) applySchema(schema)
  }, { immediate: true, deep: true })

  watch(() => props.modelValue, (val) => {
    if (val) group.value = { ...val }
  }, { deep: true })

  function applySchema(schema) {
    fields.value = schema.fields || []
    operators.value = schema.operators || []
    conditionOptions.value = schema.conditions || [
      { id: 'AND', label: '且' },
      { id: 'OR', label: '或' }
    ]
  }

  // 如果没有外部 schema，自动加载
  onMounted(async () => {
    if (!props.externalSchema) {
      try {
        const res = await getRoutingSchema()
        const data = res.data || res || {}
        applySchema(data)
      } catch (e) {
        console.error('RuleBuilder: 加载 Schema 失败', e)
      }
    }
  })

  function getOperatorsForField(fieldId) {
    const field = fields.value.find(f => f.id === fieldId)
    if (!field) return operators.value
    return operators.value.filter(op =>
      !field.operators || field.operators.includes(op.id)
    )
  }

  function getFieldInput(fieldId) {
    const field = fields.value.find(f => f.id === fieldId)
    return field?.input || 'text'
  }

  function getFieldValues(fieldId) {
    const field = fields.value.find(f => f.id === fieldId)
    return field?.values || []
  }

  function needsValue(op) {
    return op !== 'is_null' && op !== 'is_not_null'
  }

  function onFieldChange(rule) {
    // 切换字段时重置操作符和值
    const field = fields.value.find(f => f.id === rule.id)
    if (field && field.operators?.length > 0) {
      if (!field.operators.includes(rule.operator)) {
        rule.operator = field.operators[0]
      }
    }
    rule.value = ''
  }

  function addCondition() {
    const field = fields.value[0]
    group.value.rules.push({
      id: field?.id || '',
      operator: field?.operators?.[0] || 'equal',
      value: ''
    })
    emitUpdate()
  }

  function addGroup() {
    group.value.rules.push({
      condition: 'AND',
      rules: []
    })
    emitUpdate()
  }

  function removeRule(index) {
    group.value.rules.splice(index, 1)
    emitUpdate()
  }

  function updateRule(index, val) {
    group.value.rules[index] = val
    emitUpdate()
  }

  function emitUpdate() {
    emit('update:modelValue', { ...group.value })
  }
</script>

<style lang="scss" scoped>
  .rule-builder {
    width: 100%;
  }

  .rule-group {
    padding: 12px;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    background: var(--el-fill-color-lighter);
  }

  .group-header {
    margin-bottom: 10px;
  }

  .group-rules {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .rule-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 8px;
    border-radius: 6px;
    background: var(--el-bg-color);
  }

  .rule-field {
    width: 180px;
  }

  .rule-operator {
    width: 140px;
  }

  .rule-value {
    min-width: 160px;
    flex: 1;
  }

  .rule-sub-group {
    margin-left: 8px;
  }

  .group-actions {
    display: flex;
    gap: 8px;
    margin-top: 10px;
  }

  .is-disabled {
    opacity: 0.85;
    pointer-events: none;
  }
</style>
