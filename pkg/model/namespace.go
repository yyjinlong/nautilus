// Copyright @ 2021 OPS Inc.
//
// Author: Jinlong Yang
//

package model

import (
	"time"
)

type Cluster struct {
	ID       int64
	Name     string    `xorm:"varchar(50) notnull"`
	Creator  string    `xorm:"varchar(50) notnull"`
	CreateAt time.Time `xorm:"timestamp notnull created"`
}

type Namespace struct {
	ID       int64
	Name     string    `xorm:"varchar(32) notnull unique"`
	Cluster  string    `xorm:"varchar(50) notnull"`
	Creator  string    `xorm:"varchar(50) notnull"`
	CreateAt time.Time `xorm:"timestamp notnull created"`
}

func GetNamespaceByID(namespaceID int64) (*Namespace, error) {
	ns := new(Namespace)
	if has, err := SEngine.ID(namespaceID).Get(ns); err != nil {
		return nil, err
	} else if !has {
		return nil, NotFound
	}
	return ns, nil
}

func GetNamespaceByName(name string) (*Namespace, error) {
	ns := new(Namespace)
	if has, err := SEngine.Where("name=?", name).Get(ns); err != nil {
		return nil, err
	} else if !has {
		return nil, NotFound
	}
	return ns, nil
}

func GetClusterByNamespace(namespace string) (string, error) {
	ns, err := GetNamespaceByName(namespace)
	if err != nil {
		return "", err
	}
	return ns.Cluster, nil
}
