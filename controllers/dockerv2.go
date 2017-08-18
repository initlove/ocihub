package controllers

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/models"
	"github.com/initlove/ocihub/session"
	"github.com/initlove/ocihub/storage"
	"github.com/initlove/ocihub/storage/driver"
	"github.com/initlove/ocihub/utils"
)

type DockerV2Ping struct {
	beego.Controller
}

func (this *DockerV2Ping) Ping() {
	head := make(map[string]string)
	head["Content-Type"] = "application/json; charset=utf-8"
	head["Docker-Distribution-Api-Version"] = "registry/2.0"
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, "{}", head)

}

type DockerV2Tag struct {
	beego.Controller
}

func (this *DockerV2Tag) GetTagsList() {
	reponame := this.Ctx.Input.Param(":splat")
	logs.Debug("GetTagsList of '%s'.", reponame)

	repo, err := models.QueryTagsList(reponame, "docker", "v2")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get tag list of '%s'.", reponame))
		return
	} else if len(repo) == 0 {
		CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, nil, fmt.Sprintf("Cannot find tag list of '%s'.", reponame))
		return
	}

	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, repo, nil)
}

type DockerV2Manifest struct {
	beego.Controller
}

func (this *DockerV2Manifest) GetManifest() {
	reponame := this.Ctx.Input.Param(":splat")
	tags := this.Ctx.Input.Param(":tags")
	logs.Debug("GetManifest of '%s:%s'.", reponame, tags)

	data, err := storage.GetManifest(this.Ctx, reponame, tags, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to get manifest of '%s:%s'.", reponame, tags))
		} else {
			CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get manifest of '%s:%s'.", reponame, tags))
		}
		return
	}

	digest, _ := utils.DigestManifest(data)
	header := make(map[string]string)
	header["Docker-Content-Digest"] = digest
	header["Content-Length"] = fmt.Sprint(len(data))
	CTX_DATA_WRAP(this.Ctx, http.StatusOK, data, header)
}

func (this *DockerV2Manifest) PutManifest() {
	reponame := this.Ctx.Input.Param(":splat")
	tags := this.Ctx.Input.Param(":tags")
	logs.Debug("PutManifest of '%s:%s'.", reponame, tags)

	data := this.Ctx.Input.CopyBody(utils.MaxSize)
	logs.Debug("Ma <%s>", data)
	err := storage.PutManifest(this.Ctx, reponame, tags, "docker", "v2", data)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put manifest of '%s:%s'.", reponame, tags))
		return
	}

	//TODO: rollback the storage.. add error checks
	_, err = models.AddImage(reponame, tags, "docker", "v2")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to add image '%s:%s' to db.", reponame, tags))
		return
	}

	digest, _ := utils.DigestManifest(data)
	header := make(map[string]string)
	header["Docker-Content-Digest"] = digest
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, "{}", header)
}

func (this *DockerV2Manifest) DeleteManifest() {
	reponame := this.Ctx.Input.Param(":splat")
	tags := this.Ctx.Input.Param(":tags")
	logs.Debug("DeleteManifest of '%s:%s'.", reponame, tags)

	err := storage.DeleteManifest(this.Ctx, reponame, tags, "docker", "v2")
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete manifest of '%s:%s'.", reponame, tags))
		return
	}
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting manifest of '%s:%s'.", reponame, tags), nil)
}

type DockerV2Blob struct {
	beego.Controller
}

func (this *DockerV2Blob) HeadBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	digest := this.Ctx.Input.Param(":digest")
	logs.Debug("HeadBlob of '%s:%s'.", reponame, digest)

	info, err := storage.HeadBlob(this.Ctx, reponame, digest, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to head blob of '%s:%s'.", reponame, digest))
		} else {
			CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to head blob of '%s:%s'.", reponame, digest))
		}
		return
	}
	head := make(map[string]string)
	head["Content-Type"] = "application/octec-stream"
	head["Content-Length"] = fmt.Sprint(info.Size())
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, "ok", head)
}

func (this *DockerV2Blob) GetBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	digest := this.Ctx.Input.Param(":digest")
	logs.Debug("GetBlob of '%s:%s'.", reponame, digest)

	data, err := storage.GetBlob(this.Ctx, reponame, digest, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to get blob of '%s:%s'.", reponame, digest))
		} else {
			CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get blob of '%s:%s'.", reponame, digest))
		}
		return
	}
	CTX_DATA_WRAP(this.Ctx, http.StatusOK, data, nil)
}

