package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"flag"
)

func buildRDSTokenWithMFA(profile, defaultRegion, mfaARN, endpoint, dbUser string) (string, error) {
	// Profile の読み込み
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return "", fmt.Errorf("Unable to load SDK config: %w", err)
	}

	// MFA認証
	// sts による一時的な認証情報の取得
	var mfaToken string
	fmt.Fprintln(os.Stderr, "Please input MFA code:")
    fmt.Scan(&mfaToken)
	stsClient := sts.NewFromConfig(cfg)
	creds, err := stsClient.GetSessionToken(context.Background(), &sts.GetSessionTokenInput{
		TokenCode:       aws.String(mfaToken),
		SerialNumber:    aws.String(mfaARN),
		DurationSeconds: aws.Int32(3600),
	})

	if err != nil {
		return "", fmt.Errorf("Failed to get session token: %w", err)
	}

	cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		*creds.Credentials.AccessKeyId,
		*creds.Credentials.SecretAccessKey,
		*creds.Credentials.SessionToken,
	))

	// RDS の認証情報の取得
	// iam_operator のパスワード取得
	authenticationToken, err := auth.BuildAuthToken(
		context.Background(), endpoint, defaultRegion, dbUser, cfg.Credentials)
	if err != nil {
		return "", fmt.Errorf("Failed to create rds authentication token: %w", err)
	}

	return authenticationToken, nil
}

func main() {
	// 使用するAWSプロファイル名を設定
	profile := flag.String("profile", "", "AWS の Profile 名")
	mfaARN := flag.String("mfaARN", "", "mfaARN 例) arn:aws:iam::xxxxxxxx:mfa/isao-aruga")
	dbUser := flag.String("user", "", "DBのユーザー名 例) iam_operator")
	endpoint := flag.String("endpoint", "", "DBのエンドポイント 例) sms-main-xxxxxx.rds.amazonaws.com:3306")
	defaultRegion := flag.String("region", "ap-northeast-1", "利用するAWSのリージョン")
	flag.Parse()

	if *mfaARN == "" {
		fmt.Println("mfaARN is required")
		os.Exit(1)
	}

	if *dbUser == "" {
		fmt.Println("dbUser is required")
		os.Exit(1)
	}

	if *endpoint == "" {
		fmt.Println("endpoint is required")
		os.Exit(1)
	}

	token, err := buildRDSTokenWithMFA(*profile, *defaultRegion, *mfaARN, *endpoint, *dbUser)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(token)
}
