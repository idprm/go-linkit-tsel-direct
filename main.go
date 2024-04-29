package main

import (
	"time"

	"github.com/idprm/go-linkit-tsel/cmd"
	"github.com/idprm/go-linkit-tsel/internal/utils"
)

var (
	APP_TZ string = utils.GetEnv("APP_TZ")
)

func main() {

	loc, _ := time.LoadLocation(APP_TZ)
	time.Local = loc

	cmd.Execute()
}
