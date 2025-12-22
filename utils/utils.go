package utils

import "os"

func Get_file_contents(file_name string) string {

	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := os.ReadFile(file_name)
	if err != nil {
		panic(err)
	}

	return string(b)
}
