package util

func ConvertGradeAverageToString(avg float64) string {
	if avg <= 0.5 {
		return "I"
	}

	if avg > 0.5 && avg <= 1.5 {
		return "C"
	}

	if avg > 1.5 && avg <= 2.5 {
		return "B-"
	}

	if avg > 2.5 && avg <= 3.5 {
		return "B"
	}

	if avg > 3.5 && avg <= 4.5 {
		return "B+"
	}

	if avg > 4.5 && avg <= 5.5 {
		return "A-"
	}

	if avg > 5.5 && avg <= 6 {
		return "A"
	}

	return "error"
}
