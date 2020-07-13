package cmd

import "github.com/armory/plank/v3"

var pipelineID string

type PlankMock struct {
}

func (p PlankMock) GetApplication(string string) (*plank.Application, error){
	return &plank.Application{
		Name:          string,
		Email:         "",
		Description:   "",
		User:          "",
		DataSources:   nil,
		Permissions:   nil,
		Notifications: nil,
	}, nil
}

func (p PlankMock) UpdateApplicationNotifications(plank.NotificationsType, string) error{
	return nil
}

func (p PlankMock) GetApplicationNotifications(string) (*plank.NotificationsType, error){
	return nil, nil
}

func (p PlankMock) CreateApplication(*plank.Application) error{
	return nil
}

func (p PlankMock) UpdateApplication(plank.Application) error{
	return nil
}

func (p PlankMock) GetPipelines(string string) ([]plank.Pipeline, error){
	return []plank.Pipeline{{
		ID:                   pipelineID,
		Type:                 "",
		Name:                 pipelineID,
		Application:          string,
		Description:          "",
		ExecutionEngine:      "",
		Parallel:             false,
		LimitConcurrent:      false,
		KeepWaitingPipelines: false,
		Stages:               nil,
		Triggers:             nil,
		Parameters:           nil,
		Notifications:        nil,
		ExpectedArtifacts:    nil,
		LastModifiedBy:       "",
		Config:               nil,
		UpdateTs:             "",
		Locked:               nil,
	},}, nil
}

func (p PlankMock) DeletePipeline(plank.Pipeline) error{
	return nil
}

func (p PlankMock) UpsertPipeline(pipe plank.Pipeline, str string) error{
	pipelineID = pipe.Name
	return nil
}

func (p PlankMock) ResyncFiat() error{
	return nil
}

func (p PlankMock) ArmoryEndpointsEnabled() bool{
	return false
}

func (p PlankMock) EnableArmoryEndpoints(){
}
