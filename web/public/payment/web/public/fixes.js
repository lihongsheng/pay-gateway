let payEnv = {
  // 判断是否为微信环境（微信内H5/微信扫码）
  // 检测是否为微信浏览器
  isWechatBrowser() {
    const ua = navigator.userAgent.toLowerCase();
    return /micromessenger/i.test(ua)
  },
  // 检测是否为支付宝浏览器
  isAlipayBrowser() {
    const ua = navigator.userAgent.toLowerCase();
    return /alipayclient/i.test(ua)
  },
  /**
   * 判断是否为 iOS 设备（iPhone/iPad/iPod）
   * 适配所有 iOS 设备类型，包括 iPadOS（iPadOS 13+ UA 会隐藏 iPad 标识）
   */
  isIOS() {
    // 核心判断：匹配 iPhone/iPad/iPod，同时兼容 iPadOS 13+ 的特殊 UA
    const isIosDevice = /iphone|ipad|ipod/.test(this.ua);
    // iPadOS 13+ 会把自己识别为 Mac，需额外判断
    const isIpadOS = /mac os/.test(this.ua) && navigator.maxTouchPoints > 1 && !/iphone/.test(this.ua);

    return isIosDevice || isIpadOS;
  },
  /**
   * 判断是否为 Android 设备
   * 排除 Windows Phone 等伪 Android 标识，仅匹配纯 Android 设备
   */
  isAndroid() {
    // 匹配 Android 标识，同时排除 Windows Phone（部分 WP 设备 UA 含 Android 字段）
    return /android/.test(this.ua) && !/windows phone/.test(this.ua);
  },
  // 判断操作系统：ios/android/pc/other
  // 获取操作系统
  getOs() {
    if (this.isIOS()) return 'IOS'
    if (this.isAndroid()) return 'Android'
    if (/windows/i.test(this.ua)) return 'Windows'
    if (/mac os/i.test(this.ua)) return 'Mac'
    if (/linux/i.test(this.ua)) return 'Linux'
    return 'other'
  },
  // 获取浏览器类型
  getDevice() {
    const ua = navigator.userAgent.toLowerCase();
    if (this.isWechatBrowser()) return 'Wechat'
    if (this.isAlipayBrowser()) return 'Alipay'
    if (this.isMobile()) {
      if (/chrome/i.test(ua)) return 'H5'
      if (/safari/i.test(ua)) return 'H5'
      if (/firefox/i.test(ua)) return 'H5'
      return 'H5'
    } else {
      if (/chrome/i.test(ua)) return 'PC'
      if (/safari/i.test(ua)) return 'PC'
      if (/firefox/i.test(ua)) return 'PC'
      return 'PC'
    }
  },
  isMobile() {
    // 检查是否为移动设备
    const userAgent = navigator.userAgent;
    // 基本移动设备检测
    const mobileRegex = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i;
    // 检查屏幕宽度（移动端通常小于768px）
    const isSmallScreen = window.innerWidth <= 768;
    // 检查触摸支持
    const hasTouch = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
    // 综合判断：用户代理包含移动设备标识，或者屏幕较小且支持触摸
    return mobileRegex.test(userAgent) || (isSmallScreen && hasTouch);
  },
  getEnvInfo() {
    return {
      os: this.getOs(),
      ua: navigator.userAgent,
    };
  }
};
payEnv.getEnvInfo();


function getDefaultPayment() {
  if (payEnv.isWechatBrowser()) {
    return 2
  } else if (payEnv.isAlipayBrowser()) {
    return 1
  }
  return 2
}
function getDefaultPaymentProduct() {
  if (payEnv.isWechatBrowser()) {
    return 3
  } else if (payEnv.isAlipayBrowser()) {
    return 3
  }
  return 3
}

