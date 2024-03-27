package shared

import "runtime"

func IsPlatformLinux() bool {
	return runtime.GOOS == "linux"
}

func IsPlatformDarwin() bool {
	return runtime.GOOS == "darwin"
}

func IsPlatformWindows() bool {
	return runtime.GOOS == "windows"
}
