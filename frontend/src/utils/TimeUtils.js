const userLocal = navigator.languages?.[0] ?? navigator.language;

const agoFullFormat = new Intl.DateTimeFormat(userLocal ?? "en-US", {
  month: "short",
  day: "2-digit",
  year: "2-digit",
});

const agoSmallFormat = new Intl.DateTimeFormat(userLocal ?? "en-US", {
  month: "short",
  day: "2-digit",
});

/**
 * Converts the duration between now and the specified dateTime
 * string into a readable format.
 *
 * Rules:
 * If duration hours < 1, the minute difference will be used.
 * If duration hours > 24, the date will be used.
 * Otherwise the full format w/year will be used.
 *
 * Eg.
 *  Now = Jan 1, 2020 @ 12:00 UTC.
 *  DateTime = Jan 1, 2020 @ 5:00 UTC.
 *  Result = "7h"
 *
 * Eg.
 *  Now = Jan 3, 2020 @ 12:00 UTC.
 *  DateTime = Jan 1, 2020 @ 5:00 UTC.
 *  Result = "Jan 1"
 *
 * Eg.
 *  Now = Jan 3, 2021 @ 12:00 UTC.
 *  DateTime = Jan 1, 2020 @ 5:00 UTC.
 *  Result = "Jan 1, 2020"
 *
 * @param {Date} dateTime The dateTime to convert
 */
export function toAgoString(dateTime) {
  const now = new Date();

  if (now.getTime() - dateTime.getTime() < 0) {
    return "In the future";
  }

  if (now.getTime() - dateTime.getTime() < 60e3) {
    return `${Math.floor((now.getTime() - dateTime.getTime()) / 1000)}s`;
  }

  const dateDifference =
    (now.getTime() - dateTime.getTime()) / (1000 * 3600 * 24);

  const hourDifference = dateDifference * 24;

  // Difference is less than 1h //
  if (hourDifference < 1) {
    return `${Math.ceil(hourDifference * 60)}m`;
  }

  // Difference is less than 24h //
  if (hourDifference < 24) {
    return `${Math.ceil(hourDifference)}h`;
  }

  // Difference is less than 28d //
  if (Math.floor(dateDifference) <= 28) {
    return `${Math.ceil(dateDifference)}d`;
  }

  // Years are the same //
  if (now.getFullYear() === dateTime.getFullYear()) {
    return agoSmallFormat.format(dateTime);
  }

  return agoFullFormat.format(dateTime);
}
