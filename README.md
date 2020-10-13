# Hyper

Hyper is an boilerplate for Golang Gin framework

## Install

```shell
go get -u github.com/cyub/hyper
```
## Usage

```go
package main

import (
    "net/http"
	"github.com/cyub/hyper"
	"github.com/cyub/hyper/mysql"
	"github.com/cyub/hyper/queue"
	"github.com/cyub/hyper/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	app := hyper.NewApp(
		hyper.WithName("hyper-demo"),
		hyper.WithAddr(":8000"),
		hyper.WithRunMode("dev"),
	)
	// bootstrap with some component as you need
	app.BootstrapWith(
		redis.Provider(),
		mysql.Provider(),
		queue.Provider(),
	)
	// register router
	app.RegisterRouter(func(r *gin.Engine) {
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello, world"
			})
		})
	})
	app.Run()
}
```

