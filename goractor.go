package goractor

type Result interface{}
type Context interface{}

type Step interface {
	Call(ctx Context, prevResult Result) (Result, error)
	Rollback(ctx Context, lastResult Result) error
}

type Organizer struct {
	steps  []Step
	errors []error
}

func NewOrganizer() *Organizer {
	return &Organizer{
		steps:  []Step{},
		errors: []error{},
	}
}

func (o *Organizer) AddStep(step Step) {
	o.steps = append(o.steps, step)
}

func (o *Organizer) Call(ctx Context) Result {
	var (
		result Result
		err    error
	)

	for i, step := range o.steps {
		result, err = step.Call(ctx, result)

		if err != nil {
			o.addError(err)
			o.rollbackTo(i, ctx, result)
		}
	}

	return result
}

func (o *Organizer) rollbackTo(stepIndex int, ctx Context, lastResult Result) {
	for k := 0; k <= stepIndex; k++ {
		if err := o.steps[k].Rollback(ctx, lastResult); err != nil {
			o.addError(err)
		}
	}
}

func (o *Organizer) Errors() []error {
	return o.errors
}

func (o *Organizer) addError(err error) {
	o.errors = append(o.errors, err)
}
