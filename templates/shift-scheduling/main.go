// package main holds the implementation of the shift-scheduling template.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	_, err := run.Run(solver)
	if err != nil {
		log.Fatal(err)
	}
}

// schedulingProblem describes a variant of the "Nurse scheduling problem".
// There are a number of `Days` each having three shifts: morning, day, night.
// Each shift needs a worker, but not all workers need to be assigned to a
// shift. Workers can state preferences regarding their preferred shift type.
// Workers can also be unavailable for certain full days.
type schedulingProblem struct {
	Days        int           `json:"days"`
	Workers     int           `json:"workers"`
	Preferences preference    `json:"preferences"`
	Unavailable []unavailable `json:"unavailable"`
}
type preference struct {
	ShiftType []shiftTypePreference `json:"shift_type"`
}
type shiftTypePreference struct {
	Worker workerID  `json:"worker"`
	Type   shiftType `json:"type"`
}
type unavailable struct {
	Worker workerID `json:"worker"`
	Days   []int    `json:"days"`
}
type shiftType string

type workerID int

func solver(input schedulingProblem, opts store.Options) (store.Solver, error) {
	// We start with an empty schedule.
	schedule := store.New()

	// We define some helper variables that are useful later.
	nShifts := input.Days * 3
	workers := model.NewDomain(model.NewRange(1, input.Workers))
	typeMap := [3]shiftType{"morning", "day", "night"}
	prefShiftType := map[int]shiftType{}
	for _, v := range input.Preferences.ShiftType {
		prefShiftType[int(v.Worker)] = v.Type
	}

	// We track the following variables:
	// shifts holds one domain for each shift (3 per day in order morning, day,
	// night) initially it holds all workers available.
	shifts := store.Repeat(schedule, nShifts, workers)

	// The number of workers active in a schedule.
	workerCount := store.NewVar(schedule, 0)

	// We also want to maximize worker happiness by fulfilling as many
	// preferences as possible.
	happiness := store.NewVar(schedule, 0)

	// As a first step, we can now remove the workers that are unavailable for
	// certain full days. This results as a new (partial) schedule.
	changes := []store.Change{}

	for _, v := range input.Unavailable {
		for _, d := range v.Days {
			morning := (d - 1) * 3
			day := morning + 1
			night := day + 1
			changes = append(changes,
				shifts.Remove(morning, []int{int(v.Worker)}),
				shifts.Remove(day, []int{int(v.Worker)}),
				shifts.Remove(night, []int{int(v.Worker)}),
			)
		}
	}
	schedule = schedule.Apply(changes...)

	// Now we define our objective value. This is the quantity we try to
	// minimize. It is the weighted sum of workers and happiness. Happiness is
	// negative as we want to maximize it while minimizing the worker count.
	workerCountImportance := 10
	schedule = schedule.Value(func(s store.Store) int {
		return workerCountImportance*workerCount.Get(s) - happiness.Get(s)
	})

	// Now we define our constraint containing a propagator. This will remove
	// workers from domains based upon the workers in other domains.
	constraint := newBreakConstraint(shifts, 2)

	// We invoke the propagator to propagate the initial values.
	schedule.Propagate(constraint.propagate)

	if constraint.isProvenInfeasible(schedule) {
		return nil, fmt.Errorf("input can not be solved")
	}

	// Now the most important part: given a store, how can we create new
	// stores that bring us closer to a valid schedule.
	schedule = schedule.Generate(func(s store.Store) store.Generator {
		// We take the first shift for which a worker has not been yet assigned.
		// Then we generate for each worker that can be assigned to a shift
		// i a new store.
		i, ok := shifts.First(s)
		workers := shifts.Domain(s, i).Slice()

		return store.Lazy(
			func() bool {
				for j := 0; j < nShifts; j++ {
					if shifts.Domain(s, j).Empty() {
						return false
					}
				}
				return ok && len(workers) > 0
			},
			func() store.Store {
				// Pop a new worker from the list.
				worker := workers[0]
				// Note that workers is a pointer to the variable defined above.
				workers = workers[1:]

				// First we assign the worker to the current shift.
				changes := []store.Change{
					shifts.AtLeast(i, worker),
					shifts.AtMost(i, worker),
				}

				// Update happiness (i.e. are the preferences considered).
				shiftType := typeMap[i%3]
				if pref, ok := prefShiftType[worker]; ok {
					if pref == shiftType {
						changes = append(changes, happiness.Set(happiness.Get(s)+1))
					}
				} else {
					// In case the worker has no preferences we assume they are
					// happy.
					changes = append(changes, happiness.Set(happiness.Get(s)+1))
				}

				s = s.Apply(changes...)
				s = s.Propagate(constraint.propagate)

				// At last we update the worker count in case some future shifts
				// are fixed now.
				workerMap := map[int]bool{}
				for j := 0; j < nShifts; j++ {
					if value, ok := shifts.Domain(s, j).Value(); ok {
						workerMap[value] = true
					}
				}

				return s.Apply(workerCount.Set(len(workerMap)))
			})
	}).Validate(func(s store.Store) bool {
		// The store is operationally valid if exactly one worker is assigned
		// to a shift and no constraints are violated.
		if !shifts.Singleton(s) {
			return false
		}

		return !constraint.isProvenInfeasible(s)
	}).Bound(func(s store.Store) store.Bounds {
		// In order to make the search more efficient we define bounds.
		// Given a partial schedule, `Bounds` returns a lower and an upper bound
		// around the final schedule value once all workers have been assigned.
		// As a rule of thumb: the better the bounds, the faster the search.

		// The way we create new stores is that we never remove
		// workers from a fixed shift. Thus a trivial bound is
		// to assign a new worker to each unassigned shift.
		// Likewise an upper bound on happiness is that each new worker
		// is also happy.
		wcUpper := workerCount.Get(s)
		happinessUpper := happiness.Get(s)
		for i := 0; i < nShifts; i++ {
			d := shifts.Domain(s, i)
			if _, ok := d.Value(); !ok {
				wcUpper++
				happinessUpper++
			}
		}

		// For the lower bound we simply take the current worker count.
		// As we know we need at least distinct 3 workers per day, we can
		// take the maximum.
		wcLower := workerCount.Get(s)
		if wcLower < 3 {
			wcLower = 3
		}

		// To find a lower bound on happiness we assume that any shift that has
		// potential workers who have preferences cannot meet these preferences.
		happinessLower := happiness.Get(s)
		for i := 0; i < nShifts; i++ {
			d := shifts.Domain(s, i)
			increase := 1
			if d.Len() <= 1 {
				continue
			}
			for _, v := range d.Slice() {
				if _, ok := prefShiftType[v]; ok {
					// there is one worker who has a preference
					increase = 0
					break
				}
			}
			happinessLower += increase
		}
		return store.Bounds{
			Lower: workerCountImportance*wcLower - happinessUpper,
			Upper: workerCountImportance*wcUpper - happinessLower,
		}
	}).Format(format(nShifts, shifts, typeMap, happiness, workerCount))
	// A duration limit of 0 is treated as infinity. For cloud runs you need to
	// set an explicit duration limit which is why it is currently set to 10s
	// here in case no duration limit is set. For local runs there is no time
	// limitation. If you want to make cloud runs for longer than 5 minutes,
	// please contact: support@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	return schedule.Minimizer(opts), nil
}

