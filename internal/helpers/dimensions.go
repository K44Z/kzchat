package helpers

func ComputeSideWidth(w int) int {
	minWidth := 10
	percent := 0.16

	calculatedWidth := int(float64(w) * percent)
	if calculatedWidth < minWidth {
		return minWidth
	}
	maxWidth := 20
	if calculatedWidth > maxWidth {
		return maxWidth
	}

	return calculatedWidth
}

func ComputeChatWidth(w, l, r int) int {
	if w < 100 {
		padding := 6
		return w - padding
	}
	minWidth := 40
	padding := 7
	calcWidth := w - l - r - padding

	if calcWidth < minWidth && w > int(float64(minWidth)*1.5) {
		return minWidth
	}

	return calcWidth
}

func ComputeContentHeight(h int) int {
	minHeight := 10
	inputRatio := 4
	if h < 24 {
		inputRatio = 3
	}
	calcHeight := h - (h / inputRatio)
	if calcHeight < minHeight {
		return minHeight
	}
	return calcHeight
}
