package services

import (
	"gitlab.sudovi.me/erp/finance-api/data"
	"gitlab.sudovi.me/erp/finance-api/dto"
	"gitlab.sudovi.me/erp/finance-api/errors"

	"github.com/oykos-development-hub/celeritas"
	up "github.com/upper/db/v4"
)

type ArticleServiceImpl struct {
	App  *celeritas.Celeritas
	repo data.Article
}

func NewArticleServiceImpl(app *celeritas.Celeritas, repo data.Article) ArticleService {
	return &ArticleServiceImpl{
		App:  app,
		repo: repo,
	}
}

func (h *ArticleServiceImpl) CreateArticle(input dto.ArticleDTO) (*dto.ArticleResponseDTO, error) {
	data := input.ToArticle()

	id, err := h.repo.Insert(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = data.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	res := dto.ToArticleResponseDTO(*data)

	return &res, nil
}

func (h *ArticleServiceImpl) UpdateArticle(id int, input dto.ArticleDTO) (*dto.ArticleResponseDTO, error) {
	data := input.ToArticle()
	data.ID = id

	err := h.repo.Update(*data)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	data, err = h.repo.Get(id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	response := dto.ToArticleResponseDTO(*data)

	return &response, nil
}

func (h *ArticleServiceImpl) DeleteArticle(id int) error {
	err := h.repo.Delete(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return errors.ErrInternalServer
	}

	return nil
}

func (h *ArticleServiceImpl) GetArticle(id int) (*dto.ArticleResponseDTO, error) {
	data, err := h.repo.Get(id)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return nil, errors.ErrNotFound
	}
	response := dto.ToArticleResponseDTO(*data)

	return &response, nil
}

func (h *ArticleServiceImpl) GetArticleList(filter dto.ArticleFilterDTO) ([]dto.ArticleResponseDTO, *uint64, error) {
	conditionAndExp := &up.AndExpr{}
	var orders []interface{}

	if filter.InvoiceID != nil {
		conditionAndExp = up.And(conditionAndExp, &up.Cond{"invoice_id": *filter.InvoiceID})
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
	response := dto.ToArticleListResponseDTO(data)

	return response, total, nil
}
