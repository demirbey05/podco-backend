package core_test

import (
	"testing"

	"github.com/demirbey05/auth-demo/internal/core"
	"github.com/joho/godotenv"
)

var article = `Okay, here's a Medium-style article based on the provided podcast transcript. I've aimed for clarity, readability, and the kind of engaging tone you often find on the platform:

**The Fractal Intelligence of LLMs and the Looming Reality of Cost**

The world of AI is rapidly evolving, and with it, our understanding of what large language models (LLMs) can truly do. We've moved from awe at their ability to generate text to a more nuanced perspective, grappling with their limitations and the implications for the future of AI development. Recently, the buzz has centered around "fractal intelligence," a term coined to describe the seemingly unpredictable nature of these models: they work, until they don't, and it can be tough to say exactly *when* and *why* they'll falter.

But as fascinating as this may be, the field can't rest on the laurels of this nebulous understanding. The next stage of AI research must be about defining the dependable parameters of their performance.  

**Moving Past the "Fractal" Stage**

We need to move beyond the idea of LLMs as mysterious, unpredictable oracles.  Reasoning and logic have been formalized with limitations, things like "limited depth", "limited look ahead," which none of those seem to apply to LLMs at all and instead we're left with this "fractal intelligence." This points to the need for a deeper scientific inquiry, one that goes beyond the surface to quantify and understand their capacities. 

But this isn't just about abstract understanding, it's about practicality. The conversation is now turning to *cost*. Computer science isn't only about making the impossible possible, but also about the unsexy reality of expense. Think of the race to the moon. The initial efforts were about proving capability, cost was a secondary concern. Today, it's Elon Musk pushing the boundaries of space exploration because *cost* now matters. This shift applies equally to the world of AI; we're entering a phase where efficiency and value are becoming the next hurdles to clear.

**The Pitfalls of "Reasoning" by Retrieval**

Some people have begun to argue that if an LLM doesn't simply "retrieve" information, it's therefore "reasoning." This reminds one of the classic Monty Python witch trial sketch, where a series of absurd, loosely connected arguments lead to the conclusion that someone is a witch. Just because an LLM isn't pulling an answer directly from memory doesn't automatically mean it's engaging in sound reasoning, just as simply connecting random concepts doesn't automatically equate to actual logic.

**01: A New Approach?**

One of the most talked about developments in the field lately is 01, a new AI research lab emerging in Zurich.  They seem to be taking a page from the book of DeepMind by tackling reasoning tasks within LLM systems.  Specifically, they're exploring the intersection of LLM systems and search methods, inspired by the AlphaGo approach.  A key player is SensML, a compute platform optimized for AI, supporting open-source models such as Llama, with flexible pricing options that range from pay-as-you-go to constant model use.

**The Shifting Landscape Since Vienna**

Since our last deep dive into the space of LLMs at ICML in Vienna, there has been a massive shift in approaches. Back then, the focus was on the autoregressive models, where a model predicts the next word or token based on the words that have been given. Their abilities are amazing at generating creative ideas but with no guarantee of correctness.  This opened up the question of how to get LLMs to engage with more complex reasoning tasks.

**Inference Time Scaling and the Rise of "Magical Tokens"**

The field has taken some interesting turns in the development of what is often called inference time scaling, and post training techniques. This was first seen in the idea that, if an LLM produces candidates without any guarantee of correctness, perhaps we can get it to produce a huge number of candidates and then either do majority voting or self-consistency to figure out which one is right. The problem is, how do we define the right answer? The next attempt was based on an interesting observation: adding certain "magical tokens" to the prompt increased the likelihood of a correct answer.  This became known as Chain-of-Thought prompting. 

The first instance, zeroth order Chain-of-Thought, consisted of the same magical tokens being added regardless of the prompt and was inspired by human data. This led to the development of Chain-of-Thought task-specific prompting, where humans added task-specific advice to the prompts. While this does increase the time spent during inference and has shown some promise, the problem is, where are these magical tokens coming from?

**Synthetic Data and Reinforcement Learning**

Humans supplying these tokens was proving extremely costly.  Another idea was to focus on problems with systematic solvers, such as arithmetic, search, and planning and use them to produce traces of data manipulation operations.  These traces became the training data, with the hope that the LLM could then mimic the process without needing to be linked to external solvers. A few groups at Meta and Google have been following this approach, though it has the limitations that the solution is based on imitation of the trace.

The next evolution in thought is centered on the idea that the perfect tokens could be the result of reinforcement learning. Imagine an AlphaGo agent, but instead of moving pieces on a board, it's generating tokens that augment the prompt. The goal? Generate the tokens that will increase the likelihood of a correct answer.  

This involves: a smaller LLM whose entire function is to find prompt augmentation tokens, and a bigger LLM which produces a solution based on the modified prompt. This system tries different sequences, then uses the known answers to train a Q function to output tokens which will ultimately result in the correct answer. Once the training is complete, a model has been made ready for inference time. This also helps explain the costs for inference, with many more tokens being generated as the model works to improve results.

**The Evolution of Reasoning**

This new structure is radically different from the autoregressive LLMs that have been the basis for many current advancements. It goes beyond simply generating text; it now incorporates something akin to a learning mechanism to determine the best way to arrive at a solution. This approach could explain the results in the Strawberry Fields paper, where 01 outperformed state-of-the-art LLMs on blocks world and mystery domains by having something closer to reasoning as opposed to mere pattern matching.

**So, is 01 Truly "Reasoning?"**

This leaves a question: Are these systems truly "reasoning," or simply engaging in more sophisticated pattern matching? The field is still attempting to determine the difference between true reasoning and retrieval, however sophisticated.  The traditional definitions of reasoning from logic provide a better framework for comparison, but this approach does not preclude the fact that LLMs can still be put to use in safety critical scenarios.

It's easy to get caught up in the potential, or as they say "optimism vs pessimism." There are those who believe that LLMs are already showing true reasoning skills, and others who see these as simply advanced pattern recognition. The reality probably lies somewhere in between, and there may be a "blind men and the elephant" type of situation, with each perspective only looking at certain specific perspectives of the bigger picture. However, there is a widespread awareness that what may have been seen before as AI does not have the type of capabilities necessary.

**The Road Ahead**

Ultimately, the field of AI will continue to evolve as we push these new technologies. It's essential to move beyond the current fascination with mere output and look instead toward a future where AI systems can offer guarantees about their reasoning and a dependable understanding of what they can do. The future will be determined by the balance between the power of general-purpose models and the efficiency of highly specialized tools. This is a journey into a future where AI is not just a fascinating technology, but a dependable and accountable force in our world.

---
**Key takeaways from the article:**

*   **"Fractal Intelligence" is Not Enough**:  LLMs are currently too unpredictable; we need to characterize their capabilities beyond this label.
*   **Cost is King**:  The excitement of innovation is shifting to considerations of efficiency and expense.
*   **Reasoning vs. Retrieval**:  Just because a model isn't retrieving information directly doesn't mean it's truly reasoning. 
*   **The 01 Approach**:  01's model may be engaging in a reinforcement learning approach, building better performance through pseudo-actions.
*   **LLMs as Tools:** For human in the loop models, LLMs and LLRMs are great intelligence amplifiers and the next stages of exploration will be in the area of "compound systems."
*  **LLM's vs LRM's**: LLMs and LRM's are very different kinds of systems, with different strengths and weaknesses and different associated costs.

Let me know if you want any changes or further adjustments!"`

func TestGenerateQuizzesFromArticle(t *testing.T) {

	// Load environment variables same as main code
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found")
	}

	quiz, err := core.GenerateQuizzesFromArticle(article, "Turkish")
	if err != nil {
		t.Fatalf("GenerateQuizzesFromArticle failed: %v", err)
	}

	if quiz == nil {
		t.Fatal("Expected non-nil quiz response")
	}

	// Basic validation of quiz structure
	if len(quiz.Questions) != 5 {
		t.Errorf("Expected 5 questions, got %d", len(quiz.Questions))
	}

	for i, q := range quiz.Questions {
		if q.Question == "" {
			t.Errorf("Question %d has empty text", i)
		}
		if len(q.Options) != 4 {
			t.Errorf("Question %d has %d options, expected 4", i, len(q.Options))
		}
		if q.Answer < 0 || q.Answer > 3 {
			t.Errorf("Question %d has invalid answer index %d", i, q.Answer)
		}
	}
}
