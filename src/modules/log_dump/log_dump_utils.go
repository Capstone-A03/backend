package logdump

import (
	"errors"
	"math"
)

func convertVolumeToM3(volume *float64, unit *string) (*float64, error) {
	if unit == nil {
		return volume, nil
	}

	var vol float64
	switch *unit {
	case "mm3":
		vol = *volume * math.Pow(10, -9)
		return &vol, nil
	case "cm3":
		vol = *volume * math.Pow(10, -6)
		return &vol, nil
	case "dm3":
		vol = *volume * math.Pow(10, -3)
		return &vol, nil
	case "m3":
		return volume, nil
	case "dam3":
		vol = *volume * math.Pow(10, 3)
		return &vol, nil
	case "hm3":
		vol = *volume * math.Pow(10, 6)
		return &vol, nil
	case "km3":
		vol = *volume * math.Pow(10, 9)
		return &vol, nil
	default:
		return nil, errors.New("unit unknown")
	}
}
