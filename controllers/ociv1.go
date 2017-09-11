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

// OCIV1Tag defines the tag operation
type OCIV1Tag struct {
	beego.Controller
}

// GetTagsList gets the tags list
func (o *OCIV1Tag) GetTagsList() {
	reponame := o.Ctx.Input.Param(":splat")
	logs.Debug("GetTagsList of '%s'.", reponame)

	repo, err := models.QueryTagsList(reponame, "oci", "v1")
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get tag list of '%s'.", reponame))
		return
	} else if len(repo) == 0 {
		CtxErrorWrap(o.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find tag list of '%s'.", reponame))
		return
	}

	CtxSuccessWrap(o.Ctx, http.StatusOK, repo, nil)
}

// OCIV1Manifest defines the manifest operation
type OCIV1Manifest struct {
	beego.Controller
}

// GetManifest gets the manifest of 'repo:tag'
func (o *OCIV1Manifest) GetManifest() {
	reponame := o.Ctx.Input.Param(":splat")
	tags := o.Ctx.Input.Param(":tags")
	logs.Debug("GetManifest of '%s:%s'.", reponame, tags)

	data, err := storage.GetManifest(o.Ctx, reponame, tags, "oci", "v1")
	if err != nil {
		//		if err == storage.ErrorNotFound {
		//			CtxErrorWrap(o.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find manifest of '%s'.", reponame))
		//		} else {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get manifest of '%s:%s'.", reponame, tags))
		//		}
		return
	}
	CtxDataWrap(o.Ctx, http.StatusOK, data, nil)
}

// PutManifest puts the manifest of 'repo:tag'
func (o *OCIV1Manifest) PutManifest() {
	reponame := o.Ctx.Input.Param(":splat")
	tags := o.Ctx.Input.Param(":tags")
	logs.Debug("PutManifest of '%s:%s'.", reponame, tags)

	// FIXME, f, h , err
	f, _, err := o.GetFile("filename")
	// TODO: just use Writer and io.Copy to that
	data, _ := ioutil.ReadAll(f)
	f.Close()
	err = storage.PutManifest(o.Ctx, reponame, tags, "oci", "v1", data)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put manifest of '%s:%s'.", reponame, tags))
		return
	}
	CtxSuccessWrap(o.Ctx, http.StatusOK, fmt.Sprintf("Succeed in putting manifest of '%s:%s'.", reponame, tags), nil)
}

// DeleteManifest deletes the manifest of 'repo:tag'
func (o *OCIV1Manifest) DeleteManifest() {
	reponame := o.Ctx.Input.Param(":splat")
	tags := o.Ctx.Input.Param(":tags")
	logs.Debug("DeleteManifest of '%s:%s'.", reponame, tags)

	err := storage.DeleteManifest(o.Ctx, reponame, tags, "oci", "v1")
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete manifest of '%s:%s'.", reponame, tags))
		return
	}
	CtxSuccessWrap(o.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting manifest of '%s:%s'.", reponame, tags), nil)
}

// OCIV1Blob defines the blob operations
type OCIV1Blob struct {
	beego.Controller
}

// HeadBlob queries the blob info
func (o *OCIV1Blob) HeadBlob() {
	reponame := o.Ctx.Input.Param(":splat")
	digest := o.Ctx.Input.Param(":digest")
	logs.Debug("HeadBlob of '%s:%s'.", reponame, digest)

	info, err := storage.HeadBlob(o.Ctx, reponame, digest, "oci", "v1")
	if err != nil {
		//		if err == storage.ErrorNotFound {
		//			CtxErrorWrap(o.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find manifest of '%s'.", reponame))
		//		} else {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to head blob of '%s:%s'.", reponame, digest))
		//		}
		return
	}
	head := make(map[string]string)
	head["Content-Type"] = "application/octec-stream"
	head["Content-Length"] = fmt.Sprint(info.Size())
	CtxSuccessWrap(o.Ctx, http.StatusOK, "ok", head)
}

