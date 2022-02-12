package pcg

// NaiveStore stores information on the network naively by simply placing it on the local node. It generate a UUIS for
// the information and creates an information block and return information uuid
func NaiveStore(overlay *PCG, data string) string {
	uuid := overlay.AddGroup(data)
	return uuid
}

func PCGStore(overlay *PCG, data string) string {
	uuid := overlay.AddGroup(data)
	return uuid
}
