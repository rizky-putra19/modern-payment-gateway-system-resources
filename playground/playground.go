package main

import (
	"fmt"

	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
)

func main() {
	timeNow := helper.GenerateTime(0)
	timeNext := helper.GenerateTime(24)

	fmt.Println(timeNow)
	fmt.Println(timeNext)
}
