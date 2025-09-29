package utils

func ReturnOnErr(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}
