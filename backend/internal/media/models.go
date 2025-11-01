// backend/internal/media/models.go
package media

// SearchMediaParams represents query parameters for searching/filtering media
type SearchMediaParams struct {
	SearchQuery    string
	MimeTypeFilter string
	SortBy         string
	SortOrder      string
	Limit          int32
	Offset         int32
}

// VariantUploadResult holds all uploaded variant URLs and paths
type VariantUploadResult struct {
	OriginalURL  string
	LargeURL     string
	MediumURL    string
	ThumbnailURL string

	OriginalPath  string
	LargePath     string
	MediumPath    string
	ThumbnailPath string

	Width  int
	Height int
	Size   int64
}

// MediaStats contains media statistics
type MediaStats struct {
	TotalCount      int64 `json:"total_count"`
	TotalSize       int64 `json:"total_size"`
	UniqueUploaders int64 `json:"unique_uploaders"`
}
