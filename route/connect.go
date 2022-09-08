package route

import (
	"sync"

	"github.com/nextmv-io/sdk/plugin"
)

const slug = "sdk"

var connected bool

var mtx sync.Mutex

func connect() {
	if connected {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if connected {
		return
	}
	connected = true

	// Router
	plugin.Connect(slug, "RouteNewRouter", &newRouterFunc)

	// Options
	plugin.Connect(slug, "RouteStarts", &startsFunc)
	plugin.Connect(slug, "RouteEnds", &endsFunc)
	plugin.Connect(slug, "RouteCapacity", &capacityFunc)
	plugin.Connect(slug, "RouteInitializationCosts", &initializationCostsFunc)
	plugin.Connect(slug, "RoutePrecedence", &precedenceFunc)
	plugin.Connect(slug, "RouteServices", &servicesFunc)
	plugin.Connect(slug, "RouteShifts", &shiftsFunc)
	plugin.Connect(slug, "RouteWindows", &windowsFunc)
	plugin.Connect(slug, "RouteUnassigned", &unassignedFunc)
	plugin.Connect(slug, "RouteBacklogs", &backlogsFunc)
	plugin.Connect(slug, "RouteMinimize", &minimizeFunc)
	plugin.Connect(slug, "RouteMaximize", &maximizeFunc)
	plugin.Connect(slug, "RouteLimits", &limitsFunc)
	plugin.Connect(slug, "RouteLimitDistances", &limitDistancesFunc)
	plugin.Connect(slug, "RouteLimitDurations", &limitDurationsFunc)
	plugin.Connect(slug, "RouteGrouper", &grouperFunc)
	plugin.Connect(slug, "RouteValueFunctionMeasures", &valueFunctionMeasuresFunc)
	plugin.Connect(slug, "RouteTravelTimeMeasures", &travelTimeMeasuresFunc)
	plugin.Connect(slug, "RouteAttribute", &attributeFunc)
	plugin.Connect(slug, "RouteThreads", &threadsFunc)
	plugin.Connect(slug, "RouteAlternates", &alternatesFunc)
	plugin.Connect(slug, "RouteVelocities", &velocitiesFunc)
	plugin.Connect(slug, "RouteServiceGroups", &serviceGroupsFunc)
	plugin.Connect(slug, "RouteSelector", &selectorFunc)
	plugin.Connect(slug, "RouteUpdate", &updateFunc)
	plugin.Connect(slug, "RouteFilterWithRoute", &filterWithRouteFunc)
	plugin.Connect(slug, "RouteSorter", &sorterFunc)
	plugin.Connect(slug, "RouteConstraint", &constraintFunc)
	plugin.Connect(slug, "RouteFilter", &filterFunc)

	// measures
	plugin.Connect(slug, "RouteHaversineByPoint", &haversineByPointFunc)
	plugin.Connect(slug, "RouteIndexed", &indexedFunc)
	plugin.Connect(slug, "RouteConstantByPoint", &constantByPointFunc)
	plugin.Connect(slug, "RouteConstant", &constantFunc)
	plugin.Connect(slug, "RouteOverride", &overrideFunc)
	plugin.Connect(slug, "RouteScale", &scaleFunc)
	plugin.Connect(slug, "RouteLocation", &locationFunc)
}
