package controllers

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/models"
)

type OCIV1Tag struct {
	beego.Controller
}

func (this *OCIV1Tag) GetTagsList() {
	reponame := this.Ctx.Input.Param(":splat")
	logs.Debug("GetTagsList of '%s'.", reponame)

	repo, err := models.QueryTagsList(reponame, "oci", "v1")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get tag list of '%s'.", reponame))
		return
	} else if len(repo) == 0 {
		CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find tag list of '%s'.", reponame))
		return
	}

	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, repo, nil)
}

type OCIV1Manifest struct {
	beego.Controller
}

func (this *OCIV1Manifest) GetManifest() {
}

func (this *OCIV1Manifest) PutManifest() {
}

func (this *OCIV1Manifest) DeleteManifest() {
}

type OCIV1Blob struct {
	beego.Controller
}

func (this *OCIV1Blob) HeadBlob() {
}

func (this *OCIV1Blob) GetBlob() {
}

func (this *OCIV1Blob) DeleteBlob() {
}

func (this *OCIV1Blob) PostBlob() {
}

func (this *OCIV1Blob) PatchBlob() {
}

func (this *OCIV1Blob) PutBlob() {
}
