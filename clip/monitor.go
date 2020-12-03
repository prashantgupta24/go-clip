package clip

import "time"

//Monitor clipboard for changes
func Monitor(interval time.Duration, stopCh <-chan struct{}, changes chan<- string) error {
	defer close(changes)

	currentValue, err := readAll()
	if err != nil {
		return err
	}

	for {
		select {
		case <-stopCh:
			return nil
		default:
			newValue, _ := readAll()
			if newValue != currentValue {
				currentValue = newValue
				changes <- currentValue
			}
		}
		time.Sleep(interval)
	}
}
