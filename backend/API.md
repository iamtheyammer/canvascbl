# CanvasCBL Open API

The CanvasCBL Open API allows third-party applications to access grades and other user information from CanvasCBL.

It is authorized using OAuth2-- see [that section](#OAuth2) for details.

The base path is `https://api.canvascbl.com/api/v1/`.

Unless explicitly specified, the content type for all responses is `application/json`. Same for requests with POST bodies.

## OAuth2

In order to use the CanvasCBL API, you'll need an OAuth2 Client ID and Client Secret. Please contact us for these.

### Scopes

Your Client ID and Client Secret will be limited to the following scopes, with the most common being `grades` only:

| Name | Description |
| ---: | ----------: |
| `profile` | The user's entire Canvas user profile. |
| `observees` | A list of the user's observees. |
| `courses` | All of the user's Canvas courses. |
| `alignments` | Outcome alignments for a course. | 
| `rollups` | Outcome rollups for a user and a course. |
| `assignments` | A course's assignments. |
| `outcomes` | A course's outcomes. |
| `grades` | **A user's grades with the names of the courses they're for.** |
| `previous_grades` | A user's previous grades. |
| `average_course_grade` | A course's average grade. |
| `average_outcome_score` | The average score for an outcome. |

### GET /oauth2/authorize

This endpoint is the first step in the authentication flow. You should send a user to this link in their browser.

#### Query Params

| Param | Example Value | Description |
| ----: | ------------: | ----------: |
| `response` | `code` | **Required.** Must be `code`. |
| `client_id` | `d262d1d3-d969-4d48-ac1e-cfceec88b5c9` | **Required.** Your Client ID |
| `scope` | `profile,observees,grades` | **Required.** Comma-separated list of scopes you would like access to. |
| `redirect_uri` | `https://dcraft.com/oauth2/response` | **Required.** The URI where the user will be redirected after the authorization. Must match the Redirect URI on your OAuth2 Credentials. |
| `purpose` | `d.Craft` | Helps the user identify what this token is for. |

#### Response

This endpoint will forward the user to an opaque URL where they will decide whether to authorize your app.

If an error occurred, your redirect URI will be appended with the `error` query param, with one of the following values:

- `invalid_scope` - You requested a non-existent scope or your credentials don't have access to the following scope.
- `unsupported_response_type` - The `response_type` param was not recognized by the server.
- `unauthorized_client` - There's something wrong with your OAuth2 Client ID.
- `access_denied` - The user rejected your access request.
- `server_error` - There was a server error when processing your request.

If the user accepted your access request, your redirect URI will be appended with the `code` query param. Use this code
param in the [POST /oauth2/token route](#post-oauth2token).

### POST /oauth2/token

This request, when successful, returns an OAuth2 Bearer token and an OAuth2 Refresh Token for the user.

The user MUST NOT be able to see this request as it will contain your client secret.

### Body

Unlike all other CanvasCBL API requests, this body must be in the `application/x-www-form-urlencoded` format.

| Field | Example Value | Status | Description |
| ----: | ------------: | -----: | ----------: |
| `grant_type` | `code` OR `refresh_token` | Required. | Whether you want to get the token using an authorization code or a refresh token. |
| `client_id` | `d262d1d3-d969-4d48-ac1e-cfceec88b5c9` | Required. | Your OAuth2 Client ID. |
| `client_secret` | `08d1d6c9-b982-4369-8ff1-ac20c0ff0824` | Required. | Your OAuth2 Client Secret. |
| `redirect_uri` | `https://dcraft.com/oauth2/response` | Required for `grant_type` = `code`. | The same redirect URI you supplied in your authorize request. |
| `code` | `016d46c5-8f39-4763-baa0-72985b2f2977` | Required for `grant_type` = `code`. | The `code` parameter sent to your redirect URI. |
| `refresh_token` | `005f9d29-cffb-46b3-a6eb-5960bd325c80` | Required for `grant_type` = `refresh_token`. | Your refresh token. |

### Response

```json5
{
  "access_token": "fa7e1ace-6af3-4f20-9c42-778854e372f3", // woot: your token!
  "token_type": "Bearer", // always Bearer
  "expires_in": 3600, // seconds until token expiry
  "user": { "id":  123, "canvas_id": 456, "name": "David Bowie" }, // some small info about the user
  "refresh_token": "005f9d29-cffb-46b3-a6eb-5960bd325c80" // only present on grant_type = code requests.
}
```

