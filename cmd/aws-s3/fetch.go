package main

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
)

func init() {
	fetchCmd.Flags().BoolP("force", "f", false, "Force an error if no are objects found")
	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   `fetch [4mbucket[0m [4mprefix[0m [4mdestination[0m`,
	Short: "Fetch AWS S3 objects",
	Long: `
Fetches AWS S3 objects from the bucket [4mbucket[0m objects with the key
prefix "[4mprefix[0m" and downloads them to [4mdestination[0m.  Any files
in [4mdestination[0m will be overwritten.
`,
	Args: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucket := args[0]
		prefix := args[1]
		destination := args[2]

		if _, err := os.Stat(destination); os.IsNotExist(err) {
			return err
		}

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

		downloader := s3manager.NewDownloader(sess)
		for _, item := range result.Contents {
			if *item.Size > 0 {
				file, err := os.Create(path.Join(destination, path.Base(*item.Key)))
				if err != nil {
					return err
				}
				_, err = downloader.Download(file, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    item.Key,
				})
				_ = file.Close()
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}
