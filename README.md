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
"bot": { "token" : "9999:XXXXXX", "ChannelOk": -,, "ChannelKo": -1001104876899 },
"minutes2call" : 1,
"times2ok" : 60, 
"debug": false
}
```

## Crear el Bot de telegram

Per crear el bot de telegram, s'ha de demanar-li al @BotFather, /newbot i guardar-se el token

Crear dos canals i posar-li el bot com a membre i fer-los public per esbrinar el numero
(Per fer-lu public has de fer un enllaç Posar-li nom al canal)

Fent aquest crida:

curl -v "https://api.telegram.org/bot<Token Bot>/sendMessage?chat_id=@<url>&text=<msg>"
La resposta és:

{"ok":true,"result":{"message_id":2,"from":{"id":97xxxxxxxx8,"is_bot":true,"first_name":"MxxxBot","username":"MxxxxBot"},"chat":{"id":-100xxxxxxxx6,"title":"XXX","username":"<url>","type":"supergroup"},"date":1xxxxx8,"text":"<msg>"}}

Amb tota la info ja pots activar el bot

