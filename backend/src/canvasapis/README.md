# canvasapis

This package used to contain the proxy that ran CanvasCBL.

Now it just contains the CanvasOAuth2 handlers as we can't change their URIs.

## Current Routes

- `GET` `request` - forwards the user to the canvas OAuth2 URI
- `GET` `response` - Canvas OAuth2 Response URI