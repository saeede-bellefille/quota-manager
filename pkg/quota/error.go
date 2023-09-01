package quota

const (
	uniqueError Error = "Unique error"
	countError  Error = "Counter limit"
	sizeError   Error = "Size limit"
)

type Error string

func (err Error) Error() string {
	return string(err)
}
