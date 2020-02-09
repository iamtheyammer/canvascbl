export const OAUTH2CONSENT_SET_GET_CONSENT_INFO_ID =
  'OAUTH2CONSENT_SET_GET_CONSENT_INFO_ID';
export const OAUTH2CONSENT_SET_SEND_CONSENT_ID = 'OAUTH2_SET_SEND_CONSENT_ID';

export function setGetConsentInfoId(id) {
  return {
    type: OAUTH2CONSENT_SET_GET_CONSENT_INFO_ID,
    id
  };
}

export function setSendConsentId(id) {
  return {
    type: OAUTH2CONSENT_SET_SEND_CONSENT_ID,
    id
  };
}
