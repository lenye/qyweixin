package version

import (
	"fmt"
	"runtime"
)

const Binary = "201704261409"

func String(app string) string {
	return fmt.Sprintf("%s/%s (built w/%s)", app, Binary, runtime.Version())
}
