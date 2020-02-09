import {
  OAUTH2_GET_AUTHORIZED_APPS_ERROR,
  OAUTH2_GET_CONSENT_INFO_ERROR,
  OAUTH2_GOT_AUTHORIZED_APPS,
  OAUTH2_GOT_CONSENT_INFO,
  OAUTH2_REVOKE_GRANT_ERROR,
  OAUTH2_REVOKED_GRANT,
  OAUTH2_SEND_CONSENT_ERROR,
  OAUTH2_SENT_CONSENT,
  OAUTH2_SET_CONSENT_CODE
} from '../actions/oauth2';

export default function oauth2(state = {}, action) {
  switch (action.type) {
    case OAUTH2_SET_CONSENT_CODE:
      return {
        ...state,
        consentCode: action.consentCode
      };
    case OAUTH2_GET_CONSENT_INFO_ERROR:
      return {
        ...state,
        getConsentInfoError: action.err
      };
    case OAUTH2_GOT_CONSENT_INFO:
      return {
        ...state,
        consentInfo: action.info
      };
    case OAUTH2_SEND_CONSENT_ERROR:
      return {
        ...state,
        sendConsentError: action.err
      };
    case OAUTH2_SENT_CONSENT:
      return {
        ...state,
        consentRedirectUrl: action.redirectTo
      };
    case OAUTH2_GET_AUTHORIZED_APPS_ERROR:
      return {
        ...state,
        getAuthorizedAppsError: action.err
      };
    case OAUTH2_GOT_AUTHORIZED_APPS:
      return {
        ...state,
        authorizedApps: action.authorizedApps
      };
    case OAUTH2_REVOKE_GRANT_ERROR:
      return {
        ...state,
        revokeGrantErrors: {
          ...state.revokeGrantErrors,
          [action.grantId]: action.err
        }
      };
    case OAUTH2_REVOKED_GRANT:
      return {
        ...state,
        revokedGrantIds: state.revokedGrantIds
          ? [...state.revokedGrantIds, action.grantId]
          : [action.grantId]
      };
    default:
      return state;
  }
}
