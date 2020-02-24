 # Canvas Proxy API

Written in Golang for maximum speed (and so many other absolutely wonderful things!).

## CanvasCBL API

See the official API Docs - api-docs.canvascbl.com or go.canvascbl.com/docs.

## Frontend Endpoints

No query string parameters from Canvas or otherwise, except otherwise noted, are supported.

Users need a valid session to use these. OAuth2 Access Tokens will not be permitted for use on these endpoints.

### CanvasCBL+

Contains a set of APIs users need CanvasCBL+ to use. In the future, it will include things like average grades, average outcome scores and more.

Users need a valid subscription to use all endpoints except for those marked \[NS\].

- `GET` `/api/plus/session` \[NS\] - Returns info about your current session
- `DELETE` `/api/plus/session` \[NS\] - Deletes the session_string cookie, effectively logging a user out.
- `GET` `/api/plus/courses/:courseID/avg` - Returns the average grade for the course
- `GET` `/api/plus/outcomes/:outcomeID/avg` - Returns the average score for the outcome
- `GET` `/api/plus/grades/previous` - Returns the user's previous grades (5 min+ ago)

### Checkout

Users need a valid session to use all of these endpoints.

- `GET` `/api/checkout/products` - Lists all products as JSON.
- `GET` `/api/checkout/session` - Gets a Stripe checkout session ID by product ID.
  - Requires the `productId` param, received from `/api/checkout/products`.
- `POST` `/api/checkout/redeem` - Lets you redeem gift cards
  - Body should look something like this: `{ "codes": ["ABCD-EFGH-IJKL", ...] }`
- `GET` `/api/checkout/subscriptions` - Returns an array of all currently valid subscriptions the user holds.
- `DELETE` `/api/checkout/subscriptions` - Cancels your current subscription effective immediately.
- `POST` `/api/checkout/webhook` - To be used as the Stripe webhook URL. As it verifies Stripe-Signature, there is no reason to send requests here.

### Admin

These endpoints require a user's status to be `2` (admin). They're a growing set of endpoints for doing things only admins would need to do.

- `POST` `/api/admin/gift_cards` - Creates gift cards
  - Requires URL param `quantity`-- the number of gift cards to generate -- ex: `quantity=200`
  - Requires URL param `valid_for`-- the number of seconds the card will be valid for; ex: `valid_for=2629800`

## OAuth2

### CanvasCBL

CanvasCBL supports OAuth2 for other apps.

