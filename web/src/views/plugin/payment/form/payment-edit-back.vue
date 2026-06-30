<!-- payment-edit.vue -->
<template>
  <div class="payment-edit-container">
    <!-- 基础信息表单 -->
    <el-form
        ref="formRef"
        :model="formData"
        label-width="120px"
        :rules="formRules"
    >
      <el-form-item label="通道名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入通道名称" />
      </el-form-item>

      <!-- 支付渠道名称显示 -->
      <el-form-item label="支付渠道" v-if="formData.channel_name">
        <el-input v-model="formData.channel_name" disabled />
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input v-model="formData.remark" type="textarea" placeholder="请输入备注" />
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="formData.status">
          <el-radio :label="1">启用</el-radio>
          <el-radio :label="2">停用</el-radio>
        </el-radio-group>
      </el-form-item>
      <!-- max_limit 数字，单日最大支付笔数 -->
      <el-form-item label="支付限制" prop="max_limit">
        <el-input-number v-model="formData.max_limit"   placeholder="单日最大支付笔数" />
      </el-form-item>

      <!-- 动态渲染支付配置表单（基于channel_option） -->
      <el-form-item label="支付配置" prop="channel_config" class="full-width-config">
        <div class="config-form-wrapper" v-if="channelOption?.Options?.length" style="min-width: 100%;">
          <!-- 独立的动态配置表单，有自己的验证规则 -->
          <el-form
              :model="channelConfigForm"
              :rules="dynamicFormRules"
              label-width="150px"
              class="config-sub-form"
              ref="dynamicFormRef"
          >
            <el-form-item
                v-for="option in channelOption.Options"
                :key="option.name"
                :label="option.label"
                :prop="option.name"
            >
              <!-- 文本输入 -->
              <el-input
                  v-if="option.input_type === 'text'"
                  v-model="channelConfigForm[option.name]"
                  :placeholder="`请输入${option.label}`"
                  :disabled="option.readonly"
                  @blur="validateDynamicField(option.name)"
                  clearable
              />

              <!-- 密码输入 -->
              <el-input
                  v-else-if="option.input_type === 'password'"
                  v-model="channelConfigForm[option.name]"
                  type="password"
                  show-password
                  :placeholder="`请输入${option.label}`"
                  :maxlength="option.maxlength || 32"
                  @blur="validateDynamicField(option.name)"
                  clearable
              />

              <!-- 数字输入 -->
              <el-input
                  v-else-if="option.input_type === 'number'"
                  v-model.number="channelConfigForm[option.name]"
                  :placeholder="`请输入${option.label}`"
                  type="number"
                  @blur="validateDynamicField(option.name)"
                  clearable
              />

              <!-- 文本域 -->
              <el-input
                  v-else-if="option.input_type === 'textarea'"
                  v-model="channelConfigForm[option.name]"
                  type="textarea"
                  :rows="6"
                  :placeholder="`请输入${option.label}`"
                  @blur="validateDynamicField(option.name)"
                  clearable
              />

              <!-- 下拉选择 -->
              <el-select
                  v-else-if="option.input_type === 'select' && option.values && option.values.length > 0"
                  v-model="channelConfigForm[option.name]"
                  :placeholder="`请选择${option.label}`"
                  @change="validateDynamicField(option.name)"
                  clearable
              >
                <el-option
                    v-for="item in option.values"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
              </el-select>

              <!-- 单选框 -->
              <el-radio-group
                  v-else-if="option.input_type === 'radio' && option.values && option.values.length > 0"
                  v-model="channelConfigForm[option.name]"
                  @change="validateDynamicField(option.name)"
              >
                <el-radio
                    v-for="item in option.values"
                    :key="item.value"
                    :label="item.value"
                >
                  {{ item.label }}
                </el-radio>
              </el-radio-group>

              <!-- 复选框（Array类型） -->
              <el-checkbox-group
                  v-else-if="option.input_type === 'checkbox' && option.values && option.values.length > 0"
                  v-model="channelConfigForm[option.name]"
                  @change="validateDynamicField(option.name)"
              >
                <el-checkbox
                    v-for="item in option.values"
                    :key="item.value"
                    :label="item.value"
                >
                  {{ item.label }}
                </el-checkbox>
              </el-checkbox-group>

              <!-- 显示验证错误 -->
              <div v-if="validationErrors[option.name]" class="validation-error">
                {{ validationErrors[option.name] }}
              </div>
            </el-form-item>
          </el-form>
        </div>
      </el-form-item>

      <!-- 支付方式选择（整合PaymentMethodSelector功能） -->
      <el-form-item label="支付方式" prop="payment_method" class="full-width-config">
        <div class="payment-method-wrapper" v-if="paymentMethodConfig.length" style="min-width: 100%;">
          <div
              v-for="method in paymentMethodConfig"
              :key="method.method"
              class="method-card"
          >
            <h4>{{ method.label }}</h4>
            <el-checkbox-group
                v-model="selectedPaymentMethods"
                @change="handlePaymentMethodChange"
            >
              <el-checkbox
                  v-for="product in method.product"
                  :key="product.product"
                  :label="`${method.method}_${product.product}`"
                  :checked="product.used"
              >
                {{ product.label }}
              </el-checkbox>
            </el-checkbox-group>
          </div>
        </div>
      </el-form-item>

      <!-- 操作按钮 -->
      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          保存
        </el-button>
        <el-button @click="handleCancel">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { reactive, ref, watch, onMounted, computed, nextTick } from "vue";
