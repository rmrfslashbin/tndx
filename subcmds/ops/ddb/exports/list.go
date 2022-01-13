package exports

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func runListDDBExport() error {
	if opt, err := svc.db.ExportList(flags.tableArn); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("error getting export status")
		return err
	} else {
		fmt.Printf("%d Export(s) found\n", len(opt.ExportSummaries))
		for _, export := range opt.ExportSummaries {
			fmt.Printf("%s: %s\n", export.ExportStatus, *export.ExportArn)
		}

		if opt.NextToken != nil {
			fmt.Printf("\nNextToken: %s\n", *opt.NextToken)
		}
	}
	return nil
}
