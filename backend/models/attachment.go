package models

import (
	"os"
	"path/filepath"
	"time"
	
	"gorm.io/gorm"
)

// Attachment 附件模型
type Attachment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NoteID    uint           `gorm:"index;not null" json:"note_id"`
	Filename  string         `gorm:"size:255;not null" json:"filename"`
	Filepath  string         `gorm:"size:255;not null" json:"-"` // 文件存储路径，不返回给前端
	FileURL   string         `gorm:"-" json:"file_url"`          // 文件访问URL，计算属性，不存储在数据库
	Filetype  string         `gorm:"size:100" json:"filetype"`
	Filesize  int64          `json:"filesize"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Note Note `gorm:"foreignKey:NoteID" json:"-"`
}

// AfterFind 查询后自动设置文件URL
func (a *Attachment) AfterFind(tx *gorm.DB) error {
	a.FileURL = "/api/attachments/" + a.Filepath
	return nil
}

// BeforeDelete 删除前移除文件
func (a *Attachment) BeforeDelete(tx *gorm.DB) error {
	// 删除物理文件
	err := os.Remove(a.Filepath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// CreateAttachment 创建附件
func CreateAttachment(attachment *Attachment) error {
	return DB.Create(attachment).Error
}

// GetAttachmentByID 通过ID获取附件
func GetAttachmentByID(id uint) (*Attachment, error) {
	var attachment Attachment
	err := DB.First(&attachment, id).Error
	return &attachment, err
}

// GetAttachmentsByNoteID 获取笔记的所有附件
func GetAttachmentsByNoteID(noteID uint) ([]Attachment, error) {
	var attachments []Attachment
	err := DB.Where("note_id = ?", noteID).Find(&attachments).Error
	return attachments, err
}

// SaveFile 保存文件并创建附件记录
func SaveFile(file *os.File, filename string, noteID uint, fileType string) (*Attachment, error) {
	// 创建上传目录
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}
	
	// 生成文件路径
	timestamp := time.Now().Unix()
	newFilename := filepath.Join(uploadDir, string(timestamp)+"_"+filename)
	
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	
	// 创建目标文件
	dst, err := os.Create(newFilename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	
	// 复制文件内容
	/*
	if _, err = io.Copy(dst, file); err != nil {
		return nil, err
	}
	*/
	
	// 创建附件记录
	attachment := &Attachment{
		NoteID:   noteID,
		Filename: filename,
		Filepath: newFilename,
		Filetype: fileType,
		Filesize: fileInfo.Size(),
	}
	
	if err := DB.Create(attachment).Error; err != nil {
		return nil, err
	}
	
	return attachment, nil
}

// DeleteAttachment 删除附件
func DeleteAttachment(id uint) error {
	return DB.Delete(&Attachment{}, id).Error
} 