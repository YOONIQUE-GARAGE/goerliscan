package router

import (
	"github.com/gin-gonic/gin"

	ctl "rnd/goerliscan/explore/controller"
)

type Router struct {
	ct *ctl.Controller
}

func NewRouter(ct *ctl.Controller) (*Router, error) {
	r := &Router{
		ct: ct,
	}
	return r, nil
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (b *Router) Idx() *gin.Engine {
	e := gin.Default()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(CORS())

	e.GET("/", b.ct.GetAll)
	e.GET("/blocks", b.ct.GetMore)
	e.GET("/txs", b.ct.GetMore) 
	e.GET("/block/:height", b.ct.GetBlockWithHeight)
	e.GET("/tx/:hash", b.ct.GetBlockWithHash)
	
	return e
}