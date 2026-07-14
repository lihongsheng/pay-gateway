<template>
  <ele-drawer
    :size="960"
    :title="isUpdate ? '编辑路由规则' : '新建路由规则'"
    :loading="loading || detailLoading || schemaLoading"
    :body-style="{ padding: '16px 20px' }"
    :footer-style="{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }"
    v-bind="modalProps"
    :close-on-click-modal="false"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      @submit.prevent=""
    >
      <section class="edit-card">
        <div class="card-head">
          <div class="card-title">
            <span class="section-bar"></span>
            <span>基础信息</span>
            <em>填写规则的基本属性</em>
          </div>
        </div>
        <div class="card-body">
          <el-row :gutter="16">
            <el-col :xs="24" :md="12">
              <el-form-item label="规则名称" prop="ruleName">
                <el-input
                  v-model.trim="form.ruleName"
                  clearable
                  maxlength="50"
                  show-word-limit
                  placeholder="请输入规则名称，如：AE卡组美金收单规则"
                />
              </el-form-item>
            </el-col>
            <el-col :xs="24" :md="6">
              <el-form-item label="规则属性" prop="ruleAttribute">
                <el-select
                  v-model="form.ruleAttribute"
                  class="ele-fluid"
                  placeholder="请选择规则属性"
                  @change="handleAttributeChange"
                >
                  <el-option
                    v-for="item in ruleAttributeOptions"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :xs="24" :md="6">
              <el-form-item label="优先级">
                <el-input-number
                  v-model="form.priority"
                  :min="1"
                  :max="9999"
                  controls-position="right"
                  class="ele-fluid"
                />
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="规则描述" prop="ruleDesc">
                <el-input
                  v-model.trim="form.ruleDesc"
                  maxlength="500"
                  show-word-limit
                  :rows="4"
                  type="textarea"
                  placeholder="请输入规则描述，说明规则的用途和匹配逻辑..."
                />
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </section>

      <section class="edit-card">
        <div class="card-head">
          <div class="card-title">
            <span class="section-bar"></span>
            <span>规则配置</span>
            <em>设置匹配条件与执行动作</em>
          </div>
        </div>
        <div class="card-body rule-builder-section">
          <div class="builder-block">
            <div class="builder-title">
              <span>匹配条件</span>
              <em>支持多层嵌套的 AND/OR 条件组合</em>
            </div>
            <RuleBuilder
              v-model="form.rules"
              :external-schema="schema"
              :disabled="false"
            />
          </div>

          <div class="then-divider"><span>THEN</span></div>

          <div class="builder-block action-block">
            <div class="builder-title">
              <span>执行动作</span>
              <em>命中规则后的收单账号选择方式</em>
            </div>
            <div class="action-row">
              <div class="action-label">收单账号</div>
              <el-select
                :model-value="form.ruleAttribute === 'single_merchant' ? '指定为' : '自动分配'"
                disabled
                class="action-operator"
              >
                <el-option label="指定为" value="指定为" />
                <el-option label="自动分配" value="自动分配" />
              </el-select>
              <el-select
                v-if="form.ruleAttribute === 'single_merchant'"
                v-model="form.actions[0].paypalAccountNo"
                filterable
                clearable
                class="action-value"
                placeholder="请选择单商户收单账号"
              >
                <el-option
                  v-for="account in availableAccounts"
                  :key="account.accountNo"
                  :label="`${account.accountNo} / ${account.accountEmail}`"
                  :value="account.accountNo"
                />
              </el-select>
              <el-input
                v-else
                model-value="由路由引擎按账号属性、额度、风险和评分自动分配"
                disabled
                class="action-value"
              />
            </div>
          </div>
        </div>
      </section>
    </el-form>

    <template #footer>
      <btn-items
        :items="[
          { preset: 'cancel', onClick: () => handleCancel() },
          { preset: 'save', onClick: () => handleSave() }
        ]"
      />
    </template>
  </ele-drawer>
</template>

