# Inspiration

- https://python.langchain.com/
- https://github.com/tmc/langchaingo
- https://github.com/yoheinakajima/babyagi
- https://github.com/cpacker/MemGPT
- https://github.com/Significant-Gravitas/AutoGPT
- https://github.com/OpenBMB/XAgent
- https://github.com/sashabaranov/go-openai

# Description
Pure Go langchain alternative.

# Use-Cases

```

// templating (dynamic translator e.g.)
// text -> image
// voice -> text
// image -> text
// voice -> text -> image
// image -> text -> voice (awaiting for https://github.com/sashabaranov/go-openai/pull/557)
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

# TODO

## v0.1.0
-[ ] Name the project
-[ ] Name the organization
-[ ] Make readme with quickstart, description, installation, etc
-[ ] Reorganize the folders and packages
-[ ] Add examples
  -[ ] Add simple pipe example
  -[ ] Add chat example
  -[ ] Add chat multi-model example (voice=>text=>image)
  -[ ] Add templating example
-[ ] Add models support
  -[ ] Add ChatGPT support
  -[ ] Add Whisper support
  -[ ] Add Dall-E support
  -[ ] Add GPT Vision support (optional)

## Next versions
-[ ] Add support for external functions
-[ ] Add RAG example with function
-[ ] metadata (tokens used, audio dur, etc)
-[ ] Add llama support