// 通用的请求函数
async function request(config) {
  const { url, method = 'GET', data, params = {}, headers = {} } = config;
  // 构建完整URL
  let fullUrl = API_BASE_URL + url;
  // 添加查询参数
  if (params && Object.keys(params).length > 0) {
    const queryString = new URLSearchParams(params).toString();
    fullUrl += '?' + queryString;
  }
  let ContentType = 'application/json'
  if (method.toUpperCase() === 'GET') {
    // GET 请求不需要添加请求体
    ContentType = ''
  }
  // 请求配置
  const requestConfig = {
    method: method.toUpperCase(),
    headers: {
      'Content-Type': ContentType,
      'Accept': 'application/json',
      ...headers
    }
  };
  // 添加请求体
  if (data && ['POST', 'PUT', 'PATCH'].includes(method.toUpperCase())) {
    requestConfig.body = JSON.stringify(data);
  }
  try {
    const response = await fetch(fullUrl, requestConfig);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('请求失败:', error);
    throw error;
  }
}
async function getApplicationConfig(params) {
  return await request({
    url: '/public/v1/aggregate/application',
    method: 'GET',
    params: params
  })
}
async function payment(params) {
  return await request({
    url: '/public/v1/aggregate/payment',
    method: 'POST',
    data: params
  })
}
async function paymentTest(params) {
  return await request({
    url: '/public/v1/aggregate/test/payment',
    method: 'POST',
    data: params
  })
}
async function getUserOpenId(params) {
  return await request({
    url: '/public/v1/aggregate/user/auth',
    method: 'GET',
    params: params
  })
}
async function query(params) {
  return await request({
    url: '/public/v1/aggregate/payment/query',
    method: 'GET',
    params: params
  })
}
async function getRedirectUrl(params) {
  return await request({
    url: '/public/v1/aggregate/redirect',
    method: 'GET',
    params: params
  })
}

async function getTestRedirectUrl(params) {
  return await request({
    url: '/public/v1/aggregate/test/redirect',
    method: 'GET',
    params: params
  })
}

async function getCashierOrder(params) {
  return await request({
    url: '/public/v1/cashier/query',
    method: 'GET',
    params: params
  })
}

async function cashierPayment(params) {
  return await request({
    url: '/public/v1/cashier/payment',
    method: 'POST',
    data: params
  })
}

async function getCashierRedirectUrl(params) {
  return await request({
    url: '/public/v1/cashier/redirect',
    method: 'GET',
    params: params
  })
}

function hash(str, len ) {
  let crc = 0xFFFFFFFF;
  const polynomial = 0xEDB88320;
  for (let i = 0; i < str.length; i++) {
    const byte = str.charCodeAt(i);
    crc ^= byte;
    for (let j = 0; j < 8; j++) {
      if (crc & 1) {
        crc = (crc >>> 1) ^ polynomial;
      } else {
        crc >>>= 1;
      }
    }
  }
  crc ^= 0xFFFFFFFF;
  return (crc >>> 0).toString(10).padStart(len, '0').substring(0, len);
}

