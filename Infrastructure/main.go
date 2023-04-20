package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/kinesis"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {

        // Create an S3 bucket to store the Kinesis data
        s3Bucket, err := s3.NewBucket(ctx, "mydatalake", nil)
        if err != nil {
            return err
        }

        // Create a Kinesis Data Stream
        dataStream, err := kinesis.NewStream(ctx, "kinesisDataStream", &kinesis.StreamArgs{
            ShardCount: pulumi.Int(1),
        })
        if err != nil {
            return err
        }

        // Create a Lambda IAM role
        lambdaRole, err := iam.NewRole(ctx, "dataTransformLambdaRole", &iam.RoleArgs{
            AssumeRolePolicy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Action": "sts:AssumeRole",
                        "Principal": {
                            "Service": "lambda.amazonaws.com"
                        },
                        "Effect": "Allow",
                        "Sid": ""
                    }
                ]
            }`),
        })
        if err != nil {
            return err
        }

        // Create a Lambda function for data transformation
        dataTransformLambda, err := lambda.NewFunction(ctx, "dataTransformLambda", &lambda.FunctionArgs{
            Runtime: lambda.RuntimeGo1dx,
            Code:    pulumi.NewFileArchive("./lambda/bin/lambda_function.zip"),
            Handler: pulumi.String("main"),
            Role:    lambdaRole.Arn,
        })
        if err != nil {
            return err
        }

        // Create a Kinesis Firehose IAM role
        firehoseRole, err := iam.NewRole(ctx, "firehoseDeliveryStreamRole", &iam.RoleArgs{
            AssumeRolePolicy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Action": "sts:AssumeRole",
                        "Principal": {
                            "Service": "firehose.amazonaws.com"
                        },
                        "Effect": "Allow",
                        "Sid": ""
                    }
                ]
            }`),
        })
        if err != nil {
            return err
        }

        // Create a Kinesis Firehose Delivery Stream with data transformation Lambda
        firehoseStream, err := kinesis.NewFirehoseDeliveryStream(ctx, "firehoseDeliveryStream", &kinesis.FirehoseDeliveryStreamArgs{
            Destination: pulumi.String("extended_s3"),
            ExtendedS3Configuration: kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationArgs{
                BucketArn: s3Bucket.Arn,
                RoleArn:   firehoseRole.Arn,
                Prefix:    pulumi.String("stream_data/"),
            },
        })
        if err != nil {
            return err
        }

        // Stack exports
        ctx.Export("bucketName", s3Bucket.Bucket)
        ctx.Export("kinesisDataStreamName", dataStream.Name)
        ctx.Export("dataTransformLambdaName", dataTransformLambda.Name)
        ctx.Export("firehoseDeliveryStreamName", firehoseStream.Name)
        ctx.Export("lambdaRoleName", lambdaRole.Name)
        ctx.Export("firehoseRoleName", firehoseRole.Name)

        return nil
    })
}
