import request from '@/utils/request'

// 路由规则 Schema
export const getRoutingSchema = () =>
  request.get('/api/plugin/routing/schema')

// 分页查询路由规则
export const pageRoutingRules = (params) => {
  const { limit, prop, order, ...rest } = params || {};
  // ele-pro-table 传入 limit，后端使用 pageSize
  // ele-pro-table 传入 prop/order，后端使用 orderBy/orderDir
  const query = { ...rest, pageSize: limit || rest.pageSize };
  if (prop) {
    query.orderBy = prop;
    query.orderDir = order === 'ascending' ? 'asc' : order === 'descending' ? 'desc' : '';
  }
  return request.get('/api/plugin/routing/rule/search', { params: query });
}

// 获取路由规则详情
export const getRoutingRule = (id) =>
  request.get('/api/plugin/routing/rule', { params: { id } })

// 新增路由规则
export const addRoutingRule = (data) =>
  request.post('/api/plugin/routing/rule', data)

// 更新路由规则
export const updateRoutingRule = (id, data) =>
  request.post('/api/plugin/routing/rule', { ...data, id })

// 切换路由规则状态
export const toggleRoutingRuleStatus = (id, ruleStatus) =>
  request.post('/api/plugin/routing/rule/status', { id, ruleStatus })

// 获取可用的单商户收单账号列表
export const listAvailableSingleAccounts = (params) =>
  request.get('/api/plugin/routing/rule/available-accounts', { params })
