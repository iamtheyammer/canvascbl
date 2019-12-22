# Canvas Proxy API

Written in Golang for maximum speed.

## Headers

### Request

There are two headers for each request to `/api/canvas/*` (except for OAuth2):

- `X-Canvas-Token`
    - Auth token for Canvas. Required for all requests except for OAuth2.
- `X-Canvas-Subdomain` (optional-- see [env vars](#environment-variables) for more info)
    - Subdomain for Canvas-- so `hello` would make calls to `hello.instructure.com`.

### Response

Two headers are returned from every request:

- `X-Canvas-Url` is returned from every request. It contains the URL that the proxy called. Helpful for debugging. Note that this header is always `omitted` for OAuth2 calls.
- `X-Canvas-Status-Code` contains the status code returned from the proxied call to Canvas.

If an error occurred, a `502 BAD GATEWAY` status will be returned from the proxy, except for OAuth2 calls.

## Endpoints

No query string parameters from Canvas or otherwise, except otherwise noted, are supported.

### Outcomes

- `GET` `/api/canvas/outcomes/:outcomeID` - Mirror of [this](https://canvas.instructure.com/doc/api/outcomes.html#method.outcomes_api.show) Canvas endpoint.

### Users

- `GET` `/api/canvas/users/profile/self` - Mirror of [this](https://canvas.instructure.com/doc/api/users.html#method.users.api_show) Canvas endpoint, with `:id` replaced with `self`.
    - Supports a custom `generateSession` param which will generate a CanvasCBL session as a cookie (`session_string`) and as a header (`X-Session-String`). Can be used for future CanvasCBL+ calls.
- `GET` `/api/canvas/users/profile/self/observees` - Mirror of [this](https://canvas.instructure.com/doc/api/user_observees.html#method.user_observees.index) Canvas endpoint with `:user_id` replaced with `self`.
    - Requires a custom `user_id` query param for backend functions. The final requested url will have the supplied user ID instead of `self`. If you're confused, try it and look at the returned `X-Canvas-URL` header.

### Courses

- `GET` `/api/canvas/courses` - Mirror of [this](https://canvas.instructure.com/doc/api/courses.html#method.courses.index) Canvas endpoint.
- `GET` `/api/canvas/courses/:courseID/assignments` - Mirror of [this](https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.index) Canvas endpoint.
  - Supports the `include[]` query param from Canvas
- `GET` `/api/canvas/courses/:courseID/outcome_groups` - Mirror of [this](https://canvas.instructure.com/doc/api/outcome_groups.html#method.outcome_groups_api.index) Canvas endpoint.
- `GET` `/api/canvas/courses/:courseID/outcome_groups/:outcomeGroupID/outcomes` - Mirror of [this](https://canvas.instructure.com/doc/api/outcome_groups.html#method.outcome_groups_api.outcomes) Canvas endpoint.
- `GET` `/api/canvas/courses/:courseID/outcome_results` - Mirror of [this](https://canvas.instructure.com/doc/api/outcome_results.html#method.outcome_results.index) Canvas endpoint.
  - Requires the [`userId` param](#userid-param)
  - Supports the `include[]` query param from Canvas
- `GET` `/api/canvas/courses/:courseID/outcome_rollups` - Mirror of [this](https://canvas.instructure.com/doc/api/outcome_results.html#method.outcome_results.rollups) Canvas endpoint.
  - Requires the [`userId` param](#userid-param)
  - Supports the `include[]` query param from Canvas
- `GET` `/api/canvas/courses/:courseID/outcome_alignments` - Mirror of [this](https://canvas.instructure.com/doc/api/outcomes.html#method.outcomes_api.outcome_alignments) Canvas endpoint.
  - Requires the [`userId` param](#userid-param)

#### `userId` param

Sets the `user_ids[]` param to the value of the `userId` param. It should be equal to the ID of the user the token is for. This is required because students only have permission to list their own outcome results, and this endpoint defaults to listing results for all students. Ex: `userId=12345`

### ~~Tokens~~ (deprecated)

The API supports holding on to tokens. All tokens endpoints require a session.

- `GET` `/api/canvas/tokens` - List your non-expired Tokens
- `POST` `/api/canvas/tokens` - Add a token
  - Body (JSON):
    - `token` - the token to use; the API will check it before adding it-- if these checks fail it will return 400 bad request.
    - [OPTIONAL] `expiresAt` - unix epoch in seconds; when the token will expire. if the token will never expire set it to null or 0 or just leave it out.

### CanvasCBL+

Contains a set of APIs users need CanvasCBL+ to use. In the future, it will include things like average grades, average outcome scores and more.

Users need a valid subscription to use all endpoints except for those marked \[NS\].

- `GET` `/api/plus/session` \[NS\] - Returns info about your current session
- `GET` `/api/plus/courses/:courseID/avg` - Returns the average grade for the course
- `GET` `/api/plus/outcomes/:outcomeID/avg` - Returns the average score for the outcome
- `GET` `/api/plus/grades/previous` - Returns the user's previous grades (5 min+ ago)

### Checkout

- `GET` `/api/checkout/products` \[NS\] - Lists all products as JSON.
- `GET` `/api/checkout/session` - Gets a Stripe checkout session ID by product ID.
  - Requires the `productId` param, received from `/api/checkout/products`.
- `GET` `/api/checkout/subscriptions` - Returns an array of all currently valid subscriptions the user holds.
- `DELETE` `/api/checkout/subscriptions` - Cancels your current subscription effective immediately.
- `POST` `/api/checkout/webhook` - To be used as the Stripe webhook URL. As it verifies Stripe-Signature, there is no reason to send requests here.

#### Sessions

These endpoints require a session, generated from the `/api/canvas/users/profile/self` endpoint with `?generateSession=true`.
Provide it in the X-Session-String header or as a cookie (`session_string`).

## OAuth2

The backend supports proxying OAuth2 requests and responses to the frontend.

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
- `GET` `/api/canvas/oauth2/refresh_token` - Retrieves a new token based on a refresh token.
  - Params:
    - `refresh_token`: your refresh token


#### Normal (successful) Redirect URI Query String Params

In the event of a successful grant from Canvas, two query string parameters will be in the URL to the Redirect URI.

- `type` - `canvas`
- `canvas_response` - contains the JSON payload from the Canvas token grant (examples [here](https://canvas.instructure.com/doc/api/file.oauth_endpoints.html))
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

Some environment variables are required to start the proxy server.

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
