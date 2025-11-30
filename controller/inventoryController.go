package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ---------- Frontend interfaces mirrored in Go --------------
type InventoryReqBody struct {
	BrandOrCompany  string `json:"brandOrCompany"`
	State           constants.InventoryState `json:"state"`
	Item            string `json:"item"`
	NumOfBoxes		int32  `json:"numberOfBoxes" binding:"required,gt=0"`
	NumOfCartons    int32  `json:"numberOfCartons" binding:"required,gt=0"`
	PricePerCarton  int32  `json:"pricePerCarton" binding:"required,gt=0"`
	SubTotal        int32  `json:"subtotal" binding:"required,gt=0"`
}

type UpdateInventoryRequest struct {
	ID              uint                         `json:"ID" binding:"required"`
	BrandOrCompany  string                       `json:"BrandOrCompany" binding:"required"`
	Item            string                       `json:"Item" binding:"required"`
	NumOfBoxes      int32                        `json:"NumOfBoxes"`
	NumOfCartons    int32                        `json:"NumOfCartons" binding:"required"`
	PricePerCarton  int32                        `json:"PricePerCarton" binding:"required"`
	State           constants.InventoryState     `json:"State" binding:"required"`
	// SubTotal is NOT accepted from frontend — we calculate it
}

// ---------- Gin Handler -------------
func AddItemToInventoryHandler(c *gin.Context) {

	// 1. Obtain the request body
	var reqBody InventoryReqBody
	if err := c.ShouldBindBodyWithJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to obtain the request body"})
		return
	}

	// 2. Validate the field State
	if err := reqBody.State.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Wrong state"})
		return
	}

	// 3. Save the data in the DB
	inventory := model.Inventory{
		BrandOrCompany: reqBody.BrandOrCompany,
		State:          reqBody.State,        
		Item:           reqBody.Item,            
		NumOfBoxes:     reqBody.NumOfBoxes,		
		NumOfCartons:   reqBody.NumOfCartons, 
		PricePerCarton: reqBody.PricePerCarton, 
		SubTotal:       reqBody.SubTotal,
	}
	if err := initializers.DB.Create(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the data in the DB"})
		return
	}

	// 4. Send success message
	c.JSON(http.StatusOK, gin.H{"Success": "Item saved to inventory"})
}

func GetPaginatedInventoryItems(c *gin.Context) {
	// Default values
	pageNumber := 1
	pageSize := 5 // increased default for better UX
	if p := c.DefaultQuery("pageNumber", "1"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			pageNumber = n
		}
	}
	if s := c.DefaultQuery("pageSize", "10"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 && n <= 100 { // max limit
			pageSize = n
		}
	}

	offset := (pageNumber - 1) * pageSize

	// Whitelisted sorting
	sortBy := c.DefaultQuery("sortBy", "id")
	order := strings.ToUpper(c.DefaultQuery("order", "DESC"))

	// Whitelist allowed columns to prevent SQL injection
	allowedSortColumns := map[string]bool{
		"id":          true,
		"item":        true,
		"sub_total":   true,
		"created_at":  true,
		"updated_at":  true,
		"state":       true,
	}
	if !allowedSortColumns[sortBy] {
		sortBy = "id" // fallback
	}

	// Validate order
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	// Search filters
	searchTitle := strings.TrimSpace(c.Query("title"))
	searchState := strings.TrimSpace(c.Query("state"))
	
	// Base query
	query := initializers.DB.Model(&model.Inventory{})

	// Validate state if provided & Apply filters
	if searchState != "" {
		state := constants.InventoryState(searchState)
		if err := state.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state value"})
			return
		}
		query = query.Where("state = ?", state)
	}

	if searchTitle != "" {
		query = query.Where("item LIKE ?", "%"+searchTitle+"%")
	}

	// Clone query for counting and summing to avoid side effects
	countQuery := query.Session(&gorm.Session{})
	listQuery := query.Session(&gorm.Session{})

	// Apply sorting safely
	listQuery = listQuery.Clauses(clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{
				Column: clause.Column{Name: sortBy},
				Desc:   order == "DESC",
			},
		},
	})

	// Pagination
	var inventories []model.Inventory
	if err := listQuery.Offset(offset).Limit(pageSize).Find(&inventories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory items"})
		return
	}

	// Total count (filtered)
	var totalElements int64
	if err := countQuery.Count(&totalElements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count items"})
		return
	}

	// Total subtotal (sum of sub_total for filtered results)
	var totalSubtotal int64
	sumQuery := query.Session(&gorm.Session{}).Select("inventories.id")
	if err := initializers.DB.
		Table("(?) as filtered_inventories", sumQuery).
		Select("COALESCE(SUM(sub_total), 0)").
		Joins("JOIN inventories ON inventories.id = filtered_inventories.id").
		Scan(&totalSubtotal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"inventoryItems": inventories,
		"total": totalSubtotal,
		"totalElements": totalElements,
	})
}

func UpdateInventory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateInventoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
			return
		}

		// Start transaction — very important for data consistency
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Fetch existing inventory record by ID
		var existing model.Inventory
		if err := tx.First(&existing, req.ID).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Inventory record not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		// Validate that the state transition is allowed
		if existing.State != constants.InventoryStateOrdered && existing.State != constants.InventoryStateReceived {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can only update records in Ordered or Received state"})
			return
		}

		// Calculate SubTotal for the incoming data (always!)
		reqSubTotal := req.NumOfCartons * req.PricePerCarton

		// Determine next state for both existing and new record
		nextState := existing.State
		switch existing.State {
			case constants.InventoryStateOrdered:
				nextState = constants.InventoryStateReceived
			case constants.InventoryStateReceived:
				nextState = constants.InventoryStateUnpacked
		}

		if existing.NumOfCartons == req.NumOfCartons {
			// Full match → just advance state of existing record
			updates := model.Inventory{
				State:    nextState,
				SubTotal: existing.NumOfCartons * existing.PricePerCarton, // should already be correct
			}

			if err := tx.Model(&existing).Updates(updates).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing record"})
				return
			}

		} else {
			// Partial receipt → adjust existing row
			difference := existing.NumOfCartons - req.NumOfCartons
			newExistingSubTotal := difference * existing.PricePerCarton

			updateExisting := model.Inventory{
				NumOfCartons: difference,
				SubTotal:     newExistingSubTotal,
				State:        nextState,
			}

			if err := tx.Model(&existing).Updates(updateExisting).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update existing cartons"})
				return
			}

			// Create new row for the received portion
			newRecord := model.Inventory{
				BrandOrCompany: req.BrandOrCompany,
				Item:           req.Item,
				NumOfBoxes:     req.NumOfBoxes,
				NumOfCartons:   req.NumOfCartons,
				PricePerCarton: req.PricePerCarton,
				SubTotal:       reqSubTotal,
				State:          nextState,
			}

			if err := tx.Create(&newRecord).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new inventory record"})
				return
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Inventory updated successfully",
			"updated_existing": existing.ID != 0,
		})
	}
}
