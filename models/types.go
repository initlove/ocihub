package models

import (
	"github.com/astaxie/beego/orm"
)

// Repo: the stored repo, should following OCI distribution format
type Repo struct {
	Id          int    `orm:"column(id);auto"`
	Name        string `orm:"unique;column(name);size(255);null"`
	Star        int    `orm:"column(star);null"`
	DownloadNum int    `orm:"column(download_num);null"`
	Description string `orm:"column(description);null"`
}

var basicModels = []interface{}{
	new(Repo),
}

func init() {
	orm.RegisterModel(basicModels...)
}
