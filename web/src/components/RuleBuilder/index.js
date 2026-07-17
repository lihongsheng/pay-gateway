import { ref, onMounted } from 'vue'
import { getRoutingSchema } from '@/api/routing'
import RuleBuilder from './RuleBuilder.vue'

// 导出 RuleBuilder 组件
export { RuleBuilder }

// 全局缓存，避免重复请求
let cachedSchema = null
let cachedStatusOptions = null
let cachedAttributeOptions = null

/**
 * useRuleSchema - 获取路由规则 Schema（字段、操作符、状态选项等）
 * 供 RuleBuilder 组件和搜索/列表页面使用
 */
export function useRuleSchema() {
  const schema = ref(cachedSchema || {})
  const ruleStatusOptions = ref(cachedStatusOptions || [])
  const ruleAttributeOptions = ref(cachedAttributeOptions || [])
  const loading = ref(false)

  const fetchSchema = async () => {
    if (cachedSchema) {
      schema.value = cachedSchema
      ruleStatusOptions.value = cachedStatusOptions
      ruleAttributeOptions.value = cachedAttributeOptions
      return
    }

    loading.value = true
    try {
      const res = await getRoutingSchema()
      const data = res.data || res || {}
      schema.value = data
      cachedSchema = data

      // 提取状态选项，增加 type 字段用于 tag 样式
      const statusList = data.ruleStatusOptions || []
      const tagTypeMap = { 1: 'info', 2: 'success' } // 1=inactive, 2=active
      ruleStatusOptions.value = statusList.map((item) => ({
        label: item.label,
        value: typeof item.value === 'number'
          ? (item.value === 2 ? 'active' : item.value === 1 ? 'inactive' : String(item.value))
          : String(item.value),
        type: item.type || tagTypeMap[item.value] || 'info'
      }))
      cachedStatusOptions = ruleStatusOptions.value

      // 提取属性选项
      const attrList = data.ruleAttributeOptions || []
      ruleAttributeOptions.value = attrList.map((item) => ({
        label: item.label,
        value: typeof item.value === 'number'
          ? (item.value === 1 ? 'single_merchant' : item.value === 2 ? 'multi_merchant' : String(item.value))
          : String(item.value)
      }))
      cachedAttributeOptions = ruleAttributeOptions.value
    } catch (e) {
      console.error('获取规则 Schema 失败:', e)
    } finally {
      loading.value = false
    }
  }

  onMounted(fetchSchema)

  return {
    schema,
    loading,
    ruleStatusOptions,
    ruleAttributeOptions,
    refresh: fetchSchema
  }
}

/**
 * validateRuleJson - 校验规则 JSON 结构合法性
 * @param {object} ruleJson - QueryBuilder 生成的规则 JSON
 * @returns {string[]} 错误信息列表
 */
export function validateRuleJson(ruleJson) {
  const errors = []

  if (!ruleJson || !ruleJson.condition) {
    errors.push('规则必须包含条件逻辑（AND/OR）')
    return errors
  }

  if (!ruleJson.rules || !Array.isArray(ruleJson.rules) || ruleJson.rules.length === 0) {
    errors.push('规则至少需要一个条件')
    return errors
  }

  function validateNode(node, path) {
    if (node.condition) {
      // 条件组
      if (node.condition !== 'AND' && node.condition !== 'OR') {
        errors.push(`${path}: 条件逻辑必须为 AND 或 OR`)
      }
      if (!node.rules || !Array.isArray(node.rules) || node.rules.length === 0) {
        errors.push(`${path}: 条件组至少需要一个条件`)
        return
      }
      node.rules.forEach((child, i) => {
        validateNode(child, `${path}.rules[${i}]`)
      })
    } else {
      // 叶子条件
      if (!node.id) {
        errors.push(`${path}: 缺少字段 ID`)
      }
      if (!node.operator) {
        errors.push(`${path}: 缺少操作符`)
      }
      if (node.operator !== 'is_null' && node.operator !== 'is_not_null' && node.value === undefined && node.value === '') {
        errors.push(`${path}: 缺少条件值`)
      }
    }
  }

  ruleJson.rules.forEach((rule, i) => {
    validateNode(rule, `rules[${i}]`)
  })

  return errors
}

/**
 * createGroupRule - 创建一个默认的 AND 条件组
 * @returns {object} 默认规则 JSON
 */
export function createGroupRule() {
  return {
    condition: 'AND',
    rules: []
  }
}
