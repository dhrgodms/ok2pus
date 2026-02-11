package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Port int    `yaml:"port"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile("config.yaml")
	// yaml 형식의 데이터를 go의 구조체, 맵 같은 변수에 담도록 역직렬화
	// 바이트 슬라이스(data)로 분석하여 정의한 변수의 사용자가 정의한 메모리 주소(&config)에 값 하나씩 채워넣음
	if err != nil {
		fmt.Printf("Parsing Error: %v\n", err)
		return config, err
	}

	err = yaml.Unmarshal(data, &config)

	return config, err
}
