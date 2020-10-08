package dir

import "os"

func ChangeDir(dir string) (func() error, error)  {
	current, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err := os.Chdir(dir); err != nil {
		return nil, err
	}
	return func() error { return os.Chdir(current) }, nil
}
