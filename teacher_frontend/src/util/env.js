export default {
  currentVersion: process.env.REACT_APP_CURRENT_VERSION,
  defaultApiUri: process.env.REACT_APP_DEFAULT_API_URI,
  buildBranch: process.env.REACT_APP_BUILD_BRANCH || 'n/a',
  googleAnalyticsId: process.env.REACT_APP_GOOGLE_ANALYTICS_ID,
  mixpanelToken: process.env.REACT_APP_TEACHER_MIXPANEL_TOKEN,
  canvascblUrl: process.env.REACT_APP_CANVASCBL_URL,
  teacherUrl: process.env.REACT_APP_TEACHER_URL,
  accountUrl: process.env.REACT_APP_ACCOUNT_URL,
  nodeEnv: process.env.NODE_ENV
};
