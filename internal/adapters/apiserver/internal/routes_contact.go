package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/skvenkat/golang-chi-rest-api/internal/core/usecase"
	"net/http"
)

func CreateContact(uc *usecase.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &ContactToSaveRest{}
		if err := render.Bind(r, req); err != nil {
			_ = render.Render(w, r, ErrBadRequest(err))
			return
		}
		contactToSave, err := req.toModel()
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest(err))
			return
		}
		c, err := uc.AddAddrBookContact(r.Context(), contactToSave)
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest(err))
			return
		}
		resp := contactModelToRest(c)
		render.Status(r, http.StatusCreated)
		_ = render.Render(w, r, resp)
	}
}

func UpdateContact(uc *usecase.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactId := chi.URLParam(r, "contactId")
		if contactId == "" {
			_ = render.Render(w, r, NotFoundErrResponse)
			return
		}
		req := &ContactToSaveRest{}
		if err := render.Bind(r, req); err != nil {
			_ = render.Render(w, r, ErrBadRequest(err))
			return
		}
		contactToSave, err := req.toModel()
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest(err))
			return
		}
		c, found, err := uc.UpdateAddrBookContact(r.Context(), contactId, contactToSave)
		if !found {
			_ = render.Render(w, r, NotFoundErrResponse)
			return
		}
		if err != nil {
			_ = render.Render(w, r, ErrBadRequest(err))
			return
		}
		resp := contactModelToRest(c)
		_ = render.Render(w, r, resp)
	}
}

func DeleteContact(uc *usecase.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactId := chi.URLParam(r, "contactId")
		if contactId == "" {
			_ = render.Render(w, r, NotFoundErrResponse)
			return
		}
		found, err := uc.DeleteAddrBookContact(r.Context(), contactId)
		if err != nil {
			_ = render.Render(w, r, NewInternalServerErrResponse(err))
			return
		}
		if !found {
			_ = render.Render(w, r, NotFoundErrResponse)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func ListAllContacts(uc *usecase.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts, err := uc.LoadAddrBookContacts(r.Context())
		if err != nil {
			_ = render.Render(w, r, NewInternalServerErrResponse(err))
			return
		}
		contactRenderers := make([]render.Renderer, len(contacts))
		for i, contact := range contacts {
			contactRenderers[i] = contactModelToRest(contact)
		}
		if err = render.RenderList(w, r, contactRenderers); err != nil {
			_ = render.Render(w, r, ErrRender(err))
			return
		}
		render.Status(r, http.StatusOK)
	}
}

func GetContact(uc *usecase.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactId := chi.URLParam(r, "contactId")
		if contactId == "" {
			_ = render.Render(w, r, NotFoundErrResponse)
			return
		}
		c, err := uc.LoadAddrBookContactByID(r.Context(), contactId)
		if err != nil {
			_ = render.Render(w, r, NewInternalServerErrResponse(err))
			return
		}
		if c == nil {
			_ = render.Render(w, r, NotFoundErrResponse)
			return
		}
		resp := contactModelToRest(c)
		_ = render.Render(w, r, resp)
	}
}

// Bind required to properly deserialize POST body to contactToSaveRest value
func (u *ContactToSaveRest) Bind(r *http.Request) error {
	// you can perform contactToSaveRest value validation here
	// we do it in .toModel() call instead
	return nil
}

// Render required to properly serialize contactRest value into HTTP body response
func (rd *ContactRest) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}
