# Monit-telegram

Sistema automatic per fer crides web (GET de moment), cada "minutes2call" vegades i si tot va bé cada "times2ok" vegades enviar un missatge a un canal de telegram, i si va malament enviar l'error deseguida a un altre canal (potser el mateix).

## Configuració

```json
{
"urls": 
	[
		{ "url": "https://example.com", "header": null },
		{ "url": "https://github.com",
		"header": [
			{ "key": "accept", "value": "application/json" },
			{ "key": "key", "value": "value" }
		]
		},
		{ "url": "https://test.com", "header": null }
	],
"bot": { "token" : "9999:XXXXXX", "ChannelOk": -10000000000, "ChannelKo": -10000000001 },
"minutes2call" : 1,
"times2ok" : 60, 
"debug": false
}
```

