package helpers

func CheckRowsAffected(rows int64, err error) error {
	if rows == 0 {
		return err
	}

	return nil
}
