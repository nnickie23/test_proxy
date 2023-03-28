package errs

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ServiceNotAvailable = Err("service_not_available")
	ObjectNotFound      = Err("object_not_found")
)
