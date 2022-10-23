package bar

import (
	"barckend/controllers/base"
	"barckend/controllers/requests/input"
	"barckend/crud"
	"barckend/link"
	"barckend/photo"
	"barckend/storage"
)

type UploadLogoController struct {
	base.Controller
}

func (c *UploadLogoController) UploadLogo() {
	in := &input.BarIdInPathInput{}
	err := input.ParseInput(c.Ctx.Input, in)
	if err != nil {
		c.BadRequest("Invalid input")
	}
	contentTypes := c.Ctx.Request.Header["Content-Type"]
	var foundContentType string
	for _, cType := range contentTypes {
		if cType == "image/png" || cType == "image/jpeg" {
			foundContentType = cType
			break
		}
	}
	if foundContentType == "" {
		c.BadRequest("Accept jpeg and png only")
	}
	var ext storage.ImageExtension
	if foundContentType == "image/png" {
		ext = storage.PNG
	} else {
		ext = storage.JPEG
	}
	photoId, err := photo.GetPhotoStorage().SavePhoto(photo.Input{
		ByteReader: c.Ctx.Request.Body,
		Ext:        ext,
	})
	if err != nil {
		c.InternalServerError(err)
	}
	photoInfo, err := photo.GetPhotoStorage().GetPhoto(photoId)
	defer photoInfo.File.Close()
	if err != nil {
		c.InternalServerError(err)
	}
	logoUrl, err := link.GetLinkProvider().GetLinkForPhoto(c.Ctx.Request.Host, photoInfo)
	if err != nil {
		c.InternalServerError(err)
	}
	_, err = crud.GetBarCrud().UpdateBar(&crud.UpdateBar{Id: uint64(in.BarId), LogoUrl: &logoUrl}, nil)
	if err != nil {
		c.InternalServerError(err)
	}
	c.Data["json"] = logoUrl
	c.ServeJSONInternal()
}
