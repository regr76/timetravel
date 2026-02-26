package entity

import "maps"

type InMemoryRecord struct {
	ID   int               `json:"id"`
	Data map[string]string `json:"data"`
}

func (d *InMemoryRecord) Copy() Record {
	return &InMemoryRecord{
		ID:   d.ID,
		Data: maps.Clone(d.Data),
	}
}

func (d *InMemoryRecord) GetID() int {
	return d.ID
}

func (d *InMemoryRecord) SetID(id int) {
	d.ID = id
}

func (d *InMemoryRecord) GetData() map[string]string {
	return d.Data
}

func (d *InMemoryRecord) SetData(data map[string]string) {
	d.Data = data
}
