package xray

import (
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gin-gonic/gin"
	"github.com/haozzzzzzzz/go-lambda/resource"
)

func init() {
	resource.RegisterResource(resource.XRayResourceType)

	xray.Configure(xray.Config{LogLevel: "trace"})
}

func XRayGinMiddleware(strSegmentNamer string) func(*gin.Context) {
	return func(context *gin.Context) {
		xray.Handler(xray.NewFixedSegmentNamer(strSegmentNamer), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.Request = r
			context.Next()
		})).ServeHTTP(context.Writer, context.Request)
	}
}