For endpoints marked with \[P\], see the public CanvasCBL API docs at [https://go.canvascbl.com/docs](https://go.canvascbl.com/docs).

- `POST` `/api/oauth2/auth` \[P\] - Beginning auth for other apps
- `GET` `/api/oauth2/consent` - Private endpoint for getting info about a consent_token. Returns stuff like the app name and scopes.
- `PUT` `/api/oauth2/consent` - Private endpoint for confirming auth. Returns json with a redirect_to field.
- `POST` `/api/oauth2/token` \[P\] - Public endpoint for other apps getting a token
- `DELETE` `/api/oauth2/token` \[P\] - Public (but also takes sessions) endpoint for deleting an auth
- `GET` `/api/oauth2/tokens` - Private endpoint for listing tokens for the sessioned user

### ~~Google~~ (deprecated)

#### OAuth2 Endpoints

- `GET` `/api/google/oauth2/request` - Redirects to the Google OAuth2 grant page, injecting the Client, scopes and more. A user should be redirected here.
- `GET` `/api/google/oauth2/response` - OAuth2Callback for google. You shouldn't use this endpoint. This endpoint does a lot. See Successful Query String Params to learn more.

#### Successful Query String Params

The Google OAuth2 Response handler does a lot. It:

- Gets an OAuth2 token for the user
- Uses said token to pull the user's profile
- Upserts the profile into `google_users`
- Generates a session for the user, saves as cookie and is available via `X-Session-String`
- Determines whether the user has a stored Canvas token
- Redirects to the OAuth2 Callback environment variable


- `type` - `google`
- `has_token` - `true`|`false` - Whether the user has a stored token. If false, you'll probably want to show a dialog prompting the user to add a token.
- `session_string` - The session string to save as a cookie, which MUST expire in less than 2 weeks.

Note that, unlike the canvas flow, the token is never sent to the browser. It's held by the server and used when a request is made with a valid session.

#### Error Redirect URI Query String Params

Special query params are present when an error occurs during the OAuth2 response flow.

- `error` - `proxy_google_error` - if the `error` param is present, there was an error in the flow. this will never show up during a successful flow
- `error_source` - `proxy` | `google` - where the error originated
  - if `error_source` is `google`
    - [OPTIONAL] `body` - the URL encoded JSON body from the Google request that failed; pretty much for debugging
  - if `error_source` is `proxy`
    - `error_text` - a message you can show to the user about what happened, ex: `domain not allowed`

### Canvas

#### OAuth2 Endpoints

See the two below sections about Redirect URI Query String Params for handling the response data.

- `GET` `/api/canvas/oauth2/request` - Redirects the user to the Canvas OAuth2 grant page, injecting your client ID and other applicable query string params. A user would be redirected to this URL.
  - No params.
- `GET` `/api/canvas/oauth2/response` - Should be the OAuth2 response URI. Handles the error/success Canvas OAuth2 response. A user would be redirected to this URL by Canvas. You should **not** use this endpoint.
  - Params will include those from Canvas, so either `code` or `error`.

#### Normal (successful) Redirect URI Query String Params

In the event of a successful grant from Canvas, two query string parameters will be in the URL to the Redirect URI.

- `type` - `canvas`
- `name` - The user's first and last name as a string
- `subdomain` - the subdomain from the OAuth2 token grant

#### Error Redirect URI Query String Params

In the event of an OAuth2 error, the user will be redirected to the Redirect URI, however some special query params will be present.

Note that the `X-Canvas-Url` header will be present on errors, but it will just contain `omitted` as to not leak the client secret.

Those query params:

- `error` - either the error string from Canvas, or `proxy_canvas_error`
- `error_source` - `proxy`|`canvas`
- `canvas_status_code` - the status code from the proxied request to Canvas
- `body` - if Canvas returns a JSON body (currently not a possibility, but supported for future expansion on Canvas's side), this will contain the body, raw. Otherwise, this will contain `html_omitted`


## Environment Variables

Some environment variables are required to start the proxy server. This list is very, very out of date.

```sh
# Contains your Canvas OAuth2 Client ID
export CANVAS_OAUTH2_CLIENT_ID="canvasoauth2clientid"

# Contains your Canvas OAuth2 Client Secret
export CANVAS_OAUTH2_CLIENT_SECRET="canvasoauth2clientsecret"

# The subdomain that will be handling OAuth2
export CANVAS_OAUTH2_SUBDOMAIN="canvas"

# The Redirect URI. You should add your host here instead of localhost:8000, and replace http with https
export CANVAS_OAUTH2_REDIRECT_URI="http://localhost:8000/api/canvas/oauth2/response"

# Success URI-- where users will be redirected to-- see the normal query string section in the OAuth2 section.
export CANVAS_OAUTH2_SUCCESS_URI="http://localhost:3000/#/oauth2response"

# Allowed CORS origins-- should NEVER be * on a production server. Sites that are allowed to make proxied requests. Can be * to allow requests from everywhere, or be like "https://google.com" to allow requests from google.com. More: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
export CANVAS_PROXY_ALLOWED_CORS_ORIGINS="*"

# Allowed Canvas subdomains-- should NEVER be * on a production server. Also probably should match your OAuth2 Subdomain. Comma separated. Ex: "canvas,myschool" to allow canvas.instructure.com and myschool.instructure.com. Also can be * to allow all, but that will throw a warning.
export CANVAS_PROXY_ALLOWED_SUBDOMAINS="*"

# Your default Canvas subdomain for non-OAuth2 requests. Should probably match your OAuth2 subdomain and MUST be in your allowed subdomains list.
export CANVAS_PROXY_DEFAULT_SUBDOMAIN="canvas"

# Whether the proxy should serve static from the build folder. Defaults to false.
export CANVAS_PROXY_SERVE_STATIC="false"

# Database connection string
export DATABASE_DSN="postgres://postgres@localhost:5432/canvascbl"

# Stripe API Key
export STRIPE_API_KEY="sk_test_jsdfsdgeYe9myfunapikeypGnwxsfklsjdoib8jL"
```
