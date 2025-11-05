package repository

import (
	"errors"

	"github.com/jpeccia/lariharumi_croche_backend_go/config"
	"github.com/jpeccia/lariharumi_croche_backend_go/internal/model"
)

func GetLatestPromotion() (*model.Promotion, error) {
	var p model.Promotion
	err := config.DB.Order("updated_at desc").Limit(1).Find(&p).Error
	if err != nil {
		return nil, err
	}
	if p.ID == 0 {
		return nil, nil
	}
	return &p, nil
}

func SavePromotion(p *model.Promotion) error {
	if p.ID == 0 {
		return errors.New("promotion ID ausente para Save")
	}
	return config.DB.Save(p).Error
}

func CreatePromotion(p *model.Promotion) error {
	return config.DB.Create(p).Error
}
