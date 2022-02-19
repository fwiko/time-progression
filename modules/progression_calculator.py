import calendar
import datetime
import json
import math
from dataclasses import dataclass

import pytz

HOURS_IN_WEEK = 24 * 7
SECONDS_IN_DAY = 24 * 60 * 60
SECONDS_IN_HOUR = 60 * 60


@dataclass
class CalcData:
    def __init__(self, **kwargs: dict) -> None:
        """Calcdata object constructor.

        Args:
            kwargs: Keyword arguments that will be set as attributes of the dataclass.
        """
        for a, b in kwargs.items():
            setattr(self, a, b)

    def __repr__(self):
        return f'CalcData({", ".join(f"{k}: {v}" for k, v in self.__dict__.items())})'


class ProgressionCalculator:
    def __init__(
        self,
        timezone: str = "UTC",
        date: datetime.datetime = None,
    ) -> None:
        """ProgressionCalculator class used to calculate the how complete all attributes of a given datetime object are in accordance to their maximum values.

        Args:
            timezone (str, optional): The timezone to base calculations upon. Defaults to "UTC".
            date (datetime.datetime, optional): Optional datetime object if you choose to specify a custom date/time. Defaults to datetime.datetime.now(tz=pytz.UTC).
        """
        if not date:
            date = datetime.datetime.now(tz=pytz.UTC)
        if timezone not in pytz.all_timezones:
            self.timezone = "UTC"
        else:
            self.timezone = timezone
        self.date = (
            date.astimezone(pytz.timezone(self.timezone))
            if date.tzinfo is not None
            else date
        )
        self.days_in_year = 366 if calendar.isleap(date.year) else 365
        self.days_in_month = calendar.monthrange(date.year, date.month)[1]

    def year_percent(self):
        """Returns the percentage in which the year is complete. Based on the current datetime object.

        Returns:
            float: Percentage value of how much of the year has passed, based on the max seconds in a year.
        """
        return (
            (
                self.date.timetuple().tm_yday * SECONDS_IN_DAY
                + self.date.hour * SECONDS_IN_HOUR
                + self.date.minute * 60
                + self.date.second
            )
            / (self.days_in_year * SECONDS_IN_DAY)
            * 100
        )

    def month_percent(self):
        """Returns the percentage in which the month is complete. Based on the current datetime object.

        Returns:
            float: Percentage value of how much of the month has passed, based on the max seconds in a month.
        """
        return (
            (
                self.date.day * SECONDS_IN_DAY
                + self.date.hour * SECONDS_IN_HOUR
                + self.date.minute * 60
                + self.date.second
            )
            / (self.days_in_month * SECONDS_IN_DAY)
            * 100
        )

    def week_percent(self):
        """Returns the percentage in which the week is complete. Based on the current datetime object.

        Returns:
            float: Percentage value of how much of the week has passed, based on the max seconds in a week.
        """
        return (
            (
                (self.date.weekday() * 24 + self.date.hour) * SECONDS_IN_HOUR
                + self.date.minute * 60
                + self.date.second
            )
            / (HOURS_IN_WEEK * SECONDS_IN_HOUR)
            * 100
        )

    def day_percent(self):
        """Returns the percentage in which the day is complete. Based on the current datetime object.

        Returns:
            float: Percentage value of how much of the day has passed, based on the max seconds in a day.
        """
        return (
            (
                self.date.hour * SECONDS_IN_HOUR
                + self.date.minute * 60
                + self.date.second
            )
            / SECONDS_IN_DAY
            * 100
        )

    def hour_percent(self):
        """Returns the percentage in which the hour is complete. Based on the current datetime object.

        Returns:
            float: Percentage value of how much of the hour has passed, based on the max seconds in an hour.
        """
        return (self.date.minute * 60 + self.date.second) / SECONDS_IN_HOUR * 100

    def minute_percent(self):
        """Returns the percentage in which the minute is complete. Based on the current datetime object.

        Returns:
            float: Percentage value of how much of the minute has passed, based on the max seconds in a minute.
        """
        return self.date.second / 60 * 100

    @staticmethod
    def available_timezones() -> list[str]:
        """Return a list of supported timezones of the pytz library.

        Returns:
            list[str]: List of supported timezones (e.g. ["UTC", "Europe-Berlin"])
        """
        return list(map(lambda x: x.replace("/", "-"), pytz.all_timezones))

    def _timestamp(self) -> str:
        """Return a string representation of the date.".

        Returns:
            str: String representation of the date (Formatted as "YYYY-MM-DD HOUR:MINUTE:SECOND")
        """
        return self.date.strftime("%Y-%m-%d %H:%M:%S")

    def _createobject(self, data: dict) -> CalcData:
        """Returns a CalcData object with the given percentage data.

        Args:
            data (dict): Dictionary with the percentage data based on the datetime object

        Returns:
            CalcData: Object with all relevant data to the datetime object
        """
        return CalcData(
            **{
                "timezone": self.timezone,
                "timestamp": self._timestamp(),
                "data": data,
            }
        )

    def _calculate(self) -> dict:
        """Returns a dictionary with percentage complete data based on the datetime object

        Returns:
            dict: Percentage complete data calculated according to the current year, month, week, day, hour, and minute
        """
        return {
            "year": self.year_percent(),
            "month": self.month_percent(),
            "week": self.week_percent(),
            "day": self.day_percent(),
            "hour": self.hour_percent(),
            "minute": self.minute_percent(),
        }

    def simplify(self, method: str = "round") -> CalcData:
        """Returns a CalcData object with the calculated percentage data in a simplified form.

        Args:
            method (str, optional): Type of simplification to apply, round up, round down, round based on decimal. Defaults to "round".

        Returns:
            CalcData: CalcData object with the calculated percentage data.
        """
        method = method.lower()
        if method not in ("round", "ceil", "floor"):
            return self.precise()
        function = getattr(math, method) if method != "round" else round
        return self._createobject(
            dict(map(lambda x: (x[0], function(x[1])), self._calculate().items()))
        )

    def precise(self) -> CalcData:
        """Returns a CalcData object with the calculated percentage data without any simplification.

        Returns:
            CalcData: CalcData object with the calculated percentage data.
        """
        return self._createobject(self._calculate())
