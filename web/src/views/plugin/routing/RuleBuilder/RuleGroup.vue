<template>
  <div
    class="rule-group"
    :class="[
      `depth-${depth}`,
      `logic-${group.condition.toLowerCase()}`
    ]"
  >
    <!-- 左侧竖线 + AND/OR 标签 -->
    <div class="rule-group-rail">
      <div class="rail-line"></div>
      <button
        type="button"
        class="rail-badge"
        :class="[`logic-${group.condition.toLowerCase()}`]"
        :disabled="disabled"
        :title="`当前关系：${conditionLabel}，点击切换`"
        @click="toggleCondition"
      >
        {{ group.condition }}
      </button>
    </div>

    <!-- 子规则列表 -->
    <div class="rule-group-body">
      <!-- 删除整个条件组按钮 -->
      <el-button
        v-if="canRemove && !disabled"
        type="danger"
        link
        class="rule-group-remove"
        title="删除此条件组"
        @click="emit('remove')"
      >
        删除组
      </el-button>

      <template v-for="(rule, index) in group.rules" :key="index">
        <!-- 叶子条件 -->
        <RuleRow
          v-if="isLeaf(rule)"
          :rule="rule"
          :schema="schema"
          :depth="depth"
          :disabled="disabled"
          :first="index === 0"
          :can-remove="group.rules.length > 1"
          @update:rule="updateRule(index, $event)"
          @remove="removeRule(index)"
        />
        <!-- 嵌套条件组 -->
        <RuleGroup
          v-else-if="isGroup(rule)"
          :model-value="rule"
          :schema="schema"
          :depth="depth + 1"
          :max-depth="maxDepth"
          :disabled="disabled"
          :can-remove="true"
          @update:model-value="updateRule(index, $event)"
          @remove="removeRule(index)"
        />
      </template>

      <!-- 操作按钮 -->
      <div class="rule-group-actions" v-if="!disabled">
        <el-button size="small" type="primary" plain @click="addLeaf">
          + 添加条件
        </el-button>
        <el-button
          size="small"
          type="primary"
          plain
          :disabled="depth >= maxDepth"
          @click="addGroup"
        >
          + 添加条件组
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed } from 'vue';
  import RuleRow from './RuleRow.vue';
  import { isLeaf, isGroup, createLeafRule, createGroupRule, cloneRuleJson } from './rule-utils';

  const props = defineProps({
    modelValue: { type: Object, required: true },
    schema: { type: Object, required: true },
    depth: { type: Number, default: 0 },
    maxDepth: { type: Number, default: 10 },
    disabled: { type: Boolean, default: false },
    canRemove: { type: Boolean, default: false }
  });

  const emit = defineEmits(['update:modelValue', 'remove']);

  // 直接使用 props.modelValue 作为数据源，不做 reactive 拷贝
  // 这样避免了 watch + reactive + emit 的循环更新
  const group = computed(() => props.modelValue);

  const conditionLabel = computed(() =>
    group.value.condition === 'OR' ? '或' : '且'
  );

  const emitUpdate = (newGroup) => {
    emit('update:modelValue', newGroup);
  };

  const toggleCondition = () => {
    if (props.disabled) return;
    const next = cloneRuleJson(group.value);
    next.condition = next.condition === 'AND' ? 'OR' : 'AND';
    emitUpdate(next);
  };

  const addLeaf = () => {
    const next = cloneRuleJson(group.value);
    next.rules.push(createLeafRule());
    emitUpdate(next);
  };

  const addGroup = () => {
    if (props.depth >= props.maxDepth) return;
    const next = cloneRuleJson(group.value);
    next.rules.push(createGroupRule());
    emitUpdate(next);
  };

  const removeRule = (index) => {
    const next = cloneRuleJson(group.value);
    next.rules.splice(index, 1);
    // 根组至少保留一条规则，子组删空后由父级处理
    if (props.depth === 0 && next.rules.length === 0) {
      next.rules.push(createLeafRule());
    }
    emitUpdate(next);
  };

  const updateRule = (index, newRule) => {
    const next = cloneRuleJson(group.value);
    next.rules[index] = newRule;
    emitUpdate(next);
  };
</script>

