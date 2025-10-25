package ali

import (
	"github.com/QuantumNous/new-api/dto"
)



const EnableSearchModelSuffix = "-internet"

func requestOpenAI2Ali(request dto.GeneralOpenAIRequest) *dto.GeneralOpenAIRequest {
	if request.TopP >= 1 {
		request.TopP = 0.999
	} else if request.TopP <= 0 {
		request.TopP = 0.001
	}
	return &request
}