func (this *DockerV2Blob) DeleteBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	digest := this.Ctx.Input.Param(":digest")
	logs.Debug("DeleteBlob of '%s:%s'.", reponame, digest)

	err := storage.DeleteBlob(this.Ctx, reponame, digest, "docker", "v2")
	if err != nil {
		if _, ok := err.(driver.PathNotFoundError); ok {
			CTX_ERROR_WRAP(this.Ctx, http.StatusNotFound, err, fmt.Sprintf("Failed to delete blob of '%s:%s'.", reponame, digest))
		} else {
			CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to delete blob of '%s:%s'.", reponame, digest))
		}
		return
	}
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, fmt.Sprintf("Succeed in deleting blob of '%s:%s'.", reponame, digest), nil)
}

// Get session id here
func (this *DockerV2Blob) PostBlob() {
	reponame := this.Ctx.Input.Param(":splat")
	mount := this.Ctx.Input.Query("mount")
	logs.Debug("PostBlob of '%s:[%s]'.", reponame, mount)

	uuid, err := session.New(*this.Ctx, mount)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to create session to upload blob to '%s'.", reponame))
		return
	}
	header := make(map[string]string)
	header["Docker-Upload-UUID"] = uuid
	header["Range"] = "0-0"
	header["Content-Length"] = "0"
	header["Location"] = fmt.Sprintf("%s%s", this.Ctx.Input.URL(), uuid)
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusAccepted, "ok", header)
}

// real start to push
func (this *DockerV2Blob) PatchBlob() {
	var uuid string
	reponame := this.Ctx.Input.Param(":splat")
	mount := this.Ctx.Input.Query("mount")
	if mount == "" {
		uuid = this.Ctx.Input.Param(":uuid")
	} else {
		uuid = mount
	}
	// FIXME: Warn: for security reason, we should not output the uuid
	logs.Debug("PatchBlob of '%s:%s'.", reponame, uuid)
	_, err := session.Get(*this.Ctx, uuid)
	if err != nil {
		// TODO: not found error
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to get session in patching blob to '%s'.", reponame))
		return
	}

	data := this.Ctx.Input.CopyBody(utils.MaxSize)
	err = session.PutCache(*this.Ctx, uuid, data)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to cache the patched blob to '%s'.", reponame))
		return
	}

	header := make(map[string]string)
	header["Docker-Upload-UUID"] = uuid
	header["Content-Length"] = "0"
	header["Location"] = fmt.Sprintf("%s", this.Ctx.Input.URL())
	header["Range"] = fmt.Sprintf("0-%v", len(data)-1)

	CTX_SUCCESS_WRAP(this.Ctx, http.StatusNoContent, fmt.Sprintf("Succeed in patch blob to '%s'.", reponame), header)
}

// Complete the blob
func (this *DockerV2Blob) PutBlob() {
	var uuid string
	reponame := this.Ctx.Input.Param(":splat")
	mount := this.Ctx.Input.Query("mount")
	if mount == "" {
		uuid = this.Ctx.Input.Param(":uuid")
	} else {
		uuid = mount
	}
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

	err = storage.PutBlob(this.Ctx, reponame, "docker", "v2", data)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to put blob of '%s'.", reponame))
		return
	}

	err = session.Release(*this.Ctx, uuid)
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, fmt.Sprintf("Failed to release session in putting blob to '%s'.", reponame))
		return
	}

	digest := utils.GetDigest("sha256", data)
	header := make(map[string]string)
	header["Content-Length"] = "0"
	header["Content-Range"] = fmt.Sprintf("0-%v", len(data)-1)
	header["Docker-Content-Digest"] = digest
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusNoContent, fmt.Sprintf("Succeed in putting blob of '%s'.", reponame), nil)
}

type DockerV2Repo struct {
	beego.Controller
}

func (this *DockerV2Repo) GetRepoList() {
	logs.Debug("GetRepoList")

	repos, err := models.QueryReposList()
	if err != nil {
		CTX_ERROR_WRAP(this.Ctx, http.StatusInternalServerError, err, "Fail to get repos list")
		return
	}

	type cataLog struct {
		Repositories []string
	}
	var c cataLog
	c.Repositories = repos
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	CTX_SUCCESS_WRAP(this.Ctx, http.StatusOK, c, header)
}
