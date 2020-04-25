package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	app.Post("/send", func(ctx iris.Context){
		queue_name := ctx.PostValue("name")
		ctx.WriteString(queue_name)
	})

	app.Run(iris.Addr("192.168.66.71:8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
