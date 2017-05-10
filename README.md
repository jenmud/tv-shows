[![Build Status](https://travis-ci.org/jenmud/tv-shows.svg?branch=master)](https://travis-ci.org/jenmud/tv-shows)

# tv-shows
Simple JSON-based web service

## Hosting
This web service is hosted on [Heroku](https://www.heroku.com/) at the following URL https://mighty-beach-93829.herokuapp.com/

_Note: The above app name may change_

## How to use it
This web service requires that you post JSON data to the root [endpoint](https://mighty-beach-93829.herokuapp.com/).

The expected JSON format being posted is should contain a `payload` field mapping to a `array` of tv show objects.

Example of a JSON post containing one TV show:

```
{
    "payload": [
        {
            "country": "USA",
            "description": "Simmering with supernatural elements and featuring familiar and fan-favourite characters from the immensely popular drama The Vampire Diaries, it's The Originals. This sexy new series centres on the Original vampire family and the dangerous vampire/werewolf hybrid, Klaus, who returns to the magical melting pot that is the French Quarter of New Orleans, a town he helped build centuries ago.",
            "drm": true,
            "episodeCount": 1,
            "genre": "Action",
            "image": {
                "showImage": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheOriginals1280.jpg"
            },
            "language": "English",
            "nextEpisode": {
                "channel": null,
                "channelLogo": "http://catchup.ninemsn.com.au/img/player/logo_go.gif",
                "date": null,
                "html": "<br><span class=\"visit\">Visit the Official Website</span></span>",
                "url": "http://go.ninemsn.com.au/"
            },
            "primaryColour": "#df0000",
            "seasons": [
                {
                    "slug": "show/theoriginals/season/1"
                }
            ],
            "slug": "show/theoriginals",
            "title": "The Originals",
            "tvChannel": "GO!"
        }
    ]
}
```

### Supported fields for a TV show object


| Key           | Value   |
| ------------- |:-------:|
| country       | string  |
| description   | string  |
| drm           | bool    |
| episodeCount  | int     |
| genre         | string  |
| image         | map mapping key to value |
| language      | string  |
| nextEpisode   | null or map mapping key to value    |
| primaryColour | string  |
| seasons       | array of maps mapping key to value |
| slug          | string  |
| title         | string  |
| tvChannel     | string  |

## Successful POST
On a successful POST, the expected JSON response format should contain a `response` field mapping to a `array` of simplified information objects.

Example response object based on the above example:

```
{
    "response": [
        {
            "image": "http://catchup.ninemsn.com.au/img/jump-in/shows/TheOriginals1280.jpg",
            "slug": "show/theoriginals",
            "title": "The Originals"
        }
    ]
}

```
### Supported fields for a JSON response`


| Key           | Value   |
| ------------- |:-------:|
| image         | string  |
| slug          | string  |
| title         | string   |

## Errors
A 400 Bad Request HTTP status is returned if it could not parse the JSON data.

```
{
    "error": "Could not decode request: JSON parsing failed"
}
```

