package test

import (
	"go-spi/db"
	"testing"
)

func TestRedis(t *testing.T) {
	rd := db.GetRedisDbInstance()
	r, err := rd.Do("HEXISTS", "19kefei", "weijian")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("结果: ", r)
}
