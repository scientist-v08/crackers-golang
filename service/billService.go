package service

import (
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/repository"
	"github.com/scientist-v08/crackers/utils"
)

func PreviewBillService(req model.BillDetails) ([]byte, error) {
	// Generate PDF
	pdfBytes, errPdf := utils.GeneratePDF(req, 0)
	if errPdf != nil {
		return nil, errPdf
	}
	return pdfBytes, nil
}

func GenerateBillService(req model.BillDetails) error {
	tx := repository.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r:= recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := repository.AddNewBill(tx, &req); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}