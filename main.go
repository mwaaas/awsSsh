package main

import (
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"fmt"
	"os"
	"strconv"
	"os/exec"
	"github.com/urfave/cli"
	"log"
)


func main(){
	var stack, environmentName string

	supportedStacks := map[string]bool {
		"eb": true,
		"ecs": true,
	}

	app := cli.NewApp()

	app.Name = "awsSsh"

	app.Usage = "ssh to eb or ec2 cluster instance"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "stack",
			Value:       "eb",
			Usage:       "stack either eb or ecs defaults to eb",
			Destination: &stack,
		},
		cli.StringFlag{
			Name: "env",
			Value: "",
			Usage: "environment or cluster name depending of" +
				" stack that you selected. (Required)",
			Destination: &environmentName,
		},
	}

	app.Action = func(c *cli.Context) error {
		errLogger := log.New(os.Stderr, "", 0)
		if environmentName == "" {
			errLogger.Println("env name is required")
			os.Exit(1)
		}
		if ! supportedStacks[stack]{
			errLogger.Println("Stack not supported")
			os.Exit(1)
		}

		ssh(stack, environmentName)
		return nil
	}

	app.Run(os.Args)
}

func getEbInstances(client client.ConfigProvider, env_name string) ([]*elasticbeanstalk.Instance, error){
	eb := elasticbeanstalk.New(client)
	envResources, err := eb.DescribeEnvironmentResources(
		&elasticbeanstalk.DescribeEnvironmentResourcesInput{
			EnvironmentName: aws.String(env_name),
		},
	)
	if err != nil {
		return nil, err
	}else{
		instanceList := envResources.EnvironmentResources.Instances
		return instanceList, nil
	}
}

func getEcsInstances(client client.ConfigProvider, cluster_name string) ([]*elasticbeanstalk.Instance, error) {
	svc := ecs.New(client)

	listContainerInstances, err := svc.ListContainerInstances(
		&ecs.ListContainerInstancesInput{
			Cluster: aws.String(cluster_name),
		},
	)
	if err != nil {
		return nil, err
	}

	result, err := svc.DescribeContainerInstances(
		&ecs.DescribeContainerInstancesInput{
			Cluster: aws.String(cluster_name),
			ContainerInstances: listContainerInstances.ContainerInstanceArns,
		},
	)
	if err != nil {
		return nil, err
	}

	var containerInstances = make([]*elasticbeanstalk.Instance,
		len(result.ContainerInstances))
	for i, instances := range result.ContainerInstances{
		containerInstances[i] = &elasticbeanstalk.Instance{
			Id: instances.Ec2InstanceId,
		}
	}
	return containerInstances, nil
}

func getEc2Ip(client client.ConfigProvider, instanceId string) (string, error){
	ec2Service := ec2.New(client)
	instancesId := []string{instanceId}
	instanceOut, err := ec2Service.DescribeInstances(&ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(instancesId)})
	if err == nil{
		return aws.StringValue(instanceOut.Reservations[0].Instances[0].PublicIpAddress), nil
	}else {
		return "", err
	}
}

func exeCommand(program string, args ...string) {
	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin;
	cmd.Stdout = os.Stdout;
	cmd.Stderr = os.Stderr;
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func ssh(stack, env string)  {
	awsClient := session.Must(session.NewSession())
	var instanceList []*elasticbeanstalk.Instance
	var err error
	switch stack {
	case "eb":
		instanceList, err = getEbInstances(awsClient,  env)
	case "ecs":
		//instanceList, err = getEcsInstances(awsClient,  "awseb-herme-sms-transport-production-ku87knjztk")
		instanceList, err = getEcsInstances(awsClient,  "awseb-hermesSmsProduction-mfharrfyru")
	default:
		println("Unsupported stack")
		os.Exit(1)
	}

	if err != nil{
		println(err.Error())
		os.Exit(1)
	}

	for i := 0; i < len(instanceList);{
		if i == 0{
			println("Choose one of the following instance")
		}
		println(fmt.Sprintf("%d.", i), aws.StringValue(instanceList[i].Id))
		i = i + 1
	}
	var input string
	fmt.Scanln(&input)

	instanceIndex, err := strconv.Atoi(input)
	if err != nil{
		println(err.Error())
		os.Exit(1)
	}else if instanceIndex >= len(instanceList) {
		println("You have chosen invalid choice")
		os.Exit(3)
	}
	publicIp, err := getEc2Ip(awsClient, aws.StringValue(instanceList[instanceIndex].Id))

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	exeCommand("ssh", fmt.Sprintf("ec2-user@%s", publicIp))

}
