package wallet

import (
	"runtime"
	"sync"
)

// Platform represents the operating system platform
type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWindows Platform = "windows"
	PlatformMacOS   Platform = "macos"
	PlatformLinux   Platform = "linux"
	PlatformWeb     Platform = "web"
	PlatformUnknown Platform = "unknown"
)

// PlatformCapabilities represents features available on each platform
type PlatformCapabilities struct {
	SupportsBiometric    bool
	SupportsNFC          bool
	SupportsPushNotify   bool
	SupportsCamera       bool
	SupportsContacts     bool
	SupportsClipboard    bool
	SupportsFileSystem   bool
	SupportsBackgroundOps bool
}

// PlatformConfig holds platform-specific configuration
type PlatformConfig struct {
	Platform     Platform
	Version      string
	Capabilities PlatformCapabilities
	Features     map[string]bool
	mutex        sync.RWMutex
}

// NewPlatformConfig creates a new platform configuration
func NewPlatformConfig() *PlatformConfig {
	pc := &PlatformConfig{
		Platform: detectPlatform(),
		Version:  runtime.Version(),
		Features: make(map[string]bool),
	}

	pc.Capabilities = pc.getPlatformCapabilities()
	pc.setDefaultFeatures()

	return pc
}

// detectPlatform detects the current operating system platform
func detectPlatform() Platform {
	switch runtime.GOOS {
	case "darwin":
		return PlatformMacOS
	case "linux":
		return PlatformLinux
	case "windows":
		return PlatformWindows
	case "android":
		return PlatformAndroid
	case "ios":
		return PlatformIOS
	default:
		return PlatformUnknown
	}
}

// getPlatformCapabilities returns capabilities for the current platform
func (pc *PlatformConfig) getPlatformCapabilities() PlatformCapabilities {
	switch pc.Platform {
	case PlatformIOS:
		return PlatformCapabilities{
			SupportsBiometric:    true,
			SupportsNFC:          true,
			SupportsPushNotify:   true,
			SupportsCamera:       true,
			SupportsContacts:     true,
			SupportsClipboard:    true,
			SupportsFileSystem:   true,
			SupportsBackgroundOps: true,
		}
	case PlatformAndroid:
		return PlatformCapabilities{
			SupportsBiometric:    true,
			SupportsNFC:          true,
			SupportsPushNotify:   true,
			SupportsCamera:       true,
			SupportsContacts:     true,
			SupportsClipboard:    true,
			SupportsFileSystem:   true,
			SupportsBackgroundOps: true,
		}
	case PlatformWindows, PlatformMacOS, PlatformLinux:
		return PlatformCapabilities{
			SupportsBiometric:    true,  // Windows Hello, Touch ID, etc.
			SupportsNFC:          false,
			SupportsPushNotify:   true,
			SupportsCamera:       true,
			SupportsContacts:     false,
			SupportsClipboard:    true,
			SupportsFileSystem:   true,
			SupportsBackgroundOps: true,
		}
	case PlatformWeb:
		return PlatformCapabilities{
			SupportsBiometric:    false,
			SupportsNFC:          false,
			SupportsPushNotify:   true,
			SupportsCamera:       true,
			SupportsContacts:     false,
			SupportsClipboard:    true,
			SupportsFileSystem:   false,
			SupportsBackgroundOps: false,
		}
	default:
		return PlatformCapabilities{}
	}
}

// setDefaultFeatures sets default feature flags based on platform
func (pc *PlatformConfig) setDefaultFeatures() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	// Mobile-specific features
	if pc.Platform == PlatformIOS || pc.Platform == PlatformAndroid {
		pc.Features["qr_scanner"] = true
		pc.Features["nfc_payments"] = true
		pc.Features["biometric_auth"] = true
		pc.Features["push_notifications"] = true
		pc.Features["mobile_data_sync"] = true
	}

	// Desktop-specific features
	if pc.Platform == PlatformWindows || pc.Platform == PlatformMacOS || pc.Platform == PlatformLinux {
		pc.Features["advanced_trading"] = true
		pc.Features["multi_monitor"] = true
		pc.Features["hardware_wallet"] = true
		pc.Features["local_node"] = true
		pc.Features["advanced_analytics"] = true
	}

	// Universal features
	pc.Features["multi_currency"] = true
	pc.Features["exchange"] = true
	pc.Features["staking"] = true
	pc.Features["2fa"] = true
	pc.Features["backup_restore"] = true
	pc.Features["transaction_history"] = true
}

// IsFeatureEnabled checks if a feature is enabled
func (pc *PlatformConfig) IsFeatureEnabled(feature string) bool {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	enabled, exists := pc.Features[feature]
	return exists && enabled
}

// EnableFeature enables a specific feature
func (pc *PlatformConfig) EnableFeature(feature string) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.Features[feature] = true
}

// DisableFeature disables a specific feature
func (pc *PlatformConfig) DisableFeature(feature string) {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()

	pc.Features[feature] = false
}

// GetPlatformInfo returns platform information
func (pc *PlatformConfig) GetPlatformInfo() map[string]interface{} {
	pc.mutex.RLock()
	defer pc.mutex.RUnlock()

	return map[string]interface{}{
		"platform":     pc.Platform,
		"version":      pc.Version,
		"capabilities": pc.Capabilities,
		"features":     pc.Features,
		"arch":         runtime.GOARCH,
		"numCPU":       runtime.NumCPU(),
	}
}

// IsMobilePlatform checks if the platform is mobile
func (pc *PlatformConfig) IsMobilePlatform() bool {
	return pc.Platform == PlatformIOS || pc.Platform == PlatformAndroid
}

// IsDesktopPlatform checks if the platform is desktop
func (pc *PlatformConfig) IsDesktopPlatform() bool {
	return pc.Platform == PlatformWindows || pc.Platform == PlatformMacOS || pc.Platform == PlatformLinux
}

// GetRecommendedSettings returns platform-specific recommended settings
func (pc *PlatformConfig) GetRecommendedSettings() map[string]interface{} {
	settings := make(map[string]interface{})

	if pc.IsMobilePlatform() {
		settings["sync_interval"] = 300 // 5 minutes
		settings["cache_size"] = 50 * 1024 * 1024 // 50MB
		settings["max_transactions"] = 1000
		settings["enable_notifications"] = true
		settings["use_mobile_data"] = false
	} else if pc.IsDesktopPlatform() {
		settings["sync_interval"] = 60 // 1 minute
		settings["cache_size"] = 500 * 1024 * 1024 // 500MB
		settings["max_transactions"] = 10000
		settings["enable_notifications"] = true
		settings["run_local_node"] = false
	}

	return settings
}
