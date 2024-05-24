package cbhmm

import (
	"github.com/gonum/matrix/mat64"
	"log"
	"main/utils"
	"math"
)

// possible sequences of events
var transitionFromOrder = []string{
	"SS",
	"SR",
	"SI",
}

var transitionFromNoOrder = []string{
	"RS",
	"RR",
	"RI",
}

var transitionFromMaybeOrder = []string{
	"IS",
	"IR",
	"II",
}

func Cbhmm(iteration int, transitionMatrix mat64.Matrix, emissionMatrix mat64.Matrix, states []string) ([]string, float64, string) {
	// initial state
	startState := states[0]
	stateList := []string{}
	stateList = append(stateList, startState)

	i := 0

	// calculate the probability of the stateList
	prob := 1.0

	for i < iteration {
		if startState == states[2] {
			orderNextProbabilityDistribution := utils.MatrixByRow(transitionMatrix, 0)

			randomChoice := utils.SampleSet(orderNextProbabilityDistribution)
			change := transitionFromOrder[randomChoice]

			if change == "SS" {
				prob = prob * emissionMatrix.At(0, 2)
				stateList = append(stateList, states[2])
			} else if change == "SR" {
				prob = prob * emissionMatrix.At(0, 0)
				startState = states[0]
				stateList = append(stateList, states[0])
			} else {
				prob = prob * emissionMatrix.At(0, 1)
				startState = states[1]
				stateList = append(stateList, states[1])
			}

		} else if startState == states[2] {
			runNextProbabilityDistribution := utils.MatrixByRow(transitionMatrix, 1)

			randomChoice := utils.SampleSet(runNextProbabilityDistribution)
			change := transitionFromNoOrder[randomChoice]

			if change == "RR" {
				prob = prob * emissionMatrix.At(0, 0)
				stateList = append(stateList, states[0])
			} else if change == "RS" {
				prob = prob * emissionMatrix.At(0, 2)
				startState = states[0]
				stateList = append(stateList, states[2])
			} else {
				prob = prob * emissionMatrix.At(0, 1)
				startState = states[1]
				stateList = append(stateList, states[1])
			}

		} else {
			orderNextProbabilityDistribution := utils.MatrixByRow(transitionMatrix, 2)

			randomChoice := utils.SampleSet(orderNextProbabilityDistribution)
			change := transitionFromMaybeOrder[randomChoice]

			if change == "II" {
				prob = prob * emissionMatrix.At(0, 1)
				stateList = append(stateList, states[1])
			} else if change == "IS" {
				prob = prob * emissionMatrix.At(0, 2)
				startState = states[2]
				stateList = append(stateList, states[2])
			} else {
				prob = prob * emissionMatrix.At(0, 0)
				startState = states[0]
				stateList = append(stateList, states[0])
			}
		}
		i = i + 1
	}
	return stateList, prob, startState
}

func Vendor(beta float64, gamma float64, delta float64, Px float64, Py float64, Qx float64, Qy float64, storeId string) float64 {
	if (beta + gamma + delta) > 1 {
		log.Println("Vendor condition for storeId: " + storeId + " not pass")
	}
	return (beta * utils.Heaviside(Px-Py)) + (gamma * utils.Heaviside(Qx-Qy)) + (delta * utils.Heaviside(Px-Py) * utils.Heaviside(Qx-Qy))
}

func Loyalty(satisfaction float64, trust float64, percieved float64, a float64, b float64, c float64, d float64, e float64, storeId string) float64 {
	if (c + a + e) > 1 {
		log.Println("Loyalty first weights condition for storeId: " + storeId + " not pass")
	}

	if (b + d) > 1 {
		log.Println("Loyalty second weights condition for storeId: " + storeId + " not pass")
	}

	var pF = c*satisfaction + a*trust + e*percieved
	var pC = b*trust + d*satisfaction

	return math.Min((pF+pC)/pC, 1)
}

func Psychology(averageQ float64, averageP float64, Pj float64, Pi float64, Qj float64, Qi float64) float64 {
	productsIndexes := math.Max(Pi, averageP) / math.Max(Pj, averageP)
	qualityIndexes := math.Max(Qi, averageQ) / math.Max(Qj, averageQ)
	return math.Min(productsIndexes+qualityIndexes, 1)
}

func Viterbi(O []int, S []int, Tm mat64.Matrix, Em mat64.Matrix) []int {
	trellis := make([][]float64, len(S))
	pointers := make([][]int, len(S))
	bestPath := make([]int, len(S))
	var k int

	for i := 0; i < len(O); i++ {
		trellis[i] = make([]float64, 3)
	}

	for i := 0; i < len(O); i++ {
		pointers[i] = make([]int, 3)
	}

	for i := 0; i < len(S); i++ {
		trellis[i][0] = float64(O[i]) * Em.At(i, O[0])
	}
	for i := 1; i < len(O)-1; i++ {
		for j := 0; j < len(S)-1; j++ {
			argm := make([]float64, len(O))
			for c := 0; c < len(trellis); c++ {
				argm[c] = trellis[c][i-1] * Tm.At(c, j) * Em.At(j, i)
			}
			k := utils.Argmax(argm)
			trellis[j][i] = trellis[k][i-1] * Tm.At(k, j) * Em.At(j, i)
			pointers[j][i] = k
		}
	}

	k = utils.Argmax(trellis[k])
	for o := 1; o < len(O)-1; o++ {
		bestPath = append(bestPath, S[k])
		k = pointers[k][o]
	}

	return bestPath
}
