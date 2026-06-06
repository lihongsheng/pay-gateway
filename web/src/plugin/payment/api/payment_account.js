import request from '@/utils/request'

// 获取应用下的支付渠道列表
export function getAccountList(params) {
  return request({
    url: '/private/v1/payment_account/list',
    method: 'get',
    params
  })
}

// 获取应用可用的支付渠道及其配置信息
export function getChannelConfig(params) {
  return request({
    url: '/private/v1/payment_account/channel/config',
    method: 'get',
    params
  })
}

// 获取应用可用的支付渠道及其配置信息
export function getChannelPayment(params) {
  return request({
    url: '/private/v1/payment_account/channel/payment',
    method: 'get',
    params
  })
}

// 保存支付渠道配置
export function savePaymentAccount(data) {
  return request({
    url: '/private/v1/payment_account',
    method: 'post',
    data
  })
}


// 获取支付渠道详情
export function getPaymentAccountDetail(id) {
  return request({
    url: `/private/v1/payment_account`,
    method: 'get',
    params: { id: id }
  })
}

// 获取支付的渠道测试二维码
export function getChannelPaymentQrcode(appNO, accountNO) {
  return request({
    url: '/private/v1/payment_account/test/qrcode',
    method: 'get',
    params:{ app_no: appNO, account_no: accountNO}
  })
}