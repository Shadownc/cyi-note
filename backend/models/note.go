package models

import (
	"time"
	
	"gorm.io/gorm"
)

// Note 笔记模型
type Note struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Title     string         `gorm:"size:255;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Summary   string         `gorm:"size:500" json:"summary"` // AI生成的摘要
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	User        User          `gorm:"foreignKey:UserID" json:"-"`
	Tags        []*Tag        `gorm:"many2many:note_tags;" json:"tags"`
	Attachments []Attachment  `gorm:"foreignKey:NoteID" json:"attachments"`
}

// CreateNote 创建笔记
func CreateNote(note *Note) error {
	return DB.Create(note).Error
}

// GetNoteByID 通过ID获取笔记
func GetNoteByID(id uint) (*Note, error) {
	var note Note
	err := DB.Preload("Tags").Preload("Attachments").First(&note, id).Error
	return &note, err
}

// GetNotesByUserID 获取用户的所有笔记
func GetNotesByUserID(userID uint, page, pageSize int) ([]Note, int64, error) {
	var notes []Note
	var total int64
	
	query := DB.Model(&Note{}).Where("user_id = ?", userID)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").
		Preload("Tags").Preload("Attachments").Find(&notes).Error
	
	return notes, total, err
}

// SearchNotes 搜索笔记（标题、内容）
func SearchNotes(userID uint, keyword string, page, pageSize int) ([]Note, int64, error) {
	var notes []Note
	var total int64
	
	// 构建搜索条件
	query := DB.Model(&Note{}).Where("user_id = ?", userID).
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").
		Preload("Tags").Preload("Attachments").Find(&notes).Error
	
	return notes, total, err
}

// SearchNotesByTag 通过标签搜索笔记
func SearchNotesByTag(userID uint, tagID uint, page, pageSize int) ([]Note, int64, error) {
	var notes []Note
	var total int64
	
	// 构建包含标签关联的查询
	query := DB.Model(&Note{}).Where("notes.user_id = ?", userID).
		Joins("JOIN note_tags ON note_tags.note_id = notes.id").
		Where("note_tags.tag_id = ?", tagID)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("notes.created_at DESC").
		Preload("Tags").Preload("Attachments").Find(&notes).Error
	
	return notes, total, err
}

// UpdateNote 更新笔记
func UpdateNote(note *Note) error {
	return DB.Save(note).Error
}

// DeleteNote 删除笔记
func DeleteNote(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除笔记标签关联
		if err := tx.Where("note_id = ?", id).Delete(&NoteTag{}).Error; err != nil {
			return err
		}
		
		// 删除笔记附件
		if err := tx.Where("note_id = ?", id).Delete(&Attachment{}).Error; err != nil {
			return err
		}
		
		// 删除笔记
		return tx.Delete(&Note{}, id).Error
	})
}

// AddTagToNote 给笔记添加标签
func AddTagToNote(noteID, tagID uint) error {
	noteTag := NoteTag{
		NoteID: noteID,
		TagID:  tagID,
	}
	return DB.Create(&noteTag).Error
}

// RemoveTagFromNote 从笔记中移除标签
func RemoveTagFromNote(noteID, tagID uint) error {
	return DB.Where("note_id = ? AND tag_id = ?", noteID, tagID).Delete(&NoteTag{}).Error
}

// NoteTag 笔记和标签的关联表
type NoteTag struct {
	NoteID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

// TableName 指定表名
func (NoteTag) TableName() string {
	return "note_tags"
} 