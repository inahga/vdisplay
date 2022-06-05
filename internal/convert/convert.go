package convert

func BGRxToRGBA(pix []byte) {
	if len(pix)%4 != 0 {
		panic("invalid pixel buffer")
	}
	for i := 0; i < len(pix); i += 4 {
		pix[i], pix[i+2] = pix[i+2], pix[i]
	}
}
