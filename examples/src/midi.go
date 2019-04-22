package examples

import (
	"fmt"
	"pixelext/ui"

	"github.com/rakyll/portmidi"
)

func NewMidiDeviceDropDown() *ui.DropDown {
	d := ui.NewDropDown("mididevices", "basic", 150, 20, 100)
	count := portmidi.CountDevices()
	for i := 0; i < count; i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if info != nil {
			if info.IsInputAvailable {
				d.AddValue(info.Name, fmt.Sprintf("%d", i))
			}
		}
	}
	return d
}
