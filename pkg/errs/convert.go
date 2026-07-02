package errs

import "google.golang.org/grpc/status"

func GRPCtoREST(err error) (*status.Status, bool) {
	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}
	return st, true
}
