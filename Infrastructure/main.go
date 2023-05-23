package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
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
		s3Bucket, err := s3.NewBucket(ctx, "tfm-diego-datalake", &s3.BucketArgs{
			ForceDestroy: pulumi.Bool(true),
			Acl:          pulumi.String("private"),
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
				Columns: glue.CatalogTableStorageDescriptorColumnArray{
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Ip address of the client"),
						Name:    pulumi.String("ip1"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Ip address of the server"),
						Name:    pulumi.String("ip2"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Timestamp of the request"),
						Name:    pulumi.String("timestamp"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Request section"),
						Name:    pulumi.String("section"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Resource requested"),
						Name:    pulumi.String("resource"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Protocol used"),
						Name:    pulumi.String("values"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Status code of the request"),
						Name:    pulumi.String("status_code"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Size of the request"),
						Name:    pulumi.String("size"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Dash1"),
						Name:    pulumi.String("dash1"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Dash2"),
						Name:    pulumi.String("dash2"),
						Type:    pulumi.String("int"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Flags"),
						Name:    pulumi.String("flags"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Values2"),
						Name:    pulumi.String("values2"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Values3"),
						Name:    pulumi.String("values3"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Request user agent"),
						Name:    pulumi.String("user_agent"),
						Type:    pulumi.String("string"),
					},
					&glue.CatalogTableStorageDescriptorColumnArgs{
						Comment: pulumi.String("Request type"),
						Name:    pulumi.String("request"),
						Type:    pulumi.String("string"),
					},
				},
				Location: pulumi.Sprintf("s3://%s/events/", s3Bucket.ID()),
				// input format should be json
				InputFormat:  pulumi.String("org.apache.hadoop.mapred.TextInputFormat"),
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
			Runtime: pulumi.String("python3.9"),
			Name:    pulumi.String("dataTransformLambda"),
			Code:    pulumi.NewFileArchive("./lambda/lambda_function.zip"),
			Handler: pulumi.String("lambda_function.lambda_handler"),
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

		// Attach the AWSLambdaExecute policy to the Lambda role
		_, err = iam.NewRolePolicyAttachment(ctx, "lambdaPolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      lambdaRole.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AWSLambdaExecute"),
		})
		if err != nil {
			return err
		}
		// Attach the AWSLambdaBasicExecutionRole policy to the Lambda role
		_, err = iam.NewRolePolicyAttachment(ctx, "basicExecutionRole", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
			Role:      lambdaRole.Name,
		})
		if err != nil {
			return err
		}

		// Create a Cloudwatch Log Group for the Lambda Function
		logGroup, err := cloudwatch.NewLogGroup(ctx, "tfmdiegoLogGroup", &cloudwatch.LogGroupArgs{
			Name:            pulumi.Sprintf("/aws/lambda/%s", dataTransformLambda.Name),
			RetentionInDays: pulumi.Int(1),
		})
		if err != nil {
			return err
		}

		// Attach the Inline Policy for specific LogGroup access
		_, err = iam.NewRolePolicy(ctx, "allowLambdaLoggingToSpecificLogGroup", &iam.RolePolicyArgs{
			Role: lambdaRole.Name,
			Policy: pulumi.Sprintf(`{
                    "Version": "2012-10-17",
                    "Statement": [
                        {
                            "Effect": "Allow",
                            "Action": [
                                "logs:CreateLogStream",
                                "logs:PutLogEvents"
                            ],
                            "Resource": "%s"
                        }
                    ]
                }`, logGroup.Arn),
		})
		if err != nil {
			return err
		}

		// Create a Kinesis Firehose IAM role
		firehoseRole, err := iam.NewRole(ctx, "firehoseDeliveryStreamRole", &iam.RoleArgs{
			Name: pulumi.String("tfm-diego-firehose-role"),
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
		// Attach the AmazonKinesisFirehoseFullAccess policy to the Firehose role
		_, err = iam.NewRolePolicyAttachment(ctx, "firehoseRolePolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      firehoseRole.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonKinesisFirehoseFullAccess"),
		})
		if err != nil {
			return err
		}
		// Attach the AmazonKinesisFirehoseFullAccess policy to the Firehose role
		_, err = iam.NewRolePolicyAttachment(ctx, "lambdaFullPolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      firehoseRole.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AWSLambda_FullAccess"),
		})
		if err != nil {
			return err
		}
		// Attach the AmazonKinesisFirehoseFullAccess policy to the Firehose role
		_, err = iam.NewRolePolicyAttachment(ctx, "s3FullPolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      firehoseRole.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
		})
		if err != nil {
			return err
		}
		// Attach Glue CatalogRead policy to the IAM role
		lambdaExecFirehose, err := iam.NewPolicy(ctx, "lambdaExecFirehose", &iam.PolicyArgs{
			Tags: pulumi.StringMap{
				"Env":  pulumi.String("test"),
				"Name": pulumi.String("tfm-diego"),
			},
			Policy: pulumi.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Effect": "Allow",
							"Action": [
								"lambda:InvokeFunction"
							],
							"Resource": [
								"%s"
							]
						}
					]
            }`, dataTransformLambda.Arn),
		})
		if err != nil {
			return err
		}
		_, err = iam.NewRolePolicyAttachment(ctx, "lambdaExecFirehoseAttachment", &iam.RolePolicyAttachmentArgs{
			PolicyArn: lambdaExecFirehose.Arn,
			Role:      firehoseRole.Name,
		})
		if err != nil {
			return err
		}
		/////////////////////////
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

		// Create a Cloudwatch Log Group for the kinesis firehose
		firehoseLogGroup, err := cloudwatch.NewLogGroup(ctx, "firehoseLogGroup", &cloudwatch.LogGroupArgs{
			Name:            pulumi.String("/aws/kinesisfirehose/tfm-firehose-stream"),
			RetentionInDays: pulumi.Int(1),
		})
		if err != nil {
			return err
		}
		// Create a CloudWatch Log Stream
		_, err = cloudwatch.NewLogStream(ctx, "firehoseLogStream", &cloudwatch.LogStreamArgs{
			Name:         pulumi.String("tfm-firehose-stream"),
			LogGroupName: firehoseLogGroup.Name,
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
				CloudwatchLoggingOptions: &kinesis.FirehoseDeliveryStreamExtendedS3ConfigurationCloudwatchLoggingOptionsArgs{
					Enabled:       pulumi.Bool(true),
					LogGroupName:  firehoseLogGroup.Name,
					LogStreamName: pulumi.String("tfm-firehose-stream"),
				},
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
								Compression:                 pulumi.String("SNAPPY"),
								EnableDictionaryCompression: pulumi.Bool(false),
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
		ctx.Export("logGroupName", logGroup.Name)

		return nil
	})
}
