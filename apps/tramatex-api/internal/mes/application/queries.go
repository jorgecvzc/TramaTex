package application

import (
	"github.com/google/uuid"
)

type ListTasksQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type GetTaskByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type ListPositionsQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type GetPositionByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type ListServiceGroupsQuery struct {
	IsActive *bool  `form:"is_active"`
	Search   string `form:"search"`
}

type ListServiceTemplatesQuery = ListServiceGroupsQuery

type GetServiceGroupByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type GetServiceTemplateByIDQuery = GetServiceGroupByIDQuery

type ListMESWorksQuery struct {
	Status  *string `form:"status"`
	Search  string  `form:"search"`
	PartyID string  `form:"party_id"`
}

type ListWorkDefinitionsQuery = ListMESWorksQuery

type GetMESWorkByIDQuery struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

type GetWorkDefinitionByIDQuery = GetMESWorkByIDQuery

type ListOverdueMESWorksQuery struct {
	Limit int `form:"limit"`
}

type ListOverdueWorkDefinitionsQuery = ListOverdueMESWorksQuery
