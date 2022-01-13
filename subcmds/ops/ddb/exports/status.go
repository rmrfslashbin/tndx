package exports

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/message"
)

func runStatusDDBExport() error {
	if opt, err := svc.db.ExportStatus(flags.exportArn); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("error getting export status")
		return err
	} else {
		exportArnParts := strings.Split(*opt.ExportDescription.ExportArn, "/")
		exportLabel := exportArnParts[len(exportArnParts)-1]
		s, ok := os.LookupEnv("LANG")
		if !ok {
			s = "en-US"
		}
		p := message.NewPrinter(message.MatchLanguage(s))

		p.Println("Export data")
		p.Printf("Status: %s\n", opt.ExportDescription.ExportStatus)
		if opt.ExportDescription.FailureCode != nil {
			p.Printf("Failure Code: %s\n", *opt.ExportDescription.FailureCode)
			p.Printf("Failure Message: %s\n", *opt.ExportDescription.FailureMessage)
		}

		p.Printf("Path: s3://%s/%sAWSDynamoDB/%s\n", *opt.ExportDescription.S3Bucket, *opt.ExportDescription.S3Prefix, exportLabel)
		p.Printf("Manifest: s3://%s/%s\n", *opt.ExportDescription.S3Bucket, *opt.ExportDescription.ExportManifest)
		p.Printf("Billed Bytes: %d\n", *opt.ExportDescription.BilledSizeBytes)
		p.Printf("Item Count: %d\n", *opt.ExportDescription.ItemCount)
		p.Printf("ExportArn: %s\n", *opt.ExportDescription.ExportArn)
		p.Printf("Start time: %s\n", *opt.ExportDescription.StartTime)
		p.Printf("End time: %s\n", *opt.ExportDescription.EndTime)
		p.Printf("Format: %s\n", opt.ExportDescription.ExportFormat)
	}
	return nil
}
