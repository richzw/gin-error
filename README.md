
Gin Handle Error Middleware
==============

Gin Handle Error Middleware is one middleware for [Gin](https://github.com/gin-gonic/gin) framework that you could handle errors in middleware, so you could NOT do error handling within each handler.

Usage
-----

Download and install without using go module:

```shell
go get github.com/richzw/gin-error
```

Import it in your code:

```shell
import "github.com/richzw/gin-error"
```

Example
-----

- Map error to one status code

```go
var BadRequestErr = fmt.Errorf("bad request error")

func main() {
    r := gin.Default()
    r.Use(err.Error(err.NewErrMap(BadRequestErr).StatusCode(http.StatusBadRequest)))

    r.GET("/test", func(c *gin.Context) {
        c.Error(BadRequestErr)
    })

    r.Run()
}
```

- Map error to the response

```go
var BadRequestErr = fmt.Errorf("bad request error")

func main() {
    r := gin.Default()
    r.Use(err.Error(
        err.NewErrMap(BadRequestErr).Response(func(c *gin.Context) {
            c.JSON(http.StatusBadRequest, gin.H{"error": BadRequestErr.Error()})
        })))

    r.GET("/test", func(c *gin.Context) {
        c.Error(NotFoundError)
    })

    r.Run()
}
```


