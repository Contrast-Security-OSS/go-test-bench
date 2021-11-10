package common

// ConstParams are the page parameters that will not change during a run. Currently only used with the standard lib.
type ConstParams struct {
	Year      int
	Rulebar   RouteMap
	Addr      string
	Logo      string
	Framework string
}

// Parameters includes ConstParams as well as anything page-specific. Currently only used with the standard lib.
type Parameters struct {
	Name string
	ConstParams
}
