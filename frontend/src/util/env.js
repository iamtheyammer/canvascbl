export default {
  currentVersion: process.env.REACT_APP_CURRENT_VERSION,
  defaultApiUri: process.env.REACT_APP_DEFAULT_API_URI,
  googleAnalyticsId: process.env.REACT_APP_GOOGLE_ANALYTICS_ID,
  privacyPolicyUrl: process.env.REACT_APP_PRIVACY_POLICY_URL,
  stripeApiKeyPub: process.env.REACT_APP_STRIPE_API_KEY_PUB,
  upgradesPurchasableProductId: parseInt(
    process.env.REACT_APP_UPGRADES_PURCHASABLE_PRODUCT_ID
  ),
  nodeEnv: process.env.NODE_ENV
};
