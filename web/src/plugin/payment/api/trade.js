
import service from '@/utils/request'

// /private/v1/trade
export const getTrade = (params) => {
    return service({
        url: '/private/v1/trade',
        method: 'get',
        params: params
    })
}

// /private/v1/trade/search
export const getTradeList = (params) => {
    return service({
        url: '/private/v1/trade/search',
        method: 'get',
        params: params
    })
}

// post /private/v1/refund
export const refund = (params) => {
    return service({
        url: '/private/v1/refund',
        method: 'post',
        data: params
    })
}
// post /private/v1/refund/amount
export const refundAmount = (params) => {
    return service({
        url: '/private/v1/refund/amount',
        method: 'get',
        data: params
    })
}
// /private/v1/refund
export const getRefund = (params) => {
    return service({
        url: '/private/v1/refund',
        method: 'get',
        params: params
    })
}

// /private/v1/refund/search
export const getRefundList = (params) => {
    return service({
        url: '/private/v1/refund/search',
        method: 'get',
        params: params
    })
}

// /private/v1/refund/amount 可退款金额
export const getAvailableRefundAmount = (params) => {
  return service({
    url: '/private/v1/refund/amount',
    method: 'get',
    params: params
  })
}
