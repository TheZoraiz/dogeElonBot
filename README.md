# Doge Elon Bot

A bot that notifies you via desktop notifications whenever Elon Musk posts a tweet about dogecoin (if it contains the keyword "doge"). It checks Elon's recent tweets every 3 minutes.

## Usage

Executables for both Windows and Linux 64 bit achitectures can be found in the build directory, or can retrieve them from the [release.](https://github.com/TheZoraiz/dogeElonBot/releases/tag/v1.1)

For usage of the code, you will need a twitter bearer token. After obtaining it, declare a const variable with your twitter bearer token in the main package like so:

const (
	bearerToken = "Bearer <bearer_token>"
)


## Packages used:

[github.com/buger/jsonparser](https://www.github.com/buger/jsonparser)

[github.com/faiface/beep](https://www.github.com/faiface/beep)

[github.com/gen2brain/beeep](https://www.github.com/gen2brain/beeep)
