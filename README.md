## RDSIAMToken 

aws cli がない状況でも RDS IAM Token を取得したく生み出された

## 使い方

```
export AWS_ACCESS_KEY_ID="FOO"
export AWS_SECRET_ACCESS_KE="BAR"
export AWS_DEFAULT_REGION="ap-northeast-1"

./RDSIamToken -user ${RDSのユーザ名} -endpoint ${RDSのエンドポイント}:${RDSのポート} -mfaARN ${arn:aws:iam::1234567890:mfa/xxx}
```
