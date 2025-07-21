package pages

import (
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/generated"
)

type PageWithBlocks struct {
	Page   generated.Page           `json:"page"`
	Blocks []*blocks.BlockWithProps `json:"blocks"`
}

type UpdatePageRequest struct {
	Page     generated.UpdatePageParams `json:"page"`
	Blocks   []generated.Block          `json:"blocks"`
	PropsMap map[uuid.UUID]any          `json:"propsMap"`
}

type UpdatePageWithBlocksRequest struct {
	Page   generated.UpdatePageParams `json:"page"`
	Blocks []blocks.BlockWithProps    `json:"blocks"`
}

type CreatePageRequest struct {
	Page   generated.CreatePageParams `json:"page"`
	Blocks []blocks.BlockWithProps    `json:"blocks"`
}
