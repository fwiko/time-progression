import datetime
import json

import pytz
from flask import Flask, jsonify, request
from waitress import serve

import settings
from modules.progression_calculator import ProgressionCalculator

app = Flask(__name__)


@app.route("/")
def _main() -> None:
    def error_response(message):
        return jsonify({"error": message})

    data_format = request.args.get("format")
    timezone = request.args.get("timezone", settings.DEFAULT_TIMEZONE)
    if timezone != "UTC":
        if timezone not in ProgressionCalculator.available_timezones():
            return error_response("Invalid timezone")
    timezone = timezone.replace("-", "/")
    if data_format and data_format not in settings.FORMAT_METHODS:
        return error_response("Invalid format method")
    elif data_format and data_format in settings.FORMAT_METHODS:
        data = ProgressionCalculator(timezone).simplify(data_format).__dict__
    else:
        data = ProgressionCalculator(timezone).precise().__dict__
    return jsonify(data)


@app.route("/timezones")
def _timezones() -> None:
    return jsonify(ProgressionCalculator.available_timezones())


if __name__ == "__main__":
    from waitress import serve

    serve(app, port=settings.PORT)
