package common

import "strings"

// GetWrapperBinaryPath returns the BinaryPathName that should be set in SCM.
// Format: "<path-to-exe>" service <serviceName>
func GetWrapperBinaryPath(exePath, serviceName string) string {
	if !strings.HasPrefix(exePath, `"`) {
		exePath = `"` + exePath + `"`
	}
	return exePath + " service " + serviceName
}

// IsWrapperBinaryPath checks if a BinaryPathName belongs to an NSSM-Plus wrapped service.
func IsWrapperBinaryPath(binaryPathName string) bool {
	return strings.Contains(binaryPathName, " service ")
}

// ExtractServiceName extracts the service name from a wrapper BinaryPathName.
func ExtractServiceName(binaryPathName string) string {
	idx := strings.LastIndex(binaryPathName, " service ")
	if idx < 0 {
		return ""
	}
	name := binaryPathName[idx+len(" service "):]
	name = strings.TrimSpace(name)
	if strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`) {
		name = name[1 : len(name)-1]
	}
	return name
}
