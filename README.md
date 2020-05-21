# canvascbl

Closed-source fork of [iamtheyammer/canvas-grade-calculator](https://github.com/iamtheyammer/canvas-grade-calculator) which includes CanvasCBL+, a paid service with more features. It also saves all proxy responses to the database.

[![Build Status](https://travis-ci.com/iamtheyammer/canvascbl.svg?token=qQmd7eMUZpxTcqHBWHBw&branch=master)](https://travis-ci.com/iamtheyammer/canvascbl)

See the [Backend README](backend/README.md) and [Frontend README](frontend/README.md) for more information.

All images and assets have been moved to [Google Drive](https://drive.google.com/drive/u/0/folders/168p3X_pzMrbTXgtWjJJ5VQbeJ2_t1zOE).

## Current Stack

- Backend on Heroku (canvas-grade-calculator.herokuapp.com)
- Web hosting with S3 & CloudFront (canvascbl.com)
- CDN (cost savings) on the web hosting with Cloudflare
- RDS for PostgreSQL on AWS for database (connects to backend)
