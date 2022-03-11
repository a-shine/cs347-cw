package pcg

func MbToBytes(mb uint64) uint64 {
	return uint64(mb * 1024 * 1024)
}

func MaxStorage(maxMemory uint64) uint64 {
	return maxMemory / uint64(GroupStructSize)
}
