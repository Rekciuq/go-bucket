package uploadcontroller

import (
	"github.com/Rekciuq/go-bucket/ent"
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
)

type UploadController struct {
	Connection *entclient.ClientConnection
}

func (uc *UploadController) Create(uuid string) error {
	err := uc.Connection.Client.Url.Create().SetID(uuid).Exec(uc.Connection.Ctx)
	return err
}

func (uc *UploadController) IsAlreadyExists(uuid string) bool {
	_, err := uc.Connection.Client.Url.Get(uc.Connection.Ctx, uuid)
	if err != nil {
		return false
	}

	return true
}

func (uc *UploadController) IsURLUsed(urlId string) (bool, error) {
	url, err := uc.Connection.Client.Url.Get(uc.Connection.Ctx, urlId)
	if err != nil {
		return false, err
	}
	return url.IsUsed, nil
}

func (uc *UploadController) UseUrl(urlId string) error {
	err := uc.Connection.Client.Url.UpdateOneID(urlId).SetIsUsed(true).Exec(uc.Connection.Ctx)
	return err
}

func (uc *UploadController) GetUrl(urlId string) (*ent.Url, error) {
	url, err := uc.Connection.Client.Url.Get(uc.Connection.Ctx, urlId)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (uc *UploadController) CreateImage(id, path string) error {
	err := uc.Connection.Client.Image.Create().SetID(id).SetPath(path).Exec(uc.Connection.Ctx)
	return err
}
