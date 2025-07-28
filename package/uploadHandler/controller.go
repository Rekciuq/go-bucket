package uploadhandler

import (
	"github.com/Rekciuq/go-bucket/ent"
	entclient "github.com/Rekciuq/go-bucket/package/ent-client"
)

type uploadController struct {
	connection *entclient.ClientConnection
}

func (uc *uploadController) Create(uuid string) error {
	err := uc.connection.Client.Url.Create().SetID(uuid).Exec(uc.connection.Ctx)
	return err
}

func (uc *uploadController) IsAlreadyExists(uuid string) bool {
	_, err := uc.connection.Client.Url.Get(uc.connection.Ctx, uuid)
	if err != nil {
		return false
	}

	return true
}

func (uc *uploadController) IsURLUsed(urlId string) (bool, error) {
	url, err := uc.connection.Client.Url.Get(uc.connection.Ctx, urlId)
	if err != nil {
		return false, err
	}
	return url.IsUsed, nil
}

func (uc *uploadController) UseUrl(urlId string) error {
	err := uc.connection.Client.Url.UpdateOneID(urlId).SetIsUsed(true).Exec(uc.connection.Ctx)
	return err
}

func (uc *uploadController) GetUrl(urlId string) (*ent.Url, error) {
	url, err := uc.connection.Client.Url.Get(uc.connection.Ctx, urlId)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (uc *uploadController) CreateImage(id, path string) error {
	err := uc.connection.Client.Image.Create().SetID(id).SetPath(path).Exec(uc.connection.Ctx)
	return err
}
