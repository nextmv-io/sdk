package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Output struct {
	Status         string             `json:"status"`
	AssignedShifts []OutputAssignment `json:"assigned_shifts"`
	Value          float64            `json:"value"`
}

type CustomResultStatistics struct {
	NumberOfWorkers int `json:"number_of_workers"`
}

type Options struct {
	OverSupplyPenalty  float64       `json:"over_supply_penalty" default:"1000" usage:"penalty for over-supplying a demand"`
	UnderSupplyPenalty float64       `json:"under_supply_penalty" default:"500" usage:"penalty for over-supplying a demand"`
	MaxHoursPerDay     time.Duration `json:"max_hours_per_day" default:"10h" usage:"maximum number of hours per day"`
	MaxHoursPerWeek    time.Duration `json:"max_hours_per_week" default:"40h" usage:"maximum number of hours per week"`
	MinHoursPerShift   time.Duration `json:"min_hours_per_shift" default:"2h" usage:"minimum number of hours per shift"`
	MaxHoursPerShift   time.Duration `json:"max_hours_per_shift" default:"8h" usage:"maximum number of hours per shift"`
	HoursBetweenShifts time.Duration `json:"hours_between_shifts" default:"8h" usage:"minimum number of hours between shifts"`
}

// create an input struct definition that can read input.json
type Input struct {
	Workers         []Worker         `json:"workers"`
	RequiredWorkers []RequiredWorker `json:"required_workers"`
}

type Worker struct {
	Availability []Availability `json:"availability"`
	ID           int            `json:"id"`
}

type Availability struct {
	Start CustomTime `json:"start"`
	End   CustomTime `json:"end"`
}

type RequiredWorker struct {
	Id    int        `json:"id,omitempty"`
	Start CustomTime `json:"start"`
	End   CustomTime `json:"end"`
	Count int        `json:"count"`
}

func (r RequiredWorker) ID() string {
	return strconv.Itoa(int(r.Id))
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02 15:04:05"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
}

type OutputAssignment struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	DriverID int       `json:"worker_id"`
}

type Assignment struct {
	DemandsCovered []RequiredWorker `json:"demands_covered"`
	Start          time.Time        `json:"start"`
	End            time.Time        `json:"end"`
	Worker         Worker           `json:"worker"`
	Duration       time.Duration    `json:"duration"`
	Id             int              `json:"id"`
}

func (a Assignment) DurationApart(other Assignment) time.Duration {
	if a.Start.After(other.End) {
		return a.Start.Sub(other.End)
	}
	if a.End.Before(other.Start) {
		return other.Start.Sub(a.End)
	}
	return 0
}

func (a Assignment) ID() string {
	return strconv.Itoa(a.Id)
}

func GenerateInput() {
	// Generate input file
	random := rand.New(rand.NewSource(0))

	for fileCount := 0; fileCount < 5; fileCount++ {

		input := Input{}

		// create RequiredDrivers for one week
		// each RequiredDriver entry will cover 30 minutes
		// the timeframe for which entries are generated starts a week ago and ends
		// today
		start := time.Now().Add(-7 * 24 * time.Hour).Round(time.Minute * 30)
		end := time.Now().Round(time.Minute * 30)

		for start.Before(end) {
			input.RequiredWorkers = append(input.RequiredWorkers, RequiredWorker{
				Start: CustomTime{start},
				End:   CustomTime{start.Add(30 * time.Minute)},
				Count: random.Intn(20),
			})
			start = start.Add(30 * time.Minute)
		}

		// create Drivers
		// each driver has a unique id and a random performance rating between 1 and
		// 5
		// each driver has a random number of availabilities. This number of
		// availabilities varies between 1 and 10
		// each availability is between 1 and 12 hours long
		// each availability starts between 1 and 8 days ago
		for i := 0; i < 100; i++ {
			availabilityCount := random.Intn(10) + 1
			availabilities := []Availability{}
			for j := 0; j < availabilityCount; j++ {
				availability := Availability{
					Start: CustomTime{time.Now().Add(-time.Duration(random.Intn(8*24)) * time.Hour).Round(time.Minute * 30)},
				}
				availability.End = CustomTime{
					availability.Start.Add(time.Duration(random.Intn(12)) * time.Hour).Round(time.Minute * 30),
				}
				availabilities = append(availabilities, availability)
			}

			// order availabilities by start time
			for j := 0; j < len(availabilities); j++ {
				for k := j + 1; k < len(availabilities); k++ {
					if availabilities[j].Start.After(availabilities[k].Start.Time) {
						availabilities[j], availabilities[k] = availabilities[k], availabilities[j]
					}
				}
			}

			input.Workers = append(input.Workers, Worker{
				Availability: availabilities,
				ID:           i,
			})
		}

		// write input to json file
		filename := fmt.Sprintf("input_%d.json", fileCount)
		f, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewEncoder(f).Encode(input)
		if err != nil {
			log.Fatal(err)
		}
	}
}
