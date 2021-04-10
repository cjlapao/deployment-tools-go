package entities

type AzureBlobResponse struct {
	Content    interface{}         `json:"content"`
	Deleted    bool                `json:"deleted"`
	Metadata   AzureBlobMetadata   `json:"metadata"`
	Name       string              `json:"name"`
	Properties AzureBlobProperties `json:"properties"`
	Snapshot   interface{}         `json:"snapshot"`
}

type AzureBlobMetadata struct {
}

type AzureBlobProperties struct {
	AppendBlobCommittedBlockCount interface{}              `json:"appendBlobCommittedBlockCount"`
	BlobTier                      interface{}              `json:"blobTier"`
	BlobTierChangeTime            interface{}              `json:"blobTierChangeTime"`
	BlobTierInferred              bool                     `json:"blobTierInferred"`
	BlobType                      string                   `json:"blobType"`
	ContentLength                 int64                    `json:"contentLength"`
	ContentRange                  string                   `json:"contentRange"`
	ContentSettings               AzureBlobContentSettings `json:"contentSettings"`
	Copy                          AzureBlobCopy            `json:"copy"`
	CreationTime                  string                   `json:"creationTime"`
	DeletedTime                   interface{}              `json:"deletedTime"`
	Etag                          string                   `json:"etag"`
	LastModified                  string                   `json:"lastModified"`
	Lease                         AzureBlobLease           `json:"lease"`
	PageBlobSequenceNumber        interface{}              `json:"pageBlobSequenceNumber"`
	RemainingRetentionDays        interface{}              `json:"remainingRetentionDays"`
	ServerEncrypted               bool                     `json:"serverEncrypted"`
}

type AzureBlobContentSettings struct {
	CacheControl       interface{} `json:"cacheControl"`
	ContentDisposition interface{} `json:"contentDisposition"`
	ContentEncoding    interface{} `json:"contentEncoding"`
	ContentLanguage    interface{} `json:"contentLanguage"`
	ContentMd5         string      `json:"contentMd5"`
	ContentType        string      `json:"contentType"`
}

type AzureBlobCopy struct {
	CompletionTime    interface{} `json:"completionTime"`
	ID                interface{} `json:"id"`
	Progress          interface{} `json:"progress"`
	Source            interface{} `json:"source"`
	Status            interface{} `json:"status"`
	StatusDescription interface{} `json:"statusDescription"`
}

type AzureBlobLease struct {
	Duration interface{} `json:"duration"`
	State    string      `json:"state"`
	Status   string      `json:"status"`
}
