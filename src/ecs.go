package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func performDeployment(depOpts DeploymentOptions) (s string, err error) {
	var deploymentResults DeploymentResults
	svc := ecs.New(sess)

	// Get the ECS Service
	dsi := &ecs.DescribeServicesInput{
		Cluster: aws.String(depOpts.Environment),
		Services: []*string{
			aws.String(fmt.Sprintf("%s-%s", depOpts.Environment, depOpts.Application)),
		},
	}
	dso, err := svc.DescribeServices(dsi)
	if err != nil {
		return s, err
	}

	if len(dso.Failures) > 0 {
		log.Println(dso.Failures)
		return s, fmt.Errorf("unable to find service %s in cluster %s", depOpts.Application, depOpts.Environment)
	}

	log.Println(*dso.Services[0].TaskDefinition)
	log.Printf("%T", dso.Services[0])

	// Get the ECS service's full task definition
	dtdi := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: dso.Services[0].TaskDefinition,
	}
	dtdo, err := svc.DescribeTaskDefinition(dtdi)
	if err != nil {
		return s, err
	}

	for i, containerDefinition := range dtdo.TaskDefinition.ContainerDefinitions {
		repoAndVersion := strings.Split(*containerDefinition.Image, ":")
		repoAndVersion[1] = depOpts.Version
		*containerDefinition.Image = strings.Join(repoAndVersion, ":")
		dtdo.TaskDefinition.ContainerDefinitions[i] = containerDefinition
	}

	// Register new task definition
	rtdi := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: dtdo.TaskDefinition.ContainerDefinitions,
		Cpu:                  dtdo.TaskDefinition.Cpu,
		ExecutionRoleArn:     dtdo.TaskDefinition.ExecutionRoleArn,
		Family:               dtdo.TaskDefinition.Family,
		Memory:               dtdo.TaskDefinition.Memory,
		NetworkMode:          dtdo.TaskDefinition.NetworkMode,
		PlacementConstraints: dtdo.TaskDefinition.PlacementConstraints,
		TaskRoleArn:          dtdo.TaskDefinition.TaskRoleArn,
		Volumes:              dtdo.TaskDefinition.Volumes,
	}
	rtdo, err := svc.RegisterTaskDefinition(rtdi)
	if err != nil {
		return s, err
	}
	// Update the service with thte new task definition
	usi := &ecs.UpdateServiceInput{
		Cluster:                       dso.Services[0].ClusterArn,
		DeploymentConfiguration:       dso.Services[0].DeploymentConfiguration,
		DesiredCount:                  dso.Services[0].DesiredCount,
		ForceNewDeployment:            aws.Bool(true),
		HealthCheckGracePeriodSeconds: dso.Services[0].HealthCheckGracePeriodSeconds,
		NetworkConfiguration:          dso.Services[0].NetworkConfiguration,
		PlatformVersion:               dso.Services[0].PlatformVersion,
		Service:                       dso.Services[0].ServiceArn,
		TaskDefinition:                rtdo.TaskDefinition.TaskDefinitionArn,
	}
	uso, err := svc.UpdateService(usi)
	log.Println(uso)
	if err != nil {
		return s, err
	}
	deploymentResults.SuccessfullyInvoked = true
	deploymentResults.ClusterArn = *uso.Service.ClusterArn
	deploymentResults.ServiceArn = *uso.Service.ServiceArn
	deploymentResults.ServiceName = *uso.Service.ServiceName
	deploymentResults.TaskDefinition = *uso.Service.TaskDefinition

	res, err := json.Marshal(deploymentResults)
	s = string(res)
	return s, err
}

func setDesiredVersion(depOpts DeploymentOptions) error {
	svc := ssm.New(sess)

	input := &ssm.PutParameterInput{
		Name:        aws.String(fmt.Sprintf("/%s/%s/VERSION", depOpts.Environment, depOpts.Application)),
		Overwrite:   aws.Bool(true),
		Type:        aws.String("String"),
		Description: aws.String(depOpts.Description),
		Value:       aws.String(depOpts.Version),
	}
	_, err := svc.PutParameter(input)
	return err
}
