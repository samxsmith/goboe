# Obsidian Wiki

Gitlab allows you to create static sites for free with its Pages offering.

This can be public or private, which is great for your notes which you might want to share, or keep to yourself.

I publish one vault publicly as a knowledge wiki that I share.

My personal notes are private, but I can use Gitlab pages to view my Obsidian notes from my phone, and when I'm on the move.

This pairs well with obsidian-git plugin, which will push your notes regularly once setup.

## What does this give you?
A free wiki version of your vault. Also, an index page with a list of all of your notes.

## Requirements
- Your obsidian vault is a git repo
- It has a Gitlab remote to which you push. Make the repo public or private to toggle who can view it.

## Setup
Add the `.gitlab-ci.yml` to the root of your Obsidian vault.

Add the template file `_publish_template` to the root as well, to get a simple menu at the top of each note.

Push to gitlab, and your site should build.

Once the job completes you should find your site at: `YOUR_GITLAB_USER.gitlab.io/NAME_OF_GITLAB_REPO`