package entity

type Product struct {
	Barcode      string `json:"barcode"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Cost         int    `json:"cost"`
	UserID       int    `json:"user_id"`
	FileLocation string `json:"file_location"`
	CreatedAt    string `json:"created_at"`
}