<script setup>
  import { computed, reactive, ref, onUnmounted } from 'vue';
  import { EleMessage, useModal } from 'ele-admin-plus';
  import { scrollToFirstFormError } from '@/utils/common';
  import {
    addRoutingRule,
    updateRoutingRule,
    getRoutingRule,
    listAvailableSingleAccounts
  } from '@/api/routing';
  import { RuleBuilder, useRuleSchema, validateRuleJson, createGroupRule } from '@/components/RuleBuilder';

  const props = defineProps({
    /** 编辑时传入行数据，新增时为空 */
    data: Object,
    /** 保存成功后的回调 */
    onDone: Function
  });

  const { modalProps, closeModal } = useModal();
  const isUpdate = computed(() => !!props.data?.id);
  const loading = ref(false);
  const detailLoading = ref(false);
  const formRef = ref(null);
  const availableAccounts = ref([]);

  // 组件是否已卸载的标志，防止卸载后异步回调修改状态
  let unmounted = false;
  onUnmounted(() => {
    unmounted = true;
  });

  const { schema, loading: schemaLoading, ruleAttributeOptions } = useRuleSchema();

  const createAction = () => ({
    actionType: 'select_paypal_account',
    accountAttribute: 'multi_merchant',
    paypalAccountNo: ''
  });

  const form = reactive({
    id: void 0,
    ruleName: '',
    ruleDesc: '',
    ruleAttribute: 'multi_merchant',
    priority: 100,
    rules: createGroupRule(),
    actions: [createAction()]
  });

  const rules = reactive({
    ruleName: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
    ruleDesc: [{ required: true, message: '请输入规则描述', trigger: 'blur' }],
    ruleAttribute: [{ required: true, message: '请选择规则属性', trigger: 'change' }]
  });

  const loadAvailableAccounts = () => {
    listAvailableSingleAccounts({ ruleId: form.id })
      .then((data) => {
        if (unmounted) return;
        availableAccounts.value = data ?? [];
      })
      .catch((e) => {
        if (unmounted) return;
        EleMessage.error({ message: e.message, plain: true });
      });
  };

  /**
   * 旧操作符名 → 新操作符名映射
   */
  const OPERATOR_MAP = {
    eq: 'equal',
    ne: 'not_equal',
    gt: 'greater',
    gte: 'greater_or_equal',
    lt: 'less',
    lte: 'less_or_equal',
    contains: 'contains',
    is_empty: 'is_null',
    is_not_empty: 'is_not_null'
  };

  const normalizeOperator = (op) => {
    if (!op) return 'equal';
    return OPERATOR_MAP[op] || op;
  };

  const assignData = (data) => {
    form.id = data.id;
    form.ruleName = data.ruleName ?? '';
    form.ruleDesc = data.ruleDesc ?? '';
    form.ruleAttribute = data.ruleAttribute ?? 'multi_merchant';
    form.priority = data.priority ?? 100;

    // 优先使用新格式 rules（QueryBuilder JSON）
    if (data.rules && data.rules.condition) {
      form.rules = data.rules;
    } else if (data.conditionGroups?.length) {
      // 兼容旧格式：从 conditionGroups 构建 QueryBuilder JSON
      form.rules = convertLegacyToRuleJson(data);
    } else {
      form.rules = createGroupRule();
    }

    form.actions = (data.actions?.length ? data.actions : [createAction()]).map((item) => ({
      ...createAction(),
      ...item,
      accountAttribute: form.ruleAttribute
    }));
    if (!form.actions.length) {
      form.actions = [createAction()];
    }
    handleAttributeChange();
  };

  /**
   * 将旧格式 conditionGroups 转换为 QueryBuilder JSON
   * 同时将旧操作符名(eq/ne/gt等)转换为新操作符名(equal/not_equal/greater等)
   */
  const convertLegacyToRuleJson = (data) => {
    const groups = data.conditionGroups || [];
    const outerLogic = data.conditionLogic || 'AND';

    if (groups.length === 0) {
      return createGroupRule();
    }

    const convertCondition = (c) => {
      const operator = normalizeOperator(c.operator);
      let value = c.fieldValue ?? '';
      // in/not_in 操作符值应为数组
      if ((operator === 'in' || operator === 'not_in') && typeof value === 'string' && value.includes(',')) {
        value = value.split(',').map((v) => v.trim()).filter(Boolean);
      }
      return { id: c.fieldKey, operator, value };
    };

    if (groups.length === 1) {
      const group = groups[0];
      return {
        condition: group.logic || outerLogic,
        rules: (group.conditions || []).map(convertCondition)
      };
    }

    return {
      condition: outerLogic,
      rules: groups.map((group) => ({
        condition: group.logic || 'AND',
        rules: (group.conditions || []).map(convertCondition)
      }))
    };
  };

  const handleAttributeChange = () => {
    form.actions = [
      {
        ...createAction(),
        accountAttribute: form.ruleAttribute,
        paypalAccountNo:
          form.ruleAttribute === 'single_merchant'
            ? form.actions?.[0]?.paypalAccountNo || ''
            : ''
      }
    ];
    if (form.ruleAttribute === 'single_merchant') {
      loadAvailableAccounts();
    }
  };

  const validateConfig = () => {
    const errors = validateRuleJson(form.rules);
    if (errors.length > 0) {
      return errors[0];
    }
    if (form.ruleAttribute === 'single_merchant' && !form.actions[0]?.paypalAccountNo) {
      return '请选择单商户收单账号';
    }
  };

  const handleSave = () => {
    formRef.value?.validate?.((valid, invalidFields) => {
      if (!valid) {
        scrollToFirstFormError(formRef, invalidFields, '.ele-drawer-body');
        return;
      }
      const error = validateConfig();
      if (error) {
        EleMessage.error({ message: error, plain: true });
        document.querySelector('.ele-drawer-body .rule-builder-section')?.scrollIntoView?.({
          behavior: 'smooth',
          block: 'center'
        });
        return;
      }
      loading.value = true;
      const payload = {
        ruleName: form.ruleName,
        ruleDesc: form.ruleDesc,
        ruleAttribute: form.ruleAttribute,
        priority: form.priority,
        rules: form.rules,
        actions: form.actions
      };
      const saveOrUpdate = isUpdate.value
        ? updateRoutingRule(form.id, payload)
        : addRoutingRule(payload);
      saveOrUpdate
        .then((msg) => {
          if (unmounted) return;
          loading.value = false;
          EleMessage.success({ message: msg, plain: true });
          closeModal();
          props.onDone?.();
        })
        .catch((e) => {
          if (unmounted) return;
          loading.value = false;
          EleMessage.error({ message: e.message, plain: true });
        });
    });
  };

  const handleCancel = () => {
    closeModal();
  };

  // 编辑模式：加载详情数据
  if (props.data?.id) {
    detailLoading.value = true;
    getRoutingRule(props.data.id)
      .then((data) => {
        if (unmounted) return;
        assignData(data);
        detailLoading.value = false;
      })
      .catch((e) => {
        if (unmounted) return;
        detailLoading.value = false;
        EleMessage.error({ message: e.message, plain: true });
      });
  } else {
    handleAttributeChange();
  }
