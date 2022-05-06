import React from "react";

/**
 * A list of common TLDs according to me, compiled by different sources.
 * @type {string[]}
 */
export const COMMON_TLDS = [
  // Top list according to: https://www.icdsoft.com/blog/what-are-the-most-popular-tlds-domain-extensions/
  "com",
  "net",
  "org",
  "de",
  "icu",
  "uk",
  "ru",
  "info",
  "top",
  "xyz",
  "tk",
  "cn",
  "ga",
  "cf",
  "nl",

  // Common tlds we might see **our** users using //
  "dev",
  "io",
  "co",
  "gg",
  "gov",
  "mil",
  "ly",
  "app",
  "ai",
  "tv",
  "me",
  "us",
  "edu",
];

/**
 * Will attempt to parse the specified link into a URL object.
 * Will return null if parsing fails.
 *
 * @param tryLink
 * @param quickLink
 * @returns {null|URL}
 */
export function tryParseUrl(tryLink, quickLink = true) {
  let link = tryLink;

  try {
    let hadToInsertSchema = false;

    if (!link.startsWith("http")) {
      link = `http://${link}`;
      hadToInsertSchema = true;
    }

    const url = new URL(link);

    // If it's a basic link and we didn't do anything to it, it's valid //
    if (!hadToInsertSchema) {
      return url;
    }

    /*
     * What I call "quick linking" is for links that are technically invalid, but contain a
     * potential domain using a common tld.
     * For example, google.com would get caught and turned into a link, but something like google.automobile
     * would not.
     *
     * To catch something like google.automobile, you'd need to prefix the link with the proper scheme. (https://)
     * We also don't want to catch links like example.invalid/?somequery=google.com so we compare the parsed url's hostname.
     *
     * The list of COMMON_TLDs is determined by us and what our users might most likely use.
     * It's not perfect, but it'll catch most cases and it's mostly for convenience.
     *
     * For guaranteed links to parse, prefix the schema to the link.
     */
    if (
      COMMON_TLDS.find(tld => url.hostname.endsWith(tld)) &&
      hadToInsertSchema &&
      quickLink
    ) {
      return url;
    }

    return null;
  } catch (e) {
    return null;
  }
}

/**
 * Inserts anchors that it can extract from a string of text
 *
 * Supports paths after the domain.
 * Does NOT support spaces in the url nor does it support non-ascii urls
 *
 * Works by splitting text by spaces. If a link is spotted, we wrap in an anchor tag
 * and return it in the set of elements to avoid using dangerouslySetHtml
 *
 * Returns a react element containing a fragment and potential anchors.
 *
 * @param text
 * @param anchorStyle the style of the link tags.
 * @returns {React.ReactElement}
 */
export function insertAnchorElements(text, anchorStyle = {}) {
  const possibleDomainRegex = /[^\s]+\.[a-z]{2,}[^\s]*/im; // No globals. We want to only match the first instance  //
  const elements = [];

  let remainingText = text;
  let match = possibleDomainRegex.exec(remainingText);

  if (!match) {
    return text;
  }

  while (match) {
    const [link] = match; // First entry is the link text //
    const { index } = match; // Index of the link in the text //
    const afterLinkIndex = index + link.length;

    // Grab all the normal text before the link and put in the elements list //
    if (index !== 0) {
      const startText = remainingText.substring(0, index);
      elements.push(startText);
    }

    const url = tryParseUrl(link);
    if (url) {
      // This is a valid url, time to linkify and push into the elements list //
      elements.push(
        <a
          href={url.href}
          style={anchorStyle}
          target="_blank"
          rel="noopener noreferrer"
        >
          {link}
        </a>,
      );
    } else {
      // The link is invalid, treat as regular text //
      elements.push(link);
    }

    // If there's still text leftover //
    if (afterLinkIndex < remainingText.length) {
      remainingText = remainingText.substring(afterLinkIndex); // There is still some text left //
      match = possibleDomainRegex.exec(remainingText);
    } else {
      // The link was the end of the text //
      break;
    }

    if (!match) {
      // There's no links left, just push the rest and finish up //
      elements.push(remainingText);
      break;
    }
  }

  return (
    <>
      {elements.map((element, index) => (
        <React.Fragment key={index}>{element}</React.Fragment>
      ))}
    </>
  );
}
