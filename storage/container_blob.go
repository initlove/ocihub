package storage

import (
	"fmt"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/storage/driver"
	"github.com/initlove/ocihub/utils"
)

// repo is not used
func ComposeBlobPath(repo string, digest string, proto string, proto_version string) string {
	head, real := utils.Snap(digest)
	return fmt.Sprintf("%s/%s/blobs/%s/%s", proto, proto_version, head, real)
}

// TODO we need to get user in ctx, or setting in config
func HeadBlob(ctx *context.Context, repo string, digest string, proto string, proto_version string) (driver.FileInfo, error) {
	storagePath := ComposeBlobPath(repo, digest, proto, proto_version)
	logs.Debug("Head '%s'.", storagePath)

	return Driver().Stat(*ctx, storagePath)
}

// TODO we need to get user in ctx, or setting in config
func GetBlob(ctx *context.Context, repo string, digest string, proto string, proto_version string) ([]byte, error) {
	storagePath := ComposeBlobPath(repo, digest, proto, proto_version)
	logs.Debug("Get '%s'.", storagePath)

	return Driver().GetContent(*ctx, storagePath)
}

// TODO we need to get user in ctx, or setting in config
func PutBlob(ctx *context.Context, repo string, proto string, proto_version string, data []byte) error {
	digest := utils.GetDigest("sha256", data)
	storagePath := ComposeBlobPath(repo, digest, proto, proto_version)
	logs.Debug("Put '%s'.", storagePath)

	return Driver().PutContent(*ctx, storagePath, data)
}

// TODO we need to get user in ctx, or setting in config
func DeleteBlob(ctx *context.Context, repo string, digest string, proto string, proto_version string) error {
	storagePath := ComposeBlobPath(repo, digest, proto, proto_version)
	logs.Debug("Delete '%s'.", storagePath)

	return Driver().Delete(*ctx, storagePath)
}
