package repository

import (
	"errors"

	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func buildInventoryBaseQuery(filters model.InventoryFilter) *gorm.DB {
	q := initializers.DB.Model(&model.Inventory{})

	if filters.State != "" {
		q = q.Where("state = ?", filters.State)
	}
	if filters.Title != "" {
		q = q.Where("item ILIKE ?", "%"+filters.Title+"%")
	}
	return q
}

func AddItemToInventory(tx *gorm.DB, inventory model.Inventory) (bool, error) {
	db := initializers.DB
	if tx != nil {
		db = tx
	}
	if err := db.Create(&inventory).Error; err != nil {
		return false, err
	}
	return true, nil
}

func FindPaginatedInventory(offset, limit int, sortBy, order string, filters model.InventoryFilter) ([]model.Inventory, error) {
	var items []model.Inventory

	q := buildInventoryBaseQuery(filters).
		Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{Name: sortBy},
					Desc:   order == "DESC",
				},
			},
		}).
		Offset(offset).
		Limit(limit)

	err := q.Find(&items).Error
	return items, err
}

// CountInventory returns total number of matching rows
func CountInventory(filters model.InventoryFilter) (int64, error) {
	var count int64
	err := buildInventoryBaseQuery(filters).Count(&count).Error
	return count, err
}

// SumInventorySubTotal returns the sum of sub_total for the filtered rows
func SumInventorySubTotal(filters model.InventoryFilter) (int64, error) {
	var total int64

	// same technique as your original code
	sub := buildInventoryBaseQuery(filters).Select("inventories.id")

	err := initializers.DB.
		Table("(?) as filtered_inventories", sub).
		Select("COALESCE(SUM(sub_total), 0)").
		Joins("JOIN inventories ON inventories.id = filtered_inventories.id").
		Scan(&total).Error

	return total, err
}

func Begin() *gorm.DB {
	return initializers.DB.Begin()
}

func FindByInvId(tx *gorm.DB, id uint) (*model.Inventory, error) {
	db := initializers.DB
	if tx != nil {
		db = tx
	}

	var inv model.Inventory
	if err := db.First(&inv, id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func UpdateInventory(tx *gorm.DB, id uint, updates any) error {
	db := initializers.DB
	if tx != nil {
		db = tx
	}
	switch updates.(type) {
	case dto.UpdateInventoryState, dto.InventoryUpdateDiff:
		return db.Model(&model.Inventory{}).Where("id = ?", id).Updates(updates).Error
	default:
		return errors.New("Invalid update type received")
	}
}
