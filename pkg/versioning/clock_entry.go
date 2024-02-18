package versioning

import "errors"

type ClockEntry struct {
	NodeId  uint16
	Version uint64
}

func NewClockEntry(nodeId uint16, version uint64) (*ClockEntry, error) {

	if version < 1 {
		return nil, errors.New("version must be 1 of higher")
	}
	return &ClockEntry{
		NodeId:  nodeId,
		Version: version,
	}, nil
}

func (ce *ClockEntry) Incremented() *ClockEntry {
	return &ClockEntry{
		NodeId:  ce.NodeId,
		Version: ce.Version + 1,
	}
}