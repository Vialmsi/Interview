package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/Vialmsi/Interview/internal/entity"
)

type ProductStore interface {
	SaveProduct(product entity.Product) error
	RetrieveProduct(fileName string, userID int) (*entity.Product, error)
	DeleteProduct(fileName string, userID int) error
	RetrieveProductsByUserID(userID int) ([]entity.Product, error)
	CreateUserStorage(userID int) error
}

func (s *Service) SaveProduct(userID int, product entity.Product) error {
	s.logger.Info("[SaveProduct] started")

	product.UserID = userID

	exProduct, err := s.productStore.RetrieveProduct(product.Barcode, product.UserID)
	fmt.Println(exProduct)
	if exProduct != nil {
		s.logger.Info("[SaveProduct] product with this barcode already exists")
		return fmt.Errorf("product with this barcode already exists")
	}

	err = s.productStore.SaveProduct(product)
	if err != nil {
		s.logger.Errorf("[SaveProduct] error while saving product: %s", err.Error())
		return fmt.Errorf("error while saving product, try later\n%w", err)
	}

	s.logger.Info("[SaveProduct] ended")

	return nil
}

func (s *Service) RetrieveProduct(fileName string, userID int) (*entity.Product, error) {
	s.logger.Info("[RetrieveProduct] started")

	product, err := s.productStore.RetrieveProduct(fileName, userID)
	if err != nil {
		s.logger.Errorf("[RetrieveProduct] error while retrieving product: %s", err.Error())
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("no such product")
		}
		return nil, fmt.Errorf("error while retrieveing product\n%w", err)
	}

	s.logger.Info(product)
	s.logger.Info("[RetrieveProduct] ended")

	return product, nil
}

func (s *Service) DeleteProduct(fileName string, userID int) error {
	s.logger.Info("[DeleteProduct] started")

	err := s.productStore.DeleteProduct(fileName, userID)
	if err != nil {
		s.logger.Errorf("[DeleteProduct] error while deleting: %s", err)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("project with this barcode not exists")
		}
		return fmt.Errorf("error while deleting\n%w", err)
	}

	s.logger.Info("[DeleteProduct] ended")

	return nil
}

func (s *Service) RetrieveProductsByUserID(userID int) ([]entity.Product, error) {
	s.logger.Info("[RetrieveProductsByUserID] started")

	products, err := s.productStore.RetrieveProductsByUserID(userID)
	if err != nil {
		s.logger.Errorf("[RetrieveProductsByUserID] error while retrieving products: %s", err.Error())
		return nil, err
	}

	s.logger.Info(products)
	s.logger.Info("[RetrieveProductsByUserID] ended")

	return products, nil
}
