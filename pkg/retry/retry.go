package retry

import "time"

type OperationFunc func() error

func Retry(attempts int, delay time.Duration, factor float64, op OperationFunc) error {
	var err error
	f := 1.
	for attempt := 1; attempt <= attempts; attempt++ {
		err = op()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(float64(delay) * f))
		f *= factor
	}
	return err
}
