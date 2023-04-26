package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/glue"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/kinesis"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create an S3 bucket to store the Kinesis data
		s3Bucket, err := s3.NewBucket(ctx, "mydatalake", &s3.BucketArgs{
			Acl: pulumi.String("private"),
			Tags: pulumi.StringMap{
				"Environment": pulumi.String("Prod"),
				"Name":        pulumi.String("mydatalake"),
			},
		})
		if err != nil {
			return err
		}
		// Create a Kinesis Data Stream
		dataStream, err := kinesis.NewStream(ctx, "kinesisDataStream", &kinesis.StreamArgs{
			Name:       pulumi.String("tfm-stream"),
			ShardCount: pulumi.Int(1),
		})
		if err != nil {
			return err
		}

		// Create a Glue catalog database
		catalogDatabase, err := glue.NewCatalogDatabase(ctx, "awsGlueCatalogDatabase", &glue.CatalogDatabaseArgs{
			Name: pulumi.String("tfmcatalogdatabase"),
		})
		if err != nil {
			return err
		}
		// Create variables
		s3BucketName := s3Bucket.ID()
		kinesisStreamName := dataStream.Name
		// Create Glue catalog table
		catalogTable, err := glue.NewCatalogTable(ctx, "awsGlueCatalogTable", &glue.CatalogTableArgs{
			DatabaseName: catalogDatabase.Name,
			Name:         pulumi.String("tfmttable"),
			Parameters: pulumi.StringMap{
				"EXTERNAL":            pulumi.String("TRUE"),
				"parquet.compression": pulumi.String("SNAPPY"),
			},
			StorageDescriptor: &glue.CatalogTableStorageDescriptorArgs{
				Columns: glue.CatalogTableStorageDescriptorColumnArray{
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Name: pulumi.String("my_string"),
						Type: pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Name: pulumi.String("my_double"),
						Type: pulumi.String("double"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String(""),
						Name:    pulumi.String("my_date"),
						Type:    pulumi.String("date"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String(""),
						Name:    pulumi.String("my_bigint"),
						Type:    pulumi.String("bigint"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String(""),
						Name:    pulumi.String("my_struct"),
						Type:    pulumi.String("struct<my_nested_string:string>"),
					},
				},
				// Input format should be raw
				InputFormat:  pulumi.String("org.apache.hadoop.mapred.TextInputFormat"),
				Location:     pulumi.Sprintf("s3://%s/event-streams/%s", s3BucketName, kinesisStreamName),
				OutputFormat: pulumi.String("org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat"),
				SerDeInfo: &glue.CatalogTableStorageDescriptorSerDeInfoArgs{
					Name: dataStream.Name,
					Parameters: pulumi.StringMap{
						"serialization.format": pulumi.String("1"),
					},
					SerializationLibrary: pulumi.String("org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe"),
				},
			},
			TableType: pulumi.String("EXTERNAL_TABLE"),
		}, pulumi.DependsOn([]pulumi.Resource{catalogDatabase}))
		if err != nil {
			return err
		}

		// Create a Lambda IAM role
		lambdaRole, err := iam.NewRole(ctx, "dataTransformLambdaRole", &iam.RoleArgs{
			Name: pulumi.String("tfm-lambda-role"),
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
			Name:    pulumi.String("dataTransformLambda"),
			Code:    pulumi.NewFileArchive("./lambda/bin/lambda_function.zip"),
			Handler: pulumi.String("main"),
			Role:    lambdaRole.Arn,
		})
		if err != nil {
			return err
		}

		// Create a Kinesis Firehose IAM role
		firehoseRole, err := iam.NewRole(ctx, "firehoseDeliveryStreamRole", &iam.RoleArgs{
			Name: pulumi.String("firehoseDeliveryStreamRole"),
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
			Name:        pulumi.String("tfm-firehose-stream"),
			ExtendedS3Configuration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationArgs{
				RoleArn:   firehoseRole.Arn,
				BucketArn: s3Bucket.Arn,
				ProcessingConfiguration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationArgs{
					Enabled: pulumi.Bool(true),
					Processors: kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorArray{
						&kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorArgs{
							Type: pulumi.String("Lambda"),
							Parameters: kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorParameterArray{
								&kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationProcessingConfigurationProcessorParameterArgs{
									ParameterName: pulumi.String("LambdaArn"),
									ParameterValue: dataTransformLambda.Arn.ApplyT(func(arn string) (string, error) {
										return fmt.Sprintf("%v:$LATEST", arn), nil
									}).(pulumi.StringOutput),
								},
							},
						},
					},
				},
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
		ctx.Export("glueDatabaseName", catalogDatabase.Name)
		ctx.Export("glueTableName", catalogTable.Name)

		return nil
	})
}