import { ElMessage } from "element-plus";
import { accountSave } from '@/api/payment'

// 在 payment-edit.vue 中定义事件
const emit = defineEmits(["saveSuccess", "cancel"]);

// 表单引用
const formRef = ref(null);
const dynamicFormRef = ref(null); // 动态配置表单引用
const submitLoading = ref(false);

// 基础表单数据（全部下划线格式）
const formData = reactive({
  id: 0,
  name: "",
  remark: "",
  app_no: "",
  channel: "",
  channel_name: "",
  status: 1,
  channel_config: "",
  payment_method: [],
  max_limit: 0,
});

// 通道配置表单（动态渲染用，下划线格式）
const channelConfigForm = reactive({});

// 动态表单验证规则
const dynamicFormRules = reactive({});

// 验证错误信息
const validationErrors = reactive({});

// 支付方式配置（回显用）
const paymentMethodConfig = ref([]);
// 选中的支付方式（临时存储，格式：method_product）
const selectedPaymentMethods = ref([]);

// 通道配置选项（从详情接口获取）
const channelOption = ref({});

// 表单校验规则
const formRules = reactive({
  name: [{ required: true, message: "请输入通道名称", trigger: "blur" }],
  status: [{ required: true, message: "请选择状态", trigger: "change" }],
});

