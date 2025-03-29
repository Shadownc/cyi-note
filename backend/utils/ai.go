package utils

import (
	"strings"
	
	"github.com/go-ego/gse"
)

// 中文分词器
var segmenter gse.Segmenter

// 初始化分词器
func init() {
	segmenter = gse.Segmenter{}
	// 加载默认词典
	segmenter.LoadDict()
}

// ExtractKeywords 从文本中提取关键词作为标签
func ExtractKeywords(content string) ([]string, error) {
	// 简单实现：使用分词器提取关键词
	// 在实际应用中，可以使用更复杂的算法或调用外部API
	
	// 如果内容为空，返回空标签
	if strings.TrimSpace(content) == "" {
		return []string{}, nil
	}
	
	// 分词
	segments := segmenter.Cut(content, true)
	
	// 统计词频
	wordFreq := make(map[string]int)
	for _, word := range segments {
		// 过滤掉小于2个字符的词
		if len(word) < 2 {
			continue
		}
		wordFreq[word]++
	}
	
	// 按词频排序
	type wordCount struct {
		word  string
		count int
	}
	var wordCounts []wordCount
	for word, count := range wordFreq {
		wordCounts = append(wordCounts, wordCount{word, count})
	}
	
	// 简单排序
	for i := 0; i < len(wordCounts); i++ {
		for j := i + 1; j < len(wordCounts); j++ {
			if wordCounts[i].count < wordCounts[j].count {
				wordCounts[i], wordCounts[j] = wordCounts[j], wordCounts[i]
			}
		}
	}
	
	// 取前10个词作为标签
	var tags []string
	for i := 0; i < len(wordCounts) && i < 10; i++ {
		tags = append(tags, wordCounts[i].word)
	}
	
	return tags, nil
}

// GenerateSummary 生成文本摘要
func GenerateSummary(content string) (string, error) {
	// 如果内容为空，返回空摘要
	if strings.TrimSpace(content) == "" {
		return "", nil
	}
	
	// 简单实现：取内容的前200个字符作为摘要
	// 在实际应用中，可以使用NLP模型或调用外部API生成摘要
	
	// 移除所有HTML标签
	content = RemoveHTMLTags(content)
	
	// 取纯文本的前200个字符
	content = strings.TrimSpace(content)
	if len(content) > 200 {
		content = content[:200] + "..."
	}
	
	return content, nil
}

// GenerateTagSuggestions 生成标签建议（实际上是 ExtractKeywords 的别名）
func GenerateTagSuggestions(content string) ([]string, error) {
	return ExtractKeywords(content)
}

// RemoveHTMLTags 移除HTML标签
func RemoveHTMLTags(content string) string {
	// 简单实现：移除所有<>之间的内容
	inTag := false
	var result strings.Builder
	
	for _, r := range content {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	
	return result.String()
} 