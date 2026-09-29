package teardown

import "fmt"

func (e *StepError) Error() string {
	return fmt.Sprintf(stepErrorFormat, e.Step, e.Err)
}

func (e *StepError) Unwrap() error {
	return e.Err
}

func (e *UnsafeTargetError) Error() string {
	return fmt.Sprintf(unsafeTargetFormat, e.Path)
}
