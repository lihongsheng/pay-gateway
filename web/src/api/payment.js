import request from '@/utils/request'

// application
export const appList     = (params) => request.get('/api/plugin/payment/application/search', { params })
export const appDetail   = (params) => request.get('/api/plugin/payment/application', { params })
export const appCreate   = (data)   => request.post('/api/plugin/payment/application', data)
export const appUpdate   = (data)   => request.post('/api/plugin/payment/application', data)
export const appStatus   = (data)   => request.put('/api/plugin/payment/application/status', data)
export const appQrcode   = (params) => request.get('/api/plugin/payment/application/qrcode', { params })
export const appSecret   = ()       => request.get('/api/plugin/payment/application/secret')

// payment account
export const accountDetail  = (params) => request.get('/api/plugin/payment/payment_account', { params })
export const accountSave   = (data)   => request.post('/api/plugin/payment/payment_account', data)
export const accountList   = (params) => request.get('/api/plugin/payment/payment_account/list', { params })
export const channelPayment = (params) => request.get('/api/plugin/payment/payment_account/channel/payment', { params })
export const channelConfig  = (params) => request.get('/api/plugin/payment/payment_account/channel/config', { params })
export const testQrcode    = (params) => request.get('/api/plugin/payment/payment_account/test/qrcode', { params })

// trade
export const tradeDetail  = (params) => request.get('/api/plugin/payment/trade', { params })
export const tradeList    = (params) => request.get('/api/plugin/payment/trade/search', { params })
export const tradeRefund  = (data)   => request.post('/api/plugin/payment/refund', data)
export const refundDetail = (params) => request.get('/api/plugin/payment/refund', { params })
export const refundList   = (params) => request.get('/api/plugin/payment/refund/search', { params })
export const refundAmount = (params) => request.get('/api/plugin/payment/refund/amount', { params })

// dashboard
export const dashboardTotalRequest     = (params) => request.get('/api/plugin/payment/dashboard/total_request', { params })
export const dashboardIndex            = (params) => request.get('/api/plugin/payment/dashboard/index', { params })
export const dashboardMchIndex         = (params) => request.get('/api/plugin/payment/dashboard/mch-index', { params })
export const dashboardMchAppIndex      = (params) => request.get('/api/plugin/payment/dashboard/mch-app-index', { params })
export const dashboardMchAppAccountAll = (params) => request.get('/api/plugin/payment/dashboard/mch-app-account-all-index', { params })
export const dashboardSearch           = (params) => request.get('/api/plugin/payment/dashboard/search', { params })
