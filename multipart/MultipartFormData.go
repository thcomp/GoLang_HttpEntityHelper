package multipart

import (
	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

type MultipartFormData struct {
	helper *ThcompUtility.MultipartHelper
}

func (formData *MultipartFormData) EntityType() entity.HttpEntityType {
	return entity.MultipartFormData
}

func (formData *MultipartFormData) Close() error {
	if formData.helper != nil {
		return formData.helper.Close()
	} else {
		return entity.ErrNotExistFormData
	}
}

func (formData *MultipartFormData) Count() int {
	if formData.helper != nil {
		return formData.helper.Count()
	} else {
		return -1
	}
}

func (formData *MultipartFormData) GetByIndex(index int) (*ThcompUtility.FormData, error) {
	if formData.helper != nil {
		return formData.helper.GetByIndex(index)
	} else {
		return nil, entity.ErrNotExistFormData
	}
}

func (formData *MultipartFormData) GetByName(partName string) (*ThcompUtility.FormData, error) {
	if formData.helper != nil {
		return formData.helper.GetByName(partName)
	} else {
		return nil, entity.ErrNotExistFormData
	}
}

func (formData *MultipartFormData) PartNames() []string {
	if formData.helper != nil {
		return formData.helper.PartNames()
	} else {
		return nil
	}
}
