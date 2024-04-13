package services

import (
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type TaxAuthorityCodebookServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.TaxAuthorityCodebook
}

func NewTaxAuthorityCodebookServiceImpl(app *celeritas.Celeritas, repo data.TaxAuthorityCodebook) TaxAuthorityCodebookService {
	return &TaxAuthorityCodebookServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *TaxAuthorityCodebookServiceImpl) CreateTaxAuthorityCodebook(input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error) {
	data := input.ToTaxAuthorityCodebook()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToTaxAuthorityCodebookResponseDTO(*data)

	return &res, nil
}

func (h *TaxAuthorityCodebookServiceImpl) UpdateTaxAuthorityCodebook(id int, input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error) {
	data := input.ToTaxAuthorityCodebook()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToTaxAuthorityCodebookResponseDTO(*data)

	return &response, nil
}

func (h *TaxAuthorityCodebookServiceImpl) DeleteTaxAuthorityCodebook(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *TaxAuthorityCodebookServiceImpl) DeactivateTaxAuthorityCodebook(id int, active bool) error {
	err := h.repo.Deactivate(id, active)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *TaxAuthorityCodebookServiceImpl) GetTaxAuthorityCodebook(id int) (*dto.TaxAuthorityCodebookResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToTaxAuthorityCodebookResponseDTO(*data)

	return &response, nil
}

func (h *TaxAuthorityCodebookServiceImpl) GetTaxAuthorityCodebookList(filter dto.TaxAuthorityCodebookFilterDTO) ([]dto.TaxAuthorityCodebookResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	// example of making conditions
	if filter.Search != nil && *filter.Search != "" {
		likeCondition := fmt.Sprintf("%%%s%%", *filter.Search)
		search := up.Or(
			up.Cond{"title ILIKE": likeCondition},
		)
		conditionAndExp = up.And(conditionAndExp, search)
	}

	if filter.Active != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"active": *filter.Active})
	}

	if filter.SortByTitle != nil {
		if *filter.SortByTitle == "asc" {
			orders = append(orders, "-title")
		} else {
			orders = append(orders, "title")
		}
	}

	orders = append(orders, "-created_at")

	data, total, err := h.repo.GetAll(filter.Page, filter.Size, conditionAndExp, orders)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, nil, errors.ErrInternalServer
	}
	response := dto.ToTaxAuthorityCodebookListResponseDTO(data)

	return response, total, nil
}
