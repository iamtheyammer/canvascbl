# canvascbl

Closed-source fork of [iamtheyammer/canvas-grade-calculator](https://github.com/iamtheyammer/canvas-grade-calculator) which includes CanvasCBL+, a paid service with more features. It also saves all proxy responses to the database.

[![Build Status](https://travis-ci.com/iamtheyammer/canvascbl.svg?branch=master)](https://travis-ci.com/iamtheyammer/canvascbl)

See the [Backend README](backend/README.md) and [Frontend README](frontend/README.md) for more information.

You can also check out [img/](img/) for some screenshots, which are wildly out of date.

# Running on Heroku

Ready to run! Clone the repo and follow Heroku's instructions.

## Heroku Build Process

1. Heroku builds the go executable to `bin/src` (where `src` is the actual executable)
2. Heroku runs `web`, declared in the Procfile.
