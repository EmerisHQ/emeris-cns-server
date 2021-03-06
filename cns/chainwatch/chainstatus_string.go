// Code generated by "stringer -type=ChainStatus"; DO NOT EDIT.

package chainwatch

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[starting-0]
	_ = x[running-1]
	_ = x[relayerConnecting-2]
	_ = x[done-3]
}

const _ChainStatus_name = "startingrunningrelayerConnectingdone"

var _ChainStatus_index = [...]uint8{0, 8, 15, 32, 36}

func (i ChainStatus) String() string {
	if i >= ChainStatus(len(_ChainStatus_index)-1) {
		return "ChainStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ChainStatus_name[_ChainStatus_index[i]:_ChainStatus_index[i+1]]
}
