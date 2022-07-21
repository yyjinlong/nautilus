// copyright @ 2021 ops inc.
//
// author: jinlong yang
//

package publish

import (
	"fmt"

	"nautilus/golib/log"
	"nautilus/pkg/config"
	"nautilus/pkg/model"
	"nautilus/pkg/util"
)

func NewFinish() *Finish {
	return &Finish{Logid: util.UniqueID()}
}

type Finish struct {
	Logid string
}

func (f *Finish) Handle(pid int64) error {
	pipeline, err := model.GetPipeline(pid)
	if err != nil {
		return fmt.Errorf(config.DB_PIPELINE_QUERY_ERROR, pid, err)
	}

	service, err := model.GetServiceByID(pipeline.ServiceID)
	if err != nil {
		return fmt.Errorf(config.DB_SERVICE_QUERY_ERROR, err)
	}

	onlineGroup := service.DeployGroup
	deployGroup := util.GetDeployGroup(onlineGroup)
	log.ID(f.Logid).Infof("get current online_group: %s deploy_group: %s", onlineGroup, deployGroup)

	if err := model.UpdateGroup(pid, service.ID, onlineGroup, deployGroup, model.PLSuccess); err != nil {
		return fmt.Errorf(config.FSH_UPDATE_ONLINE_GROUP_ERROR, err)
	}
	log.ID(f.Logid).Infof("set current online group: %s deploy group: %s success.", onlineGroup, deployGroup)
	return nil
}
