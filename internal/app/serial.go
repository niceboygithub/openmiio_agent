package app

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
)

type SerialStats struct {
	stat struct {
		BluetoothRX uint32 `json:"bluetooth_rx,omitempty"`
		BluetoothTX uint32 `json:"bluetooth_tx,omitempty"`
		BluetoothOE uint32 `json:"bluetooth_oe,omitempty"`
		ZigbeeRX    uint32 `json:"zigbee_rx,omitempty"`
		ZigbeeTX    uint32 `json:"zigbee_tx,omitempty"`
		ZigbeeOE    uint32 `json:"zigbee_oe,omitempty"`
		IrdaRX    uint32 `json:"irda_rx,omitempty"`
		IrdaTX    uint32 `json:"irda_tx,omitempty"`
		IrdaOE    uint32 `json:"irda_oe,omitempty"`
	}
}

func (s *SerialStats) MarshalJSON() ([]byte, error) {
	switch Model {
	case ModelMGW:
		counters := readSerial("/proc/tty/driver/serial")
		if len(counters) >= 9 {
			s.stat.BluetoothTX = counters[3]
			s.stat.BluetoothRX = counters[4]
			s.stat.BluetoothOE = counters[5]
			s.stat.ZigbeeTX = counters[6]
			s.stat.ZigbeeRX = counters[7]
			s.stat.ZigbeeOE = counters[8]
		}
	case ModelE1:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 6 {
			s.stat.ZigbeeTX = counters[3]
			s.stat.ZigbeeRX = counters[4]
			s.stat.ZigbeeOE = counters[5]
		}
	case ModelMGW2:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 9 {
			s.stat.BluetoothTX = counters[6]
			s.stat.BluetoothRX = counters[7]
			s.stat.BluetoothOE = counters[8]
			s.stat.ZigbeeTX = counters[3]
			s.stat.ZigbeeRX = counters[4]
			s.stat.ZigbeeOE = counters[5]
		}
	case ModelM1S:
		counters := readSerial("/proc/tty/driver/serial")
		if len(counters) >= 6 {
			s.stat.ZigbeeTX = counters[3]
			s.stat.ZigbeeRX = counters[4]
			s.stat.ZigbeeOE = counters[5]
		}
	case ModelM2:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 9 {
			s.stat.ZigbeeTX = counters[3]
			s.stat.ZigbeeRX = counters[4]
			s.stat.ZigbeeOE = counters[5]
			s.stat.IrdaTX = counters[6]
			s.stat.IrdaRX = counters[7]
			s.stat.IrdaOE = counters[8]
		}
	case ModelM1S22, ModelG2HPro:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 6 {
			s.stat.ZigbeeTX = counters[3]
			s.stat.ZigbeeRX = counters[4]
			s.stat.ZigbeeOE = counters[5]
		}
	case ModelG2H:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 9 {
			s.stat.ZigbeeTX = counters[6]
			s.stat.ZigbeeRX = counters[7]
			s.stat.ZigbeeOE = counters[8]
		}
	case ModelG3:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 9 {
			s.stat.IrdaTX = counters[6]
			s.stat.IrdaRX = counters[7]
			s.stat.IrdaOE = counters[8]
			s.stat.ZigbeeTX = counters[9]
			s.stat.ZigbeeRX = counters[10]
			s.stat.ZigbeeOE = counters[11]
		}
	case ModelM2PoE, ModelM3:
		counters := readSerial("/proc/tty/driver/ms_uart")
		if len(counters) >= 12 {
			s.stat.ZigbeeTX = counters[3]
			s.stat.ZigbeeRX = counters[4]
			s.stat.ZigbeeOE = counters[5]
			s.stat.BluetoothTX = counters[6]
			s.stat.BluetoothRX = counters[7]
			s.stat.BluetoothOE = counters[8]
			s.stat.IrdaTX = counters[9]
			s.stat.IrdaRX = counters[10]
			s.stat.IrdaOE = counters[11]
		}
	}

	return json.Marshal(s.stat)
}

func readSerial(name string) (counters []uint32) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil
	}

	return parseSerial(b)
}

func parseSerial(b []byte) (counters []uint32) {
	for {
		// 1. Search tx start
		i := bytes.Index(b, []byte("tx:"))
		if i < 0 || i+3 > len(b) {
			return
		}
		b = b[i+3:]

		// 2. Search tx end
		i = bytes.IndexByte(b, ' ')
		if i < 0 || i+4 > len(b) {
			return
		}

		// 3. Read tx
		x, err := strconv.Atoi(string(b[:i]))
		if err != nil {
			return
		}

		tx := uint32(x)
		b = b[i+4:]

		// 4. Search rx end
		i = bytes.IndexByte(b, ' ')
		if i < 0 {
			return
		}

		// 5. Read rx
		x, err = strconv.Atoi(string(b[:i]))
		if err != nil {
			return
		}

		rx := uint32(x)

		// 6. Search line end
		i = bytes.IndexByte(b, '\n')
		if i < 0 {
			return
		}

		var oe uint32

		// 7. Search oe start
		i = bytes.Index(b[:i], []byte("oe:"))
		if i > 0 && i+3 < len(b) {
			b2 := b[i+3:]

			// 8. Search oe end
			i = bytes.IndexByte(b2, ' ')
			if i < 0 {
				return
			}

			// 9. Read oe
			x, err = strconv.Atoi(string(b2[:i]))
			if err != nil {
				return
			}

			oe = uint32(x)
		}

		counters = append(counters, tx, rx, oe)
	}
}
