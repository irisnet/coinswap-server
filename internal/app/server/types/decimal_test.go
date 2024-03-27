package types

import (
	"fmt"
	"testing"
)

func TestDecimal_Float64(t *testing.T) {
	type fields struct {
		dec Decimal
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name:   "case1",
			fields: fields{NewDecFromFloat64(0.6666)},
			want:   0.6666,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.dec.Float64(); got != tt.want {
				t.Errorf("Decimal.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClone(t *testing.T) {
	d1 := NewDecFromUint64(100)
	d2 := d1.AddUint64(100)
	d3 := d1.Add(d2).Mul(d1)

	fmt.Println(d1)
	fmt.Println(d2)
	fmt.Println(d3)
}
