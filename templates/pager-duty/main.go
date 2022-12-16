// package main holds the implementation of the pager-duty template.
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	run.Run(solver)
}

// Input for the pager duty scheduling problem. We have
// pager duty users that need to be assigned to days between the schedule start
// date and the schedule end date.
type input struct {
	ScheduleStart time.Time `json:"schedule_start"`
	ScheduleEnd   time.Time `json:"schedule_end"`
	Users         []user    `json:"users"`
}

// Users have a name, id, type, unavailable dates, and preferences.
type user struct {
	Name        string      `json:"name,omitempty"`
	ID          string      `json:"id,omitempty"`
	Type        string      `json:"type,omitempty"`
	Unavailable []time.Time `json:"unavailable,omitempty"`
	Preferences preference  `json:"preferences,omitempty"`
}

// Preference of days.
type preference struct {
	Days []time.Time `json:"days"`
}

// Provide the start, end, user, and timezone of the override to work
// with the PagerDuty API.
type override struct {
	Start    time.Time    `json:"start"`
	End      time.Time    `json:"end"`
	User     assignedUser `json:"user"`
	TimeZone string       `json:"time_zone"`
}

// An assignedUser has a name, id, and type for PagerDuty override.
type assignedUser struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

func solver(input input, opts store.Options) (store.Solver, error) {
	// We start with an empty root pagerDuty store.
	pagerDuty := store.New()

	// Next, we add a starting state to store.
	// We create a domain for each day and initialize each day to all users.
	ndays := int(input.ScheduleEnd.Sub(input.ScheduleStart).Hours()/24) + 1
	users := model.NewDomain(model.NewRange(0, len(input.Users)-1))
	days := store.Repeat(pagerDuty, ndays, users)

	// Next, we create an `assignedDays` variable to keep track of number of
	// days per user. This is so we can balance assignments across users later.
	// We initialize day length to 0 for all users.
	assignedDays := store.NewSlice[int](pagerDuty)
	for range input.Users {
		pagerDuty = pagerDuty.Apply(assignedDays.Append(0))
	}
	// We also want to maximize worker happiness by fulfilling as many
	// preferences as possible.
	happiness := store.NewSlice[int](pagerDuty)
	for range input.Users {
		pagerDuty = pagerDuty.Apply(happiness.Append(0))
	}
	// As a first step, we can now remove the users that are unavailable for
	// certain full days and apply the changes. This results in a new (partial)
	// schedule.

	// We also create maps of date -> preferred users and update store with
	// unavailable dates removed from users.
	preferenceMap := map[int][]int{}

	date := input.ScheduleStart
	dateIndex := 0
	for !date.After(input.ScheduleEnd) {
		for userIndex, user := range input.Users {
			for _, unavailable := range user.Unavailable {
				if date.Equal(unavailable) {
					pagerDuty = pagerDuty.Apply(days.Remove(dateIndex, []int{userIndex}))
				}
			}
			for _, preference := range user.Preferences.Days {
				if date == preference {
					preferenceMap[dateIndex] = append(preferenceMap[dateIndex], userIndex)
				}
			}
		}

		// If no one is available on a day, the problem is infeasible.
		if days.Domain(pagerDuty, dateIndex).Empty() {
			return nil, fmt.Errorf("problem is infeasible")
		}
		assignedUser, dayAssigned := days.Domain(pagerDuty, dateIndex).Value()
		// If there is only 1 person available, assign that person to the day.
		if dayAssigned {
			userAssignedDays := assignedDays.Get(pagerDuty, assignedUser)

			// Add 1 to their day length.
			pagerDuty = pagerDuty.Apply(
				assignedDays.Set(assignedUser, userAssignedDays+1),
			)

			// Add 1 to their happiness score if they preferred to work this day.
			if preferredUsers, ok := preferenceMap[dateIndex]; ok {
				for _, p := range preferredUsers {
					if p == assignedUser {
						pagerDuty = pagerDuty.Apply(
							happiness.Set(p, happiness.Get(pagerDuty, p)+1),
						)
						break
					}
				}
			}
		}

		// Add 1 day to the date.
		date = date.AddDate(0, 0, 1)
		dateIndex++
	}

	// Now we configure how we generate, validate, value, and format stores.
	pagerDuty = pagerDuty.Generate(func(s store.Store) store.Generator {
		// We define the method for generating child states from each parent
		// state. Each time a new store is generated - we attempt to generate more
		// stores.

		// We find the first day with 2 or more users available.
		dayIndex, ok := days.First(s)

		if !ok {
			return nil
		}

		// Create a slice of stores where we'll attempt to assign each of those
		// available users to the day.
		stores := make([]store.Store, days.Domain(s, dayIndex).Len())
		for i, user := range days.Domain(s, dayIndex).Slice() {
			// We increment the day length for the user we want to assign.
			userAssignedDays := assignedDays.Get(s, user)
			userAssignedDays++

			// We assign the user and apply changes for that assignment and the
			// day length increment.
			newStore := s.Apply(
				days.Assign(dayIndex, user),
				assignedDays.Set(user, userAssignedDays),
			)

			// If we were able to assign user to preferred day, increment
			// happiness score and apply changes.
			if preferredUsers, ok := preferenceMap[dayIndex]; ok {
				for _, p := range preferredUsers {
					if p == user {
						newStore = newStore.Apply(
							happiness.Set(p, happiness.Get(pagerDuty, p)+1),
						)
					}
				}
			}

			stores[i] = newStore
		}

		// Use an Eager generator to return all child stores we just created for
		// all users available on this day.
		return store.Eager(stores...)
	}).Validate(func(s store.Store) bool {
		// Next, we define operational validity on the store. Our plan is
		// operationally valid if all days are assigned exactly to one person.
		return days.Singleton(s)
	}).Value(func(s store.Store) int {
		// Now we define our objective value. This is the quantity we try to
		// minimize (or maximize or satisfy). This balances days assigned to each
		// user.

		// Calculate sum of squared assigned days.
		sumSquares := sumSquare(assignedDays.Slice(s))

		// Calculate minimum happiness across users.
		minHappiness := min(happiness.Slice(s))

		// Balance days between users and maximize minimum happiness.
		return sumSquares - minHappiness
	}).Format(format(assignedDays, days, input))

	// If the duration limit is unset, we set it to 10s. You can configure
	// longer solver run times here. For local runs there is no time limitation.
	// If you want to make cloud runs for longer than 5 minutes, please contact:
	// sales@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	return pagerDuty.Minimizer(opts), nil
}

