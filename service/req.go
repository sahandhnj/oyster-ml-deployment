package service

import (
	"time"

	"github.com/sahandhnj/apiclient/db"
	"github.com/sahandhnj/apiclient/types"
)

type ReqService struct {
	DBHandler      *db.DBStore
	ModelService   *ModelService
	VersionService *VersionService
}

func NewReqService(ms *ModelService, vs *VersionService, db *db.DBStore) *ReqService {
	return &ReqService{
		DBHandler:      db,
		ModelService:   ms,
		VersionService: vs,
	}
}

func (rs *ReqService) Add(modelName string, versionNumer int, t time.Time, responseTime int64) error {
	model, err := rs.ModelService.GetByName(modelName)
	if err != nil {
		return err
	}

	version, err := rs.VersionService.DBHandler.VersionService.VersionByVersionNumber(versionNumer, model.ID)
	if err != nil {
		return err
	}

	req := &types.Req{
		ID:           rs.DBHandler.ReqService.GetNextIdentifier(),
		ModelId:      model.ID,
		VersionId:    version.ID,
		Time:         t,
		ResponseTime: responseTime,
		Result:       true,
	}

	rs.DBHandler.ReqService.CreateReq(req)

	return nil
}
