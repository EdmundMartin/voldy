package cluster

import "fmt"

const (
	DEFAULT_ZONE_ID = 0
	UNSET_ZONE_ID   = -1
)

type Zone struct {
	SerialVersion int
	ZoneId        int
	ProximityList []int
}

func NewDefaultZone() *Zone {
	return &Zone{
		SerialVersion: 1,
		ZoneId:        DEFAULT_ZONE_ID,
		ProximityList: []int{},
	}
}

func NewZone(zoneId int, proximityList []int) *Zone {
	return &Zone{
		SerialVersion: 1,
		ZoneId:        zoneId,
		ProximityList: proximityList,
	}
}

func (z *Zone) String() string {
	return fmt.Sprintf("<Zone Id: %d, ProxmitiyList: %v>", z.ZoneId, z.ProximityList)
}

func (z *Zone) Equals(other *Zone) bool {

	if z.ZoneId != other.ZoneId {
		return false
	}

	idx := 0

	for idx < len(z.ProximityList) && idx < len(other.ProximityList) {

		if z.ProximityList[idx] != other.ProximityList[idx] {
			return false
		}
		idx++
	}

	return true
}

type Zones []*Zone

func (z Zones) Len() int {
	return len(z)
}

func (z Zones) Less(i, j int) bool {
	return z[i].ZoneId < z[j].ZoneId
}

func (z Zones) Swap(i, j int) {
	z[i], z[j] = z[j], z[i]
}
