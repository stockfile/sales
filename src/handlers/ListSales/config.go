package main

import "os"

type Config struct {
	Region    string
	TableName string
	IndexName string
}

var config *Config

func init() {
	config = &Config{
		Region:    os.Getenv("AWS_REGION"),
		TableName: os.Getenv("TABLE_NAME"),
		IndexName: os.Getenv("STORE_INDEX_NAME"),
	}
}

func GetConfig() *Config {
	return config
}
