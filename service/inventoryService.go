package service

import (
	"errors"
	"strings"
	"sync"

	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/repository"
	"gorm.io/gorm"
)

var (
	ErrInventoryNotFound      = errors.New("inventory record not found")
	ErrInvalidStateTransition = errors.New("can only update records in Ordered or Received state")
)

func AddItemToinventory(reqBody model.InventoryReqBody) (bool, error) {
	// 1. Create the required body to save to the DB
	inventory := model.Inventory{
		BrandOrCompany: reqBody.BrandOrCompany,
		State:          reqBody.State,        
		Item:           reqBody.Item,            
		NumOfBoxes:     reqBody.NumOfBoxes,		
		NumOfCartons:   reqBody.NumOfCartons, 
		PricePerCarton: reqBody.PricePerCarton, 
		SubTotal:       reqBody.PricePerCarton * reqBody.NumOfCartons,
	}
	// 2. Save to the DB
	var (
		isSaved bool
		isError error
	)
	isSaved, isError = repository.AddItemToInventory(nil, inventory)

	// 3. Return the result
	return isSaved, isError
}

func GetPaginatedInventoryItems(req dto.InventoryListRequest) (*dto.InventoryListResponse, error) {
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	allowedSortColumns := map[string]bool{
		"id":          true,
		"item":        true,
		"sub_total":   true,
		"created_at":  true,
		"updated_at":  true,
		"state":       true,
	}
	if !allowedSortColumns[req.SortBy] {
		req.SortBy = "id"
	}

	req.Order = strings.ToUpper(req.Order)
	if req.Order != "ASC" && req.Order != "DESC" {
		req.Order = "DESC"
	}

	// Validate state if present
	filter := model.InventoryFilter{
		Title: strings.TrimSpace(req.Title),
	}
	if req.State != "" {
		state := constants.InventoryState(req.State)
		if err := state.Validate(); err != nil {
			return nil, err
		}
		filter.State = state
	}

	offset := (req.PageNumber - 1) * req.PageSize

	// ---- Parallel DB calls (same as original) ----
	var (
		wg            sync.WaitGroup
		items         []model.Inventory
		totalElements int64
		totalSubtotal int64
		errItems      error
		errCount      error
		errSum        error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		items, errItems = repository.FindPaginatedInventory(offset, req.PageSize, req.SortBy, req.Order, filter)
	}()

	go func() {
		defer wg.Done()
		totalElements, errCount = repository.CountInventory(filter)
	}()

	go func() {
		defer wg.Done()
		totalSubtotal, errSum = repository.SumInventorySubTotal(filter)
	}()

	wg.Wait()

	if errItems != nil {
		return nil, errItems
	}
	if errCount != nil {
		return nil, errCount
	}
	if errSum != nil {
		return nil, errSum
	}

	return &dto.InventoryListResponse{
		InventoryItems: items,
		Total:          totalSubtotal,
		TotalElements:  totalElements,
	}, nil
}

func UpdateInventory(req model.UpdateInventoryRequest) error {
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

	existing, existErr := repository.FindByInvId(tx, req.ID)
	if existErr != nil {
		tx.Rollback()
		if errors.Is(existErr, gorm.ErrRecordNotFound) {
			return ErrInventoryNotFound
		}
		return existErr
	}

	if existing.State != constants.InventoryStateOrdered && existing.State != constants.InventoryStateReceived {
		tx.Rollback()
		return ErrInvalidStateTransition
	}

	reqSubTotal := req.NumOfCartons * existing.PricePerCarton
	nextState := existing.State
	switch existing.State {
	case constants.InventoryStateOrdered:
		nextState = constants.InventoryStateReceived
	case constants.InventoryStateReceived:
		nextState = constants.InventoryStateUnpacked
	}

	if existing.NumOfCartons == req.NumOfCartons {
		// Full match → just advance state
		updates := dto.UpdateInventoryState{
			State:     nextState,
			SubTotal: existing.NumOfCartons * existing.PricePerCarton,
		}
		if err := repository.UpdateInventory(tx, existing.ID, updates); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Partial receipt
		difference := existing.NumOfCartons - req.NumOfCartons
		newExistingSubTotal := difference * existing.PricePerCarton

		updateExisting := dto.InventoryUpdateDiff{
			NumOfCartons: difference,
			SubTotal:      newExistingSubTotal,
		}
		if err := repository.UpdateInventory(tx, existing.ID, updateExisting); err != nil {
			tx.Rollback()
			return err
		}

		newRecord := model.Inventory{
			BrandOrCompany: req.BrandOrCompany,
			Item:           req.Item,
			NumOfBoxes:     req.NumOfBoxes,
			NumOfCartons:   req.NumOfCartons,
			PricePerCarton: req.PricePerCarton,
			SubTotal:       reqSubTotal,
			State:          nextState,
		}
		if _, err := repository.AddItemToInventory(tx, newRecord); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}