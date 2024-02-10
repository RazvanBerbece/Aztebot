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

	embed := embed.NewEmbed()

	// Make sure that the newline characters are encoded correctly
	// embedTextData = strings.Replace(embedTextData, `\n`, "\n", -1)

	// Split the content into sections based on double newline characters ("\n\n")
	sections := strings.Split(embedTextData, "\n\n")
	for _, section := range sections {
		lines := strings.Split(section, "\n")
		if len(lines) > 0 {
			// Use the first line as the title and the rest as content
			title := lines[0]
			content := strings.Join(lines[1:], "\n")
			embed.AddField(title, content, false)
			// if idx < len(sections)-1 {
			// 	embed.AddLineBreakField()
			// }
		}
	}

	return embed
}
