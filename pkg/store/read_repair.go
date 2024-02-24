package store

import (
	"sort"
	"voldy/pkg/versioning"
)

func LatestVersion(results []*versioning.Versioned) *versioning.Versioned {
	if len(results) == 1 {
		return results[0]
	}

	sort.Sort(versioning.VersionedCollection(results))

	return results[len(results)-1]
}

func CombinedValues(latest versioning.Version, outdated versioning.Version) {
	// TODO - latest should include values not include in latest - but included on older
}
