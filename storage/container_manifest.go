package storage

import (
	"fmt"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// ComposeManifestPath composes the manifest path from the 'proto, prot version, repo and tag'
func ComposeManifestPath(repo string, tag string, proto string, proto_version string) string {
	return fmt.Sprintf("%s/%s/manifests/%s/%s", proto, proto_version, repo, tag)
}

// GetManifest gets the manifest data
// TODO we need to get user in ctx, or setting in config
func GetManifest(ctx *context.Context, repo string, tag string, proto string, proto_version string) ([]byte, error) {
	storagePath := ComposeManifestPath(repo, tag, proto, proto_version)
	logs.Debug("Get '%s'.", storagePath)

	return Driver().GetContent(*ctx, storagePath)
}

// PutManifest puts the manifest data
// TODO we need to get user in ctx, or setting in config
func PutManifest(ctx *context.Context, repo string, tag string, proto string, proto_version string, data []byte) error {
	storagePath := ComposeManifestPath(repo, tag, proto, proto_version)
	logs.Debug("Put '%s'.", storagePath)

	return Driver().PutContent(*ctx, storagePath, data)
}

// DeleteManifest deletes the manifest data
// TODO we need to get user in ctx, or setting in config
func DeleteManifest(ctx *context.Context, repo string, tag string, proto string, proto_version string) error {
	storagePath := ComposeManifestPath(repo, tag, proto, proto_version)
	logs.Debug("Delete '%s'.", storagePath)

	return Driver().Delete(*ctx, storagePath)
}
