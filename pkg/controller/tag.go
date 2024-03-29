// copyright @ 2021 ops inc.
//
// author: jinlong yang
//

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"nautilus/pkg/config"
	"nautilus/pkg/service/publish"
)

func BuildTag(c *gin.Context) {
	type params struct {
		ID      int64  `form:"pipeline_id" binding:"required"`
		Service string `form:"service" binding:"required"`
	}

	var data params
	if err := c.ShouldBind(&data); err != nil {
		ResponseFailed(c, err.Error())
		return
	}

	var (
		pid         = data.ID
		serviceName = data.Service
	)

	if err := publish.NewBuildTag(pid, serviceName); err != nil {
		log.Errorf("build tag failed: %+v", err)
		ResponseFailed(c, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

func ReceiveTag(c *gin.Context) {
	type params struct {
		ID     int64  `form:"taskid" binding:"required"`
		Module string `form:"module" binding:"required"`
		Tag    string `form:"tag" binding:"required"`
	}

	var data params
	if err := c.ShouldBind(&data); err != nil {
		ResponseFailed(c, err.Error())
		c.String(http.StatusOK, err.Error())
		return
	}

	var (
		pid    = data.ID
		module = data.Module
		tag    = data.Tag
	)

	if err := publish.NewReceiveTag(pid, module, tag); err != nil {
		log.Errorf("receive tag failed: %+v", err)
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, config.OK)
}

func ReceivePkg(c *gin.Context) {
	type params struct {
		ID     int64  `form:"taskid" binding:"required"`
		Module string `form:"module" binding:"required"`
		Pkg    string `form:"pkg" binding:"required"`
	}

	var data params
	if err := c.ShouldBind(&data); err != nil {
		ResponseFailed(c, err.Error())
		c.String(http.StatusOK, err.Error())
		return
	}

	var (
		pid    = data.ID
		module = data.Module
		pkg    = data.Pkg
	)

	if err := publish.NewReceivePkg(pid, module, pkg); err != nil {
		log.Errorf("receive compile package failed: %+v", err)
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, config.OK)
}
