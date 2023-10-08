package nextroute

import "fmt"

const (
	BalanceObjectiveModeMax    = "max"
	BalanceObjectiveModeMinMax = "minmax"
)

type BalanceObjectiveMode int

const (
	balanceObjectiveModeMax BalanceObjectiveMode = iota
	balanceObjectiveModeMinMax
)

func BalanceObjectiveModeFrom(mode string) (BalanceObjectiveMode, error) {
	switch mode {
	case BalanceObjectiveModeMax:
		return balanceObjectiveModeMax, nil
	case BalanceObjectiveModeMinMax:
		return balanceObjectiveModeMinMax, nil
	default:
		return 0, fmt.Errorf("invalid balance objective mode: %s", mode)
	}
}
