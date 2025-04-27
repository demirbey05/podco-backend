package core

import "fmt"

var rawArticlePrompt = `First, determine if this transcript contains substantive educational content suitable for creating an educational article and quiz.

If the content is NOT educational (such as movie scenes, music clips, casual conversations, or promotional material), respond ONLY with:
{
    "error": "The provided content does not appear to be educational. Please provide a transcript of educational content like a lecture, documentary, or informative podcast."
}

If the content IS educational, then you are an expert educator who deeply understands this subject. Transform this transcript into an educational article that effectively teaches the core concepts.

As an experienced teacher, you will:
1. IDENTIFY the 3-5 most important concepts or insights from the content
2. EXPLAIN these concepts clearly, as if teaching a class of engaged students
3. CONNECT ideas logically, building understanding progressively
4. ILLUSTRATE concepts with relevant examples, analogies, or applications
5. EMPHASIZE practical takeaways and "why this matters"

Structure your article with:
- A clear, informative title
- An introduction setting context and stating learning objectives
- Well-organized sections with descriptive headings
- A conclusion reinforcing key learning points

Format using appropriate markdown:
- ## Main Headings and ### Subheadings
- Bullet points for lists of related items
- **Bold** for key terms and important concepts
- Tables only when they enhance understanding

If translation is needed:
- Maintain the pedagogical clarity while adapting to the target language
- Preserve technical terminology with brief explanations where needed

Present only the educational article without any framing text or respond with the error JSON if the content is not educational.`

func generateArticlePrompt(transcript, language string) string {
	return fmt.Sprintf("%s\n\nTranscript: %s\n\nUser language: %s", rawArticlePrompt, transcript, language)
}

var rawQuizPrompt = `As an experienced educator who has just taught this material, create an assessment that effectively measures student understanding of the key concepts.

Your pedagogical approach should:
1. TEST MASTERY of the fundamental concepts rather than memorization of details
2. ASSESS different levels of understanding:
   - Basic comprehension of core ideas
   - Application of concepts to new situations
   - Analysis of relationships between concepts
3. PROVIDE questions that:
   - Are clearly worded as you would present them in class
   - Focus on what a good teacher would consider important
   - Challenge students to demonstrate true understanding
4. DESIGN thoughtful answer options that:
   - Include one clearly correct answer
   - Offer plausible distractors that reveal common misconceptions
   - Help identify gaps in understanding

Format your assessment as valid JSON:
{
    "questions": [
        {
            "question": "Clearly worded question testing an important concept",
            "options": ["Correct answer", "Plausible distractor", "Plausible distractor", "Plausible distractor"],
            "true_answer_index": 0,
            "explanation": "Brief explanation of the concept as you would explain it to a student"
        },
		{
            "question": "Clearly worded question testing an important concept",
            "options": ["Plausible distracter", "Plausible distractor", "Correct answer", "Plausible distractor"],
            "true_answer_index": 2,
            "explanation": "Brief explanation of the concept as you would explain it to a student"
        },
		{
            "question": "Clearly worded question testing an important concept",
            "options": ["Plausible distracter", "Correct answer", "Plausible distractor", "Plausible distractor"],
            "true_answer_index": 1,
            "explanation": "Brief explanation of the concept as you would explain it to a student"
        },
		{
            "question": "Clearly worded question testing an important concept",
            "options": ["Plausible distracter", "Plausible distractor", "Plausible distractor", "Correct answer"],
            "true_answer_index": 3,
            "explanation": "Brief explanation of the concept as you would explain it to a student"
        }
    ]
}

IMPORTANT: Randomize the position of correct answers across your questions. DO NOT place all correct answers in the same position (e.g., all at index 0). Deliberately vary the true_answer_index values (0, 1, 2, or 3) throughout the quiz.

If translation is needed:
- Ensure questions maintain their pedagogical clarity in the target language
- Preserve the educational value of both questions and explanations

Create 7-10 questions that collectively assess mastery of the material's most important concepts.`

func generateQuizPrompt(article, language string) string {
	return fmt.Sprintf("%s\n\nArticle: %s\n\nUser language: %s", rawQuizPrompt, article, language)
}
