package nextroute

import "fmt"

const (
	BalanceObjectiveModeMax              = "max"
	BalanceObjectiveModeMinMax           = "minmax"
	BalanceObjectiveModeTargetMin        = "targetmin"
	BalanceObjectiveModeTargetMinSquared = "targetminsquared"
)

type BalanceObjectiveMode int

const (
	BalanceObjectiveModeMaxValue BalanceObjectiveMode = iota
	BalanceObjectiveModeMinMaxValue
	BalanceObjectiveModeTargetMinValue
	BalanceObjectiveModeTargetMinSquaredValue
)

func BalanceObjectiveModeFrom(mode string) (BalanceObjectiveMode, error) {
	switch mode {
	case BalanceObjectiveModeMax:
		return BalanceObjectiveModeMaxValue, nil
	case BalanceObjectiveModeMinMax:
		return BalanceObjectiveModeMinMaxValue, nil
	case BalanceObjectiveModeTargetMin:
		return BalanceObjectiveModeTargetMinValue, nil
	case BalanceObjectiveModeTargetMinSquared:
		return BalanceObjectiveModeTargetMinSquaredValue, nil
	default:
		return 0, fmt.Errorf("invalid balance objective mode: %s", mode)
	}
}
