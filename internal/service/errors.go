package service

import "fmt"

func (e *UnsupportedActionError) Error() string {
	return fmt.Sprintf(unsupportedActionFormat, e.Name, e.Action)
}
