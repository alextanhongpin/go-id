package ids

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidID = errors.New("ids: invalid format")
	ErrBadPrefix = errors.New("ids: prefix does not match")
)

type resourceType interface {
	Prefix() string
}

type ID[T resourceType] uuid.UUID

func New[T resourceType]() ID[T] {
	return ID[T](uuid.New())
}

func (id ID[T]) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id ID[T]) Prefix() string {
	var t T
	return t.Prefix()
}

func (id ID[T]) String() string {
	var t T

	// NOTE: The usage of uuid.UUID is to avoid recursion.
	return fmt.Sprintf("%s_%s", t.Prefix(), uuid.UUID(id).String())
}

func (id ID[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ID[T]) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	prefix, uuidStr, ok := strings.Cut(s, "_")
	if !ok {
		return ErrInvalidID
	}

	var t T
	if want, got := t.Prefix(), prefix; want != got {
		return fmt.Errorf(`%w: want "%s_", got %q`, ErrBadPrefix, want, s)
	}

	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return err
	}

	*id = ID[T](uid)

	return nil
}
