/**
 * 应用管理 API 接口
 * @module api/payment/application
 */

import service from '@/utils/request'

/**
 * 获取应用列表
 * @param {Object} params 查询参数
 * @param {string} params.appNo 应用编号
 * @param {string} params.name 应用名称
 * @param {string} params.mchNo 商户编号
 * @param {number} params.status 状态
 * @param {number} params.page 页码
 * @param {number} params.pageSize 每页数量
 * @returns {Promise} 应用列表数据
 */
export const getApplicationList = (params) => {
  return service({
    url: '/private/v1/application/search',
    method: 'get',
    params: params
  })
}

/**
 * 创建应用
 * @param {Object} data 应用数据
 * @param {string} data.name 应用名称
 * @param {string} data.mchNo 商户编号
 * @param {string} data.notifyUrl 通知URL
 * @param {string} data.refundNotifyUrl 退款通知URL
 * @param {number} data.status 状态
 * @returns {Promise} 创建结果
 */
export const createApplication = (data) => {
  return service({
    url: '/private/v1/application',
    method: 'post',
    data: data
  })
}

/**
 * 更新应用
 * @param {Object} data 应用数据
 * @param {number} data.id 应用ID
 * @param {string} data.name 应用名称
 * @param {string} data.mchNo 商户编号
 * @param {string} data.notifyUrl 通知URL
 * @param {string} data.refundNotifyUrl 退款通知URL
 * @param {number} data.status 状态
 * @returns {Promise} 更新结果
 */
export const updateApplication = (data) => {
  return service({
    url: '/private/v1/application',
    method: 'post',
    data: data
  })
}

/**
 * 删除应用
 * @param {number} id 应用ID
 * @returns {Promise} 删除结果
 */
export const deleteApplication = (id) => {
  return service({
    url: `/private/v1/application/${id}`,
    method: 'delete'
  })
}

/**
 * 获取应用详情
 * @param {string} appNo 应用ID
 * @returns {Promise} 应用详情
 */
export const getApplicationDetail = (appNo) => {
  return service({
    url: `/private/v1/application`,
    method: 'get',
    params: { app_no: appNo }
  })
}

/**
 * 更新应用状态
 * @param {Object} data 状态更新数据
 * @param {number} data.id 应用ID
 * @param {number} data.status 状态
 * @returns {Promise} 状态更新结果
 */
export const updateApplicationStatus = (data) => {
  return service({
    url: '/private/v1/application/status',
    method: 'post',
    data: data
  })
}

export const getApplicationQrcode = (appNo) => {
  return service({
    url: '/private/v1/application/qrcode',
    method: 'get',
    params: { app_no: appNo }
  })
}