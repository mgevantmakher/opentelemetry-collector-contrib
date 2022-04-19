// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solacereceiver

import "errors"

// ReceiverError is an error returned from the Solace receiver.
type ReceiverError interface {
	error
	init(message string, wrapped error)
	// Message returns the error message for the API exception.
	Message() string
	//Retriable checks if error is retriable
	Retriable() bool
	// Unwrap returns the wrapped error.
	Unwrap() error

	// The error built-in interface type is the conventional interface for
	// representing an error condition, with the nil value representing no error.
	Error() string
}

type receiverError struct {
	message string
	wrapped error
}

// setErrorInfo will set the information on the error.
func (err *receiverError) init(message string, wrapped error) {
	err.message = message
	err.wrapped = wrapped
}

// Error returns the error message.
func (err *receiverError) Error() string {
	return err.message
}

// Unwrap returns the wrapped error.
func (err *receiverError) Unwrap() error {
	return err.wrapped
}

// Message returns the error message.
func (err *receiverError) Message() string {
	return err.message
}

// UnmarchallingError indicates that a AMQP messsage payload unmarshalling to Traces failed.
type UnmarshallingError struct {
	receiverError
	missingFields *[]string
}

// Message returns the error message.
func (u *UnmarshallingError) Message() string {
	return u.message
}

func (u *UnmarshallingError) Error() string {
	return u.Message()
}

// Retriable returns false, it is a fatal error.
func (u *UnmarshallingError) Retriable() bool {
	return false
}

type UnknownMessageError struct {
	UnmarshallingError
}

func (v *UnknownMessageError) Error() string {
	return v.Message()
}

func (v *UnknownMessageError) Retriable() bool {
	return false
}

type VersionIncompatibleUnmarshallingError struct {
	UnmarshallingError
}

func (v *VersionIncompatibleUnmarshallingError) Error() string {
	return v.Message()
}

// Retriable returns false, it is a fatal error.
func (v *VersionIncompatibleUnmarshallingError) Retriable() bool {
	return false
}

//// VersionIncompatible returns false, if unmarchalling problem related to incompatible version schema is detected.
//func (u *UnmarshallingError) VersionIncompatible() bool {
//	return u.versionIncompatible
//}

// AuthenticationError indicates an authentication related error occurred when connecting to a remote event broker.
// The pointer type *AuthenticationError is returned.
type AuthenticationError struct {
	receiverError
}

// ServiceUnreachableError indicates that service connection to event broker could not be established.
// The pointer type *ServiceUnreachableError is returned.
type ServiceUnreachableError struct {
	receiverError
}

// MissingConfigurationError indicates that receiver startup configuration is invalid.
// The pointer type *MissingConfigurationError is returned.
type MissingConfigurationError struct {
	receiverError
	missingSettings *[]string
}

func (m MissingConfigurationError) Message() string {
	return m.message
}

func (m MissingConfigurationError) Error() string {
	return m.Message()
}

func (m MissingConfigurationError) Retriable() bool {
	return false
}

// newError returns a new Solace error with the specified message and wrapped error.
func newError(err ReceiverError, message string, wrapped error) ReceiverError {
	err.init(message, wrapped)
	return err
}

// NewVersionIncompatibleUnmarshallingError returns a new UnmarshallingError with a flag unknownVersion set to true with the specified message and wrapped error.
func NewVersionIncompatibleUnmarshallingError(message string, wrapped error) *VersionIncompatibleUnmarshallingError {
	e := &VersionIncompatibleUnmarshallingError{}
	newError(e, message, wrapped)
	return e
}

// NewMissingConfigurationError returns a new MissingConfigurationError with a slice of missing required configuration fields, a message and wrapped error.
func NewMissingConfigurationError(message string, wrapped error, fields *[]string) *MissingConfigurationError {
	e := &MissingConfigurationError{missingSettings: fields}
	newError(e, message, wrapped)
	return e
}

func NewUnknownMessageError(message string, wrapped error) *UnknownMessageError {
	e := &UnknownMessageError{}
	newError(e, message, wrapped)
	return e
}

// NewUnmarshallingError returns a new UnmarshallingError and slice of missing required fields with the specified message and wrapped error.
func NewUnmarshallingError(message string, wrapped error, fields *[]string) *UnmarshallingError {
	e := &UnmarshallingError{missingFields: fields}
	newError(e, message, wrapped)
	return e
}

func IsVersionIncompatibleUnmarshallingError(err error) bool {
	if err == nil {
		return false
	}
	return errors.As(err, &VersionIncompatibleUnmarshallingError{})
}

func IsUnknownMessageError(err error) bool {
	if err == nil {
		return false
	}
	return errors.As(err, &UnknownMessageError{})
}
