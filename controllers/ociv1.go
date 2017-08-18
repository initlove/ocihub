package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/models"
	"github.com/initlove/ocihub/session"
	"github.com/initlove/ocihub/storage"
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
	reponame := this.Ctx.Input.Param(":splat")
	tags := this.Ctx.Input.Param(":tags")
	logs.Debug("GetManifest of '%s:%s'.", reponame, tags)

	data, err := storage.GetManifest(this.Ctx, reponame, tags, "oci", "v1")
	if err != nil {
		//		if err == storage.ErrorNotFound {
		//			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find manifest of '%s'.", reponame))
		//		} else {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get manifest of '%s:%s'.", reponame, tags))
		//		}
		return
	}
	CTX_DATA_WRAP(this.Ctx, http.StatusOK, data, nil)
}

func (this *OCIV1Manifest) PutManifest() {
	reponame := this.Ctx.Input.Param(":splat")
	tags := this.Ctx.Input.Param(":tags")
	logs.Debug("PutManifest of '%s:%s'.", reponame, tags)

	// FIXME, f, h , err
	f, _, err := this.GetFile("filename")
	// TODO: just use Writer and io.Copy to that
	data, _ := ioutil.ReadAll(f)
	f.Close()
	err = storage.PutManifest(this.Ctx, reponame, tags, "oci", "v1", data)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put manifest of '%s:%s'.", reponame, tags))
		return
	}
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in putting manifest of '%s:%s'.", reponame, tags), nil)
}

func (this *OCIV1Manifest) DeleteManifest() {
	reponame := this.Ctx.Input.Param(":splat")
	tags := this.Ctx.Input.Param(":tags")
	logs.Debug("DeleteManifest of '%s:%s'.", reponame, tags)

	err := storage.DeleteManifest(this.Ctx, reponame, tags, "oci", "v1")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete manifest of '%s:%s'.", reponame, tags))
		return
	}
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting manifest of '%s:%s'.", reponame, tags), nil)
}

type OCIV1Blob struct {
	beego.Controller
}

func (this *OCIV1Blob) HeadBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	digest := this.Ctx.Input.Param(":digest")
	logs.Debug("HeadBlob of '%s:%s'.", reponame, digest)

	info, err := storage.HeadBlob(this.Ctx, reponame, digest, "oci", "v1")
	if err != nil {
		//		if err == storage.ErrorNotFound {
		//			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find manifest of '%s'.", reponame))
		//		} else {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to head blob of '%s:%s'.", reponame, digest))
		//		}
		return
	}
	head := make(map[string]string)
	head["Content-Type"] = "application/octec-stream"
	head["Content-Length"] = fmt.Sprint(info.Size())
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, "ok", head)
}

func (this *OCIV1Blob) GetBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	digest := this.Ctx.Input.Param(":digest")
	logs.Debug("GetBlob of '%s:%s'.", reponame, digest)

	data, err := storage.GetBlob(this.Ctx, reponame, digest, "oci", "v1")
	if err != nil {
		//		if err == storage.ErrorNotFound {
		//			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find manifest of '%s'.", reponame))
		//		} else {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get blob of '%s:%s'.", reponame, digest))
		//		}
		return
	}
	CTX_DATA_WRAP(this.Ctx, http.StatusOK, data, nil)
}

func (this *OCIV1Blob) DeleteBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	digest := this.Ctx.Input.Param(":digest")
	logs.Debug("DeleteBlob of '%s:%s'.", reponame, digest)

	err := storage.DeleteBlob(this.Ctx, reponame, digest, "oci", "v1")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete blob of '%s:%s'.", reponame, digest))
		return
	}
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting blob of '%s:%s'.", reponame, digest), nil)
}

// Get session id here
func (this *OCIV1Blob) PostBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	logs.Debug("PostBlob of '%s'.", reponame)

	id, err := session.New(*this.Ctx, "")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to create session to upload blob to '%s'.", reponame))
		return
	}

	header := make(map[string]string)
	header["Content-Type"] = "text/plain"
	header["Session-Id"] = id
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusAccepted, "ok", header)
}

// real start to push
func (this *OCIV1Blob) PatchBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	uuid := this.Ctx.Input.Param(":uuid")
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PatchBlob of '%s:%s'.", reponame, uuid)

	_, err := session.Get(*this.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in patching blob to '%s'.", reponame))
		return
	}

	// FIXME, f, h , err
	f, _, err := this.GetFile("filename")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Cannot find the blob data to '%s'.", reponame))
		return
	}

	// TODO: just use Writer and io.Copy to that
	data, _ := ioutil.ReadAll(f)
	f.Close()
	err = session.PutCache(*this.Ctx, uuid, data)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to cache the patched blob to '%s'.", reponame))
		return
	}
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in patch blob to '%s'.", reponame), nil)
}

// Complete the blob
func (this *OCIV1Blob) PutBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	uuid := this.Ctx.Input.Param(":uuid")
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PutBlob of '%s:%s'.", reponame, uuid)

	_, err := session.Get(*this.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in putting blob to '%s'.", reponame))
		return
	}

	data, err := session.GetCache(*this.Ctx, uuid)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get cached session data in putting blob to '%s'.", reponame))
		return
	}

	err = storage.PutBlob(this.Ctx, reponame, "oci", "v1", data)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put blob of '%s'.", reponame))
		return
	}

	err = session.Release(*this.Ctx, uuid)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to release session in putting blob to '%s'.", reponame))
		return
	}

	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in putting blob of '%s'.", reponame), nil)
}
