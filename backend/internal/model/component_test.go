package model

import "testing"

func intPtr(v int) *int             { return &v }
func strPtr(v string) *string       { return &v }
func catPtr(v Category) *Category   { return &v }
func condPtr(v Condition) *Condition { return &v }

func TestCreateComponentInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateComponentInput
		wantErr bool
	}{
		{
			name: "valid component",
			input: CreateComponentInput{
				ItemID:         "CMP-001",
				Name:           "STM32F401 Blackpill",
				Category:       CategoryMicrocontroller,
				StockQty:       intPtr(15),
				MinStock:       intPtr(5),
				DrawerLocation: "RAK-A2",
				Condition:      ConditionGood,
			},
			wantErr: false,
		},
		{
			name: "missing item_id",
			input: CreateComponentInput{
				Name:           "DHT22",
				Category:       CategorySensor,
				StockQty:       intPtr(10),
				MinStock:       intPtr(2),
				DrawerLocation: "RAK-B1",
				Condition:      ConditionNew,
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			input: CreateComponentInput{
				ItemID:         "CMP-002",
				Name:           "DHT22",
				Category:       Category("ACTUATOR"),
				StockQty:       intPtr(10),
				MinStock:       intPtr(2),
				DrawerLocation: "RAK-B1",
				Condition:      ConditionNew,
			},
			wantErr: true,
		},
		{
			name: "negative stock_qty",
			input: CreateComponentInput{
				ItemID:         "CMP-003",
				Name:           "Jumper Wire",
				Category:       CategoryCable,
				StockQty:       intPtr(-1),
				MinStock:       intPtr(5),
				DrawerLocation: "RAK-C1",
				Condition:      ConditionGood,
			},
			wantErr: true,
		},
		{
			name: "invalid condition",
			input: CreateComponentInput{
				ItemID:         "CMP-004",
				Name:           "Resistor 10k",
				Category:       CategoryPassive,
				StockQty:       intPtr(50),
				MinStock:       intPtr(10),
				DrawerLocation: "RAK-D1",
				Condition:      Condition("BROKEN"),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestUpdateComponentInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   UpdateComponentInput
		wantErr bool
	}{
		{
			name: "valid partial update",
			input: UpdateComponentInput{
				StockQty:       intPtr(8),
				DrawerLocation: strPtr("RAK-A3"),
				Condition:      condPtr(ConditionGood),
			},
			wantErr: false,
		},
		{
			name:    "empty update payload",
			input:   UpdateComponentInput{},
			wantErr: true,
		},
		{
			name: "negative stock_qty",
			input: UpdateComponentInput{
				StockQty: intPtr(-5),
			},
			wantErr: true,
		},
		{
			name: "empty drawer_location",
			input: UpdateComponentInput{
				DrawerLocation: strPtr("   "),
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			input: UpdateComponentInput{
				Category: catPtr(Category("INVALID")),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestUploadURLInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   UploadURLInput
		wantErr bool
	}{
		{
			name:    "valid jpeg",
			input:   UploadURLInput{Filename: "stm32.jpg", ContentType: "image/jpeg"},
			wantErr: false,
		},
		{
			name:    "valid png uppercase",
			input:   UploadURLInput{Filename: "dht22.png", ContentType: "IMAGE/PNG"},
			wantErr: false,
		},
		{
			name:    "missing filename",
			input:   UploadURLInput{Filename: "", ContentType: "image/jpeg"},
			wantErr: true,
		},
		{
			name:    "unsupported content type",
			input:   UploadURLInput{Filename: "doc.pdf", ContentType: "application/pdf"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}
