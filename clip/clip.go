/*
Borrowwd from https://github.com/atotto/clipboard
*/

package clip

import (
	"os/exec"
	"time"
)

var (
	pasteCmdArgs = "pbpaste"
	copyCmdArgs  = "pbcopy"
)

func getPasteCommand() *exec.Cmd {
	return exec.Command(pasteCmdArgs)
}

func getCopyCommand() *exec.Cmd {
	return exec.Command(copyCmdArgs)
}

func readAll() (string, error) {
	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func WriteAll(text string) error {
	copyCmd := getCopyCommand()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}

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
