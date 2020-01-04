package co2mini

func decrypt(encryptedData []byte, key []byte) []byte {
	// phase1
	phase1 := make([]byte, 8)
	for i, v := range []int{2, 4, 0, 7, 1, 6, 5, 3} {
		phase1[i] = encryptedData[v]
	}

	// phase2
	phase2 := make([]byte, 8)
	for i := range phase1 {
		phase2[i] = phase1[i] ^ key[i]
	}

	// phase3
	phase3 := make([]byte, 8)
	for i := range phase2 {
		phase3[i] = (phase2[i]>>3 | phase2[(i-1+8)&7]<<5) & 0xff
	}

	// phase4
	phase4 := make([]byte, 8)
	for i, v := range []uint16{0x84, 0x47, 0x56, 0xd6, 0x07, 0x93, 0x93, 0x56} { // reverse "Htemp99e" in each half-bytes
		phase4[i] = byte((0x100 + uint16(phase3[i]) - v) & 0xff)
	}

	return phase4
}
