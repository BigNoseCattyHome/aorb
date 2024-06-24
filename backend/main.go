package main

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
)

// 运行整个项目

func main() {
	fmt.Println(config.Conf.Pod)
}
