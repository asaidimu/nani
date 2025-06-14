# [2.0.0](https://github.com/asaidimu/nani/compare/v1.0.0...v2.0.0) (2025-06-14)


* feat(ai)!: implement persistent sessions and JSON response format ([79f3414](https://github.com/asaidimu/nani/commit/79f3414a5d30ee9a53045cc7cf2f38df07881cc3))


### BREAKING CHANGES

* The internal AI model response format has changed from XML to JSON.
Any external tools or scripts directly consuming the AI's raw output will need
to be updated to parse JSON instead of XML.
The pkg/ai.AIClient interface has also been updated; consumers should review changes.

# 1.0.0 (2025-05-26)


* feat(app)!: Transform project into AI-powered TUI chat assistant ([e37c245](https://github.com/asaidimu/nani/commit/e37c24514cbeb2269fc30319237af271e0712972))


### BREAKING CHANGES

* The application has been completely redesigned from a simple 'Hello, World!' example to a full AI-powered terminal chat assistant.
Users must now obtain a Google Gemini API key and set it as the GEMINI_API_KEY environment variable before running the application.
The installation steps and usage paradigm have fundamentally changed. Refer to the updated README.md for new setup and build instructions, as well as usage documentation and keybindings.
