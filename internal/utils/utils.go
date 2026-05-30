package utils

func OctToGB(oct float64) float64 {
	return oct / (1024.0 * 1024.0 * 1024.0)
}

func OctToMo(oct float64) float64 {
	return oct / (1024.0 * 1024.0)
}
