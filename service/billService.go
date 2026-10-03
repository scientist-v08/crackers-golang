package service

import (
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/repository"
	"github.com/scientist-v08/crackers/utils"
)

func PreviewBillService(req model.BillDetails, id uint64) ([]byte, error) {
	// Generate PDF
	pdfBytes, errPdf := utils.GeneratePDF(req, id)
	if errPdf != nil {
		return nil, errPdf
	}
	return pdfBytes, nil
}

func GenerateBillService(req model.BillDetails) (uint64, error) {
	tx := repository.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	defer func() {
		if r:= recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	_id := uint64(0)
	var addBillErr error
	if _id, addBillErr = repository.AddNewBill(tx, &req); addBillErr != nil {
		tx.Rollback()
		return 0, addBillErr
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return _id, nil
}