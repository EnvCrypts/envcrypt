package types

import (
	"database/sql/driver"
	"fmt"
	"net"
)

type NullIP struct {
	IP    net.IP
	Valid bool
}

func (ip *NullIP) Scan(src any) error {
	if src == nil {
		ip.IP = nil
		ip.Valid = false
		return nil
	}

	var value string
	switch v := src.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		return fmt.Errorf("cannot scan %T into NullIP", src)
	}

	parsed := net.ParseIP(value)
	if parsed == nil {
		return fmt.Errorf("invalid IP address %q", value)
	}

	ip.IP = parsed
	ip.Valid = true
	return nil
}

func (ip NullIP) Value() (driver.Value, error) {
	if !ip.Valid {
		return nil, nil
	}
	return ip.IP.String(), nil
}
