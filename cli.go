package main

import "github.com/spf13/cobra"

var region string

func init() {
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "AWS region")
}

var rootCmd = &cobra.Command{
	Use:   "aws-s3",
	Short: "A simple AWS S3 tool",
	Long:  `A simple AWS S3 tool`,
}
