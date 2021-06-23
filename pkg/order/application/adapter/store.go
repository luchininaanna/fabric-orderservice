package adapter

type Fabric struct {
	Id       string
	Name     string
	Cost     float32
	Quantity float32
}

type StoreAdapter interface {
	GetFabrics() ([]Fabric, error)
}
