package co2mini

func checksum(plainData []byte) bool {
	return plainData[4] == 0x0d &&
		(uint16(plainData[0])+uint16(plainData[1])+uint16(plainData[2]))&0xff == uint16(plainData[3])
}
