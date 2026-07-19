<template>
  <el-card shadow="never" class="search-card">
    <el-form label-width="84px" @keyup.enter="handleSearch" @submit.prevent="">
      <el-row :gutter="16">
        <el-col :lg="6" :md="12" :sm="12" :xs="24">
          <el-form-item label="快捷日期">
            <el-select v-model="form.dateShortcut" placeholder="全部" clearable>
              <el-option label="近3天" value="last_3_days" />
              <el-option label="近7天" value="last_7_days" />
              <el-option label="近30天" value="last_30_days" />
              <el-option label="自定义" value="custom" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col v-if="form.dateShortcut === 'custom'" :lg="8" :md="12" :sm="12" :xs="24">
          <el-form-item label="日期范围">
            <el-date-picker
              v-model="form.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :lg="6" :md="12" :sm="12" :xs="24">
          <el-form-item label="规则名称">
            <el-input
              v-model.trim="form.ruleName"
              clearable
              placeholder="搜索规则名称"
            />
          </el-form-item>
        </el-col>
        <el-col :lg="6" :md="12" :sm="12" :xs="24">
          <el-form-item label="规则状态">
            <el-select v-model="form.ruleStatus" clearable placeholder="全部">
              <el-option
                v-for="item in ruleStatusOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :lg="6" :md="12" :sm="12" :xs="24">
          <el-form-item label="规则属性">
            <el-select v-model="form.ruleAttribute" clearable placeholder="全部">
              <el-option
                v-for="item in ruleAttributeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :lg="6" :md="12" :sm="12" :xs="24">
          <el-form-item label-width="0px">
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>搜索
            </el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </el-card>
</template>

<script setup>
  import { Search } from '@element-plus/icons-vue';
  import { useFormData } from './user-form-data.js';
  import { useRuleSchema } from '../RuleBuilder';

  const emit = defineEmits(['search']);

  const { ruleStatusOptions, ruleAttributeOptions } = useRuleSchema();

  const [form, resetFields] = useFormData({
    dateShortcut: '',
    dateRange: [],
    ruleName: '',
    ruleStatus: '',
    ruleAttribute: ''
  });

  const shortcutDays = {
    last_3_days: 3,
    last_7_days: 7,
    last_30_days: 30
  };

  const formatDate = (date) => {
    const year = date.getFullYear();
    const month = `${date.getMonth() + 1}`.padStart(2, '0');
    const day = `${date.getDate()}`.padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const buildParams = () => {
    if (form.dateShortcut === 'custom' && form.dateRange?.length === 2) {
      return { beginTime: form.dateRange[0], endTime: form.dateRange[1] };
    }
    const days = shortcutDays[form.dateShortcut];
    if (!days) {
      return void 0;
    }
    const end = new Date();
    const begin = new Date();
    begin.setDate(end.getDate() - days + 1);
    return { beginTime: formatDate(begin), endTime: formatDate(end) };
  };

  const handleSearch = () => {
    const where = { ...form };
    const params = buildParams();
    if (params) {
      where.params = params;
    }
    delete where.dateRange;
    emit('search', where);
  };

  const handleReset = () => {
    resetFields();
    handleSearch();
  };
</script>

<style lang="scss" scoped>
  .search-card {
    margin-bottom: 16px;
  }
</style>
