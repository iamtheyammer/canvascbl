import {
  OAUTH2CONSENT_SET_GET_CONSENT_INFO_ID,
  OAUTH2CONSENT_SET_SEND_CONSENT_ID
} from '../../actions/components/oauth2consent';

export default function (state = {}, action) {
  switch (action.type) {
    case OAUTH2CONSENT_SET_GET_CONSENT_INFO_ID:
      return {
        ...state,
        getConsentInfoId: action.id
      };
    case OAUTH2CONSENT_SET_SEND_CONSENT_ID:
      return {
        ...state,
        sendConsentId: action.id
      };
    default:
      return state;
  }
}
