package service_discovery

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"fmt"
)

func RegisteServiceToConsul(consul_address string,service_address string,check_address string){

	config := consulapi.DefaultConfig();
	config.Address=consul_address;

	client,err:=consulapi.NewClient(config);
	if err!=nil{
		logrus.Error(err)
		return;
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "battle_server_10_0_0_5_9092"
	registration.Name = "battle_server"
	registration.Port = 0;
	registration.Tags = []string{"battle_server"}
	registration.Address =service_address;

	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s", check_address);
	check.Timeout = "5s";
	check.Interval = "5s";
	registration.Check = check;
	if err:=client.Agent().ServiceRegister(registration);err!=nil{
		logrus.Error(err);
	}
}