/**
 * 收银台 API 接口
 * @module api/payment/application
 */

import service from '@/utils/request'

export const getApplication = (params) => {
    return service({
        url: '/public/v1/checkout/application',
        method: 'get',
        params: params
    })
}

export const routerAccount = (params) => {
    return service({
        url: '/public/v1/checkout/router',
        method: 'post',
        data:  params
    })
}
// .Payment_Wxpay
// .Channel_Wxpay
export const getPaymentOrder = (params) => {
    return service({
        url: '/public/v1/checkout/payment',
        method: 'get',
        params: params
    })
}

export const payment = (params) => {
    return service({
        url: '/public/v1/checkout/payment',
        method: 'post',
        data:  params
    })
}