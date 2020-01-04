package co2mini

func decrypt(encryptedData []byte, key []byte) []byte {
	// phase1
	phase1 := []byte{
		encryptedData[2],
		encryptedData[4],
		encryptedData[0],
		encryptedData[7],
		encryptedData[1],
		encryptedData[6],
		encryptedData[5],
		encryptedData[3],
	}

	// phase2
	phase2 := make([]byte, 8)
	for i, p := range phase1 {
		phase2[i] = p ^ key[i]
	}

	// phase3
	phase3 := make([]byte, 8)
	for i := range phase2 {
		phase3[i] = (phase2[i]>>3 | phase2[(i-1+8)&7]<<5) & 0xff
	}

	// phase4
	phase4 := make([]byte, 8)
	c := []uint16{0x84, 0x47, 0x56, 0xd6, 0x07, 0x93, 0x93, 0x56} // reverse "Htemp99e" in each half-bytes
	for i, p := range phase3 {
		phase4[i] = byte((0x100 + uint16(p) - c[i]) & 0xff)
	}

	return phase4
}
