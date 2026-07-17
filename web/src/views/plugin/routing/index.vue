<template>
  <ele-page>
    <el-row :gutter="16" class="routing-summary page-section">
      <el-col
        v-for="item in summaryCards"
        :key="item.key"
        :xs="24"
        :sm="12"
        :md="6"
        :lg="4"
      >
        <ele-card
          :body-style="{ padding: '16px 18px' }"
          class="summary-card"
          :class="[`is-${item.type}`, { 'is-active': activeSummaryKey === item.key }]"
          @click="item.onClick?.()"
        >
          <div class="summary-head">
            <span class="summary-dot"></span>
            <span class="summary-label">{{ item.label }}</span>
          </div>
          <div class="summary-value">{{ item.value }}</div>
          <div class="summary-desc">{{ item.desc }}</div>
        </ele-card>
      </el-col>
    </el-row>

    <routing-search @search="(where) => reload(where, 1)" />

    <ele-card class="table-section" :body-style="{ paddingTop: '8px' }">
      <ele-pro-table
        ref="tableRef"
        row-key="id"
        :columns="columns"
        :datasource="datasource"
        :show-overflow-tooltip="true"
        :highlight-current-row="true"
        v-model:selections="selections"
        :export-config="{ fileName: '路由策略' }"
        cache-key="routingRuleTableV2"
      >
        <template #toolbar>
          <btn-items
            :items="[
              {
                preset: 'add',
                title: '新增规则',
                onClick: () => handleAdd()
              }
            ]"
          />
        </template>

        <template #ruleNo="{ row }">
          <el-link type="primary" @click="handleDetail(row)">{{ row.ruleNo }}</el-link>
        </template>

        <template #ruleStatus="{ row }">
          <el-tooltip
            :content="getOptionLabel(row.ruleStatus)"
            placement="top"
          >
            <el-switch
              size="small"
              :model-value="row.ruleStatus === 'active'"
              @change="(value) => handleToggleStatus(row, value)"
            />
          </el-tooltip>
        </template>

        <template #ruleAttribute="{ row }">
          {{ getOptionLabel(row.ruleAttribute, 'ruleAttribute') }}
        </template>

        <template #statusTag="{ row }">
          <el-tag :type="getRuleStatusType(row.ruleStatus)">
            {{ getOptionLabel(row.ruleStatus) }}
          </el-tag>
        </template>

        <template #action="{ row }">
          <el-space :size="6" spacer="|">
            <el-button type="primary" link @click.stop="handleDetail(row)">
              查看
            </el-button>
            <el-button type="primary" link @click.stop="handleEdit(row)">
              编辑
            </el-button>
          </el-space>
        </template>
      </ele-pro-table>
    </ele-card>
  </ele-page>
</template>