<style lang="scss" scoped>
  .rule-group {
    position: relative;
    display: flex;
    gap: 0;
    padding: 0;
  }

  .rule-group-rail {
    position: relative;
    width: 28px;
    flex-shrink: 0;
    display: flex;
    align-items: stretch;
  }

  .rail-line {
    position: absolute;
    left: 12px;
    top: 0;
    bottom: 0;
    width: 2px;
    border-radius: 2px;
    background: var(--rule-logic-color, var(--el-color-primary-light-5));
    transition: background-color 0.2s ease;
  }

  .rail-badge {
    position: absolute;
    left: 2px;
    top: 50%;
    z-index: 1;
    display: flex;
    width: 24px;
    height: 24px;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--rule-logic-color, var(--el-color-primary));
    border-radius: 6px;
    color: #fff;
    font-size: 11px;
    font-weight: 700;
    background: var(--rule-logic-color, var(--el-color-primary));
    box-shadow: 0 2px 8px var(--rule-logic-shadow, rgb(64 158 255 / 20%));
    cursor: pointer;
    transform: translateY(-50%);
    transition:
      background-color 0.2s ease,
      transform 0.15s ease,
      box-shadow 0.15s ease;
    padding: 0;
  }

  .rail-badge:hover:not(:disabled) {
    transform: translateY(-50%) scale(1.1);
    box-shadow: 0 4px 12px var(--rule-logic-shadow, rgb(64 158 255 / 30%));
  }

  .rail-badge:disabled {
    cursor: default;
    opacity: 0.8;
  }

  .logic-and {
    --rule-logic-color: var(--el-color-primary);
    --rule-logic-shadow: rgb(64 158 255 / 20%);
    --rule-logic-bg: var(--el-color-primary-light-9);
    --rule-logic-border: var(--el-color-primary-light-7);
  }

  .logic-or {
    --rule-logic-color: var(--el-color-warning);
    --rule-logic-shadow: rgb(230 162 60 / 22%);
    --rule-logic-bg: var(--el-color-warning-light-9);
    --rule-logic-border: var(--el-color-warning-light-7);
  }

  .rule-group-body {
    flex: 1;
    min-width: 0;
    padding: 10px 12px;
    border: 1px solid var(--rule-logic-border, var(--el-border-color-light));
    border-left: 3px solid var(--rule-logic-color, var(--el-color-primary));
    border-radius: 8px;
    background: linear-gradient(
      90deg,
      var(--rule-logic-bg, var(--el-fill-color-extra-light)),
      transparent 20%
    );
    position: relative;
  }

  .rule-group-remove {
    position: absolute;
    top: 4px;
    right: 6px;
    z-index: 1;
    font-size: 12px;
  }

  .rule-group-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px dashed var(--el-border-color-lighter);
  }

  .depth-0 {
    --rule-logic-color: var(--el-color-primary);
    --rule-logic-shadow: rgb(64 158 255 / 20%);
    --rule-logic-bg: var(--el-color-primary-light-9);
    --rule-logic-border: var(--el-color-primary-light-7);
  }

  .depth-1 {
    --rule-logic-color: var(--el-color-success);
    --rule-logic-shadow: rgb(103 194 58 / 20%);
    --rule-logic-bg: var(--el-color-success-light-9);
    --rule-logic-border: var(--el-color-success-light-7);
  }

  .depth-2 {
    --rule-logic-color: var(--el-color-warning);
    --rule-logic-shadow: rgb(230 162 60 / 22%);
    --rule-logic-bg: var(--el-color-warning-light-9);
    --rule-logic-border: var(--el-color-warning-light-7);
  }

  .depth-3 {
    --rule-logic-color: var(--el-color-danger);
    --rule-logic-shadow: rgb(245 108 108 / 22%);
    --rule-logic-bg: var(--el-color-danger-light-9);
    --rule-logic-border: var(--el-color-danger-light-7);
  }

  @media (max-width: 768px) {
    .rule-group {
      flex-direction: column;
    }

    .rule-group-rail {
      width: 100%;
      height: 28px;
      flex-direction: row;
      align-items: center;
      justify-content: center;
    }

    .rail-line {
      left: 0;
      right: 0;
      top: 13px;
      bottom: auto;
      width: auto;
      height: 2px;
    }

    .rail-badge {
      position: relative;
      left: auto;
      top: auto;
      transform: none;
    }

    .rail-badge:hover:not(:disabled) {
      transform: scale(1.1);
    }
  }
</style>
