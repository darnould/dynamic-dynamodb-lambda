# dynamic-dynamodb-lambda

Work in progress, steer well clear for now.

```

## Execution Role
{
    "Description": "AWS Lambda execution role",
    "Resources": {
        "LambdaExecutionRole": {
            "Type": "AWS::IAM::Role",
            "Properties": {
                "AssumeRolePolicyDocument": {
                    "Version": "2012-10-17",
                    "Statement": [{
                        "Effect": "Allow",
                        "Principal": {
                            "Service": "lambda.amazonaws.com"
                        },
                        "Action": "sts:AssumeRole"
                    }]
                },
                "Path": "/",
                "Policies": [{
                    "PolicyName": "LambdaLogging",
                    "PolicyDocument": {
                        "Version": "2012-10-17",
                        "Statement": [{
                            "Effect": "Allow",
                            "Action": [
                                "logs:CreateLogGroup",
                                "logs:CreateLogStream",
                                "logs:PutLogEvents"
                            ],
                            "Resource": [
                                "arn:aws:logs:*:*:*"
                            ]
                        }]
                    }
                }]
            }
        }
    },
    "Outputs": {
        "LambdaExecutionRoleARN": { "Value": { "Fn::GetAtt": ["LambdaExecutionRole", "Arn"] } }
    }
}
```
