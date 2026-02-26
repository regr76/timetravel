package entity

import "maps"

type PersistentRecord struct {
	ID      int               `json:"id"`
	Version int               `json:"version"`
	Start   string            `json:"start_dt"`
	End     string            `json:"end_dt,omitempty"`
	Data    map[string]string `json:"data"`
}

func (d *PersistentRecord) Copy() Record {
	return &PersistentRecord{
		ID:      d.ID,
		Version: d.Version,
		Start:   d.Start,
		End:     d.End,
		Data:    maps.Clone(d.Data),
	}
}

func (d *PersistentRecord) GetID() int {
	return d.ID
}

func (d *PersistentRecord) SetID(id int) {
	d.ID = id
}

func (d *PersistentRecord) GetData() map[string]string {
	return d.Data
}

func (d *PersistentRecord) SetData(data map[string]string) {
	d.Data = data
}

func (d *PersistentRecord) GetVersion() int {
	return d.Version
}

func (d *PersistentRecord) SetVersion(version int) {
	d.Version = version
}

func (d *PersistentRecord) GetStart() string {
	return d.Start
}

func (d *PersistentRecord) SetStart(start string) {
	d.Start = start
}

func (d *PersistentRecord) GetEnd() string {
	return d.End
}

func (d *PersistentRecord) SetEnd(end string) {
	d.End = end
}
