import React from 'react';
import * as PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import mixpanel from 'mixpanel-browser';
import env from './env';

mixpanel.init(env.mixpanelToken);

mixpanel.set_config({
  ignore_dnt: true,
  debug: env.nodeEnv === 'development'
});

const env_check = true;

/**
 * Turns a bool into a string.
 * @param {boolean} b The boolean to stringify.
 * @returns {string} A string, like "true" or "false".
 */
const bts = b => `${b}`;

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
  profile: 'Profile',
  grades: 'Grades',
  gradeBreakdown: 'Grade Breakdown',
  upgrades: 'Upgrades',
  redeem: 'Redeem',
  settings: 'Settings',
  logout: 'Logout',
  authorize: 'OAuth2 Authorize'
};

export function pageNameFromPath(path) {
  switch (path) {
    case '/dashboard/profile':
      return pageNames.profile;
    case '/dashboard/grades':
      return pageNames.grades;
    case '/dashboard/upgrades':
      return pageNames.upgrades;
    case '/dashboard/upgrades/redeem':
      return pageNames.redeem;
    case '/dashboard/settings':
      return pageNames.settings;
    case '/dashboard/authorize':
      return pageNames.authorize;
    default:
      if (path.startsWith('/dashboard/grades/')) {
        return pageNames.gradeBreakdown;
      }
  }
}

export const vias = {
  dashboardMenu: 'Dashboard Menu',
  moreActionsSubmenu: 'More Actions Submenu',
  breadcrumb: 'Breadcrumb',
  gradesTableOpenOnCanvas: 'Grades Table Open On Canvas',
  gradesTableSeeBreakdownLink: 'Grades Table See Breakdown Link',
  gradesTableCourseName: 'Grades Table Course Name',
  gpaReportCardQuestionIcon: 'GPA Question Mark Icon (Report Card)',
  gpaTraditionalQuestionIcon: 'GPA Question Mark Icon (Traditional)',
  gpaLearnMore: 'GPA Learn More Here',
  gradeBreakdownOutcomesTableAssignmentsTableOpenOnCanvas:
    'Grade Breakdown Outcome Assignments Table Open On Canvas',
  gradeBreakdownOutcomesTableFutureAssignmentsTableOpenOnCanvas:
    'Grade Breakdown Outcome Future Assignments Table Open On Canvas',
  noCurrentSubscriptionRedeem:
    'Redeem Link On The Upgrades (No Current Subscription) Page',
  mobileGradesObserveeSwitcher: 'Mobile Grades Observee Switcher',
  mobileNavBarLogo: 'Mobile Logo At Top',
  breakdownUnavailableBackToGrades: 'Grade Breakdown Unavailable Back To Grades'
};

export const destinationNames = {
  canvas: 'Canvas',
  helpdesk: 'CanvasCBL Helpdesk',
  privacyPolicy: 'CanvasCBL Privacy Policy',
  termsOfService: 'CanvasCBL Terms of Service',
  statusPage: 'CanvasCBL Status Page',
  extension: 'CanvasCBL Extension Download Page'
};

export const destinationTypes = {
  outcome: 'Outcome',
  course: 'Course',
  assignment: 'Assignment',
  helpdesk: {
    home: 'CanvasCBL Helpdesk Home',
    gpas: 'CanvasCBL Helpdesk Article on GPAs'
  }
};

export const tabImplementations = {
  gradeCard: {
    name: 'Grade Breakdown Grade Card',
    tabNames: {
      userGrade: 'Your Grade',
      averageGrade: 'Average Grade'
    }
  },
  outcomeInfo: {
    name: 'Grade Breakdown Outcome Info Card',
    tabNames: {
      lowestOutcome: 'Lowest Outcome',
      averageOutcomeScore: 'Average Outcome Score',
      toGetAnA: 'How To Get An A',
      moreInfo: 'More Info'
    }
  }
};

export const tableNames = {
  gradeBreakdown: {
    outcomes: 'Grade Breakdown Outcomes'
  }
};

export const itemTypes = {
  outcome: 'Outcome'
};

/**
 * Tracks dashboard loads. We're defining a load as someone visiting the dashboard, whether
 * they have to reauthenticate or not. This function does four things:
 * - It identifies the user (mixpanel.identify)
 * - It sets some super properties (like subscription status)
 * - It sets people properties (to help mixpanel identify the user)
 * - It tracks a sign in event
 * @param {string} name The user's full name
 * @param {string} email The user's email
 * @param {boolean} hasValidSubscription Whether the user has a valid subscription.
 * @param {string} subscriptionStatus The user's subscription status (ex: active)
 * @param {number} userId The user's ID (from CanvasCBL, not Canvas)
 * @param {number} canvasUserId The user's Canvas ID
 * @param {number} activeUserId The active user's Canvas ID
 * @param {number|string} currentVersion The current version
 * @param {number|string} prevVersion The user's previous version
 */
