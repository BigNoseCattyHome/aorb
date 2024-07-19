package main

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
)

// 运行整个项目

func main() {
	fmt.Println(config.Conf.Consul.Addr)
}
