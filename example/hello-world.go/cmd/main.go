package main

import (
	hello "github.com/Seb-C/watch-dock/example/hello-world.go/pkg"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	for {
		logrus.Infoln(hello.GetHello())
		time.Sleep(time.Second)
	}
}
