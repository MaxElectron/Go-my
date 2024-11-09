//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func createAPIError(method string, authError *kvapi.AuthError) *kvapi.APIError {
	return &kvapi.APIError{
		Method: method,
		Err:    authError,
	}
}

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var apiError *kvapi.APIError
	var authError *kvapi.AuthError
	var conflictError *kvapi.ConflictError
	var newValue string
	var response *kvapi.GetResponse
	var err error

	keyMissed := false

	for {
		// check for key
		if !keyMissed {
			response, err = c.Get(&kvapi.GetRequest{Key: key})
		}

		// check for authentication error before anything else
		if errors.As(err, &authError) {
			return createAPIError("get", authError)
		}

	missingKey:
		keyFound := !errors.Is(err, kvapi.ErrKeyNotFound)
		// key not found
		if !keyFound {
			newValue, err = updateFn(nil)
			if err != nil {
				return err
			}
		}

		// api failed
		if errors.As(err, &apiError) && keyFound {
			continue
		}

		// key found
		if keyFound {
			newValue, err = updateFn(&response.Value)
			if err != nil {
				return err
			}
		}

		// generate new uuid
		newUUID := uuid.Must(uuid.NewV4())

		// set up the request
		setRequest := kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: uuid.UUID{},
			NewVersion: newUUID,
		}

		if response != nil {
			setRequest.OldVersion = response.Version
		}

	secondaryApiFail:
		// send the request
		_, err = c.Set(&setRequest)

		// check for authentication error before anything else
		if errors.As(err, &authError) {
			return createAPIError("set", authError)
		}

		// handle conflict at writing
		if errors.As(err, &conflictError) {
			if conflictError.ExpectedVersion == newUUID {
				return nil
			}

			response.Version = conflictError.ExpectedVersion
			continue
		}

		// key not found
		if errors.Is(err, kvapi.ErrKeyNotFound) {
			if response == nil {
				response = &kvapi.GetResponse{
					Value:   newValue,
					Version: uuid.UUID{},
				}
			} else {
				response.Value = newValue
				response.Version = uuid.UUID{}
			}

			keyMissed = true
			goto missingKey
		}

		// api failed
		if errors.As(err, &apiError) {
			goto secondaryApiFail
		}

		// unexpected fail
		return nil
	}
}
