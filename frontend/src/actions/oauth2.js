export const OAUTH2_SET_CONSENT_CODE = 'OAUTH2_SET_CONSENT_CODE';

export const OAUTH2_GET_CONSENT_INFO = 'OAUTH2_GET_CONSENT_INFO';
export const OAUTH2_GET_CONSENT_INFO_ERROR = 'OAUTH2_GET_CONSENT_INFO_ERROR';
export const OAUTH2_GOT_CONSENT_INFO = 'OAUTH2_GOT_CONSENT_INFO';

export const OAUTH2_SEND_CONSENT = 'OAUTH2_SEND_CONSENT';
export const OAUTH2_SEND_CONSENT_ERROR = 'OAUTH2_SEND_CONSENT_ERROR';
export const OAUTH2_SENT_CONSENT = 'OAUTH2_SENT_CONSENT';

export const OAUTH2_GET_AUTHORIZED_APPS = 'OAUTH2_GET_AUTHORIZED_APPS';
export const OAUTH2_GET_AUTHORIZED_APPS_ERROR =
  'OAUTH2_GET_AUTHORIZED_APPS_ERROR';
export const OAUTH2_GOT_AUTHORIZED_APPS = 'OAUTH2_GOT_AUTHORIZED_APPS';

export const OAUTH2_REVOKE_GRANT = 'OAUTH2_REVOKE_GRANT';
export const OAUTH2_REVOKE_GRANT_ERROR = 'OAUTH2_REVOKE_GRANT_ERROR';
export const OAUTH2_REVOKED_GRANT = 'OAUTH2_REVOKED_GRANT';

export function setConsentCode(consentCode) {
  return {
    type: OAUTH2_SET_CONSENT_CODE,
    consentCode
  };
}

export function getConsentInfo(id, consentCode) {
  return {
    type: OAUTH2_GET_CONSENT_INFO,
    id,
    consentCode
  };
}

export function getConsentInfoError(consentCode, err) {
  return {
    type: OAUTH2_GET_CONSENT_INFO_ERROR,
    consentCode,
    err
  };
}

export function gotConsentInfo(consentCode, info) {
  return {
    type: OAUTH2_GOT_CONSENT_INFO,
    consentCode,
    info
  };
}

export function sendConsent(id, consentCode, action) {
  return {
    type: OAUTH2_SEND_CONSENT,
    id,
    consentCode,
    action
  };
}

export function sendConsentError(consentCode, err) {
  return {
    type: OAUTH2_SEND_CONSENT_ERROR,
    consentCode,
    err
  };
}

export function sentConsent(consentCode, redirectTo) {
  return {
    type: OAUTH2_SENT_CONSENT,
    consentCode,
    redirectTo
  };
}

export function getAuthorizedApps(id) {
  return {
    type: OAUTH2_GET_AUTHORIZED_APPS,
    id
  };
}

export function getAuthorizedAppsError(err) {
  return {
    type: OAUTH2_GET_AUTHORIZED_APPS_ERROR,
    err
  };
}

export function gotAuthorizedApps(authorizedApps) {
  return {
    type: OAUTH2_GOT_AUTHORIZED_APPS,
    authorizedApps
  };
}

export function revokeGrant(id, grantId) {
  return {
    type: OAUTH2_REVOKE_GRANT,
    id,
    grantId
  };
}

export function revokeGrantError(grantId, err) {
  return {
    type: OAUTH2_REVOKE_GRANT_ERROR,
    grantId,
    err
  };
}

export function revokedGrant(grantId) {
  return {
    type: OAUTH2_REVOKED_GRANT,
    grantId
  };
}
