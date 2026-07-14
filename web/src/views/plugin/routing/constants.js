/**
 * 路由规则页面常量（仅保留向后兼容辅助函数）。
 *
 * conditionFieldOptions、operatorOptions 已迁移到后端 Schema API。
 * ruleStatusOptions、ruleAttributeOptions 也从 Schema API 获取。
 *
 * 此文件仅保留通用辅助函数，供旧代码过渡使用。
 */

/**
 * 从选项数组中获取标签
 */
export function getOptionLabel(options, value) {
  return options.find((item) => item.value === value)?.label ?? value ?? '-';
}
