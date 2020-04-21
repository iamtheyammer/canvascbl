import React from "react";
import * as PropTypes from "prop-types";
import { Link } from "react-router-dom";
import mixpanel from "mixpanel-browser";
import env from "./env";

mixpanel.init(env.mixpanelToken);

mixpanel.set_config({
  ignore_dnt: true,
  debug: env.nodeEnv === "development"
});

const env_check = true;

/**
 * Turns a bool into a string.
 * @param {boolean} b The boolean to stringify.
 * @returns {string} A string, like "true" or "false".
 */
// const bts = b => `${b}`;

const mp = {
  identify: id => {
    if (env_check) mixpanel.identify(id);
  },
  alias: id => {
    if (env_check) mixpanel.alias(id);
  },
  track: (name, props, callback) => {
    if (env_check) mixpanel.track(name, props, callback);
  },
  track_links: (query, event_name, properties) => {
    if (env_check) mixpanel.track_links(query, event_name, properties);
  },
  register: props => {
    if (env_check) mixpanel.register(props);
  },
  people: {
    set: props => {
      if (env_check) mixpanel.people.set(props);
    }
  }
};

export const pageNames = {
  profile: "Profile",
  courses: "Courses",
  courseOverview: "Course Overview"
};

const courseOverviewRegex = /^\/dashboard\/courses\/[0-9]+_[0-9]+\/overview$/;

export function pageNameFromPath(path) {
  switch (path) {
    case "/dashboard/profile":
      return pageNames.profile;
    case "/dashboard/courses":
      return pageNames.courses;
    default:
      if (courseOverviewRegex.test(path)) {
        return pageNames.courseOverview;
      }
  }
}

export const vias = {
  dashboardMenu: "Dashboard Menu",
  breadcrumb: "Breadcrumb",
  coursesCourseCard: "Courses Course Card",
  notATeacherPopup: "Not a Teacher Popup"
};

export const destinationNames = {
  courseOverview: "Course Overview",
  canvascblForStudentsAndParents: "CanvasCBL for Students and Parents",
  googleForms: "Google Forms"
};

export const destinationTypes = {
  canvascbl: "CanvasCBL",
  canvascblLogout: "CanvasCBL Logout",
  canvascblForTeachersFeedbackForm: "CanvasCBL for Teachers Feedback Form"
};

export const tabImplementations = {};

export const tableNames = {};

export const itemTypes = {};

/**
 * Tracks dashboard loads. We're defining a load as someone visiting the dashboard, whether
 * they have to reauthenticate or not. This function does four things:
 * - It identifies the user (mixpanel.identify)
 * - It sets some super properties (like subscription status)
 * - It sets people properties (to help mixpanel identify the user)
 * - It tracks a sign in event
 * @param {string} name The user's full name
 * @param {string} email The user's email
 * @param {number} userId The user's ID (from CanvasCBL, not Canvas)
 * @param {number} canvasUserId The user's Canvas ID
 * @param {number|string} currentVersion The current version
 * @param {number|string} prevVersion The user's previous version
 */
export function trackDashboardLoad(
  name,
  email,
  userId,
  canvasUserId,
  currentVersion,
  prevVersion
) {
  mp.identify(userId);

  mp.register({
    "Current Version": `${currentVersion}`
  });

  mp.people.set({
    $name: name,
    $email: email,
    "CanvasCBL User ID": userId,
    "Canvas User ID": canvasUserId
  });

  mp.track("Dashboard Load", {
    $name: name,
    $email: email,
    "CanvasCBL User ID": userId,
    "User Last Version": `${prevVersion}`
  });
}

/**
 * Tracks a page view.
 * @param {string} pageName Human-readable name of the page to track
 * @param {number} [courseId] Canvas course ID
 */
export function trackPageView(pageName, courseId) {
  mp.track("Page View", {
    "Page Name": pageName,
    "Course ID": courseId
  });
}

/**
 * Tracks a navigation from one page to another.
 * @param {string} to Where the user went. Get this from `pageNames`.
 * @param {string} via How the user got there. Get this from `vias`.
 */
export function trackNavigation(to, via) {
  mp.track("Navigation", {
    To: to,
    Via: via
  });
}

/**
 * Tracks a click to an external link.
 * @param {string} anchorId The ID of the <a> tag.
 * @param {string} destinationUrl The URL of the destination.
 * @param {string} destinationName The name of the destination. Get this from `destinationNames`. Examples: Canvas, Privacy Policy
 * @param {string} destinationType The type of the destination. Get this from `destinationTypes`. Examples: Outcome, Assignment
 * @param {string} via How the user got there. Get this from `vias`.
 */
export function trackExternalLinkClick(
  anchorId,
  destinationUrl,
  destinationName,
  destinationType,
  via
) {
  mp.track_links(`#${anchorId}`, "External Link Click", {
    "Destination URL": destinationUrl,
    "Destination Name": destinationName,
    "Destination Type": destinationType,
    Via: via
  });

  // mp.track('External Link Click', {
  //   'Destination URL': destinationUrl,
  //   'Destination Name': destinationName,
  //   'Destination Type': destinationType,
  //   Via: via
  // });
}

/**
 * Tracks an external link click that isn't an anchor tag.
 * @param {string} destinationUrl The URL of the destination.
 * @param {string} destinationName The name of the destination. Get this from `destinationNames`. Examples: Canvas, Privacy Policy
 * @param {string} destinationType The type of the destination. Get this from `destinationTypes`. Examples: Outcome, Assignment
 * @param {string} via How the user got there. Get this from `vias`.
 */
export function trackExternalLinkClickOther(
  destinationUrl,
  destinationName,
  destinationType,
  via
) {
  mp.track("External Link Click", {
    "Destination URL": destinationUrl,
    "Destination Name": destinationName,
    "Destination Type": destinationType,
    Via: via
  });
}

/**
 * Tracks a logout.
 * @param via - How the user got to logout
 */
export function trackLogout(via) {
  mp.track("Logout", {
    Via: via
  });
}

/*

Components

 */

/**
 * TrackingLink returns a link that calls `trackNavigation` when clicked.
 * @param props See the PropTypes for this function.
 * @returns {React.FunctionComponent}
 * @constructor
 */
export function TrackingLink(props) {
  const { to, pageName, via, style, children } = props;
  return (
    <Link to={to} onClick={() => trackNavigation(pageName, via)} style={style}>
      {children}
    </Link>
  );
}

TrackingLink.propTypes = {
  to: PropTypes.string.isRequired,
  pageName: PropTypes.oneOf(Object.values(pageNames)),
  via: PropTypes.oneOf(Object.values(vias)).isRequired,
  style: PropTypes.object,
  children: PropTypes.any
};
