package controller

import (
	"encoding/json"
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
func (b *Controller) GetAll(c *gin.Context) {
	logger.Debug(c.ClientIP())
	if datas, err := b.md.GetAll(); err != nil {
		logger.Debug("Can't get all datas")
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		logger.Debug("Can get all datas")
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": datas,
		})
	}
}

func (b *Controller) GetMore(c *gin.Context) {
	logger.Debug(c.ClientIP())
	// Get the URL path from the request
	path := c.Request.URL.Path
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]

	if moreDatas, err := b.md.GetMore(lastPart); err != nil {
		logger.Debug("Can't get more datas")
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		logger.Debug("Can get more datas")
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": moreDatas,
		})
	}
}

func (b *Controller) GetBlockWithHeight(c *gin.Context) {
	logger.Debug(c.ClientIP())
	height := c.Param("height")
	if len(height) <= 0 {
		logger.Debug("fail, Not Found Param")
		b.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if block, err := b.md.GetOneBlcok(height); err != nil {
		logger.Debug("Can't get one block")
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
	} else {
		logger.Debug("Can get one block")
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
		logger.Debug("fail, Not Found Param")
		b.RespError(c, nil, 400, "fail, Not Found Param", nil)
		c.Abort()
		return
	}

	if block, err := b.md.GetOneTransaction(hash); err != nil {
		logger.Debug("Can't get one transaction")
		c.JSON(http.StatusOK, gin.H{
			"res":  "fail",
			"body": err.Error(),
		})
		c.Abort()
	} else {
		logger.Debug("Can get one transaction")
		c.JSON(http.StatusOK, gin.H{
			"res":  "ok",
			"body": block,
		})
		c.Next()
	}
}
