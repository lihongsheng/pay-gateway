<template>
  <div class="rule-row" :class="[`depth-${depth}`]">
    <div class="rule-row-connector" v-if="!first"></div>
    <el-select
      :model-value="rule.id"
      placeholder="选择条件"
      class="rule-field-select"
      :disabled="disabled"
      @change="onFieldChange"
    >
      <el-option
        v-for="field in schema.fields"
        :key="field.id"
        :label="field.label"
        :value="field.id"
      />
    </el-select>

    <el-select
      :model-value="rule.operator"
      placeholder="运算符"
      class="rule-operator-select"
      :disabled="disabled"
      @change="onOperatorChange"
    >
      <el-option
        v-for="op in fieldOperators"
        :key="op.id"
        :label="op.label"
        :value="op.id"
      />
    </el-select>

    <!-- select 类型字段 + in/not_in 操作符 → 多选 -->
    <el-select
      v-if="isSelectField && isArrayOp && needsValue"
      :model-value="rule.value"
      multiple
      placeholder="选择值"
      class="rule-value-select"
      :disabled="disabled"
      @change="(val) => onValueChange(val)"
    >
      <el-option
        v-for="v in fieldValues"
        :key="v.value"
        :label="v.label"
        :value="v.value"
      />
    </el-select>

    <!-- select 类型字段 + 非 in/not_in → 单选 -->
    <el-select
      v-else-if="isSelectField && !isArrayOp && needsValue"
      :model-value="rule.value"
      placeholder="选择值"
      class="rule-value-select"
      :disabled="disabled"
      clearable
      @change="(val) => onValueChange(val)"
    >
      <el-option
        v-for="v in fieldValues"
        :key="v.value"
        :label="v.label"
        :value="v.value"
      />
    </el-select>

    <!-- radio 类型字段 -->
    <el-radio-group
      v-else-if="isRadioField && needsValue"
      :model-value="rule.value"
      :disabled="disabled"
      @change="(val) => onValueChange(val)"
    >
      <el-radio
        v-for="v in fieldValues"
        :key="v.value"
        :value="v.value"
      >
        {{ v.label }}
      </el-radio>
    </el-radio-group>

    <!-- 文本输入 -->
    <el-input
      v-else-if="needsValue"
      :model-value="rule.value"
      :placeholder="valuePlaceholder"
      :disabled="disabled"
      class="rule-value-input"
      clearable
      @update:model-value="(val) => onValueChange(val)"
    />

    <!-- is_null / is_not_null 无需值输入 -->
    <span v-else class="rule-value-placeholder">-</span>

    <el-button
      v-if="!disabled && canRemove"
      type="danger"
      link
      class="rule-remove-btn"
      @click="$emit('remove')"
    >
      删除
    </el-button>
  </div>
</template>

<script setup>
  import { computed } from 'vue';
  import {
    getOperatorsForField,
    getFieldValues,
    getFieldDef,
    operatorNeedsValue,
    operatorAcceptsArray
  } from './rule-utils';

  const props = defineProps({
    rule: { type: Object, required: true },
    schema: { type: Object, required: true },
    depth: { type: Number, default: 0 },
    disabled: { type: Boolean, default: false },
    first: { type: Boolean, default: false },
    canRemove: { type: Boolean, default: true }
  });

  const emit = defineEmits(['update:rule', 'remove']);

  // 直接从 props.rule 读取，不做 reactive 拷贝，避免循环更新
  const fieldDef = computed(() => getFieldDef(props.schema, props.rule.id));
  const fieldOperators = computed(() => getOperatorsForField(props.schema, props.rule.id));
  const fieldValues = computed(() => getFieldValues(props.schema, props.rule.id));

  const isSelectField = computed(() => fieldDef.value?.input === 'select');
  const isRadioField = computed(() => fieldDef.value?.input === 'radio');
  const isArrayOp = computed(() => operatorAcceptsArray(props.rule.operator));
  const needsValue = computed(() => operatorNeedsValue(props.rule.operator));

  const valuePlaceholder = computed(() => {
    if (!fieldDef.value) return '请输入值';
    if (fieldDef.value.type === 'integer') return '请输入整数';
    if (fieldDef.value.type === 'double') return '请输入数值';
    return '请输入值';
  });

  const onFieldChange = (newFieldId) => {
    const newRule = { ...props.rule, id: newFieldId };
    // 切换字段时重置操作符和值
    const ops = getOperatorsForField(props.schema, newFieldId);
    if (ops.length > 0 && !ops.find((op) => op.id === newRule.operator)) {
      newRule.operator = ops[0].id;
    }
    // 重置值
    if (operatorAcceptsArray(newRule.operator)) {
      newRule.value = [];
    } else {
      newRule.value = '';
    }
    emit('update:rule', newRule);
  };

  const onOperatorChange = (newOperator) => {
    const newRule = { ...props.rule, operator: newOperator };
    // 切换操作符时调整值类型
    if (operatorAcceptsArray(newOperator) && !Array.isArray(newRule.value)) {
      newRule.value = newRule.value ? [newRule.value] : [];
    } else if (!operatorAcceptsArray(newOperator) && Array.isArray(newRule.value)) {
      newRule.value = newRule.value.length > 0 ? newRule.value[0] : '';
    }
    emit('update:rule', newRule);
  };

  const onValueChange = (newValue) => {
    emit('update:rule', { ...props.rule, value: newValue });
  };
</script>

<style lang="scss" scoped>
  .rule-row {
    display: flex;
    align-items: center;
    gap: 10px;
    position: relative;
    padding: 6px 0;
  }

  .rule-row-connector {
    position: absolute;
    left: -16px;
    top: 50%;
    width: 16px;
    height: 2px;
    background: var(--rule-logic-color, var(--el-color-primary-light-5));
  }

  .rule-field-select {
    width: 160px;
    flex-shrink: 0;
  }

  .rule-operator-select {
    width: 130px;
    flex-shrink: 0;
  }

  .rule-value-select {
    min-width: 180px;
    max-width: 300px;
  }

  .rule-value-input {
    width: 180px;
    flex-shrink: 0;
  }

  .rule-value-placeholder {
    color: var(--el-text-color-placeholder);
    font-size: 13px;
    width: 60px;
    text-align: center;
  }

  .rule-remove-btn {
    flex-shrink: 0;
    min-width: 36px;
  }

  @media (max-width: 768px) {
    .rule-row {
      flex-wrap: wrap;
    }

    .rule-field-select,
    .rule-operator-select,
    .rule-value-select,
    .rule-value-input {
      width: 100%;
      max-width: none;
    }
  }
</style>
