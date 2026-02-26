package entity

// Record describes the behaviour that the rest of the application requires from a record.
type Record interface {
	Copy() Record
	GetID() int
	SetID(int)
	GetData() map[string]string
	SetData(map[string]string)
}
