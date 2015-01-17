package models

type Inventory interface {
	Init(num float64)
	Destory()
	SetStock(num float64)
	GetStock() float64
	AddStock(num float64) float64
	ReduceStock(num float64) float64
	// SetupAlert(num float64)
	// Alert() bool
}
