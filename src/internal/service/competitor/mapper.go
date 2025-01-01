package competitor

import (
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
)

type CompetitorMapper struct{}

func NewCompetitorMapper() *CompetitorMapper {
	return &CompetitorMapper{}
}

func (m *CompetitorMapper) ToDTO(entity *core_models.Competitor) *core_models.CompetitorDTO {
	if entity == nil {
		return nil
	}

	return &core_models.CompetitorDTO{
		ID:        entity.ID,
		Name:      entity.Name,
		Domain:    entity.Domain,
		URLs:      entity.URLs,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *CompetitorMapper) ToEntity(dto *core_models.CompetitorInput) *core_models.Competitor {
	if dto == nil {
		return nil
	}

	return &core_models.Competitor{
		Name:   dto.Name,
		Domain: dto.Domain,
		URLs:   dto.URLs,
	}
}
