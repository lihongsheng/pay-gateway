import { ref, computed } from 'vue';
import { getRoutingSchema as getRuleSchema } from '../api/routing';

/**
 * 全局单例 Schema 状态 — 所有组件共享，避免重复请求和卸载后的异步回调问题。
 * 只加载一次，后续调用直接返回缓存结果。
 */
let _schema = null;
let _loading = false;
let _loaded = false;
let _error = null;

const schemaRef = ref(null);
const loadingRef = ref(false);
const errorRef = ref(null);

/**
 * 确保 Schema 只加载一次。
 */
async function ensureSchema() {
  if (_loaded || _loading) return;
  _loading = true;
  loadingRef.value = true;
  _error = null;
  errorRef.value = null;
  try {
    _schema = await getRuleSchema();
    schemaRef.value = _schema;
    _loaded = true;
  } catch (e) {
    _error = e.message;
    errorRef.value = e.message;
  } finally {
    _loading = false;
    loadingRef.value = false;
  }
}

// 首次 import 时即触发加载，不依赖 onMounted
ensureSchema();

/**
 * Composable: 获取规则 Schema（全局单例，无需 onMounted）。
 */
export function useRuleSchema() {
  const schema = computed(() => schemaRef.value);
  const loading = computed(() => loadingRef.value);
  const error = computed(() => errorRef.value);

  const fields = computed(() => schema.value?.fields ?? []);
  const operators = computed(() => schema.value?.operators ?? []);
  const ruleStatusOptions = computed(() => schema.value?.ruleStatusOptions ?? []);
  const ruleAttributeOptions = computed(() => schema.value?.ruleAttributeOptions ?? []);

  const fetchSchema = async () => {
    // 强制重新加载
    _loaded = false;
    _schema = null;
    schemaRef.value = null;
    await ensureSchema();
  };

  return {
    schema,
    loading,
    error,
    fetchSchema,
    fields,
    operators,
    ruleStatusOptions,
    ruleAttributeOptions
  };
}
