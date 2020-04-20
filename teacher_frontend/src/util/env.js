export default {
  currentVersion: process.env.REACT_APP_CURRENT_VERSION,
  defaultApiUri: process.env.REACT_APP_DEFAULT_API_URI,
  buildBranch: process.env.REACT_APP_BUILD_BRANCH || "n/a",
  googleAnalyticsId: process.env.REACT_APP_GOOGLE_ANALYTICS_ID,
  defaultSubdomain: process.env.REACT_APP_DEFAULT_SUBDOMAIN,
  mixpanelToken: process.env.REACT_APP_MIXPANEL_TOKEN,
  nodeEnv: process.env.NODE_ENV
};
