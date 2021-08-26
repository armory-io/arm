package cmd

import "github.com/armory/plank/v4"

var pipelineID string

type PlankMock struct {
}

func (p PlankMock) GetApplication(string, arg1 string) (*plank.Application, error) {
	return &plank.Application{
		Name: string,
	}, nil
}

func (p PlankMock) UpdateApplicationNotifications(plank.NotificationsType, string, string) error {
	return nil
}

func (p PlankMock) GetApplicationNotifications(string, string) (*plank.NotificationsType, error) {
	return nil, nil
}

func (p PlankMock) CreateApplication(*plank.Application, string) error {
	return nil
}

func (p PlankMock) UpdateApplication(plank.Application, string) error {
	return nil
}

func (p PlankMock) GetPipelines(arg0, arg1 string) ([]plank.Pipeline, error) {
	return []plank.Pipeline{{
		ID:          "mock-" + pipelineID + "-id",
		Name:        pipelineID,
		Application: arg0,
	}}, nil
}

func (p PlankMock) DeletePipeline(plank.Pipeline, string) error {
	return nil
}

func (p PlankMock) UpsertPipeline(pipe plank.Pipeline, str, arg2 string) error {
	pipelineID = pipe.Name
	return nil
}

func (p PlankMock) ResyncFiat(string) error {
	return nil
}

func (p PlankMock) ArmoryEndpointsEnabled() bool {
	return false
}

func (p PlankMock) EnableArmoryEndpoints() {
}

func (p PlankMock) UseGateEndpoints() {
}

func (p PlankMock) UseServiceEndpoints() {
}
