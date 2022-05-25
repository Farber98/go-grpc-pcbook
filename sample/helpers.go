package sample

import "math/rand"

func randomStringFromSet(args ...string) string {
	return args[rand.Intn(len(args))]
}

func randomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomBool() bool {
	return rand.Intn(2) == 1
}
