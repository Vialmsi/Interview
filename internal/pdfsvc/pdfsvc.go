package pdfsvc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Vialmsi/Interview/internal/entity"

	"github.com/signintech/gopdf"
	"github.com/sirupsen/logrus"
)

const (
	pdfStorage       = "pdf"
	dirPermission    = 0777
	fileNameTemplate = "doc_%s_%s.pdf"

	fontFile   = "assets/Font.ttf"
	fontFamily = "Metroplex Shadow"
	fontSize   = 20

	pdfTemplate = "assets/Template.pdf"

	barcodeX = 40
	barcodeY = 80

	titleX = 40
	titleY = 185

	costX = 325
	costY = 295
)

type Store interface {
	UpdateFileLocation(ctx context.Context, fileName, barcode string, userID int) error
}

type PDFService struct {
	logger *logrus.Logger
	store  Store
}

func NewPDFService(logger *logrus.Logger, store Store) (*PDFService, error) {
	_, err := os.Open(pdfStorage)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(pdfStorage, dirPermission)
		if err != nil {
			return nil, fmt.Errorf("error while creating pdf storage: %s", err.Error())
		}
	}
	return &PDFService{
		logger: logger,
		store:  store,
	}, nil
}

func (p *PDFService) LoadPDFFromBarcode(userID int, barcode string) (string, error) {
	p.logger.Info("[LoadPDFFromBarcode] started")

	strID := strconv.Itoa(userID)

	err := p.checkUserFolder(strID)
	if err != nil {
		p.logger.Errorf("[LoadPDFFromBarcode] Error while checking user folder: %s", err)
		return "", fmt.Errorf("error, while generating pdf, try later")
	}

	dir, err := os.Open(filepath.Join(pdfStorage, strID))
	if err != nil {
		p.logger.Errorf("[LoadPDFFromBarcode] Error while open pdf storage: %s", err.Error())
		return "", err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		p.logger.Errorf("[LoadPDFFromBarcode] Error while read dir: %s", err.Error())
		return "", err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), barcode) {
			p.logger.Info("[LoadPDFFromBarcode] ended")
			return filepath.Join(pdfStorage, strID, file.Name()), nil
		}
	}

	p.logger.Info("[LoadPDFFromBarcode] ended")

	return "", errors.New("couldn't find file with this name")
}

func (p *PDFService) LoadPDFFromName(fileName string) (string, error) {
	p.logger.Info("[LoadPDFFromName] started")

	_, err := os.Stat(filepath.Join(fileName))
	if err != nil {
		p.logger.Errorf("[LoadPDFFromName] Error while check file: %s", err.Error())
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not exist")
		}
		return "", err
	}

	p.logger.Info("[LoadPDFFromName] ended")

	return fileName, nil
}

func (p *PDFService) GeneratePDF(ctx context.Context, userID int, product entity.Product) (string, error) {
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 420, H: 395}})

	err := pdf.AddTTFFont(fontFamily, fontFile)
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while load font: %s", err.Error())
		return "", fmt.Errorf("error while generating pdf, try later")
	}

	err = pdf.SetFont(fontFamily, "", fontSize)
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while setting font: %s", err.Error())
		return "", fmt.Errorf("error while generating pdf, try later")
	}

	tpl1 := pdf.ImportPage(pdfTemplate, 1, "/MediaBox")

	pdf.AddPage()

	pdf.UseImportedTemplate(tpl1, 0, 0, 420, 395)

	//Product barcode
	pdf.SetXY(barcodeX, barcodeY)
	err = pdf.Cell(nil, product.Barcode)
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while writing barcode")
		return "", fmt.Errorf("error while generating pdf, try later")
	}

	//Product name
	pdf.SetXY(titleX, titleY)
	err = pdf.Cell(nil, product.Name)
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while writing name")
		return "", fmt.Errorf("error while generating pdf, try later")
	}

	//Product cost
	pdf.SetXY(costX, costY)
	err = pdf.Cell(nil, strconv.Itoa(product.Cost))
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while writing cost")
		return "", fmt.Errorf("error while generating pdf, try later")
	}

	strID := strconv.Itoa(userID)

	fileName := filepath.Join(pdfStorage, strID, fmt.Sprintf(fileNameTemplate, product.Barcode, time.Now().Format(`02-01-2006_15:01:05`)))

	err = pdf.WritePdf(fileName)
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while saving pdf: %s", err.Error())
		return "", fmt.Errorf("error while generating pdf, try later")
	}

	err = p.store.UpdateFileLocation(ctx, fileName, product.Barcode, userID)
	if err != nil {
		p.logger.Errorf("[GeneratePDF] Error while updating file location in db: %s", err.Error())
		return "", fmt.Errorf("error,while generating pdf, try later")
	}

	return fileName, nil
}

func (p *PDFService) checkUserFolder(userID string) error {
	_, err := os.Open(filepath.Join(pdfStorage, userID))
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath.Join(pdfStorage, userID), dirPermission)
		if err != nil {
			return fmt.Errorf("error while creating pdf storage: %s", err.Error())
		}
	}
	return nil
}
