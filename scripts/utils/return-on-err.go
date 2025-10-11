package utils

func ReturnOnErr(err error) error {
	if err != nil {
		LogError("Critical error occurred", err)
		panic(err)
	}
	return nil
}
