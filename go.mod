module github.com/metalblueberry/halite-bot

go 1.12

require (
	github.com/gorilla/websocket v1.4.0
	github.com/metalblueberry/Halite-debug v0.1.1
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/sys v0.0.0-20190712062909-fae7ac547cb7 // indirect
)

replace github.com/metalblueberry/Halite-debug => ./../Halite-debug
