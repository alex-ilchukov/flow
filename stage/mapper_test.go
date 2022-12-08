package stage_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alex-ilchukov/flow/stage"
	"github.com/alex-ilchukov/flow/values"
)

type joint struct {
	ctx       context.Context
	v         []int
	w         []int
	puterr    error
	errs      []error
	reporterr error
}

func (j *joint) Ctx() context.Context {
	return j.ctx
}

func (j *joint) Get() (int, error) {
	v, err := 0, values.Over
	if len(j.v) > 0 {
		v, err, j.v = j.v[0], nil, j.v[1:]
	}

	return v, err
}

func (j *joint) Put(w int) error {
	if j.puterr == nil {
		j.w = append(j.w, w)
	}

	return j.puterr
}

func (j *joint) Report(err error) error {
	if j.reporterr == nil {
		j.errs = append(j.errs, err)
	}

	return j.reporterr
}

type carrier struct {
	limit int
	err   error
}

func (c carrier) mapper(_ context.Context, v int) (int, error) {
	if v < c.limit {
		return v * v, nil
	}

	return 0, c.err
}

func TestMapperFormWhenSuccessful(t *testing.T) {
	t.Parallel()

	j := joint{
		ctx: context.Background(),
		v:   []int{1, 2, 3, 4},
	}

	c := carrier{
		limit: 5,
	}

	m := stage.Mapper[int, int](c.mapper)
	m.Form(&j)

	switch {
	case j.errs != nil:
		t.Errorf("got errors: %v", j.errs)

	case !reflect.DeepEqual(j.w, []int{1, 4, 9, 16}):
		t.Errorf("got invalid put values: %v", j.w)
	}
}

func TestMapperFormWhenMapperIsErrorful(t *testing.T) {
	t.Parallel()

	j := joint{
		ctx: context.Background(),
		v:   []int{1, 2, 3, 4},
	}

	err := errors.New("mapper")
	c := carrier{
		limit: 3,
		err:   err,
	}

	m := stage.Mapper[int, int](c.mapper)
	m.Form(&j)

	switch {
	case !reflect.DeepEqual(j.errs, []error{err, err}):
		t.Errorf("got invalid errors: %v", j.errs)

	case !reflect.DeepEqual(j.w, []int{1, 4}):
		t.Errorf("got invalid put values: %v", j.w)
	}
}

func TestMapperFormWhenMapperIsErrorfulAndJointReportIsErrorful(t *testing.T) {
	t.Parallel()

	reporterr := errors.New("report")
	j := joint{
		ctx:       context.Background(),
		v:         []int{1, 2, 3, 4},
		reporterr: reporterr,
	}

	err := errors.New("mapper")
	c := carrier{
		limit: 3,
		err:   err,
	}

	m := stage.Mapper[int, int](c.mapper)
	m.Form(&j)

	switch {
	case j.errs != nil:
		t.Errorf("got errors: %v", j.errs)

	case !reflect.DeepEqual(j.w, []int{1, 4}):
		t.Errorf("got invalid put values: %v", j.w)
	}
}

func TestMapperFormWhenJointPutIsErrorful(t *testing.T) {
	t.Parallel()

	err := errors.New("put")
	j := joint{
		ctx:    context.Background(),
		v:      []int{1, 2, 3, 4},
		puterr: err,
	}

	c := carrier{
		limit: 3,
	}

	m := stage.Mapper[int, int](c.mapper)
	m.Form(&j)

	switch {
	case j.errs != nil:
		t.Errorf("got errors: %v", j.errs)

	case j.w != nil:
		t.Errorf("got invalid put values: %v", j.w)
	}
}
