# Goboe

Based on the very good Python package [Oboe](https://github.com/kmaasrud/oboe) which will no longer be maintained.

This package is designed to be used as a command line tool for converting your [Obsidian.md](http://obsidian.md/) notes into HTML for publishing publicly or privately on the web.

This package extends Oboe's functionality and supports all standard markdown, as well as the following Obsidian features:
- Obsidian Wiki Links
- Backlinks

## Requirements
- Go(lang)

## Install
`go install https://github.com/samxsmith/goboe/cmd/goboe/`

## How to Use
`goboe ~/Documents/my_vault -o ./public`

The first command line argument should be the root to your vault.

Then pass an `-o` flag declaring where you want the HTML output to go.

See the example folder for a way of creating a Gitlab wiki, public or private, for free using Goboe.

## Templates
You can pass an HTML template, into which each note will be inserted. This gives you full stylistic and formatting control. You pass the path to the template with the `-t` flag:
```
goboe ~/Documents/my_vault -o ./public -t ./my_template
```

The template must contain a placeholder for content: `{content}`. Each note's content will go here.

The notes are not styled at all, so that you can have full control over their appearance.

## Coming Soon
- Indexed tags
