package core

import "fmt"

var rawArticlePrompt = `I will send you a transcription of the podcast and language of the user.
Article should be in the user's language.If content language is different than user's language please translate it to user's language. But be careful in technical terms.
Can you generate medium like article from it?
In generating markdown you can use following elements :
- Headings
- Lists
- Tables
- Bold and italic text`

func generateArticlePrompt(transcript, language string) string {

	return fmt.Sprintf("%s\n\n%s\n\nlanguage of user : %s", rawArticlePrompt, transcript, language)

}

var rawQuizPrompt = `Generate a quiz based on this article. Follow these rules:
Quiz should be in the user's language.If content language is different than user's language please translate it to user's language. But be careful in technical terms.
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

func generateQuizPrompt(article, language string) string {
	return fmt.Sprintf("%s\n\n%s\n\nlanguage of user : %s", rawQuizPrompt, article, language)
}
