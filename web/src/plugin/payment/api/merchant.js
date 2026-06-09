// web/src/plugin/payment/api/merchant.js
/**
 * 商户管理 API 接口
 * @module api/payment/merchant
 */

import service from '@/utils/request'

/**
 * 获取商户列表
 * @param {Object} params 查询参数
 * @param {string} params.mchName 商户名称
 * @param {string} params.mchNo 商户编号
 * @param {number} params.status 状态
 * @param {number} params.id 商户ID
 * @param {number} params.page 页码
 * @param {number} params.pageSize 每页数量
 * @returns {Promise} 商户列表数据
 */
export const getMerchantList = (params) => {
  return service({
    url: '/private/v1/merchant/search',
    method: 'get',
    params: params
  })
}

/**
 * 获取商户详情
 * @param {number} mch_no 商户ID
 * @returns {Promise} 商户详情
 */
export const getMerchantDetail = (mch_no) => {
  return service({
    url: '/private/v1/merchant',
    method: 'get',
    params: { mch_no: mch_no }
  })
}

/**
 * 创建/更新商户
 * @param {Object} data 商户数据
 * @param {number} data.id 商户ID (可选，有ID表示更新)
 * @param {string} data.mchName 商户名称
 * @param {string} data.linker 联系人
 * @param {string} data.phone 联系电话
 * @param {string} data.email 邮箱
 * @param {number} data.status 状态
 * @param {string} data.address 地址
 * @param {string} data.reason 备注
 * @returns {Promise} 创建/更新结果
 */
export const saveMerchant = (data) => {
  return service({
    url: '/private/v1/merchant',
    method: 'post',
    data: data
  })
}

/**
 * 更新商户状态
 * @param {Object} data 状态更新数据
 * @param {number} data.id 商户ID
 * @param {number} data.status 状态
 * @returns {Promise} 状态更新结果
 */
export const updateMerchantStatus = (data) => {
  return service({
    url: '/private/v1/merchant/status',
    method: 'post',
    data: data
  })
}
