package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// ContainerRepo: the container repo, should support dockerv2, ociv1.
type ContainerRepo struct {
	Id          int    `orm:"column(id);auto"`
	Name        string `orm:"unique;column(name);size(255);null"`
	Star        int    `orm:"column(star);null"`
	DownloadNum int    `orm:"column(download_num);null"`
	Description string `orm:"column(description);null"`
}

type ContainerImage struct {
	Id     int    `orm:"column(id);auto"`
	Tag    string `orm:"column(tag);size(255);null"`
	Size   int64  `orm:"column(size);null"`
	RepoID int    `orm:"column(repo_id);null"`

	// docker, oci, rkt...
	Proto        string `orm:"column(proto);size(15);null"`
	ProtoVersion string `orm:"column(proto_version);size(15);null"`
}

var containerModels = []interface{}{
	new(ContainerRepo),
	new(ContainerImage),
}

func init() {
	orm.RegisterModel(containerModels...)
}

const (
	queryContainerTagsList = `select ci.Tag from container_image ci join container_repo cr 
	     on ci.repo_id=cr.id where cr.name=? and ci.proto=? and ci.proto_version=?`
)

func QueryTagsList(reponame string, proto string, proto_version string) ([]string, error) {
	var tags []string
	_, err := orm.NewOrm().Raw(queryContainerTagsList, reponame, proto, proto_version).QueryRows(&tags)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryTagsList] %s", err)
		return nil, err
	}

	return tags, nil
}
