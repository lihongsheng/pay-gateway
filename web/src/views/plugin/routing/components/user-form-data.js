import { reactive } from 'vue';
import { cloneDeep, set as setValue } from 'lodash-es';

/**
 * 表单数据hook
 * @param initValue 初始值
 */
export function useFormData(init) {
  let initValue = init == null ? {} : cloneDeep(init);
  /** 表单数据 */
  const form = reactive(init == null ? {} : cloneDeep(init));

  /** 重置为初始值 */
  const resetFields = (field, excludes) => {
    const keys = Object.keys(form);
    if (typeof field === 'string' && field) {
      if (keys.includes(field)) {
        form[field] = cloneDeep(initValue[field]);
      }
      return;
    }
    const initValueClone = cloneDeep(initValue);
    keys.forEach((key) => {
      if (!excludes || !excludes.includes(key)) {
        form[key] = initValueClone[key];
      }
    });
  };

  /** 赋值不改变字段 */
  const assignFields = (data, excludes) => {
    if (excludes === true) {
      initValue = data == null ? {} : cloneDeep(data);
      const initValueKeys = Object.keys(initValue);
      Object.keys(form).forEach((key) => {
        form[key] = void 0;
        if (!initValueKeys.includes(key)) {
          delete form[key];
        }
      });
      if (data != null) {
        Object.assign(form, cloneDeep(data));
      }
      return;
    }
    Object.keys(form).forEach((key) => {
      if (!excludes || !excludes.includes(key)) {
        form[key] = data?.[key];
      }
    });
  };

  /** 赋值某字段 */
  const setFieldValue = (field, value) => {
    setValue(form, field, value);
  };

  const result = [form, resetFields, assignFields, setFieldValue];
  // 支持对象解构以兼容旧版
  Object.assign(result, {
    form,
    resetFields,
    assignFields
  });
  return result;
}


/**
 * 表单校验失败时滚动到第一个错误项并聚焦
 * @param formRef 表单组件 ref 或 ref.value
 * @param invalidFields Element Plus validate 返回的错误字段
 * @param rootSelector 限定查找范围
 */
export function scrollToFirstFormError(formRef, invalidFields, rootSelector) {
  const form = formRef?.value ?? formRef;
  const firstField = Object.keys(invalidFields || {})[0];
  if (firstField) {
    form?.scrollToField?.(firstField);
  }

  setTimeout(() => {
    const root = rootSelector ? document.querySelector(rootSelector) : document;
    const errorEl = root?.querySelector?.('.el-form-item.is-error');
    errorEl?.scrollIntoView?.({ behavior: 'smooth', block: 'center' });
    errorEl?.querySelector?.('input, textarea, [tabindex]')?.focus?.();
  }, 0);
}
