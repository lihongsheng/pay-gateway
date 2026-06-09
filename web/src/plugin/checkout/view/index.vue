<template>
  <!-- 在模板中添加支付宝 JS SDK 的 script 标签 -->
  <div>
    <script src="https://gw.alipayobjects.com/as/g/h5-lib/alipayjsapi/3.1.1/alipayjsapi.min.js"></script>
  </div>
</template>

<!--<script setup>-->
<!--/**-->
<!-- * 支付环境识别工具：区分微信内/微信扫码/支付宝/浏览器，识别操作系统-->
<!-- */-->

<!--export const payEnv = {-->
<!--  // 判断是否为微信环境（微信内H5/微信扫码）-->
<!--  // 检测是否为微信浏览器-->
<!--  isWechatBrowser() {-->
<!--    const ua = navigator.userAgent.toLowerCase();-->
<!--    return /micromessenger/i.test(ua)-->
<!--  },-->
<!--  // 检测是否为支付宝浏览器-->
<!--  isAlipayBrowser() {-->
<!--    const ua = navigator.userAgent.toLowerCase();-->
<!--    return /alipayclient/i.test(ua)-->
<!--  },-->

<!--  // 判断操作系统：ios/android/pc/other-->
<!--  // 获取操作系统-->
<!--  getOS() {-->
<!--    if (this.isIOS()) return 'IOS'-->
<!--    if (this.isAndroid()) return 'Android'-->
<!--    if (/windows/i.test(this.ua)) return 'Windows'-->
<!--    if (/mac os/i.test(this.ua)) return 'Mac'-->
<!--    if (/linux/i.test(this.ua)) return 'Linux'-->
<!--    return 'other'-->
<!--  },-->

<!--  // 获取浏览器类型-->
<!--  getDevice() {-->
<!--    const ua = navigator.userAgent.toLowerCase();-->
<!--    if (this.isWechatBrowser()) return 'Wechat'-->
<!--    if (this.isAlipayBrowser()) return 'Alipay'-->
<!--    if (this.isMobile()) {-->
<!--      if (/chrome/i.test(ua)) return 'H5'-->
<!--      if (/safari/i.test(ua)) return 'H5'-->
<!--      if (/firefox/i.test(ua)) return 'H5'-->
<!--      return 'H5'-->
<!--    } else {-->
<!--      if (/chrome/i.test(ua)) return 'PC'-->
<!--      if (/safari/i.test(ua)) return 'PC'-->
<!--      if (/firefox/i.test(ua)) return 'PC'-->
<!--      return 'PC'-->
<!--    }-->
<!--    return 'H5'-->
<!--  }-->
<!--};-->
<!--// 整合环境信息（供页面使用）-->
<!--function  getEnvInfo() {-->
<!--  return {-->
<!--    isWechat: this.isWechatBrowser(),-->
<!--    isAlipay: this.isAlipayBrowser(),-->
<!--    os: this.getOs(),-->
<!--    ua: navigator.userAgent,-->
<!--  };-->
<!--}-->
<!--// 单例模式-->
<!--export const uaDetector = getEnvInfo()-->

<!--function getDefaultPayment() {-->
<!--  if (uaDetector.isWechat) {-->
<!--    return "Wxpay"-->
<!--  } else if (uaDetector.isAlipay) {-->
<!--    return "Alipay"-->
<!--  }-->
<!--  return "Wxpay"-->
<!--}-->
<!--function getDefaultPaymentProduct() {-->
<!--  if (uaDetector.isWechat) {-->
<!--    return "JSAPI"-->
<!--  } else if (uaDetector.isAlipay) {-->
<!--    return "JSAPI"-->
<!--  }-->
<!--  return "Qrcode"-->
<!--}-->

<!--function getAuthCode() {-->
<!--  return "Bearer " + localStorage.getItem("token")-->
<!--}-->
<!--/**-->
<!-- * 获取OpenId-->
<!-- * @param {string} code 微信授权code-->
<!-- * @param wxAppId-->
<!-- * @returns {Promise<string>} OpenId-->
<!-- */-->
<!--async function getWechatOpenId(code, wxAppId) {-->
<!--  if (!code) {-->
<!--    // 微信内自动获取code-->
<!--    if (uaDetector.isWechatBrowser()) {-->
<!--      const currentUrl = encodeURIComponent(window.location.href)-->
<!--      // const redirectUrl = `${window.location.origin}/payment/auth/wechat?redirect=${currentUrl}`-->
<!--      window.location.href = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${wxAppId}&redirect_uri=${currentUrl}&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect`-->
<!--    }-->
<!--  }-->
<!--  // 调用后端API获取OpenId-->
<!--  try {-->
<!--    // 这里调用你的后端API-->
<!--    // const response = await axios.get(`/api/payment/wechat/openid?code=${code}`)-->
<!--    // return response.data.openid-->
<!--    return '模拟OpenId' // 模拟返回-->
<!--  } catch (error) {-->
<!--    console.error('获取OpenId失败:', error)-->
<!--    return null-->
<!--  }-->
<!--}-->
<!--/**-->
<!-- * 获取OpenId-->
<!-- * @param {string} code 微信授权code-->
<!-- * @param appId-->
<!-- * @returns {Promise<string>} OpenId-->
<!-- */-->
<!--async function getAlipayOpenId(code, appId) {-->
<!--    if (!code) {-->
<!--      // 获取code-->
<!--      if (uaDetector.isWechatBrowser()) {-->
<!--        const currentUrl = encodeURIComponent(window.location.href)-->
<!--        window.location.href = `https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?appid=${appId}&redirect_uri=${currentUrl}&scope=auth_base`-->
<!--      }-->
<!--    }-->
<!--    // auth_code-->
<!--  // 调用后端API获取OpenId-->
<!--  try {-->
<!--    // 这里调用你的后端API-->
<!--    // const response = await axios.get(`/api/payment/wechat/openid?code=${code}`)-->
<!--    // return response.data.openid-->
<!--    return '模拟OpenId' // 模拟返回-->
<!--  } catch (error) {-->
<!--    console.error('获取OpenId失败:', error)-->
<!--    return null-->
<!--  }-->
<!--}-->
<!--</script>-->