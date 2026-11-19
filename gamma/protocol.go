package gamma

import (
	"fmt"
	"net/url"
)

// Match Gamma/ts-sdk repeated version query parameters without adding a default.
func addProtocolVersions(q url.Values, versions []string) error {
	if len(versions) > 2 {
		return fmt.Errorf("version accepts at most two values")
	}
	for _, version := range versions {
		if version != "v1" && version != "v2" {
			return fmt.Errorf("unsupported protocol version %q", version)
		}
		q.Add("version", version)
	}
	return nil
}
