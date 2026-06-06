

import service from '@/utils/request'

/**
 * 数据统计
 * @returns {Promise} 商户列表数据
 */
export const getDashboardData = (params) => {
    return service({
        url: '/private/v1/dashboard/index',
        method: 'get',
        params: params
    })
}

/**
 * 具体某个商户的数据统计
 * @returns {Promise} 商户列表数据
 */
export const getDashboardMchData = (params) => {
    return service({
        url: '/private/v1/dashboard/mch-index',
        method: 'get',
        params: params
    })
}

/**
 * 具体某个商户下某个应用的数据
 * @returns {Promise} 商户列表数据
 */
export const getDashboardMchAppData = (params) => {
    return service({
        url: '/private/v1/dashboard/mch-app-index',
        method: 'get',
        params: params
    })
}


/**
 * 具体某个商户下某个应用的所有支付账户数据
 * @returns {Promise} 商户列表数据
 */
export const getDashboardMchAppAccountData = (params) => {
    return service({
        url: '/private/v1/dashboard/mch-app-account-all-index',
        method: 'get',
        params: params
    })
}
/**
 * 具体某个商户的数据统计
 * @returns {Promise} 商户列表数据
 */
export const getDashboardTotalRequest = (params) => {
    return service({
        url: '/private/v1/dashboard/total_request',
        method: 'get',
        params: params
    })
}
/**
 * 搜索商户数据
 * @returns {Promise} 商户列表数据
 */
export const searchDashboard = (params) => {
    return service({
        url: '/private/v1/dashboard/search',
        method: 'get',
        params: params
    })
}