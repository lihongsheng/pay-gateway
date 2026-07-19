<template>
  <el-dialog
    v-model="dialogVisible"
    :width="960"
    title="路由规则详情"
    :close-on-click-modal="false"
    @close="emit('close')"
  >
    <div v-loading="loading">
      <template v-if="detail">
        <el-descriptions :column="3" border class="detail-section">
          <el-descriptions-item label="规则ID">{{ detail.ruleNo || '-' }}</el-descriptions-item>
          <el-descriptions-item label="规则名称">{{ detail.ruleName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="规则状态">
            <el-tag :type="getRuleStatusType(detail.ruleStatus)">
              {{ getOptionLabel(ruleStatusOptions, detail.ruleStatus) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="商户编号">{{ detail.mchNo || '-' }}</el-descriptions-item>
          <el-descriptions-item label="应用编号">{{ detail.appNo || '-' }}</el-descriptions-item>
          <el-descriptions-item label="规则属性">
            {{ getOptionLabel(ruleAttributeOptions, detail.ruleAttribute) }}
          </el-descriptions-item>
          <el-descriptions-item label="优先级">{{ detail.priority ?? '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ detail.createTime || '-' }}</el-descriptions-item>
          <el-descriptions-item label="规则描述">
            {{ detail.ruleDesc || '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <el-card header="匹配条件" class="detail-section" shadow="never">
          <RuleBuilder
            v-if="detail.rules && detail.rules.condition"
            :model-value="detail.rules"
            :external-schema="schema"
            :disabled="true"
          />
          <el-empty v-else description="暂无条件" :image-size="70" />
        </el-card>

        <el-card header="执行动作" class="detail-section" shadow="never">
          <el-table :data="detail.actions || []" size="small">
            <el-table-column label="动作类型" min-width="160">
              <template #default="{ row }">
                {{ row.actionType === 'select_paypal_account' ? '选择收单账号' : row.actionType || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="账号属性" min-width="150">
              <template #default="{ row }">
                {{ getOptionLabel(ruleAttributeOptions, row.accountAttribute || detail.ruleAttribute) }}
              </template>
            </el-table-column>
            <el-table-column label="PayPal账号" min-width="180">
              <template #default="{ row }">
                {{ row.paypalAccountNo || '自动分配' }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </template>
    </div>
  </el-dialog>
</template>

<script setup>
  import { reactive, ref, onUnmounted, computed } from 'vue';
  import { ElMessage } from 'element-plus';
  import { getRoutingRule } from '../api/routing';
  import { RuleBuilder, useRuleSchema } from '../RuleBuilder';

  const props = defineProps({
    visible: Boolean,
    data: Object
  });

  const emit = defineEmits(['close']);

  const dialogVisible = computed({
    get: () => props.visible,
    set: (val) => { if (!val) emit('close'); }
  });

  const { schema, ruleStatusOptions, ruleAttributeOptions } = useRuleSchema();
  const loading = ref(false);
  const detail = reactive({ ...(props.data || {}) });

  // 防止组件卸载后异步回调修改状态
  let unmounted = false;
  onUnmounted(() => {
    unmounted = true;
  });

  const getOptionLabel = (options, value) => {
    const list = Array.isArray(options) ? options : (options?.value ?? []);
    return list.find((item) => item.value === value)?.label ?? value ?? '-';
  };

  const getRuleStatusType = (value) => {
    const list = ruleStatusOptions.value ?? [];
    return list.find((item) => item.value === value)?.type ?? 'info';
  };

  const query = () => {
    if (!props.data?.id) {
      return;
    }
    loading.value = true;
    getRoutingRule(props.data.id)
      .then((data) => {
        if (unmounted) return;
        Object.assign(detail, data);
      })
      .catch((e) => {
        if (unmounted) return;
        ElMessage.error(e.message);
      })
      .finally(() => {
        if (unmounted) return;
        loading.value = false;
      });
  };

  query();
</script>

<style lang="scss" scoped>
  .detail-section {
    margin-bottom: 12px;
  }
</style>
