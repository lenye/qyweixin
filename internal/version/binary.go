package version

import (
	"fmt"
	"runtime"
)

const Binary = "0.01"

func String(app string) string {
	return fmt.Sprintf("%s/%s (built with %s %s for %s/%s)", app, Binary, runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
