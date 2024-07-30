package ManagementInterface

type IAes interface {
	Encrypt([]byte) ([]byte, error)

	Decrypt([]byte) ([]byte, error)
}
