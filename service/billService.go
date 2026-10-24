package service

import (
	"context"

	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/repository"
	"github.com/scientist-v08/crackers/utils"
)

func PreviewBillService(req dto.BillDetails, id int64) ([]byte, error) {
	// Generate PDF
	pdfBytes, errPdf := utils.GeneratePDF(req, uint64(id))
	if errPdf != nil {
		return nil, errPdf
	}
	return pdfBytes, nil
}

func GenerateBillService(req dto.BillDetails) (int64, error) {
	ctx := context.Background()

	tx, err := repository.Begin(ctx)
	if err != nil {
		return 0, err
	}

	// Rollback is safe even after a successful Commit
	defer tx.Rollback(ctx)

	id, err := repository.AddNewBill(ctx, tx, &req)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}