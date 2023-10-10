package nextroute

import "fmt"

const (
	BalanceObjectiveModeMax    = "max"
	BalanceObjectiveModeMinMax = "minmax"
)

type BalanceObjectiveMode int

const (
	BalanceObjectiveModeMaxValue BalanceObjectiveMode = iota
	BalanceObjectiveModeMinMaxValue
)

func BalanceObjectiveModeFrom(mode string) (BalanceObjectiveMode, error) {
	switch mode {
	case BalanceObjectiveModeMax:
		return BalanceObjectiveModeMaxValue, nil
	case BalanceObjectiveModeMinMax:
		return BalanceObjectiveModeMinMaxValue, nil
	default:
		return 0, fmt.Errorf("invalid balance objective mode: %s", mode)
	}
}