// 解析URL
function getUrlParams() {
  // 从URL的hash部分获取action和token
  const hash = window.location.hash.substring(1); // 移除开头的 '#'
  // https://proxy.test.jianxindianzi.com/sslab/payment/web/public/payment.html?action=oauth&action_token=PDP5VEDL6QMD42LLUA4L7W22AJXPLWWTU576P7FHSRZFPW7FWBV2HD62H7UCQMKTG2BBOXF2G27DFW2F&auth_code=456911272ad542de84cefa053edcRX99&account_no=P4bc6f2201e800&filter_account_no=&app_id=2021005199606982&source=alipay_wallet&scope=auth_base
  console.log('window.location.search:', window.location.search);
  // 如果从hash中没有获取到，尝试从查询参数中获取
  const searchParams = new URLSearchParams(window.location.search);
  const queryIndex = hash.indexOf('?');
  if (queryIndex !== -1) {
    const hashParams = new URLSearchParams(hash.slice(queryIndex + 1));
    console.log('hashParams:', hash.slice(queryIndex + 1));
    // 合并参数
    for (const [key, val] of hashParams) {
      // 追加模式：保留多值参数
      searchParams.append(key, val);
    }
  }
  return searchParams ;
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

function getBaseApi(defaultPath ) {
  const path = window.location.pathname;
  const basePath = window.location.origin;
  if (path === defaultPath) {
    return basePath;
  }
  // 如果路径是 /prefix/payment/web/public/payment.html，则取 /prefix 部分作为 API 基础路径
  const pathSegments = path.split('/');
  const defaultPathSegments = defaultPath.split('/');
  console.log('API_BASE_URL:', pathSegments);
  console.log('API_BASE_URL:', defaultPathSegments);
  // 从后往前匹配路径段
  let i = pathSegments.length - defaultPathSegments.length;
  if (i <= 0) {
    return basePath;
  }
  let match = true;
  for (let j = 1; j < defaultPathSegments.length; j++) {
    if (pathSegments[i + j] !== defaultPathSegments[j]) {
      match = false;
      console.log('API_BASE_URL:', pathSegments[i + j] , defaultPathSegments[j] , i + j, j);
      break;
    }
  }
  if (match) {
    // 获取匹配路径前面的部分作为基础路径
    return basePath + pathSegments.slice(0, i+1).join('/');
  }
  // 默认情况下返回 origin
  return basePath;
}

function closePage() {
  if (payEnv.isWechatBrowser()) {
    closeWechatH5Page();
  } else if (payEnv.isAlipayBrowser()) {
    closeAlipayH5Page();
  } else {
    window.close();
  }
}

/**
 * 关闭微信H5页面
 * 兼容不同微信版本，确保API可用
 */
function closeWechatH5Page() {
  // 方式1：优先使用WeixinJSBridge（最稳定）
  if (typeof WeixinJSBridge === 'object' && typeof WeixinJSBridge.invoke === 'function') {
    try {
      WeixinJSBridge.invoke('closeWindow', {}, function(res) {
        // 回调函数（一般无需处理）
        console.log('关闭页面结果：', res);
      });
    } catch (e) {
      // 降级处理
      window.close();
    }
  } else {
    window.close();
  }
}

/**
 * 关闭支付宝H5页面（兼容支付宝容器和普通浏览器）
 */
function closeAlipayH5Page() {
  try {
    // 检测是否在支付宝客户端内
    if (window.my && typeof my.navigateBack === 'function') {
      // 1. 尝试返回上一页（小程序/生活号内的H5）
      my.navigateBack({
        delta: 1, // 返回1级页面
        fail: function() {
          // 2. 返回失败则退出当前小程序/H5容器
          my.exitMiniProgram();
        }
      });
    } else if (window.AlipayJSBridge) {
      // 兼容旧版支付宝JSBridge
      AlipayJSBridge.call('closeWebview');
    } else {
      window.close();
    }
  } catch (e) {
    console.error('关闭页面失败:', e);
    window.close();
  }
}

function orderIsComplete(status )  {
  return !(status === 1 || status === 2 || status === 8 || status === 9 || status == 11);
}
// 实现订单号生成：时间戳+4位随机数（不足4位前面补0）+ 32位哈希值（不足32位前面补0）
function genOrderNo(prefix = 'P', token = '') {
  const timestamp = Date.now().toString();
  const randomNum = Math.floor(Math.random() * 10000).toString().padStart(4, '0');
  const rawString = `${timestamp}${randomNum}`;
  // 返回最大四位纯数字
  const hashStr = hash(token, 4);
  // 增加业务前缀，格式：前缀_时间戳+随机数+哈希
  return `${prefix}${rawString}${hashStr}`;
}

// 判断支付环境（微信/支付宝）
function judgePaymentPayEnv() {
  urlQueryParams.paymentMethod = getDefaultPayment();
  urlQueryParams.paymentProduct = getDefaultPaymentProduct();
}

let urlQueryParams = {
  encryptToken: "",
  paymentMethod: 2,
  paymentProduct: 3,
  authCode: "",
  accountNO: "",
  orderNo: "",
  action: "",
  filterAccountNoStr: "",
  openId: ""
}
// 全局变量
let applicationConfig = {
  app_no: '',
  pay_icon: '',
  pay_title: '',
  multi_channel: false
};

// 加载商户配置（根据app_no请求后端）
async function loadPaymentApplicationConfig() {
  if (urlQueryParams.action === 'hub') {
    return;
  }
  try {
    const res = await getApplicationConfig({token: urlQueryParams.encryptToken });
    if (res.code === 0) {
      applicationConfig = res.data;
    } else {
      displayError('加载商户配置失败:' + res?.msg,true, true);
    }
  } catch (err) {
    displayError('加载商户配置异常',true, true);
  }
}

// 获取用户标识（微信OpenID/支付宝OpenID）
async function loadRetryUserIdentity(pageTag = 'cashier') {
  if (urlQueryParams.action === 'hub') {
    return;
  }
  if (!urlQueryParams.accountNO) {
    displayError('未找到支付账户', true, true);
    return;
  }
  try {
    const data = await getUserOpenId({auth_code: urlQueryParams.authCode,
      account_no: urlQueryParams.accountNO,
      token: urlQueryParams.encryptToken,
      payment_method: urlQueryParams.paymentMethod,
      payment_product: urlQueryParams.paymentProduct,
      device : payEnv.getDevice()});
    if (data.code === 0) {
      urlQueryParams.openId = data.data;
    } else if (data.code === 41101 && applicationConfig.multi_channel) {
      await loadRetryRedirectUrl(urlQueryParams.orderNo,  true, pageTag);
    } else {
      displayError('授权失败',true, true);
    }
  } catch (err) {
    displayError('授权失败',true, false);
  }
}

// 重定向获取openid
async function loadRetryRedirectUrl(orderNoParams = '', retry = false, pageTag='cashier') {
  if (urlQueryParams.action === 'oauth') {
    if (!retry) {
      return;
    }
  } else if (urlQueryParams.action !== 'hub') {
    return;
  }
  try {
    if (retry) {
      if (urlQueryParams.filterAccountNoStr) {
        urlQueryParams.filterAccountNoStr = urlQueryParams.filterAccountNoStr + ',' + urlQueryParams.accountNO;
      } else {
        urlQueryParams.filterAccountNoStr = urlQueryParams.accountNO;
      }
    }
    let res = {};
    const redirectParams = {
      token: urlQueryParams.encryptToken,
      payment_method: urlQueryParams.paymentMethod,
      payment_product: urlQueryParams.paymentProduct,
      device: payEnv.getDevice(),
      filter_account_no: urlQueryParams.filterAccountNoStr,
      order_no: orderNoParams,
      retry: retry
    }
    if (pageTag === "cashier") {
      res = await getCashierRedirectUrl(redirectParams);
    } else {
      res = await getRedirectUrl(redirectParams);
    }
    if (res.code !== 0) {
      displayError('授权失败:' + res?.msg,true, false);
      return;
    }
    window.location.href = res.data;
  } catch (err) {
    displayError('授权异常：',true, false);
  }
}
// 从URL的hash部分获取action和token
function parsePaymentUrlParams(isCheckoutOrder = false ) {
  const searchParams = getUrlParams();
  urlQueryParams.action = searchParams.get('action');
  urlQueryParams.encryptToken = searchParams.get('action_token');
  if (!urlQueryParams.encryptToken) {
    throw new Error('未获取到应用信息');
  }
  if (!urlQueryParams.action) {
    throw new Error('错误：不支持的支付环境');
  }
  urlQueryParams.orderNo = searchParams.get("order_no")
  if (urlQueryParams.orderNo == 'null') {
    urlQueryParams.orderNo = ''
  }
  if ((!urlQueryParams.orderNo) && isCheckoutOrder) {
    throw new Error('未获取到订单号');
  }
  urlQueryParams.accountNO = searchParams.get('account_no');
  if (searchParams.get('filter_account_no')) {
    urlQueryParams.filterAccountNoStr = searchParams.get('filter_account_no');
  }
  if (!urlQueryParams.encryptToken) {
    urlQueryParams.encryptToken = searchParams.get('encrypt');
  }
  if (!urlQueryParams.encryptToken) {
    throw new Error('未获取到应用信息');
  }
  // 如果是hub模式，则获取授权码
  if (urlQueryParams.action === 'oauth') {
    if (payEnv.isWechatBrowser()) {
      urlQueryParams.authCode = searchParams.get('code');
    } else if (payEnv.isAlipayBrowser()) {
      urlQueryParams.authCode = searchParams.get('auth_code');
    }
    if (!urlQueryParams.authCode) {
      throw new Error('未获取到授权码');
    }
  }
}

// 实现订单号生成：时间戳+4位随机数（不足4位前面补0）+ 32位哈希值（不足32位前面补0）
function generatePaymentNo(prefix = 'P') {
  const timestamp = Date.now().toString();
  const randomNum = Math.floor(Math.random() * 10000).toString().padStart(4, '0');
  const rawString = `${timestamp}${randomNum}`;
  // 返回最大四位纯数字
  const hashStr = hash(urlQueryParams.encryptToken, 4);
  // 增加业务前缀，格式：前缀_时间戳+随机数+哈希
  return `${prefix}${rawString}${hashStr}`;
}

let orderDetail = {
  order:{
    order_no: '',
    amount: {
      total: 0
    },
    subject: '',
  },
  trade_no: '',
  create: '',
  order_subject: '',
  status: '',
}


async function loadCashierOrder() {
  if (urlQueryParams.action === 'hub') {
    return;
  }
  if (!urlQueryParams.orderNo) {
    return;
  }
  try {
    const res = await getCashierOrder({token: urlQueryParams.encryptToken, order_no: urlQueryParams.orderNo });
    if (res.code === 0) {
      orderDetail = res.data;
      if (orderIsComplete(orderDetail.status)) {
        displayError('支付完结',true, true);
      }
    } else {
      displayError('获取订单失败：' + res?.msg,true, false);
    }
  } catch (err) {
    console.error('获取商户配置失败：', err);
    displayError('获取订单失败',true, false);
  }
}
