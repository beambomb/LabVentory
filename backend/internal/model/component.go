package model

import (
	"strings"

	"labventroy/backend/internal/apperror"
)

// Category classifies laboratory electronic components.
type Category string

const (
	CategoryMicrocontroller Category = "MICROCONTROLLER"
	CategorySensor          Category = "SENSOR"
	CategoryCable           Category = "CABLE"
	CategoryPassive         Category = "PASSIVE"
)

// IsValid reports whether the category is one of the allowed constants.
func (c Category) IsValid() bool {
	switch c {
	case CategoryMicrocontroller, CategorySensor, CategoryCable, CategoryPassive:
		return true
	default:
		return false
	}
}

// Condition represents the physical condition of a component.
type Condition string

const (
	ConditionNew       Condition = "NEW"
	ConditionGood      Condition = "GOOD"
	ConditionDefective Condition = "DEFECTIVE"
)

// IsValid reports whether the condition is one of the allowed constants.
func (c Condition) IsValid() bool {
	switch c {
	case ConditionNew, ConditionGood, ConditionDefective:
		return true
	default:
		return false
	}
}

// Component represents a single item in the LabComponentInventory DynamoDB table.
type Component struct {
	ItemID         string    `json:"item_id" dynamodbav:"item_id"`
	Name           string    `json:"name" dynamodbav:"name"`
	Category       Category  `json:"category" dynamodbav:"category"`
	StockQty       int       `json:"stock_qty" dynamodbav:"stock_qty"`
	MinStock       int       `json:"min_stock" dynamodbav:"min_stock"`
	DrawerLocation string    `json:"drawer_location" dynamodbav:"drawer_location"`
	Condition      Condition `json:"condition" dynamodbav:"condition"`
	UpdatedAt      string    `json:"updated_at" dynamodbav:"updated_at"`
}

// CreateComponentInput is the request payload for POST /items.
type CreateComponentInput struct {
	ItemID         string    `json:"item_id"`
	Name           string    `json:"name"`
	Category       Category  `json:"category"`
	StockQty       *int      `json:"stock_qty"`
	MinStock       *int      `json:"min_stock"`
	DrawerLocation string    `json:"drawer_location"`
	Condition      Condition `json:"condition"`
}

// Validate checks required fields, non-negative quantities, and valid enums.
func (in *CreateComponentInput) Validate() error {
	in.ItemID = strings.TrimSpace(in.ItemID)
	in.Name = strings.TrimSpace(in.Name)
	in.DrawerLocation = strings.TrimSpace(in.DrawerLocation)

	if in.ItemID == "" {
		return apperror.Validation("item_id is required")
	}
	if in.Name == "" {
		return apperror.Validation("name is required")
	}
	if !in.Category.IsValid() {
		return apperror.Validation("category must be one of MICROCONTROLLER, SENSOR, CABLE, PASSIVE")
	}
	if in.StockQty == nil || *in.StockQty < 0 {
		return apperror.Validation("stock_qty is required and must be >= 0")
	}
	if in.MinStock == nil || *in.MinStock < 0 {
		return apperror.Validation("min_stock is required and must be >= 0")
	}
	if in.DrawerLocation == "" {
		return apperror.Validation("drawer_location is required")
	}
	if !in.Condition.IsValid() {
		return apperror.Validation("condition must be one of NEW, GOOD, DEFECTIVE")
	}
	return nil
}

// ToComponent converts a validated CreateComponentInput into a Component entity.
func (in CreateComponentInput) ToComponent(updatedAt string) Component {
	return Component{
		ItemID:         in.ItemID,
		Name:           in.Name,
		Category:       in.Category,
		StockQty:       *in.StockQty,
		MinStock:       *in.MinStock,
		DrawerLocation: in.DrawerLocation,
		Condition:      in.Condition,
		UpdatedAt:      updatedAt,
	}
}

// UpdateComponentInput is the request payload for PUT /items/{id}.
type UpdateComponentInput struct {
	Name           *string    `json:"name,omitempty"`
	Category       *Category  `json:"category,omitempty"`
	StockQty       *int       `json:"stock_qty,omitempty"`
	MinStock       *int       `json:"min_stock,omitempty"`
	DrawerLocation *string    `json:"drawer_location,omitempty"`
	Condition      *Condition `json:"condition,omitempty"`
}

// Validate ensures at least one field is provided and all supplied fields are valid.
func (in *UpdateComponentInput) Validate() error {
	if in.Name == nil &&
		in.Category == nil &&
		in.StockQty == nil &&
		in.MinStock == nil &&
		in.DrawerLocation == nil &&
		in.Condition == nil {
		return apperror.Validation("at least one updatable field is required")
	}

	if in.Name != nil {
		trimmed := strings.TrimSpace(*in.Name)
		if trimmed == "" {
			return apperror.Validation("name cannot be empty")
		}
		in.Name = &trimmed
	}
	if in.Category != nil && !in.Category.IsValid() {
		return apperror.Validation("category must be one of MICROCONTROLLER, SENSOR, CABLE, PASSIVE")
	}
	if in.StockQty != nil && *in.StockQty < 0 {
		return apperror.Validation("stock_qty must be >= 0")
	}
	if in.MinStock != nil && *in.MinStock < 0 {
		return apperror.Validation("min_stock must be >= 0")
	}
	if in.DrawerLocation != nil {
		trimmed := strings.TrimSpace(*in.DrawerLocation)
		if trimmed == "" {
			return apperror.Validation("drawer_location cannot be empty")
		}
		in.DrawerLocation = &trimmed
	}
	if in.Condition != nil && !in.Condition.IsValid() {
		return apperror.Validation("condition must be one of NEW, GOOD, DEFECTIVE")
	}
	return nil
}
