package exports

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/tndx/pkg/database"
	"github.com/sirupsen/logrus"
)

func runStartDDBExport() error {
	opt, err := svc.db.ExportTable(&database.TableExportRequest{
		ExportFormat: flags.format,
		S3Bucket:     flags.s3Bucket,
		S3Prefix:     flags.s3Prefix,
		TableArn:     flags.tableArn,
		ExportTime:   time.Now(),
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("error exporting table")
		return err
	}
	spew.Dump(opt.ExportDescription)

	exportLabel := strings.Split(*opt.ExportDescription.ExportArn, "/")[:1]
	fmt.Println("Export data")
	fmt.Printf("Status: %s\n", opt.ExportDescription.ExportStatus)
	if opt.ExportDescription.FailureCode != nil {
		fmt.Printf("Failure Code: %s\n", *opt.ExportDescription.FailureCode)
		fmt.Printf("Failure Message: %s\n", *opt.ExportDescription.FailureMessage)
	}

	fmt.Printf("Path: s3://%s/%sAWSDynamoDB/%s\n", *opt.ExportDescription.S3Bucket, *opt.ExportDescription.S3Prefix, exportLabel)
	fmt.Printf("ExportArn: %s\n", *opt.ExportDescription.ExportArn)
	fmt.Printf("Start time: %s\n", *opt.ExportDescription.StartTime)
	fmt.Printf("End time: %s\n", *opt.ExportDescription.EndTime)
	fmt.Printf("Format: %s\n", opt.ExportDescription.ExportFormat)
	return nil
}
