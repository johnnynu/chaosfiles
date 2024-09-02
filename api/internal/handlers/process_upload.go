package handlers

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/johnnynu/agreatchaos/api/internal/db"
)

func ProcessUpload(ctx context.Context, s3Event events.S3Event) error {
    for _, record := range s3Event.Records {
        key := record.S3.Object.Key // fileID
        size := record.S3.Object.Size

        // Fetch file metadata
        file, err := db.GetFile(ctx, key)
        if err != nil {
            log.Printf("Error fetching file metadata: %v", err)
            return err
        }

        // Update file metadata
        file.FileSize = size
		file.UpdatedAt = time.Now().Format(time.RFC3339)

        err = db.UpdateFile(ctx, *file)
        if err != nil {
            log.Printf("Error updating file metadata: %v", err)
            return err
        }

        log.Printf("Successfully processed upload for file: %s", key)
    }

    return nil
}