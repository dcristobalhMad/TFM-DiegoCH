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
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
		})
		if err != nil {
			return err
		}
		// Create a Kinesis Data Stream
		dataStream, err := kinesis.NewStream(ctx, "kinesisDataStream", &kinesis.StreamArgs{
			Name:       pulumi.String("tfm-stream"),
			ShardCount: pulumi.Int(1),
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
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
		// Create Glue catalog table
		catalogTable, err := glue.NewCatalogTable(ctx, "awsGlueCatalogTable", &glue.CatalogTableArgs{
			DatabaseName: catalogDatabase.Name,
			Name:         pulumi.String("tfmttable"),
			Description:  pulumi.String("An example Glue Catalog Table with output in Parquet format"),
			TableType:    pulumi.String("EXTERNAL_TABLE"),
			Parameters: pulumi.StringMap{
				"EXTERNAL":                      pulumi.String("TRUE"),
				"parquet.compression":           pulumi.String("SNAPPY"),
				"projection.enabled":            pulumi.String("true"),
				"projection.date.type":          pulumi.String("date"),
				"projection.date.format":        pulumi.String("yyyy-MM-dd"),
				"projection.date.range":         pulumi.String("2022-10-01,NOW"),
				"projection.date.interval":      pulumi.String("1"),
				"projection.date.interval.unit": pulumi.String("DAYS"),
				"storage.location.template":     pulumi.Sprintf("s3://%s/events/date=$${date}", s3Bucket.ID()),
			},
			StorageDescriptor: &glue.CatalogTableStorageDescriptorArgs{
				Location:     pulumi.Sprintf("s3://%s/events/", s3Bucket.ID()),
				InputFormat:  pulumi.String("org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat"),
				OutputFormat: pulumi.String("org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat"),
				Compressed:   pulumi.Bool(false),
				SerDeInfo: &glue.CatalogTableStorageDescriptorSerDeInfoArgs{
					Name: pulumi.String("events"),
					Parameters: pulumi.StringMap{
						"serialization.format": pulumi.String("1"),
					},
					SerializationLibrary: pulumi.String("org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe"),
				},
			},
		},
			pulumi.DependsOn([]pulumi.Resource{catalogDatabase}))
		if err != nil {
			return err
		}

		// Create a Lambda IAM role
		lambdaRole, err := iam.NewRole(ctx, "dataTransformLambdaRole", &iam.RoleArgs{
			Name: pulumi.String("tfm-lambda-role"),
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
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
			Timeout: pulumi.Int(60),
			Role:    lambdaRole.Arn,
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
		})
		if err != nil {
			return err
		}

		// Create a Kinesis Firehose IAM role
		firehoseRole, err := iam.NewRole(ctx, "firehoseDeliveryStreamRole", &iam.RoleArgs{
			Name: pulumi.String("firehoseDeliveryStreamRole"),
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
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

		// Attach Glue CatalogRead policy to the IAM role
		readGluePolicy, err := iam.NewPolicy(ctx, "myReadGluePolicy", &iam.PolicyArgs{
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
			Policy: pulumi.Sprintf(`{
              "Version": "2012-10-17",
              "Statement": [
                {
                  "Action": [
                    "glue:GetDatabase",
                    "glue:GetTable"
                  ],
                  "Resource": [
                    "arn:aws:glue:*:*:*"
                  ],
                  "Effect": "Allow"
                }
              ]
            }`),
		})
		if err != nil {
			return err
		}
		_, err = iam.NewRolePolicyAttachment(ctx, "myReadGluePolicyAttachment", &iam.RolePolicyAttachmentArgs{
			PolicyArn: readGluePolicy.Arn,
			Role:      firehoseRole.Name,
		})
		if err != nil {
			return err
		}

		putDataS3, err := iam.NewPolicy(ctx, "putDataS3", &iam.PolicyArgs{
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
			Policy: pulumi.Sprintf(`{
              "Version": "2012-10-17",
              "Statement": [
                {
                  "Action": [
                    "s3:AbortMultipartUpload",
					"s3:GetBucketLocation",
					"s3:GetObject",
					"s3:ListBucket",
					"s3:ListBucketMultipartUploads",
					"s3:PutObject"
                  ],
                  "Resource": [
                    "%s"
                  ],
                  "Effect": "Allow"
                }
              ]
            }`, s3Bucket.Arn),
		})
		if err != nil {
			return err
		}
		// Attach the policy to put data to s3 bucket
		_, err = iam.NewRolePolicyAttachment(ctx, "putDataS3PolicyAttachment", &iam.RolePolicyAttachmentArgs{
			PolicyArn: putDataS3.Arn,
			Role:      firehoseRole.Name,
		})
		if err != nil {
			return err
		}

		// Attach the policy to read the Kinesis Stream.
		_, err = iam.NewRolePolicyAttachment(ctx, "kinesisStreamRolePolicyAttachment", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonKinesisReadOnlyAccess"),
			Role:      firehoseRole.Name,
		})
		if err != nil {
			return err
		}

		// Create a Kinesis Firehose Delivery Stream with data transformation Lambda
		firehoseStream, err := kinesis.NewFirehoseDeliveryStream(ctx, "firehoseDeliveryStream", &kinesis.FirehoseDeliveryStreamArgs{
			Destination: pulumi.String("extended_s3"),
			Name:        pulumi.String("tfm-firehose-stream"),
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
			KinesisSourceConfiguration: &kinesis.FirehoseDeliveryStreamKinesisSourceConfigurationArgs{
				KinesisStreamArn: dataStream.Arn,
				RoleArn:          firehoseRole.Arn, // Replace with your IAM role ARN
			},
			ExtendedS3Configuration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationArgs{
				RoleArn:           firehoseRole.Arn,
				BucketArn:         s3Bucket.Arn,
				BufferSize:        pulumi.Int(128),
				BufferInterval:    pulumi.Int(60),
				CompressionFormat: pulumi.String("UNCOMPRESSED"),
				Prefix:            pulumi.String("events/date=!{timestamp:yyyy}-!{timestamp:MM}-!{timestamp:dd}/"),
				ErrorOutputPrefix: pulumi.String("events_error/year=!{timestamp:yyyy}/month=!{timestamp:MM}/day=!{timestamp:dd}/hour=!{timestamp:HH}/!{firehose:error-output-type}/"),
				DataFormatConversionConfiguration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationArgs{
					Enabled: pulumi.Bool(true),
					InputFormatConfiguration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationInputFormatConfigurationArgs{
						Deserializer: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationInputFormatConfigurationDeserializerArgs{
							HiveJsonSerDe: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationInputFormatConfigurationDeserializerHiveJsonSerDeArgs{
								TimestampFormats: pulumi.StringArray{
									pulumi.String("yyyy-MM-dd'T'HH:mm:ss.SSSSSS"),
								},
							},
						},
					},
					OutputFormatConfiguration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationOutputFormatConfigurationArgs{
						Serializer: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationOutputFormatConfigurationSerializerArgs{
							ParquetSerDe: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationOutputFormatConfigurationSerializerParquetSerDeArgs{
								Compression: pulumi.String("SNAPPY"),
							},
						},
					},
					SchemaConfiguration: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationDataFormatConversionConfigurationSchemaConfigurationArgs{
						CatalogId:    pulumi.String(""), // empty string means current account
						DatabaseName: catalogDatabase.Name,
						RoleArn:      firehoseRole.Arn,
						TableName:    catalogTable.Name,
						Region:       pulumi.String("us-east-1"),
						VersionId:    pulumi.String("LATEST"),
					},
				},
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
		ctx.Export("glueTableNameX", catalogTable.Name)

		return nil
	})
}