export function trackDashboardLoad(
  name,
  email,
  hasValidSubscription,
  subscriptionStatus,
  userId,
  canvasUserId,
  activeUserId,
  currentVersion,
  prevVersion
) {
  mp.identify(userId);

  mp.register({
    'Subscription Status': subscriptionStatus,
    'Current Version': `${currentVersion}`,
    'Active User ID': activeUserId
  });

  mp.people.set({
    $name: name,
    $email: email,
    'CanvasCBL User ID': userId,
    'Canvas User ID': canvasUserId,
    'Has Valid Subscription': hasValidSubscription,
    'Subscription Status': subscriptionStatus
  });

  mp.track('Dashboard Load', {
    $name: name,
    $email: email,
    'Subscription Status': subscriptionStatus,
    'CanvasCBL User ID': userId,
    'User Last Version': `${prevVersion}`
  });
}

/**
 * Tracks a page view.
 * @param {string} pageName Human-readable name of the page to track
 * @param {number} [courseId] Canvas course ID
 */
export function trackPageView(pageName, courseId) {
  mp.track('Page View', {
    'Page Name': pageName,
    'Course ID': courseId
  });
}

/**
 * Tracks a navigation from one page to another.
 * @param {string} to Where the user went. Get this from `pageNames`.
 * @param {string} via How the user got there. Get this from `vias`.
 */
export function trackNavigation(to, via) {
  mp.track('Navigation', {
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
  mp.track_links(`#${anchorId}`, 'External Link Click', {
    'Destination URL': destinationUrl,
    'Destination Name': destinationName,
    'Destination Type': destinationType,
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
 * Tracks a notification status toggle.
 * @param {boolean} notificationStatus The new status
 * @param {string} notificationTypeShortName
 */
export function trackNotificationStatusToggle(
  notificationStatus,
  notificationTypeShortName
) {
  mp.track('Notification Status Toggle', {
    // Notification Status should be submitted as a string
    'Notification Status': bts(notificationStatus),
    'Notification Type Short Name': notificationTypeShortName
  });
}

/**
 * Tracks a row expansion on a table.
 * @param {string} tableName The name of the table. Get this from `tableNames`.
 * @param {number} expandedItemId The ID of the expanded item.
 * @param {string} expandedItemType The type of the expanded item. Get this from `itemTypes`.
 * @param {boolean} expansionStatus Whether the row was opened (true) or closed (false).
 * @param {number} courseId The ID of the course, if applicable.
 */
export function trackTableRowExpansion(
  tableName,
  expandedItemId,
  expandedItemType,
  expansionStatus,
  courseId
) {
  mp.track('Table Row Expansion', {
    'Table Name': tableName,
    'Expanded Item ID': expandedItemId,
    'Expanded Item Type': expandedItemType,
    'Expansion Status': expansionStatus,
    'Course ID': courseId
  });
}

/**
 * Tracks tab changes.
 * @param {string} containerName The name of the thing with the tabs. Get this from `tabImplementations`.
 * @param {string} newTabName The name of the tab that was just selected. Get this from `tabImplementations`.
 */
export function trackTabChange(containerName, newTabName) {
  mp.track('Tab Change', {
    'Container Name': containerName,
    'New Tab Name': newTabName
  });
}

/**
 * Tracks an active user change.
 * @param activeUserId The new active user ID.
 * @param via How the user changed active users. Get this from `vias`.
 */
export function trackActiveUserChange(activeUserId, via) {
  mp.register({ 'Active User ID': activeUserId });

  mp.track('Active User Change', {
    Via: via
  });
}

/**
 * Tracks a user's OAuth2 Authorize decision.
 * @param {number} oAuth2CredentialId The ID of the OAuth2 Credential.
 * @param {boolean} didAuthorize Whether the user authorized the app.
 * @param {string} consentCode The consent code for the interaction.
 * @param {string[]} scopes The scopes of the OAuth2 Request.
 */
export function trackOAuth2Decision(
  oAuth2CredentialId,
  didAuthorize,
  consentCode,
  scopes
) {
  mp.track('OAuth2 Decision', {
    'OAuth2 Credential ID': oAuth2CredentialId,
    'Did Authorize App': bts(didAuthorize),
    'Consent Code': consentCode,
    Scopes: scopes
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
