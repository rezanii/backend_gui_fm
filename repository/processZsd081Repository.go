package repository

import (
	"backend_gui/dto"
	"backend_gui/models"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProcessZsd081Repository interface {
	InsertDataZsd081(detail models.FpmDetailInputZsd081) error
	GetHistoryZsd081FileByName(fileName string) dto.HistoryZsd081
	SaveDataHistoryZsd081(data models.FpmHistoryFileZsd081) int
	GetFileZsd081Process(fileName string) dto.FileProcessZsd081
	FlagFileZsd081Process(fileName string) error
	UpdateStatusZsd081(data models.FpmHistoryFileZsd081) error
}

type mapProcessZsd081RepositoryCon struct {
	mapProcessFileRepositoryCon *gorm.DB
}

func (m mapProcessZsd081RepositoryCon) UpdateStatusZsd081(data models.FpmHistoryFileZsd081) error {
	errUpdate:= m.mapProcessFileRepositoryCon.Save(&data)
	if errUpdate.Error != nil{
		log.Error("update sts zsd081 data ",data, "with error :"+errUpdate.Error.Error())
		return fmt.Errorf("%s", "Failed to update status ", data.FileName)
	}
	return  nil
}

func (m mapProcessZsd081RepositoryCon) FlagFileZsd081Process(fileName string) error {
	err:= m.mapProcessFileRepositoryCon.Save(&models.FpmListFileProcess{
		FileName: fileName,
	})
	if err.Error != nil{
		log.Error("Error insert flag zsd081 with data ", fileName)
		return fmt.Errorf("%s", "Failed to flag "+fileName)
	}
	return nil
}

func (m mapProcessZsd081RepositoryCon) GetFileZsd081Process(fileName string) dto.FileProcessZsd081 {
	fileProcessModels := models.FpmListFileProcess{}
	m.mapProcessFileRepositoryCon.Where("file_name=?", fileName).First(&fileProcessModels)
	var fileProcessDto dto.FileProcessZsd081
	fileProcessDto.ID = fileProcessModels.ID
	fileProcessDto.FileName = fileProcessModels.FileName
	return fileProcessDto

}

func (m mapProcessZsd081RepositoryCon) InsertDataZsd081(detail models.FpmDetailInputZsd081) error {
	err := m.mapProcessFileRepositoryCon.Save(&detail)
	if err.Error != nil {
		log.Error("Error insert zsd081 with data ", detail, "with error :"+err.Error.Error())
		return fmt.Errorf("%s", "Failed to save data zsd081")
	}
	return nil

}

func (m mapProcessZsd081RepositoryCon) SaveDataHistoryZsd081(data models.FpmHistoryFileZsd081) int {
	id := 0
	err := m.mapProcessFileRepositoryCon.Save(&data)
	if err.Error != nil {
		log.Error("error save history zsd081 with data ", data)
		return id
	}
	id = data.ID
	return id

}

func (m mapProcessZsd081RepositoryCon) GetHistoryZsd081FileByName(fileName string) dto.HistoryZsd081 {
	var data models.FpmHistoryFileZsd081
	m.mapProcessFileRepositoryCon.Where("file_name=?", fileName).First(&data)
	dataHistory := dto.HistoryZsd081{
		ID:         data.ID,
		UploadBy:   data.UploadBy,
		FileName:   data.FileName,
		StatusFile: data.StatusFile,
	}
	return dataHistory
}

func InstanceProcessZsd081Repository(db *gorm.DB) ProcessZsd081Repository {
	return &mapProcessZsd081RepositoryCon{
		mapProcessFileRepositoryCon: db,
	}
}
