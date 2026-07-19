<template>
  <div class="rule-builder" v-loading="schemaLoading">
    <RuleGroup
      v-if="resolvedSchema"
      :model-value="modelValue"
      :schema="resolvedSchema"
      :depth="0"
      :max-depth="maxDepth"
      :disabled="disabled"
      @update:model-value="onUpdate"
    />
    <el-empty v-else-if="!schemaLoading" description="规则 Schema 加载失败" :image-size="60">
      <el-button type="primary" size="small" @click="fetchSchema">重试</el-button>
    </el-empty>
  </div>
</template>

<script setup>
  import { computed } from 'vue';
  import RuleGroup from './RuleGroup.vue';
  import { useRuleSchema } from './useRuleSchema';
  import { createGroupRule } from './rule-utils';

  const props = defineProps({
    modelValue: {
      type: Object,
      default: () => createGroupRule()
    },
    disabled: {
      type: Boolean,
      default: false
    },
    maxDepth: {
      type: Number,
      default: 10
    },
    externalSchema: {
      type: Object,
      default: null
    }
  });

  const emit = defineEmits(['update:modelValue']);

  const { schema: apiSchema, loading: schemaLoading, fetchSchema } = useRuleSchema();

  const resolvedSchema = computed(() => props.externalSchema || apiSchema.value);

  const onUpdate = (newVal) => {
    emit('update:modelValue', newVal);
  };
</script>

<style lang="scss" scoped>
  .rule-builder {
    min-height: 60px;
  }
</style>
