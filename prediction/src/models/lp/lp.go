package lp

import (
	"main/utils"
	"math"
)

func Elp(predictedPeriod int8, order int, inputVector []float64, pctWageP float64) []float64 {
	var predictVisitorsP []float64
	var a []float64
	var shift int
	var lshift int

	predictionVector := make([]float64, len(inputVector))
	copy(predictionVector, inputVector)

	if len(inputVector) > order+1 {
		for i := 1; i < int(predictedPeriod); i++ {
			if math.Remainder(float64(i), 3) == 0 || i == 1 {
				a, _ = utils.LevinsonDurbin(predictionVector, order)
				shift, lshift = utils.CalculateShift(predictionVector, order)
			}
			addValueP := 0.0
			for lp := 1; lp <= order; lp++ {
				addValueP += predictionVector[len(predictionVector)-lp] * a[lp]
				if len(predictionVector) > shift+lp {
					addValueP += predictionVector[(len(predictionVector)-shift)-lp] * a[lp]
				}
				if len(predictionVector) > lshift+lp {
					addValueP += predictionVector[(len(predictionVector)-lshift)-lp] * a[lp]
				}
			}
			predictVisitorsP = append(predictVisitorsP, math.Ceil(addValueP))
			predictionVector = append(predictionVector, math.Ceil(addValueP))
		}
	}
	return predictVisitorsP
}
