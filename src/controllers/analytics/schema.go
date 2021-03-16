package analytics

type rowResponseContentCount struct {
	Count uint `json:"count"`
}

type responseContentCount struct {
	Article   rowResponseContentCount `json:"article"`
	Handcraft rowResponseContentCount `json:"handcraft"`
	Culinary  rowResponseContentCount `json:"culinary"`
	Lodging   rowResponseContentCount `json:"lodging"`
	Travel    rowResponseContentCount `json:"travel"`
}
