package competitor

import "github.com/wizenheimer/iris/internal/domain/models"

type CompetitorMapper struct{}

func NewCompetitorMapper() *CompetitorMapper {
	return &CompetitorMapper{}
}

func (m *CompetitorMapper) ToDTO(entity *models.Competitor) *models.CompetitorDTO {
	if entity == nil {
		return nil
	}

	return &models.CompetitorDTO{
		ID:        entity.ID,
		Name:      entity.Name,
		Domain:    entity.Domain,
		URLs:      entity.URLs,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (m *CompetitorMapper) ToEntity(dto *models.CompetitorInput) *models.Competitor {
	if dto == nil {
		return nil
	}

	return &models.Competitor{
		Name:   dto.Name,
		Domain: dto.Domain,
		URLs:   dto.URLs,
	}
}
