package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/RazvanBerbece/Aztebot/pkg/shared/embed"
)

func GetTextFromFile(filepath string) string {
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

func GetLongEmbedFromStaticData(embedTextData string) *embed.Embed {

	// Make sure that the newline characters are encoded correctly
	embedTextData = strings.Replace(embedTextData, `\n`, "\n", -1)

	embed := embed.NewEmbed()

	// Split the content into sections based on double newlines
	sections := strings.Split(embedTextData, "\n\n")
	for idx, section := range sections {
		lines := strings.Split(section, "\n")
		if len(lines) > 0 {
			// Use the first line as the title and the rest as content
			title := lines[0]
			content := strings.Join(lines[1:], "\n")
			// Replace formatting character sequences
			formattedTitle := strings.Replace(title, `\br`, "\n", -1)
			formattedContent := strings.Replace(content, `\br`, "\n", -1)
			// Add section to embed
			embed.AddField(formattedTitle, formattedContent, false)
			// If not last section, add a new line for aesthetic purposes
			if idx < len(sections)-1 {
				embed.AddLineBreakField()
			}
		}
	}

	return embed
}

func GetSectionContent(section string) (string, error) {

	// The section content is the second substring from first newline till the last newline
	newlineIndex := strings.IndexByte(section, '\n')

	if newlineIndex != -1 {
		// Extract the substring starting from the character right after the newline character
		content := section[newlineIndex+1:]
		fmt.Println("Substring after the first newline character ->", content, " | ")
		return content, nil
	} else {
		return "", fmt.Errorf("no newline character found in the string")
	}

}
