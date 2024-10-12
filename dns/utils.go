package dns

import (
	"encoding/binary"
)

var btoi16 = binary.BigEndian.Uint16
var btoi32 = binary.BigEndian.Uint32
