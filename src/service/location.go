package service

import (
	"dev.chaiyapluek.cloud.final.backend/src/dto"
	appError "dev.chaiyapluek.cloud.final.backend/src/pkg/error"
	"dev.chaiyapluek.cloud.final.backend/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocationService interface {
	GetAllLocation() ([]*dto.LocationResponse, error)
	GetLocationById(id string) (*dto.SingleLocationResponse, error)
	GetLocationMenu(locationId, menuId string) (*dto.MenuResponse, error)
}

type locationServiceImpl struct {
	locationRepository repository.LocationRepository
}

func NewLocationService(locationRepo repository.LocationRepository) LocationService {
	return &locationServiceImpl{
		locationRepository: locationRepo,
	}
}

func (s *locationServiceImpl) GetAllLocation() ([]*dto.LocationResponse, error) {
	locations, err := s.locationRepository.GetAllLocation()
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	resp := []*dto.LocationResponse{}
	for idx := range locations {
		resp = append(resp, &dto.LocationResponse{
			Id:   locations[idx].Id.Hex(),
			Name: locations[idx].Name,
		})
	}
	return resp, nil
}

func (s *locationServiceImpl) GetLocationById(id string) (*dto.SingleLocationResponse, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid location id format")
	}

	location, menus, err := s.locationRepository.GetLocationById(&oid)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	if location == nil {
		return nil, appError.NewErrNotFound("location not found")
	}

	resp := &dto.SingleLocationResponse{
		Id:    location.Id.Hex(),
		Name:  location.Name,
		Menus: []*dto.MenuResponse{},
	}
	for idx := range menus {
		resp.Menus = append(resp.Menus, &dto.MenuResponse{
			Id:             menus[idx].Id.Hex(),
			Name:           menus[idx].Name,
			Description:    menus[idx].Description,
			Price:          menus[idx].Price,
			IconImage:      menus[idx].IconImage,
			ThumbnailImage: menus[idx].ThumbnailImage,
		})
	}

	return resp, nil
}

func (s *locationServiceImpl) GetLocationMenu(locationId, menuId string) (*dto.MenuResponse, error) {
	oLocationId, err := primitive.ObjectIDFromHex(locationId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid location id format")
	}

	oMenuId, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return nil, appError.NewErrBadRequest("invalid menu id format")
	}

	menu, err := s.locationRepository.GetMenuItmes(&oLocationId, &oMenuId)
	if err != nil {
		return nil, appError.NewErrInternalServerError(err.Error())
	}

	if menu == nil {
		return nil, appError.NewErrNotFound("menu not found")
	}

	steps := []*dto.StepResponse{}
	for idx := range menu.Steps {
		step := menu.Steps[idx]
		options := []*dto.OptionResponse{}
		for idx := range step.Options {
			options = append(options, &dto.OptionResponse{
				Name:  step.Options[idx].Name,
				Value: step.Options[idx].Value,
				Price: step.Options[idx].Price,
			})
		}
		steps = append(steps, &dto.StepResponse{
			Name:        step.Name,
			Description: step.Description,
			Type:        step.Type,
			Required:    step.Required,
			Min:         step.Min,
			Max:         step.Max,
			Options:     options,
		})
	}

	menuResp := &dto.MenuResponse{
		Id:             menu.Id.Hex(),
		Name:           menu.Name,
		Description:    menu.Description,
		Price:          menu.Price,
		IconImage:      menu.IconImage,
		ThumbnailImage: menu.ThumbnailImage,
		Steps:          steps,
	}

	return menuResp, nil
}
