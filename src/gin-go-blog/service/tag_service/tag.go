package tag_service

import (
	"encoding/json"
	"gin-blog/models"
	"gin-blog/pkg/export"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/service/cache_service"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	tag := map[string]interface{}{
		"modified_by": t.ModifiedBy,
		"name":        t.Name,
	}

	if t.State >= 0 {
		tag["state"] = t.State
	}

	return models.EditTag(t.ID, tag)
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]*models.Tag, error) {
	var (
		tags, cacheTags []*models.Tag
	)

	cache := cache_service.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()

	if gredis.Exist(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())

	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["delete_on"] = 0

	if t.State >= 0 {
		maps["state"] = t.State
	}

	if t.Name != "" {
		maps["name"] = t.Name
	}

	return maps
}

func (t *Tag) CleanAll() error {
	return models.CleanAllTag()
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	log.Println("tag_service tag.go Export()")
	log.Println(tags)

	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titiles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titiles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, tag := range tags {
		values := []string{
			strconv.Itoa(tag.ID),
			tag.Name,
			tag.CreatedBy,
			strconv.Itoa(tag.CreatedOn),
			tag.ModifiedBy,
			strconv.Itoa(tag.ModifiedOn),
		}
		log.Println(values)
		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tag-" + time + ".xlsx"

	fullpath := export.GetExcelFullPath() + filename
	err = file.Save(fullpath)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (t *Tag) Import(r io.Reader) error {

	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("标签信息")
	for i, row := range rows {
		if i > 0 {
			var data []string
			for _, colCell := range row {
				data = append(data, colCell)
			}

			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}
