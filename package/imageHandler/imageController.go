package imagehandler

import (
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
)

type ImageController struct {
	Connection *entclient.ClientConnection
}

func (ic *ImageController) GetImage(imageId string) (string, error) {
	image, err := ic.Connection.Client.Image.Get(ic.Connection.Ctx, imageId)
	if err != nil {
		return "", err
	}
	return image.Path, nil
}
