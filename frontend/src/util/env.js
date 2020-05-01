export default {
  currentVersion: process.env.REACT_APP_CURRENT_VERSION,
  defaultApiUri: process.env.REACT_APP_DEFAULT_API_URI,
  googleAnalyticsId: process.env.REACT_APP_GOOGLE_ANALYTICS_ID,
  privacyPolicyUrl: process.env.REACT_APP_PRIVACY_POLICY_URL,
  termsOfServiceUrl: process.env.REACT_APP_TERMS_OF_SERVICE_URL,
  stripeApiKeyPub: process.env.REACT_APP_STRIPE_API_KEY_PUB,
  buildBranch: process.env.REACT_APP_BUILD_BRANCH || 'n/a',
  upgradesPurchasableProductId: parseInt(
    process.env.REACT_APP_UPGRADES_PURCHASABLE_PRODUCT_ID
  ),
  defaultSubdomain: process.env.REACT_APP_DEFAULT_SUBDOMAIN,
  mixpanelToken: process.env.REACT_APP_MIXPANEL_TOKEN,
  teacherUrl: process.env.REACT_APP_TEACHER_URL,
  canvascblUrl: process.env.REACT_APP_CANVASCBL_URL,
  accountUrl: process.env.REACT_APP_ACCOUNT_URL,
  nodeEnv: process.env.NODE_ENV
};
