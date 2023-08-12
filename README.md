Simple go application that lets you do chatgpt prompts from command line.
By default uses gpt-3.5-turbo model.

Windows batch files pipes the response to glow (required). https://github.com/charmbracelet/glow

gpt.bat runs main.go


**main.go** without 2nd argument (after the prompt) default behaviour.
![example](prompt.jpg)


**main.go** with 2nd argument (after the prompt) greater than zero lets you use different behaviour,
which can be defined in the personality variable.
![example2](prompt2.jpg)


gptm.bat runs memory.go


**memory.go** same behaviour flags, but saves conversation history locally to "message_history" file
located in application directory, it allows to make follow up questions.
![example3](prompt3.jpg)
