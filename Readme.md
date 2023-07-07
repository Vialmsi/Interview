# PDF_Server

This server may store users and products.

Also it may generate pdf with price using pdf template.

This service support next routes:

### [POST] /api/auth/signup - user registration

### [POST] /api/auth/signin - user login

### [GET] /api/products/:barcode - get product by barcode(id)

### [DELETE] /api/products/:barcode - delete product by barcode(id)

### [POST] /api/products - save product

### [GET] /api/products - get all user products

this method allow pagination and you can use limit and page like this:

#### localhost:8081/api/products?limit=5&page=1

### [GET] /api/prices/:barcode - get product by barcode (if product not exist this route will generate it)

### [GET] /api/prices - get product by filename

#### Example: localhost:8081/api/prices?path=pdf/10/doc_foo_07-06-2023_03:35:00.pdf
You can find location of files in products info after generating price for it

## For run
### docker-compose up --build