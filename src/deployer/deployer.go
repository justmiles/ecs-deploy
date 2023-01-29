package deployer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/mitchellh/copystructure"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
)

// PerformDeployment initiates an ECS deployment by
//
//	setting desired version in SSM Parameter Store /<env>/<app>/VERSION
//	bumping the image version in task definition
//	registering new task definition with the ECS service
func PerformDeployment(depOpts DeploymentOptions) (s string, err error) {
	var deploymentResults DeploymentResults
	// Set the desired application version
	if !depOpts.DryRun {
		err = setDesiredVersion(depOpts)
		if err != nil {
			return s, err
		}
	}

	var svc *ecs.ECS

	if depOpts.Role != "" {
		creds := stscreds.NewCredentials(sess, depOpts.Role)
		svc = ecs.New(sess, &aws.Config{Credentials: creds})
	} else {
		svc = ecs.New(sess)
	}

	// Get the ECS Service
	dsi := &ecs.DescribeServicesInput{
		Cluster: aws.String(depOpts.Environment),
		Services: []*string{
			aws.String(depOpts.Application),
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

	// Get the ECS service's full task definition
	dtdi := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: dso.Services[0].TaskDefinition,
	}
	dtdo, err := svc.DescribeTaskDefinition(dtdi)
	if err != nil {
		return s, err
	}

	// Deep copy to preserve original container definitions for diff
	copyContainerDefs, err := copystructure.Copy(dtdo.TaskDefinition.ContainerDefinitions)
	if err != nil {
		fmt.Printf("Error performing deep copy of container definitions: %v", err)
		os.Exit(1)
	}
	desiredContainerDefinitions, ok := copyContainerDefs.([]*ecs.ContainerDefinition)
	if !ok {
		fmt.Printf("Error converting interface to ecs.ContainerDefinition: %v", err)
		os.Exit(1)
	}

	// Update only first container definition
	repoAndVersion := strings.Split(*desiredContainerDefinitions[0].Image, ":")
	if len(repoAndVersion) == 1 {
		repoAndVersion = append(repoAndVersion, depOpts.Version)
	} else {
		repoAndVersion[1] = depOpts.Version
	}

	// Update only the first contianer image version - ignore sidecar containers assuming they are defined second, third, and so on.
	*desiredContainerDefinitions[0].Image = strings.Join(repoAndVersion, ":")

	// Register new task definition
	rtdi := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    desiredContainerDefinitions,
		Cpu:                     dtdo.TaskDefinition.Cpu,
		ExecutionRoleArn:        dtdo.TaskDefinition.ExecutionRoleArn,
		Family:                  dtdo.TaskDefinition.Family,
		Memory:                  dtdo.TaskDefinition.Memory,
		NetworkMode:             dtdo.TaskDefinition.NetworkMode,
		PlacementConstraints:    dtdo.TaskDefinition.PlacementConstraints,
		TaskRoleArn:             dtdo.TaskDefinition.TaskRoleArn,
		Volumes:                 dtdo.TaskDefinition.Volumes,
		RequiresCompatibilities: dtdo.TaskDefinition.RequiresCompatibilities,
	}

	if depOpts.RefreshSecrets {
		var ssmClient *ssm.SSM
		if depOpts.Role != "" {
			creds := stscreds.NewCredentials(sess, depOpts.Role)
			ssmClient = ssm.New(sess, &aws.Config{Credentials: creds})
		} else {
			ssmClient = ssm.New(sess)
		}

		pageNum := 0
		containerSecrets := []*ecs.Secret{}
		err := ssmClient.GetParametersByPathPages(&ssm.GetParametersByPathInput{
			Path: aws.String(depOpts.SecretsPrefix),
		},
			func(page *ssm.GetParametersByPathOutput, lastPage bool) bool {
				pageNum++
				for _, v := range page.Parameters {

					ss := strings.Split(*v.Name, "/")
					s := ss[len(ss)-1]

					secret := &ecs.Secret{
						Name:      aws.String(s),
						ValueFrom: v.ARN,
					}
					containerSecrets = append(containerSecrets, secret)
				}
				return pageNum <= 100
			})

		if err != nil {
			fmt.Printf("Error refreshing ssm params: %v", err)
			os.Exit(1)
		}

		desiredContainerDefinitions[0].Secrets = containerSecrets

	}

	// Diff task image version
	diff := NewDiff("container", *dtdo.TaskDefinition.ContainerDefinitions[0].Name)
	for index, currentContainerDef := range dtdo.TaskDefinition.ContainerDefinitions {
		diff.AddChange("image", *currentContainerDef.Image, *desiredContainerDefinitions[index].Image)

		// Diff task secrets
		for _, x := range currentContainerDef.Secrets {
			found := false
			for _, y := range desiredContainerDefinitions[index].Secrets {
				// TODO diff secrets
				if *x.Name == *y.Name {
					diff.AddChange(*x.Name, *x.ValueFrom, *y.ValueFrom)
					found = true
				}
			}

			if !found {
				diff.AddChange(*x.Name, *x.ValueFrom, "")
			}
		}

		for _, y := range desiredContainerDefinitions[index].Secrets {
			found := false
			for _, x := range currentContainerDef.Secrets {
				if *x.Name == *y.Name {
					found = true
				}
			}
			if !found {
				diff.AddChange(*y.Name, "", *y.ValueFrom)
			}
		}

		if len(diff.changes) > 0 {
			fmt.Println(diff)
		}
	}

	if depOpts.DryRun {
		os.Exit(0)
	}

	rtdo, err := svc.RegisterTaskDefinition(rtdi)
	if err != nil {
		return s, err
	}
	// Update the service with the new task definition
	usi := &ecs.UpdateServiceInput{
		Cluster:                 dso.Services[0].ClusterArn,
		DeploymentConfiguration: dso.Services[0].DeploymentConfiguration,
		DesiredCount:            dso.Services[0].DesiredCount,
		ForceNewDeployment:      aws.Bool(true),
		NetworkConfiguration:    dso.Services[0].NetworkConfiguration,
		PlatformVersion:         dso.Services[0].PlatformVersion,
		Service:                 dso.Services[0].ServiceArn,
		TaskDefinition:          rtdo.TaskDefinition.TaskDefinitionArn,
	}
	// If HealthCheckGracePeriodSeconds == 0 (Default), assume that the previous definition did not include a health check.
	if *dso.Services[0].HealthCheckGracePeriodSeconds != 0 {
		usi.HealthCheckGracePeriodSeconds = dso.Services[0].HealthCheckGracePeriodSeconds
	}

	uso, err := svc.UpdateService(usi)
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

// PerformReDeployment initiates an ECS re-deployment
func PerformReDeployment(depOpts DeploymentOptions) (s string, err error) {
	var deploymentResults DeploymentResults

	var svc *ecs.ECS

	if depOpts.Role != "" {
		creds := stscreds.NewCredentials(sess, depOpts.Role)
		svc = ecs.New(sess, &aws.Config{Credentials: creds})
	} else {
		svc = ecs.New(sess)
	}

	// Get the ECS Service
	dsi := &ecs.DescribeServicesInput{
		Cluster: aws.String(depOpts.Environment),
		Services: []*string{
			aws.String(depOpts.Application),
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

	uso, err := svc.UpdateService(&ecs.UpdateServiceInput{
		ForceNewDeployment: aws.Bool(true),
		Cluster:            dso.Services[0].ClusterArn,
		Service:            dso.Services[0].ServiceArn,
	})
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

func WaitForDeployment(depOpts DeploymentOptions) (err error) {

	var svc *ecs.ECS

	if depOpts.Role != "" {
		creds := stscreds.NewCredentials(sess, depOpts.Role)
		svc = ecs.New(sess, &aws.Config{Credentials: creds})
	} else {
		svc = ecs.New(sess)
	}

	err = svc.WaitUntilServicesStableWithContext(aws.BackgroundContext(), &ecs.DescribeServicesInput{
		Cluster: aws.String(depOpts.Environment),
		Services: []*string{
			aws.String(depOpts.Application),
		},
	}, request.WithWaiterMaxAttempts(depOpts.MaxAttempts))

	if err != nil {
		return err
	}

	return nil
}

func setDesiredVersion(depOpts DeploymentOptions) error {
	var svc *ssm.SSM

	if depOpts.Role != "" {
		creds := stscreds.NewCredentials(sess, depOpts.Role)
		svc = ssm.New(sess, &aws.Config{Credentials: creds})
	} else {
		svc = ssm.New(sess)
	}

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
