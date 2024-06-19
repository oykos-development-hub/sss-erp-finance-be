package services

import (
	"context"
	"fmt"

	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	newErrors "gitlab.sudovi.me/erp/finance-api/pkg/errors"

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

func (h *TaxAuthorityCodebookServiceImpl) CreateTaxAuthorityCodebook(ctx context.Context, input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error) {
	data := input.ToTaxAuthorityCodebook()

	id, err := h.repo.Insert(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo tax authority codebook insert")
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo tax authority codebook get")
	}

	res := dto.ToTaxAuthorityCodebookResponseDTO(*data)

	return &res, nil
}

func (h *TaxAuthorityCodebookServiceImpl) UpdateTaxAuthorityCodebook(ctx context.Context, id int, input dto.TaxAuthorityCodebookDTO) (*dto.TaxAuthorityCodebookResponseDTO, error) {
	data := input.ToTaxAuthorityCodebook()
	data.ID = id

	err := h.repo.Update(ctx, *data)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo tax authority codebook update")
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo tax authority codebook get")
	}

	response := dto.ToTaxAuthorityCodebookResponseDTO(*data)

	return &response, nil
}

func (h *TaxAuthorityCodebookServiceImpl) DeleteTaxAuthorityCodebook(ctx context.Context, id int) error {
	err := h.repo.Delete(ctx, id)
	if err != nil {
		return newErrors.Wrap(err, "repo tax authority codebook delete")
	}

	return nil
}

func (h *TaxAuthorityCodebookServiceImpl) DeactivateTaxAuthorityCodebook(ctx context.Context, id int, active bool) error {
	err := h.repo.Deactivate(ctx, id, active)
	if err != nil {
		return newErrors.Wrap(err, "repo tax authority codebook deactivate")
	}

	return nil
}

func (h *TaxAuthorityCodebookServiceImpl) GetTaxAuthorityCodebook(id int) (*dto.TaxAuthorityCodebookResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		return nil, newErrors.Wrap(err, "repo tax authority codebook get")
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
		return nil, nil, newErrors.Wrap(err, "repo tax authority codebook get all")
	}
	response := dto.ToTaxAuthorityCodebookListResponseDTO(data)

	return response, total, nil
}
