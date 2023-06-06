// copyright @ 2021 ops inc.
//
// author: jinlong yang
//

package publish

import (
	"fmt"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"nautilus/pkg/config"
	"nautilus/pkg/model"
	"nautilus/pkg/util/cm"
)

func NewBuildImage() *BuildImage {
	return &BuildImage{}
}

type BuildImage struct{}

func (bi *BuildImage) Handle(pid int64, service string) error {
	pipeline, err := model.GetPipeline(pid)
	if err != nil {
		return fmt.Errorf(config.IMG_QUERY_PIPELINE_ERROR, err)
	}

	if err := bi.checkStatus(pipeline.Status); err != nil {
		return err
	}

	if err := model.CreatePhase(pid, model.PHASE_DEPLOY, model.PHASE_IMAGE, model.PHProcess); err != nil {
		log.Errorf("create pipeline: %d image phase error: %s", pid, err)
		return err
	}

	updateList, err := model.FindUpdateInfo(pid)
	if err != nil {
		return fmt.Errorf(config.IMG_QUERY_UPDATE_ERROR, err)
	}

	_, curPath, _, _ := runtime.Caller(1)
	var (
		mainPath   = filepath.Dir(filepath.Dir(filepath.Dir(curPath)))
		scriptPath = filepath.Join(mainPath, "script")
		changes    []string
		retains    []string
	)

	for _, item := range updateList {
		if err := model.CreateImage(pid, item.CodeModule); err != nil {
			return fmt.Errorf(config.IMG_CREATE_IMAGE_INFO_ERROR, err)
		}

		codeModule, err := model.GetCodeModuleInfo(item.CodeModule)
		if err != nil {
			return fmt.Errorf(config.TAG_QUERY_UPDATE_ERROR, err)
		}
		lang := codeModule.Language
		repo := codeModule.RepoAddr

		param := fmt.Sprintf("%s/makeimg -s %s -m %s -l %s -a %s -t %s -i %d",
			scriptPath, service, item.CodeModule, lang, repo, item.CodeTag, pid)
		log.Infof("makeimg command: %s", param)
		if !CallRealtimeOut(param, nil) {
			return fmt.Errorf(config.IMG_BUILD_FAILED)
		}
		changes = append(changes, item.CodeModule)
	}

	// 获取未变更的模块(服务所有模块-当前变更的模块)
	totals, err := model.FindServiceCodeModules(service)
	if err != nil {
		return fmt.Errorf(config.DB_QUERY_MODULE_BINDING_ERROR, err)
	}

	for _, item := range totals {
		codeModule := item.CodeModule.Name
		if cm.In(codeModule, changes) {
			continue
		}
		retains = append(retains, codeModule)
	}
	log.Infof("build image pipeline: %d fetch unchange code modules: %v", pid, retains)

	for _, codeModule := range retains {
		image, err := model.QueryLatestSuccessModuleImage(service, codeModule)
		if err != nil {
			log.Errorf(config.DB_IMAGE_CREATE_OR_UPDATE_ERROR, err)
			return err
		}
		imageURL := image.ImageURL
		imageTag := image.ImageTag

		if err := model.CreateOrUpdatePipelineImage(pid, service, codeModule, imageURL, imageTag); err != nil {
			return err
		}
		log.Infof("build image pipeline: %d record latest module: %s image: %s:%s success", pid, codeModule, imageURL, imageTag)
	}

	if err := model.UpdatePhase(pid, model.PHASE_DEPLOY, model.PHASE_IMAGE, model.PHSuccess); err != nil {
		log.Errorf("update pipeline: %d image phase error: %s", pid, err)
	}
	return nil
}

func (bi *BuildImage) checkStatus(status int) error {
	statusList := []int{
		model.PLSuccess,
		model.PLFailed,
		model.PLRollbackSuccess,
		model.PLRollbackFailed,
		model.PLTerminate,
	}
	if cm.Ini(status, statusList) {
		return fmt.Errorf(config.IMG_BUILD_FINISHED)
	}
	return nil
}

func UpdateImageInfo(pid int64, module, imageURL, imageTag string) error {
	return model.UpdateImage(pid, module, imageURL, imageTag)
}
