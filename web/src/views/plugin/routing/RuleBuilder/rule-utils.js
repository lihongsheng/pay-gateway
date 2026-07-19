/**
 * 规则引擎工具函数
 */

/**
 * 创建空叶子条件
 */
export function createLeafRule(overrides = {}) {
  return {
    id: '',
    operator: 'equal',
    value: '',
    ...overrides
  };
}

/**
 * 创建空条件组
 */
export function createGroupRule(overrides = {}) {
  return {
    condition: 'AND',
    rules: [createLeafRule()],
    ...overrides
  };
}

/**
 * 判断规则节点是否为叶子条件
 */
export function isLeaf(rule) {
  return rule != null && !Array.isArray(rule) && 'id' in rule && !('condition' in rule);
}

/**
 * 判断规则节点是否为条件组
 */
export function isGroup(rule) {
  return rule != null && !Array.isArray(rule) && 'condition' in rule && 'rules' in rule;
}

/**
 * 深拷贝规则 JSON
 */
export function cloneRuleJson(ruleJson) {
  return JSON.parse(JSON.stringify(ruleJson));
}

/**
 * 验证规则 JSON 完整性（前端轻量验证）
 * 返回错误信息数组，空数组表示通过
 */
export function validateRuleJson(ruleJson) {
  const errors = [];

  if (!ruleJson || !ruleJson.condition) {
    errors.push('缺少条件逻辑（AND/OR）');
    return errors;
  }

  if (!Array.isArray(ruleJson.rules) || ruleJson.rules.length === 0) {
    errors.push('请至少配置一条条件');
    return errors;
  }

  function validateNode(node, path) {
    if (isGroup(node)) {
      if (!['AND', 'OR'].includes(node.condition)) {
        errors.push(`${path}: 条件逻辑必须为 AND 或 OR`);
      }
      if (!Array.isArray(node.rules) || node.rules.length === 0) {
        errors.push(`${path}: 条件组不能为空`);
      } else {
        node.rules.forEach((rule, i) => validateNode(rule, `${path}.rules[${i}]`));
      }
    } else if (isLeaf(node)) {
      if (!node.id) {
        errors.push(`${path}: 请选择条件字段`);
      }
      const noValueOps = ['is_null', 'is_not_null'];
      if (!noValueOps.includes(node.operator) && (node.value === '' || node.value == null)) {
        errors.push(`${path}: 请输入条件值`);
      }
    } else {
      errors.push(`${path}: 无效的规则节点`);
    }
  }

  ruleJson.rules.forEach((rule, i) => validateNode(rule, `rules[${i}]`));

  return errors;
}

/**
 * 根据 schema 获取字段可用的操作符列表
 */
export function getOperatorsForField(schema, fieldId) {
  if (!schema?.fields || !fieldId) return [];
  const field = schema.fields.find((f) => f.id === fieldId);
  if (!field) return [];

  if (field.operators && field.operators.length > 0) {
    return (schema.operators || []).filter((op) => field.operators.includes(op.id));
  }

  return schema.operators || [];
}

/**
 * 根据 schema 获取字段的值选项
 */
export function getFieldValues(schema, fieldId) {
  if (!schema?.fields || !fieldId) return [];
  const field = schema.fields.find((f) => f.id === fieldId);
  return field?.values || [];
}

/**
 * 根据 schema 获取字段定义
 */
export function getFieldDef(schema, fieldId) {
  if (!schema?.fields || !fieldId) return null;
  return schema.fields.find((f) => f.id === fieldId) || null;
}

/**
 * 判断操作符是否需要值输入
 */
export function operatorNeedsValue(operator) {
  return !['is_null', 'is_not_null'].includes(operator);
}

/**
 * 判断操作符是否接受数组值
 */
export function operatorAcceptsArray(operator) {
  return ['in', 'not_in'].includes(operator);
}
