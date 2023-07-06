package errors

import "errors"

var (
	UserAlreadyExistError         = errors.New("user with this login already exist")
	ProductAlreadyExistError      = errors.New("product with this barcode already exist")
	NoProductToDeleteError        = errors.New("no such product to delete")
	IncorrectLoginOrPasswordError = errors.New("incorrect login or password")
	NoSuchProductError            = errors.New("no such product")
	FileNotExistError             = errors.New("file not exist")
	PDFNotExistError              = errors.New("couldn't find file with this name")
	TokenExpiredError             = errors.New("token is expired, pls login again")
)
