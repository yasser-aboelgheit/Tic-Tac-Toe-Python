package metric

import (
	"fmt"

)

func StartCollector(tags []string) (func(), error) {
	if err := Run(DefaultAddr(), WithTags(tags)); err != nil {
		return nil, fmt.Errorf("error starting metrics collector: %v", err)
	}

	return Stop, nil
}
