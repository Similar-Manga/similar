/*
 * MangaDex API
 *
 * MangaDex is an ad-free manga reader offering high-quality images!  Here is some generic stuff about the API  # Authentication  You can login with `/auth/login` endpoints. It will return a JWT that remains for 15min and that have a 4h refresh token.  # Rating limits  We are using rating limits in order to avoid too much calls on our endpoints, here is how is it configured:  | Endpoint                | Calls            | Time frame                | |-------------------------|------------------|---------------------------| | `/account/create` | 1 | 60 minutes | | `/account/activate/{code}` | 30 | 60 minutes | | `/account/activate/resend`, `/account/recover`, `/account/recover/{code}` | 5 | 60 minutes | | `/auth/login` | 30 | 60 minutes | | `/auth/refresh` | 30 | 60 minutes | | `/chapter/{id}/read` | 300              | 10 minutes                | | `/upload/begin`, `/upload/{id}`, `/upload/{id}/commit` | 30               | 1 minute                  | | `PUT /chapter/{id}` | 10               | 1 minute                  | | `DELETE /chapter/{id}` | 10               | 1 minute                  | | `POST /manga` | 10               | 60 minutes                | | `PUT /manga/{id}` | 10               | 1 minute                  | | `DELETE /manga/{id}` | 10               | 10 minutes                | | `POST /group` | 10               | 60 minutes                | | `PUT /group/{id}` | 10               | 1 minute                  | | `DELETE /group/{id}` | 10               | 10 minutes                | | `POST /author` | 10               | 60 minutes                | | `PUT /author` | 10               | 1 minutes                 | | `DELETE /author/{id}` | 10               | 10 minutes                | | `POST /captcha/solve` | 10 | 10 minutes |  You can get details about how your rate limit is going by reading following headers:  | Header                  | Description                                              | |-------------------------|-----------------------------------------------------------| | X-RateLimit-Limit       | Number of max requests allowed in the current time period | | X-RateLimit-Remaining   | Number of remaining requests in the current time period   | | X-RateLimit-Retry-After | Timestamp of end of current time period as UNIX timestamp |  # Captchas  Some endpoints may require captchas to proceed, in order to slow down automated malicious traffic. Regular users might see a couple of captchas, based on the frequency of write requests or on certain endpoints like user signup.  Once an endpoint decides that a captcha needs to be solved, a 403 Forbidden response will be returned, where the error title is `captcha_required_exception`. The sitekey needed for recaptcha to function is provided in both the `X-Captcha-Sitekey` header field, as well as in the error context, specified as the `siteKey` parameter.  The captcha result of the client can either be passed into the repeated original request with the `X-Captcha-Result` header or alternatively the `POST /captcha/solve` endpoint can be called to solve this captcha. The time a solved captcha is remembered varies across different endpoints and can also be influenced by individual client behavior.  Authentication is not required for the `POST /captcha/solve` endpoint, captchas are tracked separately for client ip and user id. If you are logged in, you want to send the session token so you validate the captcha for your client ip and user id at the same time, but it is not required.  # Chapter pages processing  ## Pages processing  When you fetch a chapter response you'll have 4 fields that you need for pages processing:  | Field                       | Type     | Description        | |-----------------------------|----------|--------------------| | `data.id`                   | `string` | API identifier     | | `data.attributes.hash`      | `string` | MD@H identifier    | | `data.attributes.data`      | `array`  | High quality pages | | `data.attributes.dataSaver` | `array`  | Low quality pages  |  From this point you miss one thing: a MD@H backend server to get images from, to get that make a request to `GET /at-home/server/{data.id}`, it will return the server url to use.  Then to build the pages, you have to build your url as following: `https://{md@h server node}/data/{data.attributes.hash}/{data.attributes.data}`  Or if you want to use the low quality files: `https://{md@h server node}/data-saver/{data.attributes.hash}/{data.attributes.dataSaver}`  Here is an example of what the url could looks like:  | Field                       | Value | |-----------------------------|-------| | `data.id`                   | `000002b1-e8de-4281-9781-8e81e869f579` | | `data.attributes.hash`      | `caad0c22434276b9e3e56a78fe2e7993` | | `data.attributes.data`      | `[\"x1-a87ae6522fa5c244fd76985c7d953ccf3975bec66ce9b8e813549e642b38a47a.png\", ...]` | | `data.attributes.dataSaver` | `[\"x1-a1d3047dfccd77b3117a86ccf19a9c5403e09baec6a78893ed1d3825d2c71256.jpg\", ...]` |  As a \"fake\" MD@H node we'll use `https://s2.mangadex.org/` server.  So for high quality we'll have an URL like that: https://s2.mangadex.org/data/caad0c22434276b9e3e56a78fe2e7993/x1-a87ae6522fa5c244fd76985c7d953ccf3975bec66ce9b8e813549e642b38a47a.png  And for low quality: https://s2.mangadex.org/data-saver/caad0c22434276b9e3e56a78fe2e7993/x1-a1d3047dfccd77b3117a86ccf19a9c5403e09baec6a78893ed1d3825d2c71256.jpg  ## Report  In order to make everything works well, we keep statistics over MD@H nodes and how they perform. In order to keep theses statistics you have to post data for each page you fetch from a MD@H node.  Here is an example: ```curl POST https://api.mangadex.network/report {   \"url\": \"https://s2.mangadex.org/data/caad0c22434276b9e3e56a78fe2e7993/x1-a87ae6522fa5c244fd76985c7d953ccf3975bec66ce9b8e813549e642b38a47a.png\",   \"success\": true,   \"bytes\": 800000, // size of the loaded image   \"duration\": 213, // miliseconds to load the image   \"cached\": false, // X-Cache header of the MDAH node == 'HIT' ? } ```  # Static data  ## Manga publication demographic  | Value            | Description               | |------------------|---------------------------| | shonen           | Manga is a Shonen         | | shoujo           | Manga is a Shoujo         | | josei            | Manga is a Josei          | | seinen           | Manga is a Seinen         |  ## Manga status  | Value            | Description               | |------------------|---------------------------| | ongoing          | Manga is still going on   | | completed        | Manga is completed        | | hiatus           | Manga is paused           | | abandoned        | Manga has been abandoned  |  ## Manga reading status  | Value            | |------------------| | reading          | | on_hold          | | plan\\_to\\_read   | | dropped          | | re\\_reading      | | completed        |  ## Manga content rating  | Value            | Description               | |------------------|---------------------------| | safe             | Safe content              | | suggestive       | Suggestive content        | | erotica          | Erotica content           | | pornographic     | Pornographic content      |  ## CustomList visibility  | Value            | Description               | |------------------|---------------------------| | public           | CustomList is public      | | private          | CustomList is private     |  ## Relationship types  | Value            | Description                    | |------------------|--------------------------------| | manga            | Manga resource                 | | chapter          | Chapter resource               | | author           | Author resource                | | artist           | Author resource (drawers only) | | scanlation_group | ScanlationGroup resource       | | tag              | Tag resource                   | | user             | User resource                  | | custom_list      | CustomList resource            |
 *
 * API version: 5.0.0
 * Contact: mangadexstaff@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Manga struct {
	Id string `json:"id,omitempty"`
	Type_ string `json:"type,omitempty"`
	Attributes *MangaAttributes `json:"attributes,omitempty"`
}
