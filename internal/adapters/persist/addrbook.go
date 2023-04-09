package persist

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/persist/internal/cache"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/persist/internal/mapper"
	"github.com/skvenkat/golang-chi-rest-api/internal/adapters/persist/internal/repo"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/app"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/model"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/outport"
)

type addrBookAdapter struct {
	repo             *repo.AddrBookRepo
	contactByIdCache cache.ContactByIdPartition
}

func NewAddrBookAdapter(
	p outport.Persistence,
	c outport.Cache,
) outport.AddrBook {
	return &addrBookAdapter{
		repo:             repo.NewAddrBookRepo(p.DB()),
		contactByIdCache: cache.RegisterContactByID(c),
	}
}

func (a *addrBookAdapter) LoadAllContacts(ctx context.Context) ([]*model.Contact, error) {
	all, err := a.repo.SelectAllContacts(ctx)
	if err != nil {
		return nil, err
	}
	return lo.Map(all, func(item *repo.ContactWithPhonesEntity, _ int) *model.Contact {
		return mapper.ContactEntityToModel(item)
	}), nil
}

func (a *addrBookAdapter) LoadContactByID(ctx context.Context, ID string) (*model.Contact, error) {
	if cachedContact := a.contactByIdCache.Get(ctx, ID); cachedContact != nil {
		return cachedContact, nil
	}
	repoID, err := mapper.ModelIdToRepoId(ID)
	if err != nil {
		panic(fmt.Sprintf("invalid contact ID: %s", ID))
	}
	entity, err := a.repo.SelectContactByID(ctx, repoID)
	if err != nil || entity == nil {
		return nil, err
	}
	m := mapper.ContactEntityToModel(entity)
	a.contactByIdCache.Set(ctx, m)
	return m, nil
}

func (a *addrBookAdapter) AddContact(ctx context.Context, c *model.ContactToSave) (*model.Contact, error) {
	entity := mapper.ContactToSaveModelToEntity(c)
	entity, err := a.repo.AddContact(ctx, entity)
	if err != nil {
		return nil, err
	}
	contact := mapper.ContactEntityToModel(entity)
	a.contactByIdCache.Set(ctx, contact)
	return contact, nil
}

func (a *addrBookAdapter) UpdateContact(ctx context.Context, ID string, c *model.ContactToSave) (*model.Contact, error) {
	var err error
	entity := mapper.ContactToSaveModelToEntity(c)
	entity.ID, err = mapper.ModelIdToRepoId(ID)
	if err != nil {
		app.Logger(ctx).Debugln("error parsing id:", ID)
		return nil, nil
	}
	found, err := a.repo.UpdateContact(ctx, entity)
	if err == nil {
		if !found {
			return nil, nil
		}
		entity, err = a.repo.SelectContactByID(ctx, entity.ID)
		if err != nil || entity == nil {
			return nil, err
		}
		contact := mapper.ContactEntityToModel(entity)
		a.contactByIdCache.Set(ctx, contact)
		return contact, nil
	}
	return nil, err
}

func (a *addrBookAdapter) DeleteContact(ctx context.Context, ID string) (found bool, err error) {
	repoID, err := mapper.ModelIdToRepoId(ID)
	if err != nil {
		app.Logger(ctx).Debugln("error parsing id:", ID)
		return false, nil // no error is needed, we assume that record does not exist
	}
	found, err = a.repo.DeleteContact(ctx, repoID)
	if found {
		a.contactByIdCache.Del(ctx, ID)
	}
	return
}
