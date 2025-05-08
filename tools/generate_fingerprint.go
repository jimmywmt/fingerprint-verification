package tools

import (
	"os"
	"strings"
)

func GenerateFingerprint() (string, error) {
	cpuModel, err := readFieldFromFile("/proc/cpuinfo", "model name")
	if err != nil {
		return "", err
	}
	cpuVendor, err := readFieldFromFile("/proc/cpuinfo", "vendor_id")
	if err != nil {
		return "", err
	}
	uuidBytes, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
	if err != nil {
		return "", err
	}
	biosUUID := strings.TrimSpace(string(uuidBytes))

	fingerprint := cpuVendor + "::" + cpuModel + "::" + biosUUID
	return fingerprint, nil
}