// We define a min helper function to use in our Value calculations.
func min(array []int) int {
	min := math.MaxInt
	for _, value := range array {
		if min > value {
			min = value
		}
	}
	return min
}

// We define a sumSquare helper function to use in our Value calculations.
func sumSquare(array []int) int {
	ss := 0
	for _, value := range array {
		ss += value * value
	}
	return ss
}

// format returns a function to format the solution output.
func format(
	assignedDays store.Slice[int],
	days store.Domains,
	input input,
) func(s store.Store) any {
	return func(s store.Store) any {
		// Next, we define the output format for our schedule.
		// We want to structure our output in a way that the PagerDuty API
		// understands.

		values, ok := days.Values(s)

		if !ok {
			return "No schedule found"
		}
		overrides := []override{}
		for v, nameIndex := range values {
			assignedUser := assignedUser{
				Name: input.Users[nameIndex].Name,
				ID:   input.Users[nameIndex].ID,
				Type: input.Users[nameIndex].Type,
			}
			overrides = append(overrides, override{
				Start:    input.ScheduleStart.AddDate(0, 0, v),
				End:      input.ScheduleStart.AddDate(0, 0, v+1),
				User:     assignedUser,
				TimeZone: "UTC",
			})
		}

		return map[string]any{
			"overrides":       overrides,
			"min_days_worked": min(assignedDays.Slice(s)),
		}
	}
}
