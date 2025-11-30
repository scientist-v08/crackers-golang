package constants

import "fmt"

type InventoryState string

const (
	InventoryStateOrdered  InventoryState = "Ordered"
	InventoryStateReceived InventoryState = "Received"
	InventoryStateUnpacked InventoryState = "Unpacked"
)

func (s InventoryState) Validate() error {
	switch s {
	case InventoryStateOrdered, InventoryStateReceived, InventoryStateUnpacked:
		return nil
	default:
		return fmt.Errorf("invalid state: %s", s)
	}
}