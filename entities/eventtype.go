package entities

import (
	"database/sql/driver"
	"errors"
)

type EventType string

func (et *EventType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan EventType")
	}

	*et = EventType(str)

	return nil
}

func (et *EventType) Value() (driver.Value, error) {
	return et, nil
}
