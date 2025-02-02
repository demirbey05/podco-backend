package core

import "fmt"

var rawArticlePrompt = `I will send you a transcription of the podcast. 
Can you generate medium like article from it?
In generating markdown you can use following elements :
- Headings
- Lists
- Tables
- Bold and italic text`

func generateArticlePrompt(transcript string) string {

	return fmt.Sprintf("%s\n\n%s", rawArticlePrompt, transcript)

}

var rawQuizPrompt = `Generate a quiz based on this article. Follow these rules:
1. Each question must have exactly 4 options
2. Options must be plausible distractors
3. True answer index must be 0-3
4. Output must be valid JSON matching this format:
{
	"questions": [
		{
			"question": "...",
			"options": ["a", "b", "c", "d"],
			"true_answer_index": 0
		}
	]
}

The number of questions is up to you.
`

func generateQuizPrompt(article string) string {
	return fmt.Sprintf("%s\n\n%s", rawQuizPrompt, article)
}
