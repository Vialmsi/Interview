package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vialmsi/Interview/internal/entity"
	customErrors "github.com/Vialmsi/Interview/internal/errors"
)

type ProductStore interface {
	SaveProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	RetrieveProduct(ctx context.Context, productID string, userID int) (*entity.Product, error)
	DeleteProduct(ctx context.Context, productID string, userID int) error
	RetrieveProductsByUserID(ctx context.Context, userID, limit, page int) ([]entity.Product, int, error)
}

func (s *Service) SaveProduct(ctx context.Context, product *entity.Product, userID int) (*entity.Product, error) {
	s.logger.Info("[SaveProduct] started")

	product.UserID = userID

	exProduct, err := s.store.RetrieveProduct(ctx, product.Barcode, userID)
	fmt.Println(exProduct)
	if exProduct.Barcode != "" {
		s.logger.Info("[SaveProduct] product with this barcode already exists")
		return nil, customErrors.ProductAlreadyExistError
	}

	product, err = s.store.SaveProduct(ctx, product)
	if err != nil {
		s.logger.Errorf("[SaveProduct] error while saving product: %s", err.Error())
		return nil, fmt.Errorf("error while saving product, try later\n%w", err)
	}

	s.logger.Info("[SaveProduct] ended")

	return product, nil
}

func (s *Service) RetrieveProduct(ctx context.Context, barcode string, userID int) (*entity.Product, error) {
	s.logger.Info("[RetrieveProduct] started")

	product, err := s.store.RetrieveProduct(ctx, barcode, userID)
	if err != nil {
		s.logger.Errorf("[RetrieveProduct] error while retrieving product: %s", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NoSuchProductError
		}
		return nil, fmt.Errorf("error while retrieveing product\n%w", err)
	}

	s.logger.Info(product)
	s.logger.Info("[RetrieveProduct] ended")

	return product, nil
}

func (s *Service) DeleteProduct(ctx context.Context, barcode string, userID int) error {
	s.logger.Info("[DeleteProduct] started")

	err := s.store.DeleteProduct(ctx, barcode, userID)
	if err != nil {
		s.logger.Errorf("[DeleteProduct] error while deleting: %s", err)
		if errors.As(err, &customErrors.NoProductToDeleteError) {
			return customErrors.NoProductToDeleteError
		}
		return fmt.Errorf("error while deleting\n%w", err)
	}

	s.logger.Info("[DeleteProduct] ended")

	return nil
}

func (s *Service) RetrieveProductsByUserID(ctx context.Context, userID, limit, page int) ([]entity.Product, int, error) {
	s.logger.Info("[RetrieveProductsByUserID] started")

	products, count, err := s.store.RetrieveProductsByUserID(ctx, userID, limit, page)
	if err != nil {
		s.logger.Errorf("[RetrieveProductsByUserID] error while retrieving products: %s", err.Error())
		return nil, 0, err
	}

	s.logger.Info(products)
	s.logger.Info("[RetrieveProductsByUserID] ended")

	return products, count, nil
}
