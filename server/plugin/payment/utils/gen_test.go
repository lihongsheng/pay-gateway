package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateID(t *testing.T) {

	fmt.Println(fmt.Sprintf("ID:%d", GenDeviceID.Generate()))
	// 1265610445407232
	m := ParseID(GenDeviceID.Generate())
	fmt.Println(time.UnixMilli(1758357486530 + m["timestamp_diff"]))
	fmt.Println(fmt.Sprintf("ID:%d", HashStringTo32("1265610445407232")))
}

func TestHashStringTo32(t *testing.T) {
	fmt.Println(HashCrc("1265610445407232", 4))
}

func TestGenAesHexStr(t *testing.T) {
	fmt.Println(AesGCMEncryptHex([]byte(fmt.Sprintf(`{"a":"%s"}`, "A4966c7e946c00")), []byte("17K9mP2n08rT4vY9")))
	fmt.Println(AesGCMDecryptHex("HN2DTUPZ6B7N3LEXPXUMJD5BE6TES7FGA26TQWHPMOUW37WPBV7XU45ZEJZICF2G6U7PH2F3K3OKUD4G", []byte("17K9mP2n08rT4vY9")))
}

func TestGenMd5Meta(t *testing.T) {
	fmt.Println(GenMd5Meta("yZuUPBYY0SPlOhY0c28luzyTg2qHVcnA", "A4966c7e946c00"))
}
