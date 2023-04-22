package main

func isWifiTrusted(configTrusted map[string]string, wifiName string) bool {
	for _, wt := range configTrusted {
		if wt == wifiName {
			return true
		}
	}

	return false
}
