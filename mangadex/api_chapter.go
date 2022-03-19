
/*
 * MangaDex API
 *
 * MangaDex is an ad-free manga reader offering high-quality images!  This document details our API as it is right now. It is in no way a promise to never change it, although we will endeavour to publicly notify any major change.  # Acceptable use policy  Usage of our services implies acceptance of the following: - You **MUST** credit us - You **MUST** credit scanlation groups if you offer the ability to read chapters - You **CANNOT** run ads or paid services on your website and/or apps  These may change at any time for any and no reason and it is up to you check for updates from time to time.  # Security issues  If you believe you found a security issue in our API, please check our [security.txt](/security.txt) to get in touch privately.  # Authentication  You can login with the `/auth/login` endpoint. On success, it will return a JWT that remains valid for 15 minutes along with a session token that allows refreshing without re-authenticating for 1 month.  # Rate limits  The API enforces rate-limits to protect our servers against malicious and/or mistaken use. The API keeps track of the requests on an IP-by-IP basis. Hence, if you're on a VPN, proxy or a shared network in general, the requests of other users on this network might affect you.  At first, a **global limit of 5 requests per second per IP address** is in effect.  > This limit is enforced across multiple load-balancers, and thus is not an exact value but rather a lower-bound that we guarantee. The exact value will be somewhere in the range `[5, 5*n]` (with `n` being the number of load-balancers currently active). The exact value within this range will depend on the current traffic patterns we are experiencing.  On top of this, **some endpoints are further restricted** as follows:  | Endpoint                           | Requests per time period    | Time period in minutes | |------------------------------------|-----------------------------|------------------------| | **Account management**                                                                    | | `GET    /account/available         | 60                          | 60                     | | `POST   /account/create`           | 5                           | 60                     | | `POST   /account/activate/{code}`  | 30                          | 60                     | | `POST   /account/activate/resend`  | 5                           | 60                     | | `POST   /account/recover`          | 5                           | 60                     | | `POST   /account/recover/{code}`   | 5                           | 60                     | | `POST   /user/email`               | 1                           | 60                     | | `DELETE /user/{id}`                | 5                           | 60                     | | **Authentication**                                                                        | | `POST   /auth/login`               | 30                          | 60                     | | `POST   /auth/refresh`             | 60                          | 60                     | | **Manga**                                                                                 | | `POST   /manga`                    | 10                          | 60                     | | `PUT    /manga/{id}`               | 10                          | 60                     | | `DELETE /manga/{id}`               | 10                          | 10                     | | `POST   /manga/draft/{id}/commit`  | 10                          | 60                     | | **Authors**                                                                               | | `POST   /author`                   | 10                          | 60                     | | `PUT    /author`                   | 10                          | 1                      | | `DELETE /author/{id}`              | 10                          | 10                     | | **Covers**                                                                                | | `POST   /cover`                    | 100                         | 10                     | | `PUT    /cover/{id}`               | 100                         | 10                     | | `DELETE /cover/{id}`               | 10                          | 10                     | | **Chapters**                                                                              | | `POST   /chapter/{id}/read`        | 300                         | 10                     | | `PUT    /chapter/{id}`             | 10                          | 1                      | | `DELETE /chapter/{id}`             | 10                          | 1                      | | **Scanlation groups**                                                                     | | `POST   /group`                    | 10                          | 60                     | | `PUT    /group/{id}`               | 10                          | 1                      | | `DELETE /group/{id}`               | 10                          | 10                     | | **CDN** (i.e.: MangaDex@Home)                                                             | | `GET    /at-home/server/{id}`      | 40                          | 1                      | | **Uploads**                                                                               | | -> Sessions                                                                               | | `GET    /upload`                   | 30                          | 1                      | | `POST   /upload/begin`             | 20  (shared)                | 1                      | | `POST   /upload/begin/{id}`        | 20  (shared)                | 1                      | | `POST   /upload/{id}/commit`       | 10                          | 1                      | | `DELETE /upload/{id}`              | 30                          | 1                      | | -> Files                                                                                  | | `POST   /upload/{id}`              | 250 (shared)                | 1                      | | `DELETE /upload/{id}/{id}`         | 250 (shared)                | 1                      | | `DELETE /upload/{id}/batch`        | 250 (shared)                | 1                      | | **reCaptcha**                                                                             | | `POST   /captcha/solve`            | 10                          | 10                     | | **Reports**                                                                               | | `POST   /report`                   | 10                          | 1                      | | `GET    /report`                   | 10                          | 1                      |  Calling these endpoints will further provide details via the following headers about your remaining quotas:  | Header                    | Description                                                                 | |---------------------------|-----------------------------------------------------------------------------| | `X-RateLimit-Limit`       | Maximal number of requests this endpoint allows per its time period         | | `X-RateLimit-Remaining`   | Remaining number of requests within your quota for the current time period  | | `X-RateLimit-Retry-After` | Timestamp of the end of the current time period, as UNIX timestamp          |  # Result Limit  Most of our listing endpoints will return a maximum number of total results that is currently capped at 10.000 items. Beyond that you will not receive any more items no matter how far you paginate and the results will become empty instead. This is for performance reasons and a limitation we will not lift.  Note that the limit is applied to a search query and list endpoints with or without any filters are search queries. If you need to retrieve more items, use filters to narrow down your search.  # Reference Expansion  Reference Expansion is a feature of certain endpoints where relationships of a resource are expanded with their attributes, which reduces the amount of requests that need to be sent to the API to retrieve a complete set of data.  It works by appending a list of includes to the query with the type names from the relationships. If an endpoint supports this feature is indicated by the presence of the optional `includes` parameter.  ## Example  To fetch a specific manga with `author`, `artist` and `cover_art` expanded, you can send the following request: `GET /manga/f9c33607-9180-4ba6-b85c-e4b5faee7192?includes[]=author&includes[]=artist&includes[]=cover_art`. You will now find the objects attributes inside the returned relationships array.  ## Note  Your current user needs `*.view` permission on each type of reference you want to expand. Guests have most of these permissions and for logged-in user accounts you can check `GET /auth/check` to list all permissions you have been granted. For example, to be able to expand `cover_art`, you need to have been granted the `cover.view` permission, for `author` and `artist` types you need the `author.view` permission and so on. If a reference can't be expanded, the request will still succeed and no error indication will be visible.  # Captchas  Some endpoints may require captchas to proceed, in order to slow down automated malicious traffic. Legitimate users might also be affected, based on the frequency of write requests or due certain endpoints being particularly sensitive to malicious use, such as user signup.  Once an endpoint decides that a captcha needs to be solved, a 403 Forbidden response will be returned, with the error code `captcha_required_exception`. The sitekey needed for recaptcha to function is provided in both the `X-Captcha-Sitekey` header field, as well as in the error context, specified as `siteKey` parameter.  The captcha result of the client can either be passed into the repeated original request with the `X-Captcha-Result` header or alternatively to the `POST /captcha/solve` endpoint. The time a solved captcha is remembered varies across different endpoints and can also be influenced by individual client behavior.  Authentication is not required for the `POST /captcha/solve` endpoint, captchas are tracked both by client ip and logged in user id. If you are logged in, you want to send the session token along, so you validate the captcha for your client ip and user id at the same time, but it is not required.  # Reading a chapter using the API  ## Retrieving pages from the MangaDex@Home network  A valid [MangaDex@Home network](https://mangadex.network) page URL is in the following format: `{server-specific base url}/{temporary access token}/{quality mode}/{chapter hash}/{filename}`  There are currently 2 quality modes: - `data`: Original upload quality - `data-saver`: Compressed quality  You only need the chapter ID to start getting a chapter pages. Now by fetching `GET /at-home/server/{ chapterId }` you'll get all required fields to compute MangaDex@Home page URLs:  | Field                | Type     | Description                       | |----------------------|----------|-----------------------------------| | `.chapter.hash`      | `string` | MangaDex@Home Chapter Hash        | | `.chapter.data`      | `array`  | data quality mode filenames       | | `.chapter.dataSaver` | `array`  | data-saver quality mode filenames |  The full URL is the constructed as follows ``` { server .baseUrl }/{ quality mode }/{ chapter .chapter.hash }/{ chapter .chapter.{ quality mode }.[*] }  Examples  data quality: https://abcdefg.hijklmn.mangadex.network:12345/some-token/data/e199c7d73af7a58e8a4d0263f03db660/x1-b765e86d5ecbc932cf3f517a8604f6ac6d8a7f379b0277a117dc7c09c53d041e.png        base url: https://abcdefg.hijklmn.mangadex.network:12345/some-token   quality mode: data   chapter hash: e199c7d73af7a58e8a4d0263f03db660       filename: x1-b765e86d5ecbc932cf3f517a8604f6ac6d8a7f379b0277a117dc7c09c53d041e.png   data-saver quality: https://abcdefg.hijklmn.mangadex.network:12345/some-token/data-saver/e199c7d73af7a58e8a4d0263f03db660/x1-ab2b7c8f30c843aa3a53c29bc8c0e204fba4ab3e75985d761921eb6a52ff6159.jpg        base url: https://abcdefg.hijklmn.mangadex.network:12345/some-token   quality mode: data-saver   chapter hash: e199c7d73af7a58e8a4d0263f03db660       filename: x1-ab2b7c8f30c843aa3a53c29bc8c0e204fba4ab3e75985d761921eb6a52ff6159.jpg ```  If the server you have been assigned fails to serve images, you are allowed to call the `/at-home/server/{ chapter id }` endpoint again to get another server.  Whether successful or not, **please do report the result you encountered as detailed below**. This is so we can pull the faulty server out of the network.  # Manga Creation  Similar to how the Chapter Upload works, after a Mangas creation with the Manga Create endpoint it is in a \"draft\" state, needs to be submitted (committed) and get either approved or rejected by Staff.  The purpose of this is to force at least one CoverArt uploaded in the original language for the Manga Draft and to discourage troll entries. You can use the list-drafts endpoint to investigate the status of your submitted Manga Drafts. Rejected Drafts are occasionally cleaned up at an irregular interval. You can edit Drafts at any time, even after they have been submitted.  Because a Manga that is in the draft state is not available through the search, there are slightly different endpoints to list or retrieve Manga Drafts, but outside from that the schema is identical to a Manga that is published.  # Language Codes & Localization  To denote Chapter Translation language, translated fields such as Titles and Descriptions, the API expects a 2-letter language code in accordance with the [ISO 639-1 standard](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes). Additionally, some cases require the [5-letter extension](https://en.wikipedia.org/wiki/IETF_language_tag) if the alpha-2 code is not sufficient to determine the correct sub-type of a language, in the style of $language-$region, e.g. `zh-hk` or `pt-br`.  Because there is no standardized method of denoting romanized translations, we chose to append the `-ro` suffix. For example the romanized version of `五等分の花嫁` is `5Toubun no Hanayome` or `Gotoubun no Hanayome`. Both would have the `ja-ro` language code, alternative versions are inserted as alternative titles. This is a clear distinction from the localized `en` translation `The Quintessential Quintuplets`  Notable exceptions are in the table below, otherwise ask a staff member if unsure.  | alpha-5 | Description            | |---------|------------------------| | `zh`    | Simplified Chinese     | | `zh-hk` | Traditional Chinese    | | `pt-br` | Brazilian Portugese    | | `es`    | Castilian Spanish      | | `es-la` | Latin American Spanish | | `ja-ro` | Romanized Japanese     | | `ko-ro` | Romanized Korean       | | `zh-ro` | Romanized Chinese      |  # Report  In order to keep track of the health of the servers in the network and to improve the quality of service and reliability, we ask that you call the MangaDex@Home report endpoint after each image you retrieve, whether successfully or not.  It is a `POST` request against `https://api.mangadex.network/report` and expects the following payload with our example above:  | Field                       | Type       | Description                                                                         | |-----------------------------|------------|-------------------------------------------------------------------------------------| | `url`                       | `string`   | The full URL of the image                                                           | | `success`                   | `boolean`  | Whether the image was successfully retrieved                                        | | `cached `                   | `boolean`  | `true` iff the server returned an `X-Cache` header with a value starting with `HIT` | | `bytes`                     | `number`   | The size in bytes of the retrieved image                                            | | `duration`                  | `number`   | The time in miliseconds that the complete retrieval (not TTFB) of this image took   |  Examples herafter.  **Success:** ```json POST https://api.mangadex.network/report Content-Type: application/json  {   \"url\": \"https://abcdefg.hijklmn.mangadex.network:12345/some-token/data/e199c7d73af7a58e8a4d0263f03db660/x1-b765e86d5ecbc932cf3f517a8604f6ac6d8a7f379b0277a117dc7c09c53d041e.png\",   \"success\": true,   \"bytes\": 727040,   \"duration\": 235,   \"cached\": true } ```  **Failure:** ```json POST https://api.mangadex.network/report Content-Type: application/json  {  \"url\": \"https://abcdefg.hijklmn.mangadex.network:12345/some-token/data/e199c7d73af7a58e8a4d0263f03db660/x1-b765e86d5ecbc932cf3f517a8604f6ac6d8a7f379b0277a117dc7c09c53d041e.png\",  \"success\": false,  \"bytes\": 25,  \"duration\": 235,  \"cached\": false } ```  While not strictly necessary, this helps us monitor the network's healthiness, and we appreciate your cooperation towards this goal. If no one reports successes and failures, we have no way to know that a given server is slow/broken, which eventually results in broken image retrieval for everyone.  # Retrieving Covers from the API  ## Construct Cover URLs  ### Source (original/best quality)  `https://uploads.mangadex.org/covers/{ manga.id }/{ cover.filename }`<br/> The extension can be png, jpeg or gif.  Example: `https://uploads.mangadex.org/covers/8f3e1818-a015-491d-bd81-3addc4d7d56a/4113e972-d228-4172-a885-cb30baffff97.jpg`  ### <=512px wide thumbnail  `https://uploads.mangadex.org/covers/{ manga.id }/{ cover.filename }.512.jpg`<br/> The extension is always jpg.  Example: `https://uploads.mangadex.org/covers/8f3e1818-a015-491d-bd81-3addc4d7d56a/4113e972-d228-4172-a885-cb30baffff97.jpg.512.jpg`  ### <=256px wide thumbnail  `https://uploads.mangadex.org/covers/{ manga.id }/{ cover.filename }.256.jpg`<br/> The extension is always jpg.  Example: `https://uploads.mangadex.org/covers/8f3e1818-a015-491d-bd81-3addc4d7d56a/4113e972-d228-4172-a885-cb30baffff97.jpg.256.jpg`  ## ℹ️ Where to find Cover filename ?  Look at the [Get cover operation](#operation/get-cover) endpoint to get Cover information. Also, if you get a Manga resource, you'll have, if available a `covert_art` relationship which is the main cover id.  # Static data  ## Manga publication demographic  | Value            | Description               | |------------------|---------------------------| | shounen          | Manga is a Shounen        | | shoujo           | Manga is a Shoujo         | | josei            | Manga is a Josei          | | seinen           | Manga is a Seinen         |  ## Manga status  | Value            | Description               | |------------------|---------------------------| | ongoing          | Manga is still going on   | | completed        | Manga is completed        | | hiatus           | Manga is paused           | | cancelled        | Manga has been cancelled  |  ## Manga reading status  | Value            | |------------------| | reading          | | on_hold          | | plan\\_to\\_read   | | dropped          | | re\\_reading      | | completed        |  ## Manga content rating  | Value            | Description               | |------------------|---------------------------| | safe             | Safe content              | | suggestive       | Suggestive content        | | erotica          | Erotica content           | | pornographic     | Pornographic content      |  ## CustomList visibility  | Value            | Description               | |------------------|---------------------------| | public           | CustomList is public      | | private          | CustomList is private     |  ## Relationship types  | Value            | Description                    | |------------------|--------------------------------| | manga            | Manga resource                 | | chapter          | Chapter resource               | | cover_art        | A Cover Art for a manga `*`    | | author           | Author resource                | | artist           | Author resource (drawers only) | | scanlation_group | ScanlationGroup resource       | | tag              | Tag resource                   | | user             | User resource                  | | custom_list      | CustomList resource            |  `*` Note, that on manga resources you get only one cover_art resource relation marking the primary cover if there are more than one. By default this will be the latest volume's cover art. If you like to see all the covers for a given manga, use the cover search endpoint for your mangaId and select the one you wish to display.  ## Manga links data  In Manga attributes you have the `links` field that is a JSON object with some strange keys, here is how to decode this object:  | Key   | Related site  | URL                                                                                           | URL details                                                    | |-------|---------------|-----------------------------------------------------------------------------------------------|----------------------------------------------------------------| | al    | anilist       | https://anilist.co/manga/`{id}`                                                               | Stored as id                                                   | | ap    | animeplanet   | https://www.anime-planet.com/manga/`{slug}`                                                   | Stored as slug                                                 | | bw    | bookwalker.jp | https://bookwalker.jp/`{slug}`                                                                | Stored as \"series/{id}\"                                       | | mu    | mangaupdates  | https://www.mangaupdates.com/series.html?id=`{id}`                                            | Stored as id                                                  | | nu    | novelupdates  | https://www.novelupdates.com/series/`{slug}`                                                  | Stored as slug                                                | | kt    | kitsu.io      | https://kitsu.io/api/edge/manga/`{id}` or https://kitsu.io/api/edge/manga?filter[slug]={slug} | If integer, use id version of the URL, otherwise use slug one  | | amz   | amazon        | N/A                                                                                           | Stored as full URL                                             | | ebj   | ebookjapan    | N/A                                                                                           | Stored as full URL                                             | | mal   | myanimelist   | https://myanimelist.net/manga/{id}                                                            | Store as id                                                    | | cdj   | CDJapan       | N/A                                                                                           | Stored as full URL                                             | | raw   | N/A           | N/A                                                                                           | Stored as full URL, untranslated stuff URL (original language) | | engtl | N/A           | N/A                                                                                           | Stored as full URL, official english licenced URL              |  ## Manga related enum  This data is used in the \"related\" field of a Manga relationships  | Value             | Description |-------------------|------------------- | monochrome        | A monochrome variant of this manga | colored           | A colored variant of this manga | preserialization  | The original version of this manga before its official serialization | serialization     | The official serialization of this manga | prequel           | The previous entry in the same series | sequel            | The next entry in the same series | main_story        | The original narrative this manga is based on | side_story        | A side work contemporaneous with the narrative of this manga | adapted_from      | The original work this spin-off manga has been adapted from | spin_off          | An official derivative work based on this manga | based_on          | The original work this self-published derivative manga is based on | doujinshi         | A self-published derivative work based on this manga | same_franchise    | A manga based on the same intellectual property as this manga | shared_universe   | A manga taking place in the same fictional world as this manga | alternate_story   | An alternative take of the story in this manga | alternate_version | A different version of this manga with no other specific distinction  # Chapter Upload  ## In A Nutshell  To upload a chapter, you need to start an upload-session, upload files to this session and once done, commit the session with a ChapterDraft. Uploaded Chapters will generally be put into a queue for staff approval and image processing before it is available to the public.  ## Limits  - 1 Active Upload Session per user. Before opening a second session, either commit or abandon your current session - 10 files per one PUT request is max - 500 files per upload session is max - 20 MB max uploaded session filesize - 150 MB max total sum of all uploaded session filesizes - Allowed file extensions: jpg, jpeg, png, gif - Image must fit into max resolution of 10000x10000 px  ## Example  You need to be logged in for any upload operation. Which Manga you're allowed to upload to and which contributing Scanlation Groups you're free to credit depend on your individual user state.  To start an upload session, we pick the manga-id we want to upload to and the group-ids we have to credit. We use the official test manga `f9c33607-9180-4ba6-b85c-e4b5faee7192` and the group \"Unknown\" with id `145f9110-0a6c-4b71-8737-6acb1a3c5da4`. If no group can be credited, we would not send any group-id at all, otherwise we can credit up to 5 groups. Note that crediting all involved groups is mandatory, doing otherwise will lead to a rejection of your upload.  The first step is optional, but because only one upload session is allowed per user, we check if we have any open upload sessions by doing `GET /upload`. We expect a 404 response with error detail 'No upload session found'.  Next step is to begin our upload session. We send a `POST /upload/begin` with json data. (If you want to edit an existing chapter, append the chapter id after it `POST /upload/begin/db99d333-76e9-4e66-9c97-4831c43ac96c` with its version as the json payload)  Request: ```json {   'manga' => 'f9c33607-9180-4ba6-b85c-e4b5faee7192',   'groups': [     '145f9110-0a6c-4b71-8737-6acb1a3c5da4'   ] } ``` Response: ```json {   'result': 'ok',   'data': {     'id': '113b7724-dcc2-4fbc-968f-9d775fcb1cd6',     'type': 'upload_session',     'attributes': {       'isCommitted': false,       'isProcessed': false,       'isDeleted': false     },     'relationships': [       {         'id': '41ce3e1a-8325-45b5-af8e-06aaf648a0df',         'type': 'user'       },       {         'id': 'f9c33607-9180-4ba6-b85c-e4b5faee7192',         'type': 'manga'       },       {         'id': '145f9110-0a6c-4b71-8737-6acb1a3c5da4',         'type': 'scanlation_group'       }     ]   } } ```  the `data.id` is what you want to store because you will need it for the following steps. We will refer to it as the `uploadSessionId` from here on out.  Remember the `GET /upload` request from the beginning? Try it again and you will see that it will return the same uploadSessionId. You can only have one upload session per user until you commit or abandon it, which makes it easy for you to continue uploading at a later time.  Now that we have a `uploadSessionId`, we can upload images. Any .jpg, .jpeg, .png or .gif files are fine, archives like .zip, .cbz or .rar are not. You will have to extract those archives beforehand if you want to make this work.  For each file, send a POST request like `POST /upload/{uploadSessionId}` with the image data. FormData seems to work best with `Content-Type: multipart/form-data; boundary=boundary`, mileage might vary depending on your programming language. Join our discord and ask for advice if in doubt.  You can upload a number of files in a single request (currently max. 10). The response body will be successful (response.result == 'ok') but might also contain errors if one or more files failed to validate. It's up to you to handle bad uploads and retry or reupload as you see fit. Successful uploads will be returned in the data array as type `upload_session_file`  A successful response could look like this: ```json {   'result': 'ok',   'errors': [],   'data': [     {       'id': '12cc211a-c3c3-4f64-8493-f26f9b98c6f6',       'type': 'upload_session_file',       'attributes': {         'originalFileName': 'testimage1.png',         'fileHash': 'bbf9b9548ee4605c388acb09e8ca83f625e5ff8e241f315eab5291ebd8049c6f',         'fileSize': 18920,         'mimeType': 'image/png',         'version': 1       }     }   ] } ``` Store the data[{index}].id attribute as the `uploadSessionFileId`, this will be the unique identifier for the file you just uploaded.  If you change your mind and want to remove a previously uploaded image, you can send a request like `DELETE /upload/{uploadSessionId}/{uploadSessionFileId}`, expecting a response like ```json {   'response': 'ok' } ```  Finally you can commit your upload session. We opened with a manga-id and group-ids but what we actually want is to upload a chapter. For that we have to build a payload consisting of two things: a chapterDraft and a pageOrder. The payload will look similar to the following:  ```json {   'chapterDraft': {     'volume': '1',     'chapter': '2.5',     'title': 'Read Online',     'translatedLanguage': 'en'   },   'pageOrder': [       '12cc211a-c3c3-4f64-8493-f26f9b98c6f6'   ] } ```  the `chapterDraft` is the chapter data you would like to create, the pageOrder is an ordered list of uploadSessionFileIds you uploaded earlier.  Order didnt matter, now it does. Any files you uploaded but do not specify in this pageOrder array will be deleted.  An example response is: ```json {   'result': 'ok',   'data': {     'id': '14d4639b-5a8f-4f42-a277-b222412930ca',     'type': 'chapter',     'attributes': {       'volume': '1',       'chapter': '2.5',       'title': 'Read Online',       'translatedLanguage': 'en',       'publishAt': null,       'createdAt': '2021-06-16T00:40:22+00:00',       'updatedAt': '2021-06-16T00:40:22+00:00',       'version': 1     },     'relationships': [       {         'id': '145f9110-0a6c-4b71-8737-6acb1a3c5da4',         'type': 'scanlation_group'       },       {         'id': 'f9c33607-9180-4ba6-b85c-e4b5faee7192',         'type': 'manga'       },       {         'id': '41ce3e1a-8325-45b5-af8e-06aaf648a0df',         'type': 'user'       }     ]   } } ```  You just uploaded a chapter. Congratz!  The returned chapter has empty data and dataSaver attributes where otherwise the pages would be. This is because the image processing happens asynchroniously. Depending on how many chapters need to be processed at a given time, it might take a while for this to be updated.  The first time you upload a chapter in a language you didn't upload before, it will have to be approved by staff. Until both imageprocessing and possible approval have happened, your chapter will be held back and not appear on the website or be found in the list and search endpoints.  As mentioned in the beginning, to edit a chapter use the `POST /upload/begin/{chapterId}` endpoint [`begin-edit-session`](#operation/begin-edit-session) with the current chapter version as the json POST-body payload, and the UploadSession will come pre-filled with the remote existing UploadSessionFiles which you can reorder, remove, upload new images and commit your changes afterward as if it was a new chapter.  # Bugs and questions  ## I have a question  You may join [our Discord](https://discord.gg/mangadex)'s #dev-talk-api channel to ask questions or for help.  However we're all busy so please read the docs first, then a second time, or try searching in the channel. Then ask away if you can't figure it out.  ## I found a bug  Please read the docs carefully and **triple-check** the request body you're actually sending to us. If you're sure you found a bug, then congrats and please report it to us so we can fix it!  Every HTTP response from our services has a `X-Request-ID` header whose value is a UUID. Please log it in your client and provide it to us. If your request has a body, please also provide the JSON body you sent (you did check it after all, right?).  If you're not sure how do do this, your client could look something similar to this pseudocode: ``` var httpResponse = httpClient.execute(httpRequest); if (httpResponse.status >= 500) { # feel free to also log it for 4XXs   logger.error(```     Request ID: ${httpResponse.headers['X-Request-ID]'}      Request:     ${httpRequest.body}      Response:     ${httpResponse.body}   ```); } ```
 *
 * API version: 5.5.7
 * Contact: support@mangadex.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package mangadex

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"github.com/antihax/optional"
)

// Linger please
var (
	_ context.Context
)

type ChapterApiService service
/*
ChapterApiService Delete Chapter
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Chapter ID
@return Response
*/
func (a *ChapterApiService) DeleteChapterId(ctx context.Context, id string) (Response, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue Response
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/chapter/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		if localVarHttpResponse.StatusCode == 200 {
			var v Response
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 400 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 403 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 404 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
/*
ChapterApiService Chapter list
Chapter list. If you want the Chapters of a given Manga, please check the feed endpoints.
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param optional nil or *ChapterApiGetChapterOpts - Optional Parameters:
     * @param "Limit" (optional.Int32) - 
     * @param "Offset" (optional.Int32) - 
     * @param "Ids" (optional.Interface of []string) -  Chapter ids (limited to 100 per request)
     * @param "Title" (optional.String) - 
     * @param "Groups" (optional.Interface of []string) - 
     * @param "Uploader" (optional.Interface of Uploader) - 
     * @param "Manga" (optional.Interface of string) - 
     * @param "Volume" (optional.Interface of Volume) - 
     * @param "Chapter" (optional.Interface of Chapter) - 
     * @param "TranslatedLanguage" (optional.Interface of []string) - 
     * @param "OriginalLanguage" (optional.Interface of []string) - 
     * @param "ExcludedOriginalLanguage" (optional.Interface of []string) - 
     * @param "ContentRating" (optional.Interface of []string) - 
     * @param "ExcludedGroups" (optional.Interface of []string) - 
     * @param "ExcludedUploaders" (optional.Interface of []string) - 
     * @param "IncludeFutureUpdates" (optional.String) - 
     * @param "CreatedAtSince" (optional.String) - 
     * @param "UpdatedAtSince" (optional.String) - 
     * @param "PublishAtSince" (optional.String) - 
     * @param "OrderCreatedAt" (optional.String) - 
     * @param "OrderUpdatedAt" (optional.String) - 
     * @param "OrderPublishAt" (optional.String) - 
     * @param "OrderReadableAt" (optional.String) - 
     * @param "OrderVolume" (optional.String) - 
     * @param "OrderChapter" (optional.String) - 
     * @param "Includes" (optional.Interface of []string) - 
@return ChapterList
*/

type ChapterApiGetChapterOpts struct {
    Limit optional.Int32
    Offset optional.Int32
    Ids optional.Interface
    Title optional.String
    Groups optional.Interface
    Uploader optional.Interface
    Manga optional.Interface
    Volume optional.Interface
    Chapter optional.Interface
    TranslatedLanguage optional.Interface
    OriginalLanguage optional.Interface
    ExcludedOriginalLanguage optional.Interface
    ContentRating optional.Interface
    ExcludedGroups optional.Interface
    ExcludedUploaders optional.Interface
    IncludeFutureUpdates optional.String
    CreatedAtSince optional.String
    UpdatedAtSince optional.String
    PublishAtSince optional.String
    OrderCreatedAt optional.String
    OrderUpdatedAt optional.String
    OrderPublishAt optional.String
    OrderReadableAt optional.String
    OrderVolume optional.String
    OrderChapter optional.String
    Includes optional.Interface
}

func (a *ChapterApiService) GetChapter(ctx context.Context, localVarOptionals *ChapterApiGetChapterOpts) (ChapterList, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ChapterList
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/chapter"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.Limit.IsSet() {
		localVarQueryParams.Add("limit", parameterToString(localVarOptionals.Limit.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Offset.IsSet() {
		localVarQueryParams.Add("offset", parameterToString(localVarOptionals.Offset.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Ids.IsSet() {
		localVarQueryParams.Add("ids[]", parameterToString(localVarOptionals.Ids.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.Title.IsSet() {
		localVarQueryParams.Add("title", parameterToString(localVarOptionals.Title.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Groups.IsSet() {
		localVarQueryParams.Add("groups[]", parameterToString(localVarOptionals.Groups.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.Uploader.IsSet() {
		localVarQueryParams.Add("uploader", parameterToString(localVarOptionals.Uploader.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Manga.IsSet() {
		localVarQueryParams.Add("manga", parameterToString(localVarOptionals.Manga.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Volume.IsSet() {
		localVarQueryParams.Add("volume[]", parameterToString(localVarOptionals.Volume.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Chapter.IsSet() {
		localVarQueryParams.Add("chapter", parameterToString(localVarOptionals.Chapter.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.TranslatedLanguage.IsSet() {
		localVarQueryParams.Add("translatedLanguage[]", parameterToString(localVarOptionals.TranslatedLanguage.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.OriginalLanguage.IsSet() {
		localVarQueryParams.Add("originalLanguage[]", parameterToString(localVarOptionals.OriginalLanguage.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.ExcludedOriginalLanguage.IsSet() {
		localVarQueryParams.Add("excludedOriginalLanguage[]", parameterToString(localVarOptionals.ExcludedOriginalLanguage.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.ContentRating.IsSet() {
		localVarQueryParams.Add("contentRating[]", parameterToString(localVarOptionals.ContentRating.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.ExcludedGroups.IsSet() {
		localVarQueryParams.Add("excludedGroups[]", parameterToString(localVarOptionals.ExcludedGroups.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.ExcludedUploaders.IsSet() {
		localVarQueryParams.Add("excludedUploaders[]", parameterToString(localVarOptionals.ExcludedUploaders.Value(), "multi"))
	}
	if localVarOptionals != nil && localVarOptionals.IncludeFutureUpdates.IsSet() {
		localVarQueryParams.Add("includeFutureUpdates", parameterToString(localVarOptionals.IncludeFutureUpdates.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CreatedAtSince.IsSet() {
		localVarQueryParams.Add("createdAtSince", parameterToString(localVarOptionals.CreatedAtSince.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.UpdatedAtSince.IsSet() {
		localVarQueryParams.Add("updatedAtSince", parameterToString(localVarOptionals.UpdatedAtSince.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.PublishAtSince.IsSet() {
		localVarQueryParams.Add("publishAtSince", parameterToString(localVarOptionals.PublishAtSince.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.OrderCreatedAt.IsSet() {
		localVarQueryParams.Add("order[createdAt]", parameterToString(localVarOptionals.OrderCreatedAt.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.OrderUpdatedAt.IsSet() {
		localVarQueryParams.Add("order[updatedAt]", parameterToString(localVarOptionals.OrderUpdatedAt.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.OrderPublishAt.IsSet() {
		localVarQueryParams.Add("order[publishAt]", parameterToString(localVarOptionals.OrderPublishAt.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.OrderReadableAt.IsSet() {
		localVarQueryParams.Add("order[readableAt]", parameterToString(localVarOptionals.OrderReadableAt.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.OrderVolume.IsSet() {
		localVarQueryParams.Add("order[volume]", parameterToString(localVarOptionals.OrderVolume.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.OrderChapter.IsSet() {
		localVarQueryParams.Add("order[chapter]", parameterToString(localVarOptionals.OrderChapter.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Includes.IsSet() {
		localVarQueryParams.Add("includes", parameterToString(localVarOptionals.Includes.Value(), "multi"))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		if localVarHttpResponse.StatusCode == 200 {
			var v ChapterList
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 400 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 403 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
/*
ChapterApiService Get Chapter
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Chapter ID
 * @param optional nil or *ChapterApiGetChapterIdOpts - Optional Parameters:
     * @param "Includes" (optional.Interface of []string) - 
@return ChapterResponse
*/

type ChapterApiGetChapterIdOpts struct {
    Includes optional.Interface
}

func (a *ChapterApiService) GetChapterId(ctx context.Context, id string, localVarOptionals *ChapterApiGetChapterIdOpts) (ChapterResponse, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ChapterResponse
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/chapter/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.Includes.IsSet() {
		localVarQueryParams.Add("includes[]", parameterToString(localVarOptionals.Includes.Value(), "multi"))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		if localVarHttpResponse.StatusCode == 200 {
			var v ChapterResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 403 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 404 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
/*
ChapterApiService Update Chapter
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param contentType
 * @param id Chapter ID
 * @param optional nil or *ChapterApiPutChapterIdOpts - Optional Parameters:
     * @param "Body" (optional.Interface of ChapterEdit) -  The size of the body is limited to 32KB.
@return ChapterResponse
*/

type ChapterApiPutChapterIdOpts struct {
    Body optional.Interface
}

func (a *ChapterApiService) PutChapterId(ctx context.Context, contentType string, id string, localVarOptionals *ChapterApiPutChapterIdOpts) (ChapterResponse, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Put")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ChapterResponse
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/chapter/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	localVarHeaderParams["Content-Type"] = parameterToString(contentType, "")
	// body params
	if localVarOptionals != nil && localVarOptionals.Body.IsSet() {
		
		localVarOptionalBody:= localVarOptionals.Body.Value()
		localVarPostBody = &localVarOptionalBody
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		if localVarHttpResponse.StatusCode == 200 {
			var v ChapterResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 400 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 403 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		if localVarHttpResponse.StatusCode == 404 {
			var v ErrorResponse
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
