# canvascbl

Closed-source fork of [iamtheyammer/canvas-grade-calculator](https://github.com/iamtheyammer/canvas-grade-calculator) which includes CanvasCBL+, a paid service with more features. It also saves all proxy responses to the database.

[![Build Status](https://travis-ci.com/iamtheyammer/canvascbl.svg?token=qQmd7eMUZpxTcqHBWHBw&branch=master)](https://travis-ci.com/iamtheyammer/canvascbl)

See the [Backend README](backend/README.md) and [Frontend README](frontend/README.md) for more information.

You can also check out [img/](img/) for some screenshots, which are wildly out of date.

## Current Stack

- Backend on Heroku (canvas-grade-calculator.herokuapp.com)
- Web hosting with GitHub Pages (canvascbl.com)
- SSLification on the web hosting with Cloudflare (canvascbl.com)
- RDS for PostgreSQL on AWS for database (connects to backend)
