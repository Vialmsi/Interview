package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Vialmsi/Interview/internal/entity"
)

const (
	fileStorage   = "files"
	jsonExtension = ".json"
	// filePermission - readable by all the user groups, but writable by the user only
	filePermission = 0644
	// dirPermission - full permissions
	dirPermission = 0777
)

type ProductStore struct {
}

func NewProductStore() (*ProductStore, error) {
	_, err := os.Open(fileStorage)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fileStorage, dirPermission)
		if err != nil {
			return nil, fmt.Errorf("error while creating file storage: %s", err.Error())
		}
	}
	return &ProductStore{}, nil
}

// SaveProduct func which allows to save product into dir
func (p *ProductStore) SaveProduct(product entity.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	strID := strconv.Itoa(product.UserID)

	err = os.WriteFile(filepath.Join(fileStorage, strID, product.Barcode)+jsonExtension, data, filePermission)
	if err != nil {
		return err
	}

	return err
}

// RetrieveProduct func which allows to get product from dir using file name
func (p *ProductStore) RetrieveProduct(fileName string, userID int) (*entity.Product, error) {
	strID := strconv.Itoa(userID)
	f, err := os.Open(filepath.Join(fileStorage, strID, fileName+jsonExtension))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, fi.Size())
	_, err = f.Read(data)
	if err != nil {
		return nil, err
	}

	product := &entity.Product{}

	err = json.Unmarshal(data, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct func which allows to delete product from dir
func (p *ProductStore) DeleteProduct(fileName string, userID int) error {
	strID := strconv.Itoa(userID)
	return os.Remove(filepath.Join(fileStorage, strID, fileName) + jsonExtension)
}

// RetrieveProductsByUserID func which allows to get every user's product
func (p *ProductStore) RetrieveProductsByUserID(userID int) ([]entity.Product, error) {
	strID := strconv.Itoa(userID)
	dir, err := os.Open(filepath.Join(fileStorage, strID))
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	products := make([]entity.Product, len(files))

	for i, file := range files {
		product, err := p.RetrieveProduct(strings.Trim(file.Name(), jsonExtension), userID)
		if err != nil {
			return nil, err
		}
		products[i] = *product
	}

	return products, nil
}

func (p *ProductStore) CreateUserStorage(userID int) error {
	strID := strconv.Itoa(userID)

	return os.Mkdir(filepath.Join(fileStorage, strID), dirPermission)
}
