package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Vialmsi/Interview/internal/entity"
)

const (
	FileStorage    = "files"
	FilePermission = 0644
	JSONExtension  = ".json"
)

type ProductStore struct {
}

func NewProductStore() *ProductStore {
	return &ProductStore{}
}

// SaveProduct func which allows to save product into dir
func (p *ProductStore) SaveProduct(fileName, userDir string, product entity.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(FileStorage, userDir, fileName)+JSONExtension, data, FilePermission)
	if err != nil {
		return err
	}

	return err
}

// RetrieveProduct func which allows to get product from dir using file name
func (p *ProductStore) RetrieveProduct(fileName string, userDir string) (*entity.Product, error) {
	f, err := os.Open(filepath.Join(FileStorage, userDir, fileName))
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
func (p *ProductStore) DeleteProduct(fileName string, userDir string) error {
	return os.Remove(filepath.Join(FileStorage, userDir, fileName) + JSONExtension)
}

// RetrieveProductsByUserID func which allows to get every user's product
func (p *ProductStore) RetrieveProductsByUserID(userDir string) ([]entity.Product, error) {
	dir, err := os.Open(filepath.Join(FileStorage, userDir))
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
		product, err := p.RetrieveProduct(file.Name(), userDir)
		if err != nil {
			return nil, err
		}
		products[i] = *product
	}

	return products, nil
}
