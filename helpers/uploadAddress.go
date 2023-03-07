package helpers

import "fmt"

func UploadAddress(server string) string {
	return fmt.Sprintf("https://%s.gofile.io/uploadFile", server)
}