// output defines the output format of a store.
type output struct {
	// Shifts describe the assigned workers to days and shift type.
	Shifts []shift `json:"shifts"`
	// Happiness is the number of times a preference of a worker has been met.
	// If a worker has no preferences, then it is assumed they are happy.
	Happiness int `json:"happiness"`
	// Workers shows the number of total workers assigned.
	Workers int `json:"workers"`
}

type shift struct {
	Day    int       `json:"day"`
	Worker workerID  `json:"worker"`
	Type   shiftType `json:"type"`
}

type constraint interface {
	propagate(store.Store) []store.Change
	isProvenInfeasible(s store.Store) bool
}

type breakConstraint struct {
	shifts        store.Domains
	breakDuration int
}

// newBreakConstraint create returns a new breakConstraint.
func newBreakConstraint(shifts store.Domains, breakDuration int) constraint {
	return &breakConstraint{
		shifts:        shifts,
		breakDuration: breakDuration,
	}
}

// Propagate removes values from the neighboring domains if a domain
// has been assigned to single value. That is, when a worker has been assigned
// to a shift it can no longer be in the two shifts before and after.
func (t *breakConstraint) propagate(s store.Store) []store.Change {
	changes := make([]store.Change, 0)

	for i := 0; i < t.shifts.Len(s); i++ {
		domain := t.shifts.Domain(s, i)

		if v, singleton := domain.Value(); singleton {
			for j := 1; j <= t.breakDuration; j++ {
				if i-j > 0 {
					changes = t.remove(s, i-j, v, changes)
				}
				if i+j < t.shifts.Len(s) {
					changes = t.remove(s, i+j, v, changes)
				}
			}
		}
	}

	return changes
}

// isProvenInfeasible returns true if the values in the domains can never
// result in an acceptable solution assuming only more values can be removed
// from the domains if the search would continue. That is, if we see two shifts
// with the same worker assigned in any sequence of 3 shifts it is proven to
// be infeasible. Also, if we see a shift for which there are no more values
// in the domain it is proven infeasible.
func (t *breakConstraint) isProvenInfeasible(s store.Store) bool {
	for i := 0; i < t.shifts.Len(s); i++ {
		domain := t.shifts.Domain(s, i)
		// being empty can never result in a solution if we assume we can only
		// remove more values (which there aren't any)
		if domain.Empty() {
			return true
		}
		if w1, singleton := domain.Value(); singleton {
			for j := 1; j <= t.breakDuration; j++ {
				if i-j > 0 && t.checkFeasible(s, w1, i-j) {
					return true
				}
				if i+j < t.shifts.Len(s) && t.checkFeasible(s, w1, i+j) {
					return true
				}
			}
		}
	}
	return false
}

func (t *breakConstraint) remove(
	s store.Store,
	domain,
	value int,
	changes []store.Change,
) []store.Change {
	if t.shifts.Domain(s, domain).Contains(value) {
		return append(changes, t.shifts.Remove(domain, []int{value}))
	}
	return changes
}

func (t *breakConstraint) checkFeasible(s store.Store, w1, domain int) bool {
	if w2, isSingleton := t.shifts.Domain(s, domain).Value(); isSingleton {
		if w1 == w2 {
			return true
		}
	}
	return false
}

// format returns a function to format the solution output.
func format(
	nShifts int,
	shifts store.Domains,
	typeMap [3]shiftType,
	happiness store.Var[int],
	workerCount store.Var[int],
) func(s store.Store) any {
	return func(s store.Store) any {
		outputShifts := make([]shift, 0, nShifts)
		for i, v := range shifts.Slices(s) {
			day := i/3 + 1
			shiftType := typeMap[i%3]
			outputShifts = append(outputShifts, shift{
				Worker: workerID(v[0]),
				Day:    day,
				Type:   shiftType,
			})
		}
		return output{
			Shifts:    outputShifts,
			Happiness: happiness.Get(s),
			Workers:   workerCount.Get(s),
		}
	}
}
