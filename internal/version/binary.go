package version

import (
	"fmt"
	"runtime"
)

const Binary = "0.01"

func String(app string) string {
	return fmt.Sprintf("%s/%s (built w/%s %s %s)", app, Binary, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
