// Copyright @ 2021 OPS Inc.
//
// Author: Jinlong Yang
//

package model

import (
	"time"

	"xorm.io/xorm"
)

type Service struct {
	ID            int64
	Name          string    `xorm:"varchar(32) notnull unique"`
	Namespace     string    `xorm:"varchar(32) notnull"`
	ImageAddr     string    `xorm:"varchar(500) notnull"`
	QuotaCPU      string    `xorm:"varchar(20)"`
	QuotaMaxCPU   string    `xorm:"varchar(20)"`
	QuotaMem      string    `xorm:"varchar(20)"`
	QuotaMaxMem   string    `xorm:"varchar(20)"`
	Replicas      int32     `xorm:"int"`
	Configmap     string    `xorm:"text"`
	ReserveTime   int       `xorm:"int"`
	Port          int       `xorm:"int"`
	ContainerPort int       `xorm:"int"`
	OnlineGroup   string    `xorm:"varchar(20) notnull"`
	DeployGroup   string    `xorm:"varchar(20) notnull"`
	MultiPhase    bool      `xorm:"bool"`
	Lock          string    `xorm:"varchar(100) notnull"`
	RD            string    `xorm:"varchar(50) notnull"`
	OP            string    `xorm:"varchar(50) notnull"`
	CreateAt      time.Time `xorm:"timestamp notnull created"`
	UpdateAt      time.Time `xorm:"timestamp notnull updated"`
}

type CodeModule struct {
	ID       int64
	Name     string    `xorm:"varchar(50) notnull"`
	Language string    `xorm:"varchar(20) notnull"`
	RepoName string    `xorm:"varchar(10) notnull"`
	RepoAddr string    `xorm:"varchar(200)"`
	CreateAt time.Time `xorm:"timestamp notnull created"`
	UpdateAt time.Time `xorm:"timestamp notnull updated"`
}

type ModuleBinding struct {
	ID           int64
	ServiceID    int64     `xorm:"int notnull"`
	CodeModuleID int64     `xorm:"int notnull"`
	CreateAt     time.Time `xorm:"timestamp notnull created"`
	UpdateAt     time.Time `xorm:"timestamp notnull updated"`
}

type BindingUnionQuery struct {
	Service       `xorm:"extends"`
	CodeModule    `xorm:"extends"`
	ModuleBinding `xorm:"extends"`
}

func BindingSession() *xorm.Session {
	return SEngine.Table("service").Alias("s").
		Join("INNER", []string{"module_binding", "b"}, "s.id = b.service_id").
		Join("INNER", []string{"code_module", "m"}, "m.id = b.code_module_id")
}

func GetServiceInfo(name string) (*Service, error) {
	service := new(Service)
	if has, err := SEngine.Where("name = ?", name).Get(service); err != nil {
		return nil, err
	} else if !has {
		return nil, NotFound
	}
	return service, nil
}

func GetServiceByID(serviceID int64) (*Service, error) {
	service := new(Service)
	if has, err := SEngine.ID(serviceID).Get(service); err != nil {
		return nil, err
	} else if !has {
		return nil, NotFound
	}
	return service, nil
}

func GetCodeModuleInfo(module string) (*CodeModule, error) {
	codeModule := new(CodeModule)
	if has, err := SEngine.Where("name=?", module).Get(codeModule); err != nil {
		return nil, err
	} else if !has {
		return nil, NotFound
	}
	return codeModule, nil
}

func GetCodeModuleInfoByID(moduleID int64) (*CodeModule, error) {
	codeModule := new(CodeModule)
	if has, err := SEngine.ID(moduleID).Get(codeModule); err != nil {
		return nil, err
	} else if !has {
		return nil, NotFound
	}
	return codeModule, nil
}

func FindServiceCodeModules(service string) ([]BindingUnionQuery, error) {
	bindings := make([]BindingUnionQuery, 0)
	if err := BindingSession().Where("s.name = ?", service).Find(&bindings); err != nil {
		return nil, err
	}
	return bindings, nil
}

func UpdateConfigMap(name string, pair string) error {
	service := new(Service)
	service.Configmap = pair
	if affected, err := MEngine.Where("name = ?", name).Cols("configmap").Update(service); err != nil {
		return err
	} else if affected == 0 {
		return NotFound
	}
	return nil
}
