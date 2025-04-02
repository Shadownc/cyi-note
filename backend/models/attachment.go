package models

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"gorm.io/gorm"
)

// Attachment 附件模型
type Attachment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NoteID    *uint          `gorm:"index;null" json:"note_id"` // 修改为指针类型，允许为NULL
	UserID    uint           `gorm:"index;not null;default:1" json:"user_id"` // 为现有记录设置默认值
	Filename  string         `gorm:"size:255;not null" json:"filename"`
	Filepath  string         `gorm:"size:255;not null" json:"-"` // 文件存储路径，不返回给前端
	FileURL   string         `gorm:"-" json:"file_url"`          // 文件访问URL，计算属性，不存储在数据库
	Filetype  string         `gorm:"size:100" json:"filetype"`
	Filesize  int64          `json:"filesize"`
	IsTemp    bool           `gorm:"default:false" json:"is_temp"` // 是否是临时附件
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Note Note `gorm:"foreignKey:NoteID" json:"-"`
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// AfterFind 查询后自动设置文件URL
func (a *Attachment) AfterFind(tx *gorm.DB) error {
	// 使用ID而不是文件路径构建URL
	a.FileURL = "/api/attachments/" + fmt.Sprintf("%d", a.ID)
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
	// 使用IS NOT NULL条件确保只查询已关联到笔记的附件
	err := DB.Where("note_id = ?", noteID).Find(&attachments).Error
	return attachments, err
}

// SaveFile 保存文件并创建附件记录
func SaveFile(file *os.File, filename string, noteID uint, fileType string, userID uint, isTemp bool) (*Attachment, error) {
	// 创建上传目录
	var uploadDir string
	if isTemp {
		// 临时文件保存在临时目录
		uploadDir = filepath.Join("uploads", "temp", fmt.Sprintf("%d", userID))
	} else {
		// 普通附件保存在uploads目录
		uploadDir = filepath.Join("uploads", fmt.Sprintf("%d", userID))
	}
	
	// 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}
	
	// 生成文件路径
	timestamp := time.Now().UnixNano()
	randomSuffix := fmt.Sprintf("%x", timestamp)[:8] // 使用时间戳的前8位十六进制作为随机后缀
	fileExt := filepath.Ext(filename)
	if fileExt == "" {
		// 如果没有扩展名，尝试从Content-Type推断
		fileExt = getExtensionFromContentType(fileType)
	}
	
	// 安全化文件名，只保留基本部分
	safeName := fmt.Sprintf("%d_%d_%s%s", 
		userID, timestamp, randomSuffix, fileExt)
	newFilepath := filepath.Join(uploadDir, safeName)
	
	// 创建目标文件
	dst, err := os.Create(newFilepath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()
	
	// 复制文件内容
	if _, err = io.Copy(dst, file); err != nil {
		// 如果复制失败，删除已创建的文件
		os.Remove(newFilepath)
		return nil, fmt.Errorf("复制文件内容失败: %v", err)
	}
	
	// 获取文件信息，包括大小
	fileInfo, err := os.Stat(newFilepath)
	if err != nil {
		os.Remove(newFilepath)
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}
	
	// 创建附件记录
	attachment := &Attachment{
		UserID:   userID,
		Filename: filename,
		Filepath: newFilepath,
		Filetype: fileType,
		Filesize: fileInfo.Size(),
		IsTemp:   isTemp,
	}
	
	// 只有在非临时附件且提供了noteID的情况下，才设置NoteID
	if !isTemp && noteID > 0 {
		attachment.NoteID = &noteID
	}
	
	if err := DB.Create(attachment).Error; err != nil {
		// 如果创建失败，删除已创建的文件
		os.Remove(newFilepath)
		return nil, fmt.Errorf("创建附件记录失败: %v", err)
	}
	
	return attachment, nil
}

// getExtensionFromContentType 根据Content-Type猜测文件扩展名
func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "application/pdf":
		return ".pdf"
	case "text/plain":
		return ".txt"
	case "text/markdown":
		return ".md"
	case "application/json":
		return ".json"
	default:
		return ""
	}
}

// DeleteAttachment 删除附件
func DeleteAttachment(id uint) error {
	return DB.Delete(&Attachment{}, id).Error
}

// DateGroup 日期分组结构
type DateGroup struct {
	Date        string `json:"date"`
	Count       int    `json:"count"`
	DisplayDate string `json:"displayDate"`
}

