package kernel

import (
	"fmt"
	"syscall"
	"time"

	"github.com/alibaba/pouch/pkg/exec"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// VersionInfo holds information about the kernel.
type VersionInfo struct {
	Kernel int    // Version of the kernel (e.g. 4.1.2-generic -> 4)
	Major  int    // Major part of the kernel version (e.g. 4.1.2-generic -> 1)
	Minor  int    // Minor part of the kernel version (e.g. 4.1.2-generic -> 2)
	Flavor string // Flavor of the kernel version (e.g. 4.1.2-generic -> generic)
}

// String returns the kernel version's string format.
func (k *VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", k.Kernel, k.Major, k.Minor, k.Flavor)
}

// GetKernelVersion returns the kernel version info.
func GetKernelVersion() (*VersionInfo, error) {

	var (
		kernel, major, minor int
		flavor               string
	)

	_, stdout, _, err := exec.Run(0, "uname", "-r")
	if err != nil {
		return nil, errors.Wrap(err, "failed to run command uname -r")
	}

	parsed, _ := fmt.Sscanf(stdout, "%d.%d.%d-%s", &kernel, &major, &minor, &flavor)
	if parsed < 3 {
		return nil, fmt.Errorf("Can't parse kernel version, release: %s" + stdout)
	}

	return &VersionInfo{
		Kernel: kernel,
		Major:  major,
		Minor:  minor,
		Flavor: flavor,
	}, nil
}

var (
	kernelVersionCacheInfo *VersionInfo
)

func GetKernelVersionByVarCache() (*VersionInfo, error) {

	if kernelVersionCacheInfo == nil {

		var (
			kernel, major, minor int
			flavor               string
		)

		_, stdout, _, err := exec.Run(0, "uname", "-r")
		if err != nil {
			return nil, errors.Wrap(err, "failed to run command uname -r")
		}

		parsed, _ := fmt.Sscanf(stdout, "%d.%d.%d-%s", &kernel, &major, &minor, &flavor)
		if parsed < 3 {
			return nil, fmt.Errorf("Can't parse kernel version, release: %s" + stdout)
		}

		kernelVersionCacheInfo = &VersionInfo{
			Kernel: kernel,
			Major:  major,
			Minor:  minor,
			Flavor: flavor,
		}
	}

	return kernelVersionCacheInfo, nil
}

var (
	kernelVersionCacheTTLInfo        *VersionInfo
	kernelVersionCacheTTLLastUpdated int64 = 0
	kernelVersionCacheTTLTimeMax     int64 = 60
)

func GetKernelVersionByVarCacheWithTTL() (*VersionInfo, error) {

	updated := time.Now().Unix()

	if kernelVersionCacheTTLInfo == nil ||
		updated-kernelVersionCacheTTLLastUpdated > kernelVersionCacheTTLTimeMax {

		var (
			kernel, major, minor int
			flavor               string
		)

		_, stdout, _, err := exec.Run(0, "uname", "-r")
		if err != nil {
			return nil, errors.Wrap(err, "failed to run command uname -r")
		}

		parsed, _ := fmt.Sscanf(stdout, "%d.%d.%d-%s", &kernel, &major, &minor, &flavor)
		if parsed < 3 {
			return nil, fmt.Errorf("Can't parse kernel version, release: %s" + stdout)
		}

		kernelVersionCacheTTLInfo = &VersionInfo{
			Kernel: kernel,
			Major:  major,
			Minor:  minor,
			Flavor: flavor,
		}
		kernelVersionCacheTTLLastUpdated = updated
	}

	return kernelVersionCacheTTLInfo, nil
}

func GetKernelVersionByUnix() (*VersionInfo, error) {
	var (
		kernel, major, minor int
		flavor               string
	)
	buf := unix.Utsname{}
	err := unix.Uname(&buf)
	if err != nil {
		return nil, err
	}
	release := string(buf.Release[:])
	parsed, _ := fmt.Sscanf(release, "%d.%d.%d-%s", &kernel, &major, &minor, &flavor)
	if parsed < 3 {
		return nil, fmt.Errorf("Can't parse kernel version, release: %s" + release)
	}

	return &VersionInfo{
		Kernel: kernel,
		Major:  major,
		Minor:  minor,
		Flavor: flavor,
	}, nil
}

func GetKernelVersionBySyscall() (*VersionInfo, error) {
	var (
		kernel, major, minor int
		flavor               string
	)
	buf := syscall.Utsname{}
	err := syscall.Uname(&buf)
	if err != nil {
		return nil, err
	}
	release := charsToString(buf.Release[:])
	parsed, _ := fmt.Sscanf(release, "%d.%d.%d-%s", &kernel, &major, &minor, &flavor)
	if parsed < 3 {
		return nil, fmt.Errorf("Can't parse kernel version, release: %s" + release)
	}

	return &VersionInfo{
		Kernel: kernel,
		Major:  major,
		Minor:  minor,
		Flavor: flavor,
	}, nil
}

func charsToString(ca []int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}
