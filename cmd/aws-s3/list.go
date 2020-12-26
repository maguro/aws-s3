package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().BoolP("force", "f", false, "Force an error if no objects are found")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   `list [4mbucket[0m [4mprefix[0m`,
	Short: "List AWS S3 objects",
	Long: `Lists AWS S3 objects from the bucket [4mbucket[0m objects with the key
prefix "[4mprefix[0m".
`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		prefix := args[1]

		sess, err := session.NewSession(&aws.Config{
			Region:                         &region,
			DisableRestProtocolURICleaning: aws.Bool(true),
		})
		if err != nil {
			return err
		}

		svc := s3.New(sess)

		result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		})
		if err != nil {
			return err
		}

		forced, _ := cmd.Flags().GetBool("force")
		if len(result.Contents) == 0 && forced {
			return fmt.Errorf("no objects found for s3://%s/%s", bucket, prefix)
		}

		for _, item := range result.Contents {
			if *item.Size > 0 {
				println(*item.Key)
			}
		}

		return nil
	},
}
