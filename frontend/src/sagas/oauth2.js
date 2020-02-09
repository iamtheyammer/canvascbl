import { put, takeLeading, takeEvery, all } from 'redux-saga/effects';
import {
  OAUTH2_GET_CONSENT_INFO,
  OAUTH2_SEND_CONSENT,
  getConsentInfoError,
  gotConsentInfo,
  sendConsentError,
  sentConsent,
  gotAuthorizedApps,
  getAuthorizedAppsError,
  OAUTH2_GET_AUTHORIZED_APPS,
  revokedGrant,
  revokeGrantError,
  OAUTH2_REVOKE_GRANT
} from '../actions/oauth2';
import makeOauth2Request from '../util/oauth2/makeOauth2Request';
import { endLoading, startLoading } from '../actions/loading';

function* getConsentInfo({ id, consentCode }) {
  yield put(startLoading(id));
  try {
    const consentInfoReq = yield makeOauth2Request('consent', {
      consent_code: consentCode
    });
    yield put(gotConsentInfo(consentCode, consentInfoReq.data));
  } catch (e) {
    console.log(e.response.data);
    yield put(
      getConsentInfoError(
        consentCode,
        e.response ? e.response.data : "can't connect"
      )
    );
  }
  yield put(endLoading(id));
}

function* sendConsent({ id, consentCode, action }) {
  yield put(startLoading(id));
  try {
    const consentReq = yield makeOauth2Request(
      'consent',
      { action, consent_code: consentCode },
      'put'
    );
    yield put(sentConsent(consentCode, consentReq.data.redirect_to));
  } catch (e) {
    yield put(
      sendConsentError(
        consentCode,
        e.response ? e.response.data : "can't connect"
      )
    );
  }
  yield put(endLoading(id));
}

function* getAuthorizedApps({ id }) {
  yield put(startLoading(id));
  try {
    const appsReq = yield makeOauth2Request('tokens');
    yield put(gotAuthorizedApps(appsReq.data));
  } catch (e) {
    yield put(
      getAuthorizedAppsError(e.response ? e.response.data : "can't connect")
    );
  }
  yield put(endLoading(id));
}

function* revokeGrant({ id, grantId }) {
  yield put(startLoading(id));
  try {
    yield makeOauth2Request('token', { token_id: grantId }, 'delete');
    yield put(revokedGrant(grantId));
  } catch (e) {
    yield put(
      revokeGrantError(grantId, e.response ? e.response.data : "can't connect")
    );
  }
  yield put(endLoading(id));
}

function* watcher() {
  yield takeLeading(OAUTH2_GET_CONSENT_INFO, getConsentInfo);
  yield takeLeading(OAUTH2_SEND_CONSENT, sendConsent);
  yield takeLeading(OAUTH2_GET_AUTHORIZED_APPS, getAuthorizedApps);
  yield takeEvery(OAUTH2_REVOKE_GRANT, revokeGrant);
}

export default function* oauth2RootSaga() {
  yield all([watcher()]);
}
