package controllers

import (
	"github.com/astaxie/beego"
)

type OCIV1Tag struct {
	beego.Controller
}

func (this *OCIV1Tag) GetTagsList() {
	this.Ctx.Output.Body([]byte("tag list"))
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