// 初始化动态表单验证规则
const initDynamicFormRules = () => {
  if (!channelOption.value?.Options) return;

  // 清空现有规则
  Object.keys(dynamicFormRules).forEach(key => {
    delete dynamicFormRules[key];
  });

  // 为每个选项创建验证规则
  channelOption.value.Options.forEach(option => {
    const rules = [];
    const fieldName = option.name;

    // 必填验证
    if (option.require) {
      rules.push({
        required: true,
        message: `${option.label}是必填项`,
        trigger: ['blur', 'change'],
      });
    }

    // 根据 validate_type 创建验证规则
    switch(option.validate_type) {
      case 'Int':
        rules.push({
          validator: (rule, value, callback) => {
            if (!value && !option.require) {
              callback();
              return;
            }
            if (value !== undefined && value !== null && value !== '') {
              const numValue = Number(value);
              if (isNaN(numValue) || !Number.isInteger(numValue)) {
                callback(new Error(`${option.label}必须是整数`));
              } else {
                callback();
              }
            } else {
              callback();
            }
          },
          trigger: ['blur', 'change'],
        });
        break;

      case 'Email':
        rules.push({
          type: 'email',
          message: `${option.label}格式不正确`,
          trigger: ['blur', 'change'],
        });
        break;

      case 'Phone':
        rules.push({
          pattern: /^1[3-9]\d{9}$/,
          message: `${option.label}必须是11位手机号码`,
          trigger: ['blur', 'change'],
        });
        break;

      case 'Domain':
        rules.push({
          pattern: /^(?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$/,
          message: `${option.label}格式不正确`,
          trigger: ['blur', 'change'],
        });
        break;

      case 'Url':
        rules.push({
          type: 'url',
          message: `${option.label}格式不正确`,
          trigger: ['blur', 'change'],
        });
        break;

      case 'RsaPublic':
        rules.push({
          validator: (rule, value, callback) => {
            if (!value && !option.require) {
              callback();
              return;
            }
            if (value) {
              // 支持两种RSA公钥格式：PKCS#1 和 PKIX
              const pkcs1Pattern = /-----BEGIN RSA PUBLIC KEY-----[\s\S]*?-----END RSA PUBLIC KEY-----/;
              const pkixPattern = /-----BEGIN PUBLIC KEY-----[\s\S]*?-----END PUBLIC KEY-----/;
              if (!pkcs1Pattern.test(value) && !pkixPattern.test(value)) {
                callback(new Error(`${option.label}格式不正确，必须是有效的RSA公钥`));
              } else {
                callback();
              }
            } else {
              callback();
            }
          },
          trigger: ['blur', 'change'],
        });
        break;

      case 'RsaPrivate':
        rules.push({
          validator: (rule, value, callback) => {
            if (!value && !option.require) {
              callback();
              return;
            }
            if (value) {
              if (!value.includes('-----BEGIN PRIVATE KEY-----') ||
                  !value.includes('-----END PRIVATE KEY-----')) {
                callback(new Error(`${option.label}格式不正确，必须是有效的RSA私钥`));
              } else {
                callback();
              }
            } else {
              callback();
            }
          },
          trigger: ['blur', 'change'],
        });
        break;

      case 'RsaCert':
        rules.push({
          validator: (rule, value, callback) => {
            if (!value && !option.require) {
              callback();
              return;
            }
            if (value) {
              if (!value.includes('-----BEGIN CERTIFICATE-----') ||
                  !value.includes('-----END CERTIFICATE-----')) {
                callback(new Error(`${option.label}格式不正确，必须是有效的证书`));
              } else {
                callback();
              }
            } else {
              callback();
            }
          },
          trigger: ['blur', 'change'],
        });
        break;

      case 'Reg':
        if (option.validate_reg && option.validate_reg.trim() !== '') {
          try {
            const pattern = new RegExp(option.validate_reg);
            rules.push({
              validator: (rule, value, callback) => {
                if (!value && !option.require) {
                  callback();
                  return;
                }
                if (value && !pattern.test(value)) {
                  callback(new Error(`${option.label}格式不正确`));
                } else {
                  callback();
                }
              },
              trigger: ['blur', 'change'],
            });
          } catch (e) {
            console.warn(`正则表达式无效: ${option.validate_reg}`, e);
          }
        }
        break;
    }

    // 将规则添加到动态表单规则中
    dynamicFormRules[fieldName] = rules;
  });
};

// 验证动态配置字段
const validateDynamicField = async (fieldName) => {
  if (!dynamicFormRef.value) return;

  try {
    await dynamicFormRef.value.validateField(fieldName, (error) => {
      if (error && error.length > 0) {
        validationErrors[fieldName] = error[0].message;
      } else {
        delete validationErrors[fieldName];
      }
    });
  } catch (error) {
    console.log(`字段 ${fieldName} 验证通过`);
    delete validationErrors[fieldName];
  }
};

// 验证整个动态表单
const validateDynamicForm = async () => {
  if (!dynamicFormRef.value || !channelOption.value?.Options) {
    return true;
  }

  try {
    await dynamicFormRef.value.validate();
    return true;
  } catch (error) {
    console.error('动态表单验证失败:', error);

    // 收集错误信息
    const errors = [];
    if (error && typeof error === 'object') {
      Object.keys(error).forEach(field => {
        if (error[field] && error[field].length > 0) {
          errors.push(error[field][0].message);
        }
      });
    }

    if (errors.length > 0) {
      ElMessage.error({
        message: `请检查以下配置项：\n${errors.join('\n')}`,
        duration: 5000,
      });
    }

    return false;
  }
};

// 处理支付方式变更
const handlePaymentMethodChange = () => {
  // 转换为提交格式：[{"method": "Wxpay", "product": "JSAPI"}]
  formData.payment_method = selectedPaymentMethods.value
      .filter(item => item) // 过滤空值
      .map((item) => {
        const [method, product] = item.split("_");
        return { method, product };
      });
};

// 初始化表单数据
const initForm = (data) => {
  console.log('PaymentEdit: 接收到初始化数据', data);

  try {
    // 重置所有响应式数据
    resetFormData();

    // 基础信息赋值
    if (data) {
      formData.id = Number(data.id) || 0;
      formData.name = data.name || "";
      formData.remark = data.remark || "";
      formData.app_no = data.app_no || "";
      formData.channel = data.channel || "";
      formData.channel_name = data.channel_name || "";
      formData.status = Number(data.status) || 1;
      formData.payment_method = Array.isArray(data.payment_method) ? data.payment_method : [];
      formData.max_limit = Number(data.max_limit) || 0;
      // 处理 channel_option
      channelOption.value = data.channel_option || {};
      console.log('channelOption:', channelOption.value);

      // 处理 channel_config - 使用安全的解析方式
      if (data.channel_config) {
        try {
          const configObj = typeof data.channel_config === 'string'
              ? JSON.parse(data.channel_config)
              : data.channel_config;

          // 清空现有配置
          Object.keys(channelConfigForm).forEach(key => {
            delete channelConfigForm[key];
          });

          // 如果有 channel_option 的 Options，使用它们初始化
          if (channelOption.value?.Options && Array.isArray(channelOption.value.Options)) {
            channelOption.value.Options.forEach(option => {
              channelConfigForm[option.name] = configObj[option.name] || "";
            });
          } else {
            // 否则直接赋值
            Object.assign(channelConfigForm, configObj);
          }

          formData.channel_config = JSON.stringify(configObj);
        } catch (error) {
          console.error('解析 channel_config 失败:', error);
          formData.channel_config = "{}";
        }
      } else {
        // 初始化空的配置项
        if (channelOption.value?.Options && Array.isArray(channelOption.value.Options)) {
          channelOption.value.Options.forEach(option => {
            switch (option.type) {
              case 'String':
                channelConfigForm[option.name] = option?.default || "";
                break ;
              case 'Int':
                  channelConfigForm[option.name] = option?.default || 0;
                break ;
              case'Bool':
                if (option.default === 'true') {
                  channelConfigForm[option.name] = true;
                } else if (option.default === 'false') {
                  channelConfigForm[option.name] = false;
                }
                channelConfigForm[option.name] = option?.default || false;
                break ;
              default:
                  channelConfigForm[option.name] = option?.default || "";
            }
          });
        }
        formData.channel_config = "{}";
      }

      // 处理 payment_method_config
      paymentMethodConfig.value = Array.isArray(data.payment_method_config)
          ? data.payment_method_config
          : [];

      // 初始化选中的支付方式
      selectedPaymentMethods.value = [];
      if (Array.isArray(formData.payment_method) && formData.payment_method.length > 0) {
        formData.payment_method.forEach(item => {
          if (item.method && item.product) {
            selectedPaymentMethods.value.push(`${item.method}_${item.product}`);
          }
        });
      } else if (paymentMethodConfig.value.length > 0) {
        // 从 paymentMethodConfig 中初始化
        paymentMethodConfig.value.forEach(method => {
          if (method.product && Array.isArray(method.product)) {
            method.product.forEach(product => {
              if (product.used) {
                selectedPaymentMethods.value.push(`${method.method}_${product.product}`);
              }
            });
          }
        });
      }

      // 确保支付方式数据同步
      handlePaymentMethodChange();

      // 初始化动态表单验证规则
      nextTick(() => {
        initDynamicFormRules();
      });
    }

    console.log('初始化后的表单数据:', {
      formData,
      channelConfigForm,
      selectedPaymentMethods: selectedPaymentMethods.value,
      channelOption: channelOption.value
    });

  } catch (error) {
    console.error("初始化表单失败:", error);
    ElMessage.error("初始化表单失败");
  }
};

// 重置表单数据
const resetFormData = () => {
  // 重置 formData
  Object.assign(formData, {
    id: 0,
    name: "",
    remark: "",
    app_no: "",
    channel: "",
    channel_name: "",
    status: 1,
    channel_config: "",
    payment_method: [],
    max_limit: 0,
  });

  // 重置 channelConfigForm
  Object.keys(channelConfigForm).forEach(key => {
    delete channelConfigForm[key];
  });

  // 重置其他响应式数据
  paymentMethodConfig.value = [];
  selectedPaymentMethods.value = [];
  channelOption.value = {};

  // 重置验证相关数据
  Object.keys(dynamicFormRules).forEach(key => {
    delete dynamicFormRules[key];
  });
  Object.keys(validationErrors).forEach(key => {
    delete validationErrors[key];
  });
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) {
    console.error('formRef 不存在');
    return;
  }

  try {
    console.log('开始表单验证...');

    // 验证基础表单
    await formRef.value.validate();

    // 验证动态配置表单
    const dynamicValid = await validateDynamicForm();
    if (!dynamicValid) {
      return;
    }

    // 检查动态配置的必填项
    if (channelOption.value?.Options) {
      const missingRequired = [];
      channelOption.value.Options.forEach(option => {
        if (option.require) {
          const value = channelConfigForm[option.name];
          if (value === undefined || value === null || value === '') {
            missingRequired.push(option.label);
          }
        }
      });

      if (missingRequired.length > 0) {
        ElMessage.error(`以下必填项未填写：${missingRequired.join('、')}`);
        return;
      }
    }

    console.log('表单验证通过，准备提交数据');
    submitLoading.value = true;

    // 整合channel_config（下划线格式）
    try {
      // 确保 channelConfigForm 是有效的对象
      const configToSave = { ...channelConfigForm };
      // 移除空值
      Object.keys(configToSave).forEach(key => {
        if (configToSave[key] === undefined || configToSave[key] === null) {
          configToSave[key] = "";
        }
      });
      formData.channel_config = JSON.stringify(configToSave);
    } catch (error) {
      console.error('序列化 channel_config 失败:', error);
      formData.channel_config = "{}";
    }

    // 确保支付方式数据已更新
    handlePaymentMethodChange();

    // 构造提交数据（严格匹配要求格式）
    const submitData = {
      id: Number(formData.id) || 0,
      name: formData.name || "",
      remark: formData.remark || "",
      app_no: formData.app_no || "",
      channel: formData.channel || "",
      channel_name: formData.channel_name || "",
      status: Number(formData.status) || 1,
      channel_config: formData.channel_config,
      payment_method: Array.isArray(formData.payment_method) ? formData.payment_method : [],
      max_limit: Number(formData.max_limit) || 0,
    };

    console.log('提交的数据:', submitData);

    // 调用保存接口
    const res = await accountSave(submitData);
    if (res.code === 0) {
      ElMessage.success("保存成功");
      // 通知父组件刷新列表
      emit("saveSuccess");
      // 关闭编辑弹窗/页面
      handleCancel();
    } else {
      ElMessage.error(res.msg || "保存失败");
    }
  } catch (error) {
    console.error("提交失败:", error);
    // 判断是否是验证错误
    if (error instanceof Error && error.name !== 'ValidationError') {
      ElMessage.error("保存失败");
    }
  } finally {
    submitLoading.value = false;
  }
};

// 取消操作
const handleCancel = () => {
  emit("cancel");
};

// 暴露方法给父组件
defineExpose({
  initForm
});
</script>

<style scoped>
.payment-edit-container {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
}

.config-form-wrapper {
  margin-top: 10px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.config-sub-form {
  margin-top: 10px;
}

.validation-error {
  color: #f56c6c;
  font-size: 12px;
  margin-top: 4px;
  line-height: 1.2;
}

.method-card {
  margin-bottom: 15px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.method-card h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  color: #303133;
}

.el-checkbox {
  margin-right: 15px;
  margin-bottom: 8px;
}

/* 必填项标识 */
:deep(.el-form-item.is-required .el-form-item__label:before) {
  content: '*';
  color: #f56c6c;
  margin-right: 4px;
}

/* 验证错误样式 */
:deep(.el-form-item.is-error .el-input__inner),
:deep(.el-form-item.is-error .el-textarea__inner) {
  border-color: #f56c6c;
}

:deep(.el-form-item.is-error .el-input__inner:focus),
:deep(.el-form-item.is-error .el-textarea__inner:focus) {
  border-color: #f56c6c;
  box-shadow: 0 0 0 2px rgba(245, 108, 108, 0.2);
}
</style>