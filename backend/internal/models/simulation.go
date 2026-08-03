package models

type SimulationRequest struct{
	Count int `json:"count" binding:"required,min=1,max=1000"`
}