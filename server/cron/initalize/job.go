package initalize

import (
  "context"
  "github.com/lihongsheng/pay-gateway/cron"
  log2 "github.com/lihongsheng/pay-gateway/log"
  job1 "github.com/lihongsheng/pay-gateway/plugin/payment/cron"
  "time"
)

func GetCronJobs() []cron.Job {
  return []cron.Job{
    {
      Spec:        "0 30 1 * * *", // 每天凌晨 1 点执行
      JobName:     "删除时间记录",
      Job:         job1.Job.DeleteEventRecord,
      MaxExecTime: 2 * time.Minute,
    },
    {
      Spec:        "0 10 2 * * 1", // 每周一凌晨2点执行
      JobName:     "删除请求日志",
      Job:         job1.Job.DeletePaymentRequestLog,
      MaxExecTime: 2 * time.Minute,
    },
    {
      Spec:        "0 15 2 * * *", // 每天凌晨2点执行
      JobName:     "删除通知日志",
      Job:         job1.Job.DeleteNotifyRecord,
      MaxExecTime: 2 * time.Minute,
    },
    // 每五秒执行一次
    {
      Spec:        "@every 5s",
      JobName:     "通知三方支付订单",
      Job:         job1.Job.CronPaymentPushNotifyHandler,
      MaxExecTime: time.Hour,
    },
    // 每五秒执行一次
    {
      Spec:        "@every 5s",
      JobName:     "通知三方退款订单",
      Job:         job1.Job.CronRefundPushNotifyHandler,
      MaxExecTime: time.Hour,
    },
    {
      Spec:        "0 0 1 * * *", // 每天凌晨1点执行
      JobName:     "删除通知记录",
      Job:         job1.Job.DeleteNotifyRecord,
      MaxExecTime: 2 * time.Minute,
    },
    {
      Spec:        "@every 30m", // 每12个小时执行一次
      JobName:     "删除订单过期记录",
      Job:         job1.Job.DeleteExpireRecord,
      MaxExecTime: 2 * time.Minute,
    },
    {
      Spec:        "@every 10m", // 每10分钟执行一次
      JobName:     "支付订单过期处理",
      Job:         job1.Job.PaymentExpireRecordJob,
      MaxExecTime: time.Hour,
    },
    {
      Spec:        "0 30 2 1 * *", // 每月凌晨2点执行
      JobName:     "支付订单过期处理",
      Job:         job1.Job.DeleteExpireRouterRecord,
      MaxExecTime: time.Hour,
    },
  }
}

func TestJob(ctx context.Context) error {
  log2.Info("TestJob", "time", time.Now())
  return nil
}