// GetBlob gets the blob of a certain digest
func (o *OCIV1Blob) GetBlob() {
	reponame := o.Ctx.Input.Param(":splat")
	digest := o.Ctx.Input.Param(":digest")
	logs.Debug("GetBlob of '%s:%s'.", reponame, digest)

	data, err := storage.GetBlob(o.Ctx, reponame, digest, "oci", "v1")
	if err != nil {
		//		if err == storage.ErrorNotFound {
		//			CtxErrorWrap(o.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find manifest of '%s'.", reponame))
		//		} else {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get blob of '%s:%s'.", reponame, digest))
		//		}
		return
	}
	CtxDataWrap(o.Ctx, http.StatusOK, data, nil)
}

// DeleteBlob deletes the blob of a certain digest
func (o *OCIV1Blob) DeleteBlob() {
	reponame := o.Ctx.Input.Param(":splat")
	digest := o.Ctx.Input.Param(":digest")
	logs.Debug("DeleteBlob of '%s:%s'.", reponame, digest)

	err := storage.DeleteBlob(o.Ctx, reponame, digest, "oci", "v1")
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete blob of '%s:%s'.", reponame, digest))
		return
	}
	CtxSuccessWrap(o.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting blob of '%s:%s'.", reponame, digest), nil)
}

// PostBlob starts to post a blob and get an uuid in return
// It is just a mimic of docker post blob
func (o *OCIV1Blob) PostBlob() {
	reponame := o.Ctx.Input.Param(":splat")
	logs.Debug("PostBlob of '%s'.", reponame)

	id, err := session.New(*o.Ctx, "")
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to create session to upload blob to '%s'.", reponame))
		return
	}

	header := make(map[string]string)
	header["Content-Type"] = "text/plain"
	header["Session-Id"] = id
	CtxSuccessWrap(o.Ctx, http.StatusAccepted, "ok", header)
}

// PatchBlob starts to upload a blob
// It is just a mimic of docker patch blob
func (o *OCIV1Blob) PatchBlob() {
	reponame := o.Ctx.Input.Param(":splat")
	uuid := o.Ctx.Input.Param(":uuid")
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PatchBlob of '%s:%s'.", reponame, uuid)

	_, err := session.Get(*o.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in patching blob to '%s'.", reponame))
		return
	}

	// FIXME, f, h , err
	f, _, err := o.GetFile("filename")
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusBadRequest, err, fmt.Sprintf("Cannot find the blob data to '%s'.", reponame))
		return
	}

	// TODO: just use Writer and io.Copy to that
	data, _ := ioutil.ReadAll(f)
	f.Close()
	err = session.PutCache(*o.Ctx, uuid, data)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to cache the patched blob to '%s'.", reponame))
		return
	}
	CtxSuccessWrap(o.Ctx, http.StatusOK, fmt.Sprintf("Succeed in patch blob to '%s'.", reponame), nil)
}

// PutBlob marks the blob uploading status to done
// It is just a mimic of docker putblob
func (o *OCIV1Blob) PutBlob() {
	reponame := o.Ctx.Input.Param(":splat")
	uuid := o.Ctx.Input.Param(":uuid")
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PutBlob of '%s:%s'.", reponame, uuid)

	_, err := session.Get(*o.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in putting blob to '%s'.", reponame))
		return
	}

	data, err := session.GetCache(*o.Ctx, uuid)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get cached session data in putting blob to '%s'.", reponame))
		return
	}

	err = storage.PutBlob(o.Ctx, reponame, "oci", "v1", data)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put blob of '%s'.", reponame))
		return
	}

	err = session.Release(*o.Ctx, uuid)
	if err != nil {
		CtxErrorWrap(o.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to release session in putting blob to '%s'.", reponame))
		return
	}

	CtxSuccessWrap(o.Ctx, http.StatusOK, fmt.Sprintf("Succeed in putting blob of '%s'.", reponame), nil)
}
