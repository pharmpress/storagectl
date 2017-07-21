package command

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/urfave/cli"
)

// NewAzureUploadCommand is upload implementation fro s3
func NewAzureUploadCommand() cli.Command {
	return cli.Command{
		Name:  "upload",
		Usage: "upload blob into a container",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "container", Usage: "Container Name"},
			cli.StringFlag{Name: "blob", Usage: "Blob name"},
		},
		Action: func(c *cli.Context) {
			azureUploadCommandFunc(c)
		},
	}
}

func azureUploadCommandFunc(c *cli.Context) error {
	accountName := c.GlobalString("account-name")
	accountKey := c.GlobalString("account-key")
	containerName := c.String("container")
	blobName := c.String("blob")
	fileToUpload := "/file_to_upload"

	if len(c.Args()) != 0 {
		fileToUpload = c.Args()[0]
	}

	client, err := storage.NewBasicClient(accountName, accountKey)

	if err != nil {
		fmt.Println(err)
		return err
	}

	blobService := client.GetBlobService()

	fileReader, err := os.Open(fileToUpload)
	defer fileReader.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	fi, err := fileReader.Stat()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Uploading " + fileToUpload)
	size := uint64(fi.Size())
	var chunkSize = storage.MaxBlobBlockSize

	if size > uint64(chunkSize) {
		err = createBlockBlobFromReaderWithChunks(blobService, containerName, blobName, size, fileReader, chunkSize)
	} else {
		err = blobService.CreateBlockBlobFromReader(containerName, blobName, size, fileReader, nil)
	}

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func createBlockBlobFromReaderWithChunks(blobService storage.BlobStorageClient, container, name string, size uint64, inputSourceReader io.Reader, chunkSize int) error {
	blocksLen := uint64(0)
	//blockList, err := blobService.GetBlockList(container, name, storage.BlockListTypeAll)
	//if err != nil {
	//	return err
	//}
	//blocksLen = len(blockList.CommittedBlocks)

	amendList := []storage.Block{}
	//for _, v := range blockList.CommittedBlocks {
	//	amendList = append(amendList, storage.Block{v.Name, storage.BlockStatusCommitted})
	//}
	chunk := make([]byte, chunkSize)
	total := 0
	for {
		n, err := inputSourceReader.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		atomic.AddUint64(&blocksLen, 1)
		blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%011d\n", blocksLen)))

		data := chunk[:n]
		total += n
		err = blobService.PutBlock(container, name, blockID, data)
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Progress %.2f%%", (float64(total)/float64(size))*100.0))
		amendList = append(amendList, storage.Block{ID: blockID, Status: storage.BlockStatusUncommitted})
	}
	err := blobService.PutBlockList(container, name, amendList)
	return err
}
