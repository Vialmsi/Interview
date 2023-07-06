package store

import (
	"context"
	"fmt"

	"github.com/Vialmsi/Interview/internal/entity"
	customErrors "github.com/Vialmsi/Interview/internal/errors"
)

// SaveProduct func which allows to save product into dir
func (s *Store) SaveProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	query := `INSERT INTO products(barcode,name,description,cost,user_id) VALUES($1,$2,$3,$4,$5)
				RETURNING barcode,name,description,cost,user_id,file_location,created_at`

	err := s.db.QueryRowContext(ctx, query,
		product.Barcode,
		product.Name,
		product.Desc,
		product.Cost,
		product.UserID).Scan(
		&product.Barcode,
		&product.Name,
		&product.Desc,
		&product.Cost,
		&product.UserID,
		&product.FileLocation,
		&product.CreatedAt)

	if isUniqueViolation(err) {
		return nil, customErrors.ProductAlreadyExistError
	}

	return product, err
}

// RetrieveProduct func which allows to get product from dir using file name
func (s *Store) RetrieveProduct(ctx context.Context, productID string, userID int) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	query := `SELECT barcode,name,description,cost,user_id 
				FROM products 
				WHERE barcode=$1 AND user_id=$2
				ORDER BY products.created_at`

	var product entity.Product

	err := s.db.QueryRowContext(ctx, query, productID, userID).Scan(
		&product.Barcode,
		&product.Name,
		&product.Desc,
		&product.Cost,
		&product.UserID)

	return &product, err
}

// DeleteProduct func which allows to delete product from dir
func (s *Store) DeleteProduct(ctx context.Context, productID string, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	query := `DELETE FROM products WHERE barcode=$1 AND user_id=$2`

	res, err := s.db.ExecContext(ctx, query, productID, userID)

	affectedRows, err := res.RowsAffected()
	if affectedRows == 0 {
		return customErrors.NoProductToDeleteError
	}

	return err
}

// RetrieveProductsByUserID func which allows to get every user's product
func (s *Store) RetrieveProductsByUserID(ctx context.Context, userID, limit, page int) ([]entity.Product, int, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	if page <= 0 {
		page = s.defaultPage
	}
	if limit <= 0 {
		limit = s.defaultLimit
	}

	page = (page - 1) * limit

	fmt.Println(page, limit)

	query := `SELECT barcode, name, description, cost, user_id, file_location, created_at
				FROM products 
				WHERE user_id = $1
				ORDER BY created_at DESC
				LIMIT $2 OFFSET $3;`

	rows, err := s.db.QueryContext(ctx, query, userID, limit, page)
	if err != nil {
		return nil, 0, err
	}

	var products []entity.Product
	var count int

	for rows.Next() {
		product := entity.Product{}

		err := rows.Scan(
			&product.Barcode,
			&product.Name,
			&product.Desc,
			&product.Cost,
			&product.UserID,
			&product.FileLocation,
			&product.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, product)
	}

	query = `SELECT COUNT(*) FROM products WHERE user_id=$1`

	err = s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return products, count, err
}

func (s *Store) UpdateFileLocation(ctx context.Context, fileName, barcode string, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	query := `UPDATE products SET file_location=$1 WHERE barcode=$2 AND user_id=$3`

	_, err := s.db.ExecContext(ctx, query, fileName, barcode, userID)

	return err
}
