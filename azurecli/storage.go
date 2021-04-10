package azurecli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/cjlapao/common-go/helper"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-06-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest"
)

type StorageAccountKey struct {
	Name  string
	Value string
}

type StorageAccountClient struct {
	Authorizer           autorest.Authorizer
	StorageAccountClient *storage.AccountsClient
	StorageAccountKeys   *[]StorageAccountKey
}

func CreateStorageAccountClient() *StorageAccountClient {
	cli := StorageAccountClient{}

	return &cli
}

func (cli *StorageAccountClient) Test() {

}

var storageKeys *[]storage.AccountKey

func getStorageAccountClient() (*storage.AccountsClient, error) {
	logger.Debug("Trying to get the Azure storage account client")
	storageAccountsClient := storage.NewAccountsClient(ctx.SubscriptionID)
	authorizeToken, err := GetClientAuthorizer()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}
	storageAccountsClient.Authorizer = authorizeToken
	storageAccountsClient.AddToUserAgent(UserAgent)
	logger.Debug("Finished getting the Azure storage account client")
	return &storageAccountsClient, nil
}

func getStorageAccountKeys() error {
	logger.Debug("Trying to get the Azure storage account keys")
	if storageKeys != nil {
		logger.Debug("Found cached Azure storage account key")
		return nil
	}

	if ctx.Storage.AccountName == "" {
		err := errors.New("The Azure context is missing the account name")
		logger.LogError(err)
		return err
	}

	localCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli, err := getStorageAccountClient()
	result, err := cli.ListKeys(localCtx, ctx.ResourceGroup, ctx.Storage.AccountName)
	if err != nil {
		logger.LogError(err)
		return err
	}

	storageKeys = result.Keys
	keys := *result.Keys

	if len(keys) >= 1 {
		ctx.Storage.PrimaryAccountKey = *keys[0].Value
		if len(keys) == 2 {
			ctx.Storage.SecondaryAccountKey = *keys[1].Value
		}
	}

	logger.Debug("Finished getting and caching the Azure storage account keys")
	return nil
}

func getServiceURL() (*azblob.ServiceURL, error) {
	logger.Debug("Trying to get Azure storage blob serviceURL")

	err := getStorageAccountKeys()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	if ctx.Storage.AccountName == "" || ctx.Storage.PrimaryAccountKey == "" || ctx.Storage.ContainerName == "" {
		err := errors.New("The Azure storage context is missing mandatory parameters")
		logger.Error(err.Error())
		return nil, err
	}

	credentials, err := azblob.NewSharedKeyCredential(ctx.Storage.AccountName, ctx.Storage.PrimaryAccountKey)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	p := azblob.NewPipeline(credentials, azblob.PipelineOptions{})
	url, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", ctx.Storage.AccountName))

	serviceURL := azblob.NewServiceURL(*url, p)
	logger.Debug("Finished getting Azure serviceURL")
	return &serviceURL, nil
}

func getContainerURL() (*azblob.ContainerURL, error) {
	logger.Debug("Trying to get Azure storage blob containerURL")
	if ctx.Storage.ContainerName == "" {
		err := errors.New("Container name cannot be null")
		logger.Error(err.Error())
		return nil, err
	}

	serviceURL, err := getServiceURL()
	if err != nil {
		return nil, err
	}
	containerURL := serviceURL.NewContainerURL(ctx.Storage.ContainerName)

	logger.Debug("Finished getting Azure storage blob containerURL")
	return &containerURL, nil
}

func CreateContainer() error {
	localCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := getContainerURL()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = c.Create(localCtx, azblob.Metadata{}, azblob.PublicAccessContainer)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func DeleteContainer() error {
	localCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := getContainerURL()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_, err = c.Delete(localCtx, azblob.ContainerAccessConditions{})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func ListFilesInContainer(next string) ([]azblob.BlobItemInternal, error) {
	logger.Debug("Starting to list files in %v", ctx.Storage.ContainerName)
	marker := azblob.Marker{}
	if next != "" {
		marker.Val = &next
	}

	items := make([]azblob.BlobItemInternal, 0)
	containerURL, err := getContainerURL()
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	localCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := containerURL.ListBlobsFlatSegment(localCtx,
		marker,
		azblob.ListBlobsSegmentOptions{
			Details: azblob.BlobListingDetails{
				Snapshots: true,
			},
		})

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	for _, item := range response.Segment.BlobItems {
		items = append(items, item)
	}

	if *response.NextMarker.Val != "" {
		nextItems, err := ListFilesInContainer(*response.NextMarker.Val)
		if err != nil {
			return items, err
		}
		if len(nextItems) > 0 {
			for _, item := range nextItems {
				items = append(items, item)
			}
		}
	}

	return items, nil
}

// UploadBlob Uploads a blob to a container
func UploadBlob() error {
	logger.Debug("Starting upload of blob file %v", ctx.Storage.FileName)
	if ctx.Storage.FileName == "" || !helper.FileExists(ctx.Storage.FileName) {
		err := fmt.Errorf("File %v does not exist or was not defined", ctx.Storage.FileName)
		logger.LogError(err)
		return err
	}

	containerURL, err := getContainerURL()
	if err != nil {
		logger.LogError(err)
		return err
	}

	logger.Command("Trying to upload file %v from container %v in account %v", ctx.Storage.FileName, ctx.Storage.ContainerName, ctx.Storage.AccountName)
	localCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	uploadData, err := helper.ReadFromFile(ctx.Storage.FileName)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(uploadData)
	blobURL := containerURL.NewBlockBlobURL(ctx.Storage.FileName)
	_, err = blobURL.Upload(localCtx, reader, azblob.BlobHTTPHeaders{ContentType: "application/octet-stream"}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.DefaultAccessTier, nil, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		logger.LogError(err)
		return err
	}
	logger.Success("Finished uploading file %v to container %v in account %v", ctx.Storage.FileName, ctx.Storage.ContainerName, ctx.Storage.AccountName)
	return nil
}

// DownloadBlob Downloads a blob to a container
func DownloadBlob() error {
	logger.Debug("Starting download of blob file %v", ctx.Storage.FileName)
	if ctx.Storage.FileName == "" {
		err := fmt.Errorf("File %v does not exist or was not defined", ctx.Storage.FileName)
		logger.LogError(err)
		return err
	}

	containerURL, err := getContainerURL()
	if err != nil {
		logger.LogError(err)
		return err
	}

	logger.Command("Trying to download file %v from container %v in account %v", ctx.Storage.FileName, ctx.Storage.ContainerName, ctx.Storage.AccountName)
	localCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	blobURL := containerURL.NewBlockBlobURL(ctx.Storage.FileName)
	get, err := blobURL.Download(localCtx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		logger.LogError(err)
		return err
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(azblob.RetryReaderOptions{})
	downloadedData.ReadFrom(reader)
	reader.Close()

	helper.WriteToFile(downloadedData.String(), ctx.Storage.FileName)

	logger.Success("Finished downloading file %v to container %v in account %v", ctx.Storage.FileName, ctx.Storage.ContainerName, ctx.Storage.AccountName)
	return nil
}
