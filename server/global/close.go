package global

func Close() {
  if Redis != nil {
    Redis.Close()
  }
  if KafkaProduct != nil {
    KafkaProduct.Close()
  }

}
