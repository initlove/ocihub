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
	RepoId int    `orm:"column(repo_id);null"`

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
	queryContainerReposList = `select name from container_repo order by id asc`
	queryContainerImage     = `select * from container_image 
	     where repo_id=? and tag=? and proto=? and proto_version=? limit 1`
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

func QueryReposList() ([]string, error) {
	var names []string
	_, err := orm.NewOrm().Raw(queryContainerReposList).QueryRows(&names)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryNamesList] %s", err)
		return nil, err
	}

	return names, nil
}

func AddRepo(reponame string) (*ContainerRepo, error) {
	repo := &ContainerRepo{}

	if err := orm.NewOrm().QueryTable("container_repo").
		Filter("Name__exact", reponame).One(repo); err == nil {
		logs.Debug("[AddRepo] repo '%s' is exist.", reponame)
		return repo, nil
	} else if err != orm.ErrNoRows {
		logs.Error("[AddRepo] fail to find repo '%s': %v", reponame, err)
		return nil, err
	}

	repo.Name = reponame
	if _, err := orm.NewOrm().Insert(repo); err != nil {
		logs.Error("[AddRepo] fail to insert repo '%s': %v", reponame, err)
		return nil, err
	}

	return repo, nil
}

func QueryImage(repoid int, tag string, proto string, proto_version string) (*ContainerImage, error) {
	var images []ContainerImage

	_, err := orm.NewOrm().Raw(queryContainerImage, repoid, tag, proto, proto_version).QueryRows(&images)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryImage] %v", err)
		return nil, err
	}

	if len(images) == 0 {
		logs.Debug("[QueryImage] cannot find the row.")
		return nil, nil
	}

	return &images[0], nil
}

//TODO: lots of rollback
func AddImage(reponame string, tags string, proto string, proto_version string) (*ContainerImage, error) {
	repo, err := AddRepo(reponame)
	if err != nil {
		return nil, err
	}

	if img, err := QueryImage(repo.Id, tags, proto, proto_version); err != nil {
		logs.Error("[AddImage] %v", err)
		return nil, err
	} else if img != nil {
		// Already exist, TODO: update info?
		logs.Debug("[AddImage] image is already exist")
		return img, nil
	}

	image := &ContainerImage{}
	image.RepoId = repo.Id
	image.Tag = tags
	image.Proto = proto
	image.ProtoVersion = proto_version
	if _, err := orm.NewOrm().Insert(image); err != nil {
		logs.Error("[AddImage] fail to insert image '%s:%s': %v", reponame, tags, err)
		return nil, err
	}
	return image, nil
}
