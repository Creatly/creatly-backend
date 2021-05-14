package main

import "github.com/zhashkevych/courses-backend/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
