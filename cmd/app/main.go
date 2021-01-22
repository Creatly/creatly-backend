package main

import "github.com/zhashkevych/courses-backend/internal/app"

const configPath = "configs/main"

func main() {
	app.Run(configPath)
}