<script setup>
  import { computed, ref } from 'vue';
  import { ElMessageBox } from 'element-plus';
  import { EleMessage, useModal } from 'ele-admin-plus';
  import {
    pageRoutingRules,
    toggleRoutingRuleStatus
  } from '@/api/routing';
  import RoutingSearch from './components/routing-search.vue';
  import { useRuleSchema } from '@/components/RuleBuilder';

  defineOptions({ name: 'RoutingRule' });

  const { openModal } = useModal();
  const tableRef = ref(null);
  const selections = ref([]);
  const summary = ref({});
  const activeSummaryKey = ref('all');

  const { ruleStatusOptions, ruleAttributeOptions } = useRuleSchema();

  const getOptionLabel = (value, type = 'ruleStatus') => {
    const options = type === 'ruleAttribute' ? ruleAttributeOptions.value : ruleStatusOptions.value;
    return options?.find((item) => item.value === value)?.label ?? value ?? '-';
  };

  const getRuleStatusType = (value) => {
    return ruleStatusOptions.value?.find((item) => item.value === value)?.type ?? 'info';
  };

  const columns = ref([
    { type: 'selection', columnKey: 'selection', width: 50, align: 'center', fixed: 'left' },
    { type: 'index', columnKey: 'index', width: 50, align: 'center', fixed: 'left' },
    {
      prop: 'ruleNo',
      label: '规则ID',
      minWidth: 130,
      fixed: 'left',
      slot: 'ruleNo'
    },
    { prop: 'mchNo', label: '商户编号', minWidth: 140 },
    { prop: 'appNo', label: '应用编号', minWidth: 140 },
    { prop: 'ruleName', label: '规则名称', minWidth: 180 },
    {
      columnKey: 'ruleStatusSwitch',
      prop: 'ruleStatus',
      label: '状态开关',
      width: 100,
      align: 'center',
      slot: 'ruleStatus',
      formatter: (row) => getOptionLabel(row.ruleStatus)
    },
    {
      columnKey: 'ruleStatusTag',
      prop: 'ruleStatus',
      label: '规则状态',
      width: 110,
      align: 'center',
      slot: 'statusTag',
      formatter: (row) => getOptionLabel(row.ruleStatus)
    },
    {
      prop: 'ruleAttribute',
      label: '规则属性',
      width: 130,
      slot: 'ruleAttribute',
      formatter: (row) => getOptionLabel(row.ruleAttribute, 'ruleAttribute')
    },
    { prop: 'ruleDesc', label: '规则描述', minWidth: 220 },
    { prop: 'priority', label: '优先级', width: 90, align: 'center' },
    { prop: 'createTime', label: '创建时间', width: 180, align: 'center' },
    {
      columnKey: 'action',
      label: '操作',
      width: 150,
      align: 'center',
      slot: 'action',
      fixed: 'right',
      hideInPrint: true,
      hideInExport: true
    }
  ]);

  const selectSummary = (key, where) => {
    activeSummaryKey.value = key;
    reload(where, 1);
  };

  const summaryCards = computed(() => [
    {
      key: 'all',
      label: '规则总数',
      value: summary.value.total ?? 0,
      desc: '全部规则',
      type: 'primary',
      onClick: () => selectSummary('all', {})
    },
    {
      key: 'active',
      label: '激活中',
      value: summary.value.active ?? 0,
      desc: '参与路由',
      type: 'success',
      onClick: () => selectSummary('active', { ruleStatus: 'active' })
    },
    {
      key: 'inactive',
      label: '未激活',
      value: summary.value.inactive ?? 0,
      desc: '允许编辑',
      type: 'info',
      onClick: () => selectSummary('inactive', { ruleStatus: 'inactive' })
    },
    {
      key: 'single',
      label: '单商户规则',
      value: summary.value.singleMerchant ?? 0,
      desc: '绑定账号',
      type: 'warning',
      onClick: () => selectSummary('single', { ruleAttribute: 'single_merchant' })
    },
    {
      key: 'multi',
      label: '多商户规则',
      value: summary.value.multiMerchant ?? 0,
      desc: '自动分配',
      type: 'danger',
      onClick: () => selectSummary('multi', { ruleAttribute: 'multi_merchant' })
    }
  ]);

  const datasource = async ({ pages, where, orders }) => {
    const data = await pageRoutingRules({ ...where, ...orders, ...pages });
    summary.value = data.summary ?? {};
    return {
      rows: data.rows ?? [],
      total: data.total ?? data.rows?.length ?? 0
    };
  };

  const reload = (where, page) => {
    tableRef.value?.reload?.({ where, page });
  };

  const handleToggleStatus = (row, value) => {
    const nextStatus = value ? 'active' : 'inactive';
    const label = value ? '激活' : '停用';
    ElMessageBox.confirm(
      `确认${label}路由规则「${row.ruleName}」？`,
      '系统提示',
      { type: 'warning', draggable: true }
    )
      .then(() => {
        const loading = EleMessage.loading({ message: '请求中..', plain: true });
        toggleRoutingRuleStatus(row.id, nextStatus)
          .then(() => {
            loading.close();
            EleMessage.success({ message: '规则状态已更新', plain: true });
            reload();
          })
          .catch((e) => {
            loading.close();
            EleMessage.error({ message: e.message, plain: true });
            reload();
          });
      })
      .catch(() => {
        reload();
      });
  };

  /** 新增规则 */
  const handleAdd = () => {
    openModal({
      custom: true,
      asyncComponent: () => import('./components/routing-edit.vue'),
      componentProps: { onDone: () => reload() }
    });
  };

  /** 编辑规则 */
  const handleEdit = (row) => {
    if (row?.ruleStatus === 'active') {
      EleMessage.warning({ message: '当前状态下不支持编辑', plain: true });
      return;
    }
    openModal({
      custom: true,
      asyncComponent: () => import('./components/routing-edit.vue'),
      componentProps: { data: row, onDone: () => reload() }
    });
  };

  /** 查看详情 */
  const handleDetail = (row) => {
    openModal({
      custom: true,
      asyncComponent: () => import('./components/routing-detail.vue'),
      componentProps: { data: row }
    });
  };
</script>

<style lang="scss" scoped>
  .page-section {
    margin-bottom: 16px;
  }

  .table-section {
    margin-bottom: 0;
  }

  .routing-summary {
    --summary-primary: var(--el-color-primary);
    --summary-success: var(--el-color-success);
    --summary-info: var(--el-color-info);
    --summary-warning: var(--el-color-warning);
    --summary-danger: var(--el-color-danger);
  }

  .summary-card {
    position: relative;
    margin-bottom: 0 !important;
    cursor: pointer;
    overflow: hidden;
    border: 1px solid transparent;
    transition:
      border-color 0.2s ease,
      box-shadow 0.2s ease,
      background-color 0.2s ease;
  }

  .routing-summary :deep(.ele-card) {
    margin-bottom: 0 !important;
  }

  .summary-card::after {
    position: absolute;
    right: -26px;
    bottom: -36px;
    width: 92px;
    height: 92px;
    border-radius: 50%;
    background: var(--summary-color);
    opacity: 0.06;
    content: '';
  }

  .summary-card:hover,
  .summary-card.is-active {
    border-color: var(--summary-color);
    box-shadow: var(--el-box-shadow-light);
  }

  .summary-card.is-primary {
    --summary-color: var(--summary-primary);
  }

  .summary-card.is-success {
    --summary-color: var(--summary-success);
  }

  .summary-card.is-info {
    --summary-color: var(--summary-info);
  }

  .summary-card.is-warning {
    --summary-color: var(--summary-warning);
  }

  .summary-card.is-danger {
    --summary-color: var(--summary-danger);
  }

  .summary-head {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .summary-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--summary-color);
  }

  .summary-label {
    color: var(--el-text-color-secondary);
    font-size: 13px;
  }

  .summary-value {
    margin-top: 10px;
    color: var(--el-text-color-primary);
    font-size: 24px;
    font-weight: 700;
    line-height: 1.2;
  }

  .summary-desc {
    margin-top: 6px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }
</style>
