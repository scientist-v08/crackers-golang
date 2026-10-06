package service

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/repository"
)

var (
	ErrInventoryNotFound      = errors.New("inventory record not found")
	ErrInvalidStateTransition = errors.New("can only update records in Ordered or Received state")
)

func AddItemToinventory(reqBody dto.InventoryReqBody) (bool, error) {
	ctx := context.Background()

	inventory := dto.Inventory{
		BrandOrCompany: reqBody.BrandOrCompany,
		State:          reqBody.State,
		Item:           reqBody.Item,
		NumOfBoxes:     reqBody.NumOfBoxes,
		NumOfCartons:   reqBody.NumOfCartons,
		PricePerCarton: reqBody.PricePerCarton,
		SubTotal:       reqBody.PricePerCarton * reqBody.NumOfCartons,
	}

	_, err := repository.AddItemToInventory(ctx, nil, inventory)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetPaginatedInventoryItems(req dto.InventoryListRequest) (*dto.InventoryListResponse, error) {
	ctx := context.Background()

	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// Note: sorting is now hardcoded to id DESC in the query.
	// We keep the validation only for future flexibility.
	filter := dto.InventoryFilter{
		Title: strings.TrimSpace(req.Title),
	}
	if req.State != "" {
		state := constants.InventoryState(req.State)
		if err := state.Validate(); err != nil {
			return nil, err
		}
		filter.State = state
	}

	offset := int32((req.PageNumber - 1) * req.PageSize)
	limit := int32(req.PageSize)

	var (
		wg            sync.WaitGroup
		items         []db.Inventory
		totalElements int64
		totalSubtotal int64
		errItems      error
		errCount      error
		errSum        error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		items, errItems = repository.FindPaginatedInventory(ctx, offset, limit, filter)
	}()

	go func() {
		defer wg.Done()
		totalElements, errCount = repository.CountInventory(ctx, filter)
	}()

	go func() {
		defer wg.Done()
		totalSubtotal, errSum = repository.SumInventorySubTotal(ctx, filter)
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

	// Convert db.Inventory → model.Inventory (so existing DTO still works)
	modelItems := make([]dto.Inventory, 0, len(items))
	for _, item := range items {
		modelItems = append(modelItems, dto.Inventory{
			ID:             uint64(item.ID),
			BrandOrCompany: item.BrandOrCompany.String,
			State:          constants.InventoryState(item.State.String),
			Item:           item.Item.String,
			NumOfBoxes:     item.NumOfBoxes.Int32,
			NumOfCartons:   item.NumOfCartons.Int32,
			PricePerCarton: item.PricePerCarton.Int32,
			SubTotal:       item.SubTotal.Int32,
		})
	}

	return &dto.InventoryListResponse{
		InventoryItems: modelItems,
		Total:          totalSubtotal,
		TotalElements:  totalElements,
	}, nil
}

func UpdateInventory(req dto.UpdateInventoryRequest) error {
	ctx := context.Background()

	tx, err := repository.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // safe even after Commit

	existing, err := repository.FindByInvId(ctx, tx, int64(req.ID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInventoryNotFound
		}
		return err
	}

	currentState := constants.InventoryState(existing.State.String)

	if currentState != constants.InventoryStateOrdered && currentState != constants.InventoryStateReceived {
		return ErrInvalidStateTransition
	}

	reqSubTotal := req.NumOfCartons * existing.PricePerCarton.Int32
	nextState := currentState

	switch currentState {
	case constants.InventoryStateOrdered:
		nextState = constants.InventoryStateReceived
	case constants.InventoryStateReceived:
		nextState = constants.InventoryStateUnpacked
	}

	if existing.NumOfCartons.Int32 == req.NumOfCartons {
		// Full match → just advance state
		updates := dto.UpdateInventoryState{
			State:    nextState,
			SubTotal: existing.NumOfCartons.Int32 * existing.PricePerCarton.Int32,
		}
		if err := repository.UpdateInventory(ctx, tx, existing.ID, updates); err != nil {
			return err
		}
	} else {
		// Partial receipt
		difference := existing.NumOfCartons.Int32 - req.NumOfCartons
		newExistingSubTotal := difference * existing.PricePerCarton.Int32

		updateExisting := dto.InventoryUpdateDiff{
			NumOfCartons: difference,
			SubTotal:     newExistingSubTotal,
		}
		if err := repository.UpdateInventory(ctx, tx, existing.ID, updateExisting); err != nil {
			return err
		}

		newRecord := dto.Inventory{
			BrandOrCompany: req.BrandOrCompany,
			Item:           req.Item,
			NumOfBoxes:     req.NumOfBoxes,
			NumOfCartons:   req.NumOfCartons,
			PricePerCarton: req.PricePerCarton,
			SubTotal:       reqSubTotal,
			State:          nextState,
		}
		if _, err := repository.AddItemToInventory(ctx, tx, newRecord); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}