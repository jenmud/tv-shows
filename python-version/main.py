# pylint: disable=invalid-name
"""
Simple flask web server.
"""
import logging
from flask import Flask, request, jsonify


LOGGER = logging.getLogger(__name__)
app = Flask(__name__)


def filter_tvShows_for_drm(data):
    """
    Filter data for all shows that have their DRM enabled.

    :param data: Data that your are filtering for shows their DRM enabled.
    :type data: iterable of :class:`dict`
    :returns: All shows with their DRM enabled.
    :rtype: iterable of :class:`dict`
    """
    for each in data:
        drm = each.get("drm", False)
        if drm is True:
            yield each


def filter_tvShows_with_episodes(data):
    """
    Filter data for all shows that have a episode count greater
    then 0.

    :param data: Data that your are filtering for shows with 1
        or more episodes.
    :type data: iterable of :class:`dict`
    :returns: All shows with 1 or more episodes.
    :rtype: iterable of :class:`dict`
    """
    for each in data:
        if each.get("episodeCount", 0) > 0:
            yield each


def make_json_response(data):
    """
    Take some data and returns a dict.

    :param data: Data that you are making a response with.
    :type data: iterable of :class:`dict`
    :returns: A JSON reponse containing a `response` field which is a
        iterable of :class:`dict`.
    :rtype: :class:`dict` containing a :class:`dict`
    """
    response = []

    for each in data:
        response.append(
            {
                "image": each.get("image", {}).get("showImage"),
                "slug": each["slug"],
                "title": each["title"],
            }
        )

    return {
        "response": response,
    }


@app.route("/", methods=["GET", "POST"])
def index():
    """
    Index entrypoint which accepts POST JSON data with a `payload` field
    mapping to a array of shows.

    :returns: JSON repsonse of simplified show information.
    :rtype: :class:`flask.Response`
    """
    try:
        shows = request.json.get("payload", [])
        shows = filter_tvShows_for_drm(shows)
        shows = filter_tvShows_with_episodes(shows)
        return jsonify(make_json_response(shows))
    except Exception:  # pylint: disable=broad-except
        LOGGER.exception("JSON parsing error")
        return jsonify(
            {
                "error": "Could not decode request: JSON parsing failed",
            }
        )


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    app.run(debug=True)