// GetAttachmentsByDate 按日期分组获取附件
func GetAttachmentsByDate(userID uint, page, pageSize int, fileType string) ([]Attachment, []DateGroup, int64, error) {
	var attachments []Attachment
	var total int64
	
	// 基础查询：获取属于用户的所有附件
	query := DB.Model(&Attachment{}).
		Joins("JOIN notes ON notes.id = attachments.note_id").
		Where("notes.user_id = ? AND attachments.note_id IS NOT NULL", userID)
	
	// 应用文件类型过滤（如果有）
	if fileType != "" {
		query = query.Where("attachments.filetype LIKE ?", fileType+"%")
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}
	
	// 获取分页数据
	if err := query.
		Order("attachments.created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&attachments).Error; err != nil {
		return nil, nil, 0, err
	}
	
	// 查询日期分组数据
	var dateGroups []DateGroup
	var rows *sql.Rows
	var err error
	
	// 检查数据库类型，根据不同类型使用不同的SQL语法
	dbType := DB.Dialector.Name()
	
	// 构建日期查询
	var dateQuery string
	params := []interface{}{userID}
	
	if dbType == "sqlite" {
		// SQLite 日期函数
		dateQuery = `
			SELECT 
				strftime('%Y-%m-%d', attachments.created_at) as date,
				COUNT(attachments.id) as count
			FROM attachments
			JOIN notes ON notes.id = attachments.note_id
			WHERE notes.user_id = ? AND attachments.deleted_at IS NULL
                AND attachments.note_id IS NOT NULL
		`
	} else if dbType == "mysql" || dbType == "postgres" {
		// MySQL/PostgreSQL 日期函数
		dateQuery = `
			SELECT 
				DATE(attachments.created_at) as date,
				COUNT(attachments.id) as count
			FROM attachments
			JOIN notes ON notes.id = attachments.note_id
			WHERE notes.user_id = ? AND attachments.deleted_at IS NULL
                AND attachments.note_id IS NOT NULL
		`
	} else {
		// 兜底方案：直接返回时间戳，后续处理
		dateQuery = `
			SELECT 
				cast(attachments.created_at as text) as date,
				COUNT(attachments.id) as count
			FROM attachments
			JOIN notes ON notes.id = attachments.note_id
			WHERE notes.user_id = ? AND attachments.deleted_at IS NULL
                AND attachments.note_id IS NOT NULL
		`
	}
	
	if fileType != "" {
		dateQuery += " AND attachments.filetype LIKE ?"
		params = append(params, fileType+"%")
	}
	
	// 分组和排序
	if dbType == "sqlite" {
		dateQuery += " GROUP BY strftime('%Y-%m-%d', attachments.created_at) ORDER BY date DESC"
	} else if dbType == "mysql" || dbType == "postgres" {
		dateQuery += " GROUP BY DATE(attachments.created_at) ORDER BY date DESC"
	} else {
		dateQuery += " GROUP BY cast(attachments.created_at as text) ORDER BY date DESC"
	}
	
	// 执行查询
	rows, err = DB.Raw(dateQuery, params...).Rows()
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var group DateGroup
		var date string
		var count int
		
		if err := rows.Scan(&date, &count); err != nil {
			continue
		}
		
		// 标准化日期格式
		// 如果返回的是完整日期时间，提取日期部分
		if len(date) > 10 && strings.Contains(date, "T") {
			date = strings.Split(date, "T")[0]
		} else if len(date) > 10 {
			// 如果是其他格式的日期时间，只保留前10个字符（YYYY-MM-DD）
			date = date[:10]
		}
		
		// 尝试解析日期
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			// 尝试解析其他可能的日期格式
			formats := []string{"2006-01-02", "01/02/2006", "02-01-2006"}
			for _, format := range formats {
				t, err = time.Parse(format, date)
				if err == nil {
					break
				}
			}
			
			// 如果还是无法解析，使用当前日期
			if err != nil {
				t = time.Now()
				date = t.Format("2006-01-02")
			}
		}
		
		group.Date = date
		group.Count = count
		group.DisplayDate = formatDisplayDate(t)
		
		dateGroups = append(dateGroups, group)
	}
	
	return attachments, dateGroups, total, nil
}

// formatDisplayDate 格式化显示日期
func formatDisplayDate(t time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	
	if t.Year() == today.Year() && t.Month() == today.Month() && t.Day() == today.Day() {
		return "今天"
	} else if t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day() {
		return "昨天"
	} else if t.Year() == now.Year() {
		return t.Format("01月02日")
	}
	
	return t.Format("2006年01月02日")
}

// UpdateAttachment 更新附件
func UpdateAttachment(attachment *Attachment) error {
	return DB.Save(attachment).Error
} 