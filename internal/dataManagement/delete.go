package dataManagement

import "os"

func (r *Repository) DeleteFolder(filepath string) error {
	return os.RemoveAll(filepath)
}
