// Package metrics allows to send metrics values.
package metric

// MustInit checks whether result can be initialized.
func MustInit[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}

	return result
}
