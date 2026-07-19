<template>
  <div class="routing-page">
    <el-row :gutter="16" class="routing-summary page-section">
      <el-col
        v-for="item in summaryCards"
        :key="item.key"
        :xs="24"
        :sm="12"
        :md="6"
        :lg="4"
      >
        <el-card
          :body-style="{ padding: '16px 18px' }"
          class="summary-card"
          :class="[`is-${item.type}`, { 'is-active': activeSummaryKey === item.key }]"
          shadow="hover"
          @click="item.onClick?.()"
        >
          <div class="summary-head">
            <span class="summary-dot"></span>
            <span class="summary-label">{{ item.label }}</span>
          </div>
          <div class="summary-value">{{ item.value }}</div>
          <div class="summary-desc">{{ item.desc }}</div>
        </el-card>
      </el-col>
    </el-row>

    <routing-search @search="(where) => reload(where, 1)" />

    <el-card class="table-section" :body-style="{ paddingTop: '8px' }" shadow="never">
      <div class="table-toolbar">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>新增规则
        </el-button>
      </div>

      <el-table
        ref="tableRef"
        v-loading="tableLoading"
        :data="tableData"
        row-key="id"
        :show-overflow-tooltip="true"
        :highlight-current-row="true"
        @selection-change="selections = $event"
      >
        <el-table-column type="selection" width="50" align="center" fixed="left" />
        <el-table-column type="index" label="#" width="50" align="center" fixed="left" />
        <el-table-column prop="ruleNo" label="规则ID" min-width="130" fixed="left">
          <template #default="{ row }">
            <el-link type="primary" @click="handleDetail(row)">{{ row.ruleNo }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="mchNo" label="商户编号" min-width="140" />
        <el-table-column prop="appNo" label="应用编号" min-width="140" />
        <el-table-column prop="ruleName" label="规则名称" min-width="180" />
        <el-table-column label="状态开关" width="100" align="center">
          <template #default="{ row }">
            <el-tooltip :content="getOptionLabel(row.ruleStatus)" placement="top">
              <el-switch
                size="small"
                :model-value="row.ruleStatus === 'active'"
                @change="(value) => handleToggleStatus(row, value)"
              />
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="规则状态" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getRuleStatusType(row.ruleStatus)">
              {{ getOptionLabel(row.ruleStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="规则属性" width="130">
          <template #default="{ row }">
            {{ getOptionLabel(row.ruleAttribute, 'ruleAttribute') }}
          </template>
        </el-table-column>
        <el-table-column prop="ruleDesc" label="规则描述" min-width="220" />
        <el-table-column prop="priority" label="优先级" width="90" align="center" />
        <el-table-column prop="createTime" label="创建时间" width="180" align="center" />
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click.stop="handleDetail(row)">查看</el-button>
            <el-button type="primary" link @click.stop="handleEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.limit"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="onPageSizeChange"
          @current-change="onPageChange"
        />
      </div>
    </el-card>

    <!-- 编辑抽屉 -->
    <routing-edit
      v-if="editVisible"
      :visible="editVisible"
      :data="editData"
      @close="editVisible = false"
      @done="onEditDone"
    />

    <!-- 详情弹窗 -->
    <routing-detail
      v-if="detailVisible"
      :visible="detailVisible"
      :data="detailData"
      @close="detailVisible = false"
    />
  </div>
</template>

<script setup>
  import { computed, ref, onMounted } from 'vue';
  import { ElMessage, ElMessageBox } from 'element-plus';
  import { Plus } from '@element-plus/icons-vue';
  import {
    pageRoutingRules,
    toggleRoutingRuleStatus
  } from './api/routing';
  import RoutingSearch from './components/routing-search.vue';
  import RoutingEdit from './components/routing-edit.vue';
  import RoutingDetail from './components/routing-detail.vue';
  import { useRuleSchema } from './RuleBuilder';

  defineOptions({ name: 'RoutingRule' });

  const tableRef = ref(null);
  const selections = ref([]);
  const summary = ref({});
  const activeSummaryKey = ref('all');
  const tableLoading = ref(false);
  const tableData = ref([]);
  const currentWhere = ref({});
  const pagination = ref({ page: 1, limit: 10, total: 0 });

  // 编辑抽屉
  const editVisible = ref(false);
  const editData = ref(null);

  // 详情弹窗
  const detailVisible = ref(false);
  const detailData = ref(null);

  const { ruleStatusOptions, ruleAttributeOptions } = useRuleSchema();

  const getOptionLabel = (value, type = 'ruleStatus') => {
    const options = type === 'ruleAttribute' ? ruleAttributeOptions.value : ruleStatusOptions.value;
    return options?.find((item) => item.value === value)?.label ?? value ?? '-';
  };

  const getRuleStatusType = (value) => {
    return ruleStatusOptions.value?.find((item) => item.value === value)?.type ?? 'info';
  };

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

  const fetchData = async () => {
    tableLoading.value = true;
    try {
      const { page, limit } = pagination.value;
      const data = await pageRoutingRules({
        ...currentWhere.value,
        page,
        limit
      });
      tableData.value = data.rows ?? [];
      pagination.value.total = data.total ?? data.rows?.length ?? 0;
      summary.value = data.summary ?? {};
    } catch (e) {
      ElMessage.error(e.message);
    } finally {
      tableLoading.value = false;
    }
  };

  const reload = (where, page) => {
    if (where !== undefined) currentWhere.value = where;
    if (page !== undefined) pagination.value.page = page;
    fetchData();
  };

  const onPageSizeChange = (size) => {
    pagination.value.limit = size;
    pagination.value.page = 1;
    fetchData();
  };

  const onPageChange = (page) => {
    pagination.value.page = page;
    fetchData();
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
        const loading = ElMessage({ message: '请求中..', type: 'loading' });
        toggleRoutingRuleStatus(row.id, nextStatus)
          .then(() => {
            loading.close();
            ElMessage.success('规则状态已更新');
            reload();
          })
          .catch((e) => {
            loading.close();
            ElMessage.error(e.message);
            reload();
          });
      })
      .catch(() => {
        reload();
      });
  };

  /** 新增规则 */
  const handleAdd = () => {
    editData.value = null;
    editVisible.value = true;
  };

  /** 编辑规则 */
  const handleEdit = (row) => {
    if (row?.ruleStatus === 'active') {
      ElMessage.warning('当前状态下不支持编辑');
      return;
    }
    editData.value = row;
    editVisible.value = true;
  };

  /** 查看详情 */
  const handleDetail = (row) => {
    detailData.value = row;
    detailVisible.value = true;
  };

  /** 编辑完成回调 */
  const onEditDone = () => {
    editVisible.value = false;
    reload();
  };

  onMounted(() => {
    fetchData();
  });
</script>

<style lang="scss" scoped>
  .routing-page {
    min-height: 100%;
  }

  .page-section {
    margin-bottom: 16px;
  }

  .table-section {
    margin-bottom: 0;
  }

  .table-toolbar {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    padding-bottom: 12px;
  }

  .table-pagination {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
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
