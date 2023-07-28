package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rnd/goerliscan/explore/logger"
	"rnd/goerliscan/explore/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewCTL(rep *model.Model) (*Controller, error) {
	r := &Controller{
		md: rep,
	}

	return r, nil
}

func (b *Controller) Check(c *gin.Context) {
	b.RespOK(c, 0)
}

func (b *Controller) RespOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func (b *Controller) RespError(c *gin.Context, body interface{}, status int, err ...interface{}) {
	bytes, _ := json.Marshal(body)

	fmt.Println("Request error", "path", c.FullPath(), "body", bytes, "status", status, "error", err)

	c.JSON(status, gin.H{
		"Error":  "Request Error",
		"path":   c.FullPath(),
		"body":   bytes,
		"status": status,
		"error":  err,
	})
	c.Abort()
}

func (b *Controller) GetOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

// MainPage
func (b *Controller) GetAllInfo(c *gin.Context) {
	logger.Debug(c.ClientIP())
	if blocks, err := b.md.GetAll(); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": blocks,
		})
	}
}

func (b *Controller) GetMore(c *gin.Context) {
	logger.Debug(c.ClientIP())
	// Get the URL path from the request
	path := c.Request.URL.Path
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]

	if block, err := b.md.GetMore(lastPart); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": block,
		})
	}
}


func (b *Controller) GetBlockWithHeight(c *gin.Context) {
	logger.Debug(c.ClientIP())
	height := c.Param("height")
	if len(height) <= 0 {
		b.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if block, err := b.md.GetOneBlcok(height); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": block,
		})
	}
}

func (b *Controller) GetBlockWithHash(c *gin.Context) {
	logger.Debug(c.ClientIP())

	hash := c.Param("hash")
	if len(hash) <= 0 {
		b.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if block, err := b.md.GetOneTransaction(hash); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": block,
		})
		c.Next()
	}
}