</script>

<style lang="scss" scoped>
  .edit-card {
    margin-bottom: 16px;
    border-radius: 12px;
    background: var(--el-bg-color);
    box-shadow: 0 8px 24px rgb(15 23 42 / 5%);
    transition:
      box-shadow 0.2s ease,
      transform 0.2s ease;
  }

  .edit-card:hover,
  .edit-card:focus-within {
    box-shadow: 0 12px 30px rgb(15 23 42 / 8%);
    transform: translateY(-1px);
  }

  .card-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 18px 20px 4px;
  }

  .card-title {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--el-text-color-primary);
    font-size: 15px;
    font-weight: 700;
  }

  .card-title em {
    color: var(--el-text-color-placeholder);
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
  }

  .section-bar {
    width: 4px;
    height: 20px;
    border-radius: 4px;
    background: var(--el-color-primary);
  }

  .card-body {
    padding: 12px 20px 20px;
  }

  .rule-builder-section {
    padding-top: 22px;
  }

  .builder-block {
    position: relative;
  }

  .builder-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 14px;
    color: var(--el-text-color-primary);
    font-weight: 600;
  }

  .builder-title em {
    color: var(--el-text-color-secondary);
    font-size: 12px;
    font-style: normal;
    font-weight: 400;
  }

  .then-divider {
    display: flex;
    align-items: center;
    max-width: 920px;
    margin: 28px 0 24px;
    color: var(--el-text-color-placeholder);
    font-size: 12px;
    font-weight: 700;
  }

  .then-divider::before,
  .then-divider::after {
    flex: 1;
    height: 1px;
    background: var(--el-border-color-lighter);
    content: '';
  }

  .then-divider span {
    padding: 0 14px;
  }

  .action-row {
    display: flex;
    max-width: 720px;
    align-items: center;
    gap: 12px;
  }

  .action-label {
    width: 112px;
    color: var(--el-text-color-regular);
  }

  .action-operator {
    width: 140px;
  }

  .action-value {
    width: 320px;
  }

  :deep(.el-form-item__label) {
    color: var(--el-text-color-regular);
    font-weight: 600;
  }

  @media (max-width: 768px) {
    .card-head {
      align-items: flex-start;
      flex-direction: column;
    }

    .action-row {
      align-items: stretch;
      flex-direction: column;
    }

    .action-label,
    .action-operator,
    .action-value {
      width: 100%;
    }
  }
</style>
