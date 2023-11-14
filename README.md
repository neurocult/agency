# Inspiration

- https://python.langchain.com/
- https://github.com/tmc/langchaingo
- https://github.com/yoheinakajima/babyagi
- https://github.com/cpacker/MemGPT
- https://github.com/Significant-Gravitas/AutoGPT
- https://github.com/OpenBMB/XAgent
- https://github.com/sashabaranov/go-openai

# Use-Cases

```
// easy:
// continue text
// rewrite text
// templating (dynamic translator e.g.)
// voice -> text
// text -> image
// image -> text
// voice -> text -> image
// image -> text -> voice
// function call (print user input message)

// medium:
// chat (use assistant API?)
// chat with voice from user
// chat with voice from both user and system
// text -> summarize with small gpt3 -> answer with gpt4
// image -> image (edit image / variations)
// sequential multiple function calls
// concurrent multiple function calls
// conditional function calls (depending on user's input)

// hard:
// RAG - (file-system, vector database, relational database)
// fine-tuning - (make it answer in a desired form)
// phone-call - (similar to how chatGPT works)
// group-function calls - (like google meet works, maybe as a google meet member)

// extra-hard (autonomous agents)
// find a job for me
// help me with my coding pet-project
// help me with my book
// be my life-copilot (read my notes, give me insights)
```