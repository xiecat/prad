package main

import "log"

const banner = `
╱╱╱╱╱╱╱╱╱╱╭╮
╱╱╱╱╱╱╱╱╱╱┃┃
╭━━┳━┳━━┳━╯┃
┃╭╮┃╭┫╭╮┃╭╮┃
┃╰╯┃┃┃╭╮┃╰╯┃
┃╭━┻╯╰╯╰┻━━╯
┃┃
╰╯ v0.0.1
`

func showBanner() {
	log.Printf("%s\n", banner)
}
