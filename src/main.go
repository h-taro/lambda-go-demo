package main

import (
	"encoding/csv"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Hoge struct {
	Name string `json: "name"`
	Age  string `json: "age"`
}

func hander(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	session := session.New()
	s3Client := s3.New(session)

	bucketName := "qiitahoge"
	objectKey := "hoge.csv"

	object, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	defer object.Body.Close()
	reader := csv.NewReader(object.Body)
	rows, err := reader.ReadAll()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	var response []Hoge
	for i, v := range rows {
		if i == 0 {
			continue
		}

		name := v[0]
		age := v[1]

		hoge := Hoge{
			Name: name,
			Age:  age,
		}

		response = append(response, hoge)
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(hander)
}
