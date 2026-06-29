### 项目介绍
   本系统支持多商户多应用的聚合支付系统。在支付渠道支持微信、支付宝、富友支付等。支付渠道动态配置，新增支付渠道只需按照规范设置对应配置即可。
### 项目特点
1. 聚合支付系统，基于多商户多应用的聚合支付系统。
2. 适配多支付渠道，新增支付渠道按照 sdk 规范适配即可。
3. 支持单应用配置多微信账户支付。
4. 支付失败，收银台自动轮询重试支付。
5. 在应用支持多微信账户应用,开启路由模式，自动屏蔽微信限制支付，限制获取用户openid的微信账户

### 业务流程

#### 聚合码支付
##### 非支持多微信账户应用
```mermaid
sequenceDiagram
actor User
participant PayPage as 支付H5页面（二维码）
participant PaySys as 支付系统
participant WxPay as 微信/支付宝..支付渠道

    %% 1. 生成聚合二维码
    Note over PaySys: 1. 生成聚合二维码
    PaySys->>PayPage: 生成聚合二维码

    %% 2. 用户扫码
    User->>PayPage: 2. 微信扫码打开页面
    PayPage->>PayPage: 加载页面基础信息


    %% 3. 用户输入金额
    User->>PayPage: 3. 输入金额，点击支付（携带UA等信息）

    %% 4. 支付系统处理
    PayPage->>PaySys: 4. 提交支付请求
    Note over PaySys: 5. 调用微信等支付接口
    PaySys->>WxPay: 5. 调用统一下单
    WxPay->>PaySys: 返回prepay_id等参数
    Note over PaySys: 6. 记录支付单&日志
    PaySys->>PayPage: 7. 返回微信支付参数

    %% 8. 唤起支付
    PayPage->>WxPay: 8. 唤起微信支付
    User->>WxPay: 9. 用户输入密码/确认支付

    %% 10. 支付成功回调
    WxPay->>PaySys: 10. 异步回调支付成功
    Note over PaySys: 11. 验签、更新订单状态

    %% 11. 页面跳转成功
    PayPage->>PaySys: 监听支付结果
    PayPage->>User: 12. 跳转到支付成功页
```
##### 支持多微信账户应用
```mermaid
sequenceDiagram
    actor User
    participant PayPage as 支付H5页面（二维码）
    participant PaySys as 支付系统
    participant WxPay as 微信/支付宝..支付渠道

    %% 1. 生成聚合二维码
    Note over PaySys: 1. 生成聚合二维码
    PaySys->>PayPage: 生成聚合二维码

    %% 2. 用户扫码
    User->>PayPage: 2. 微信扫码打开页面
    PayPage->>PayPage: 加载页面基础信息


    %% 3. 用户输入金额
    User->>PayPage: 3. 输入金额，点击支付（携带UA等信息）

    %% 4. 支付系统处理
    PayPage->>PaySys: 4. 提交支付请求
    Note over PaySys: 5. 调用微信等支付接口
    PaySys->>WxPay: 5. 调用统一下单
    WxPay->>PaySys: 返回支付失败
    Note over PaySys: 6. 记录支付单&日志
    PaySys->>PaySys: 7. 判断是否可以重试
    PaySys->>PayPage: 9. 返回支付结果

    %% 8. 支付重试
    PayPage->>PayPage: 9.1 判断是否是可以重试错误，是刷新页面获取另一个微信账户
    PayPage->>PayPage: 9.2 判断是否刷新后的重试，自动唤起支付。

    Note over PaySys: 10. 调用微信等支付接口
    PaySys->>WxPay: 11. 调用统一下单
    WxPay->>PaySys: 返回prepay_id等参数
    Note over PaySys: 12. 记录支付单&日志
    PaySys->>PayPage: 13. 返回微信支付参数

    %% 11. 页面跳转成功
    PayPage->>PaySys: 监听支付结果
    PayPage->>User: 12. 跳转到支付成功页
```
#### 收银台支付
##### 非多微信账户应用
```mermaid
sequenceDiagram
    actor User
    participant mch as     商户系统
    participant PayPage as 支付收银台页面
    participant PaySys as  支付系统
    participant WxPay as   微信/支付宝..支付渠道

    %% 1. 生成支付单
    Note over mch: 1. 生成订单
    mch->>PaySys:  2. 创建支付单
    Note over PaySys: 记录支付单信息
    PaySys->>mch:  返回支付单信息以及收银台页面


    %% 2. 引导用户支付
    mch->>User: 3. 引导用户打卡收银台页面
    User->>PayPage: 加载收银台页面


    %% 3. 加载订单信息
    PayPage->>PaySys: 4. 加载订单标题及金额等信息
    PayPage->>PayPage: 渲染页面
    %% 4. 支付
    User->>PayPage: 5. 提交支付请求
    PayPage->>PaySys: 6. 接受请求
    Note over PaySys: 10. 调用微信等支付接口
    PaySys->>WxPay: 11. 调用统一下单
    WxPay->>PaySys: 返回prepay_id等参数
    Note over PaySys: 12. 更新支付单&日志
    PaySys->>PayPage: 13. 返回微信支付参数

    %% 11. 页面跳转成功
    PayPage->>PaySys: 监听支付结果
    PayPage->>User: 12. 跳转到支付成功页
```
##### 多微信账户应用
重试流程和聚合码支付一致
#### 接口支付
接口支付对于多微信账户应用，需要先获取微信账户信息，再进行支付。因此商户系统在调用支付接口的时候需要传递在支付系统中支付渠道的编码以及商户系统获取到的openid。如果不是多微信账户应用，则不需要传递支付渠道的编码。
```mermaid
sequenceDiagram
    actor User
    participant mch as     商户系统
    participant PayPage as 支付收银台页面
    participant PaySys as  支付系统
    participant WxPay as   微信/支付宝..支付渠道

    %% 1. 生成支付单
    User->>mch: 1. 生成订单
     User->>mch: 2. 发起支付
    mch->>PaySys:  2. 调用支付
    PaySys->>PaySys: 接受请求，记录支付单&日志
    PaySys->>WxPay: 3. 调用统一下单
    WxPay->>PaySys: 返回prepay_id等参数
    Note over PaySys: 6. 更新支付单&日志
    PaySys->>mch: 7. 返回微信支付参数
    mch->>mch: 监听支付结果
    mch->>User: 12. 跳转到支付成功页
```
### 功能模块
1. 商户管理
   1. 商户信息管理
   2. 商户状态管理
   3. 商户用户管理
2. 应用管理
   1. 应用信息管理
   2. 应用聚合二维码生成
   3. 应用下支付配置
3. 支付
   1. 支付列表
   2. 支付详情
   3. 支付退款
   4. 支付超时关闭
   5. 支付退款发起
4. 支付路由
   1. 路由条件 监听到获取用户openid受到微信限制
   2. 路由条件 监听到支付受到微信限制（如达到最大收款次数）
   3. 路由条件 应用支付账户配置的日最大收款次数
5. 退款
   1. 退款列表
   2. 退款详情
6. 数据看板
7. 菜单管理等

### 模块截图
#### 商户管理
![商户管理](/docs/img/mch.png)
#### 应用管理
![应用列表](/docs/img/application_list.png)
![应用编辑](/docs/img/application_edit.png)
![支付通道列表](/docs/img/channel_list.png)
#### 数据看板
![dashboard](/docs/img/dashboard.png)
