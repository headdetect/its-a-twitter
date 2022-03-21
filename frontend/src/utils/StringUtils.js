/**
 * A version of toFixed that doesn't round.
 *
 * For example if you have:
 *  9.999.toFixed(1) -> "10.0". Which is wrong af
 *
 * This will return the following:
 *   toFixedNoRounding(9.9999, 1) -> "9.9" which is good and makes sense
 * 
 * @param num
 * @param dec
 * @returns {string}
 */
export function toFixedNoRounding(num, dec = 2) {
  const calcDec = 10 ** dec;
  return `${Math.trunc(num * calcDec) / calcDec}`;
}

/**
 * Returns the shorthand version of the specified number.
 * Numbers below 10k can have a single decimal.
 * 
 * Eg:
 *  - 1040830 => 1M
 *  - 3049 => 3K
 *  - 5108 => 5.1K
 *  - 9999 => 9.9k
 * @param {*} value 
 * @returns 
 */
export function shortHandNumber(value) {
  if (value >= 1000000000) {
    return `${Math.floor(value / 1000000000)}B`;
  }

  if (value >= 1000000) {
    return `${Math.floor(value / 1000000)}M`;
  }

  if (value >= 10000) {
    return `${Math.floor(value / 1000)}K`;
  }

  // Below 10k gets one decimal place //
  if (value >= 1000) {
    return `${toFixedNoRounding(value / 1000, 1).replace(/\.0$/, "")}K`;
  }

  return value;
}