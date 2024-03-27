package types

import (
	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/types"
)

var (
	ZeroDec = Decimal{types.Decimal{Big: new(decimal.Big).SetUint64(0)}}
)

type Decimal struct {
	dec types.Decimal
}

func Wrap(dec types.Decimal) Decimal {
	return Decimal{types.NewDecimal(dec.Big)}
}

func NewDecFromUint64(d1 uint64) Decimal {
	dec := types.Decimal{Big: new(decimal.Big).SetUint64(d1)}
	return Decimal{dec}
}

func NewDecFromFloat64(d1 float64) Decimal {
	dec := types.Decimal{Big: new(decimal.Big).SetFloat64(d1)}
	return Decimal{dec}
}

func (d Decimal) Unwrap() types.Decimal {
	return d.dec
}

func (d Decimal) Clone() Decimal {
	big := new(decimal.Big)
	big.Copy(d.dec.Big)
	return Decimal{types.Decimal{Big: big}}
}

func (d Decimal) Add(d1 Decimal) Decimal {
	src := d.Clone()
	dst := d1.Clone()
	src.dec.Add(src.dec.Big, dst.dec.Big)
	return src
}

func (d Decimal) AddUint64(d1 uint64) Decimal {
	copy := d.Clone()
	copy.dec.Add(copy.dec.Big, new(decimal.Big).SetUint64(d1))
	return copy
}

func (d Decimal) AddFloat64(d1 float64) Decimal {
	copy := d.Clone()
	copy.dec.Add(copy.dec.Big, new(decimal.Big).SetFloat64(d1))
	return copy
}

func (d Decimal) Sub(d1 Decimal) Decimal {
	copy := d.Clone()
	copy.dec.Sub(copy.dec.Big, d1.dec.Big)
	return copy
}

func (d Decimal) SubUint64(d1 uint64) Decimal {
	copy := d.Clone()
	copy.dec.Sub(copy.dec.Big, new(decimal.Big).SetUint64(d1))
	return copy
}

func (d Decimal) SubFloat64(d1 float64) Decimal {
	copy := d.Clone()
	copy.dec.Sub(copy.dec.Big, new(decimal.Big).SetFloat64(d1))
	return copy
}

func (d Decimal) Mul(d1 Decimal) Decimal {
	copy := d.Clone()
	copy.dec.Mul(copy.dec.Big, d1.dec.Big)
	return copy
}

func (d Decimal) MulUint64(d1 uint64) Decimal {
	copy := d.Clone()
	copy.dec.Mul(copy.dec.Big, new(decimal.Big).SetUint64(d1))
	return copy
}

func (d Decimal) MulFloat64(d1 float64) Decimal {
	copy := d.Clone()
	copy.dec.Mul(copy.dec.Big, new(decimal.Big).SetFloat64(d1))
	return copy
}

func (d Decimal) Quo(d1 Decimal) Decimal {
	copy := d.Clone()
	copy.dec.Quo(copy.dec.Big, d1.dec.Big)
	return copy
}

func (d Decimal) QuoUint64(d1 uint64) Decimal {
	copy := d.Clone()
	copy.dec.Quo(copy.dec.Big, new(decimal.Big).SetUint64(d1))
	return copy
}

func (d Decimal) QuoFloat64(d1 float64) Decimal {
	copy := d.Clone()
	copy.dec.Quo(copy.dec.Big, new(decimal.Big).SetFloat64(d1))
	return copy
}

func (d Decimal) Uint64() uint64 {
	ui, _ := d.dec.Uint64()
	return ui
}

func (d Decimal) Float64() float64 {
	ui, _ := d.dec.Float64()
	return ui
}

func (d Decimal) String() string {
	return d.dec.String()
}
