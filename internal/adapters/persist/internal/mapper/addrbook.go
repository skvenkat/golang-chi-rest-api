package mapper

import (
	"fmt"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/persist/internal/repo"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/model"

	"github.com/samber/lo"
)

func ContactEntityToModel(e *repo.ContactWithPhonesEntity) *model.Contact {
	phones := make([]*model.ContactPhone, len(e.Phones))
	for i, ph := range e.Phones {
		phones[i] = &model.ContactPhone{
			PhoneType:   phoneTypeEntityToModel(ph.PhoneType),
			PhoneNumber: ph.PhoneNumber,
		}
	}
	return &model.Contact{
		ID:        RepoIdToModelId(e.ID),
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Phones:    phones,
	}
}

func ContactToSaveModelToEntity(m *model.ContactToSave) *repo.ContactWithPhonesEntity {
	return &repo.ContactWithPhonesEntity{
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Phones: lo.Map(m.Phones, func(item *model.ContactPhoneToSave, _ int) *repo.PhoneEntity {
			return &repo.PhoneEntity{
				PhoneType:   phoneTypeModelToEntity(item.PhoneType),
				PhoneNumber: item.PhoneNumber,
			}
		}),
	}
}

func phoneTypeEntityToModel(phoneType string) model.ContactPhoneType {
	switch phoneType {
	case "mobile":
		return model.ContactPhoneTypeMobile
	case "home":
		return model.ContactPhoneTypeHome
	case "work":
		return model.ContactPhoneTypeWork
	default:
		panic(fmt.Sprintf("unexpected database mobile phone type: %s", phoneType))
	}
}

func phoneTypeModelToEntity(phoneType model.ContactPhoneType) string {
	switch phoneType {
	case model.ContactPhoneTypeMobile:
		return "mobile"
	case model.ContactPhoneTypeHome:
		return "home"
	case model.ContactPhoneTypeWork:
		return "work"
	default:
		panic(fmt.Sprintf("unexpected model phone type: %s", phoneType))
	}
}
