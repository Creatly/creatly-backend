package main

import "github.com/zhashkevych/creatly-backend/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
