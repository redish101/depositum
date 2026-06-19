package common

import (
	"errors"
	"strconv"
)

var ErrIDNotUint = errors.New("ID 必须为整数")

func ReadID(idStr string) (uint, error) {
	v64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, ErrIDNotUint
	}

	v := uint(v64)

	return v, nil
}
