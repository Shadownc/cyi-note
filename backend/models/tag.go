package models

import (
	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;uniqueIndex;not null" json:"name"`
	
	// 关联
	Notes []*Note `gorm:"many2many:note_tags;" json:"-"`
}

// TagWithCount 包含笔记数量的标签结构
type TagWithCount struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	NoteCount int64  `json:"noteCount"`
}

// CreateTag 创建标签
func CreateTag(tag *Tag) error {
	return DB.Create(tag).Error
}

// GetTagByID 通过ID获取标签
func GetTagByID(id uint) (*Tag, error) {
	var tag Tag
	err := DB.First(&tag, id).Error
	return &tag, err
}

// GetTagByName 通过名称获取标签
func GetTagByName(name string) (*Tag, error) {
	var tag Tag
	err := DB.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// GetAllTags 获取所有标签
func GetAllTags() ([]Tag, error) {
	var tags []Tag
	err := DB.Find(&tags).Error
	return tags, err
}

// GetAllTagsWithCount 获取所有标签，并包含每个标签关联的笔记数量
func GetAllTagsWithCount(userID uint) ([]TagWithCount, error) {
	var tagsWithCount []TagWithCount
	
	// SQL查询：获取标签及其关联的笔记数量
	err := DB.Raw(`
		SELECT t.id, t.name, COUNT(nt.note_id) as note_count
		FROM tags t
		LEFT JOIN note_tags nt ON t.id = nt.tag_id
		LEFT JOIN notes n ON nt.note_id = n.id AND n.deleted_at IS NULL
		WHERE (n.user_id = ? OR n.id IS NULL)
		GROUP BY t.id, t.name
		ORDER BY t.name
	`, userID).Scan(&tagsWithCount).Error
	
	return tagsWithCount, err
}

// GetTagNoteCount 获取标签关联的笔记数量
func GetTagNoteCount(tagID uint, userID uint) (int64, error) {
	var count int64
	
	err := DB.Model(&Note{}).
		Joins("JOIN note_tags ON note_tags.note_id = notes.id").
		Where("note_tags.tag_id = ? AND notes.user_id = ?", tagID, userID).
		Count(&count).Error
	
	return count, err
}

// GetTagsByNoteID 获取笔记的所有标签
func GetTagsByNoteID(noteID uint) ([]Tag, error) {
	var tags []Tag
	err := DB.Model(&Tag{}).
		Joins("JOIN note_tags ON note_tags.tag_id = tags.id").
		Where("note_tags.note_id = ?", noteID).
		Find(&tags).Error
	return tags, err
}

// DeleteTag 删除标签
func DeleteTag(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除标签与笔记的关联
		if err := tx.Where("tag_id = ?", id).Delete(&NoteTag{}).Error; err != nil {
			return err
		}
		
		// 删除标签
		return tx.Delete(&Tag{}, id).Error
	})
}

// GetOrCreateTag 获取或创建标签
func GetOrCreateTag(name string) (*Tag, error) {
	var tag Tag
	
	// 尝试查找标签
	if err := DB.Where("name = ?", name).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 标签不存在，创建新标签
			tag = Tag{Name: name}
			if err := DB.Create(&tag).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	return &tag, nil
} 