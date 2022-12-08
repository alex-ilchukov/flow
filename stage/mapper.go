package stage

import "context"

// Mapper is type of generalized functions, which map values of type V to
// values of type W within the provided context. The functions can return an
// error value.
type Mapper[V, W any] func(context.Context, V) (W, error)

// Form maps every incoming value from joint j with help of mapper m to new
// value and sends the new value to the joint. It reports to joint on mapping
// errors, but it just breaks its routine on any joint error.
func (m Mapper[V, W]) Form(j Joint[V, W]) {
	for {
		v, err := j.Get()
		if err != nil {
			return
		}

		w, err := m(j.Ctx(), v)
		if err != nil {
			err = j.Report(err)
			if err != nil {
				return
			}

			continue
		}

		err = j.Put(w)
		if err != nil {
			return
		}
	}
}

var _ Former[int, int] = Mapper[int, int](
	func(context.Context, int) (int, error) { return 0, nil },
)
